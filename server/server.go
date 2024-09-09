package server

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/yuanzix/hello/gen"
)

type Connection struct {
	pb.UnimplementedChatRoomServer
	stream pb.ChatRoom_JoinServer

	id     string
	roomID string
	error  chan error
}

type Server struct {
	pb.UnimplementedChatRoomServer
	Chatrooms map[string][]*Connection

	mu sync.RWMutex
}

func (p *Server) Join(jReq *pb.JoinRequest, stream pb.ChatRoom_JoinServer) error {
	conn := Connection{
		stream: stream,
		id:     jReq.User.Id,
		roomID: jReq.RoomID,
		error:  make(chan error),
	}

	p.mu.Lock()
	_, ok := p.Chatrooms[conn.roomID]
	if !ok {
		p.Chatrooms[conn.roomID] = []*Connection{&conn}
	} else {
		p.Chatrooms[conn.roomID] = append(p.Chatrooms[conn.roomID], &conn)

	}
	p.mu.Unlock()

	log.Printf("%v joined chatroom %v.\n", conn.id, conn.roomID)

	go func() {
		<-stream.Context().Done()
		log.Printf("User %v disconnected from chatroom %v.\n", conn.id, conn.roomID)
		conn.error <- stream.Context().Err()
	}()

	err := <-conn.error

	p.mu.Lock()
	for i, curr := range p.Chatrooms[conn.roomID] {
		if curr.id == conn.id {
			p.Chatrooms[conn.roomID] = append(p.Chatrooms[conn.roomID][:i], p.Chatrooms[conn.roomID][i+1:]...)
		}
		log.Printf("%v left chatroom %v.\n", conn.id, conn.roomID)
	}

	if len(p.Chatrooms[conn.roomID]) == 0 {
		delete(p.Chatrooms, conn.roomID)
		log.Printf("chatroom %v deleted because it was empty.\n", conn.roomID)
	}

	p.mu.Unlock()

	return err
}

func (p *Server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Exit, error) {
	if _, ok := p.Chatrooms[msg.RoomID]; !ok {
		return &pb.Exit{
			Success: false,
			Message: fmt.Sprintf("chatroom with id %v does not exist.\n", msg.RoomID),
		}, nil
	}

	wait := sync.WaitGroup{}

	done := make(chan int)

	p.mu.RLock()
	for _, conn := range p.Chatrooms[msg.RoomID] {
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
	p.mu.RUnlock()

	return &pb.Exit{
		Success: true,
	}, nil
}
