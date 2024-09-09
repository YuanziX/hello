package server

import (
	"context"
	"log"
	"sync"

	pb "github.com/yuanzix/hello/gen"
)

type Connection struct {
	pb.UnimplementedChatRoomServer
	stream pb.ChatRoom_JoinServer

	id    string
	error chan error
}

type Pool struct {
	pb.UnimplementedChatRoomServer
	Connection []*Connection
}

func (p *Pool) Join(jReq *pb.JoinRequest, stream pb.ChatRoom_JoinServer) error {
	conn := Connection{
		stream: stream,
		id:     jReq.User.Id,
		error:  make(chan error),
	}

	p.Connection = append(p.Connection, &conn)
	log.Printf("%v joined the chat.\n", conn.id)

	err := <-conn.error

	for i, curr := range p.Connection {
		if curr.id == conn.id {
			p.Connection = append(p.Connection[:i], p.Connection[i+1:]...)
		}
	}
	return err
}

func (p *Pool) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Exit, error) {
	wait := sync.WaitGroup{}

	done := make(chan int)

	for _, conn := range p.Connection {
		wait.Add(1)

		go func(msg *pb.Message, conn *Connection) {
			defer wait.Done()

			err := conn.stream.Send(msg)
			log.Printf("Sending message to: %v from %v\n", conn.id, msg.Id)

			if err != nil {
				log.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
				conn.error <- err
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &pb.Exit{}, nil
}
