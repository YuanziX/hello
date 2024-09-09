# Hello - Chatroom Server

This project implements a simple chatroom using Go and gRPC. Users can join chatrooms, send messages to other users in the room, and disconnect from the chatroom when done. When all users leave a room, the room is deleted.

## Features

-   **Join Chatroom**: Users can join a chatroom by providing their user ID and the room ID.
-   **Send Messages**: Users in a chatroom can send messages to all other users in the room.
-   **Auto-Disconnect**: When a user disconnects from the chatroom, they are automatically removed from the room.
-   **Room Deletion**: If all users leave a chatroom, the room is deleted.

## Architecture

The chat server uses gRPC for communication between clients and the server. A `Server` struct is used to manage active chatrooms and handle user connections, while each user connection is represented by the `Connection` struct.

-   `Server`:
    -   Manages multiple chatrooms and user connections.
    -   Ensures thread-safe access to chatrooms using a read/write mutex.
-   `Connection`:
    -   Represents a user in a chatroom, including the gRPC stream for sending messages and an error channel for managing disconnections.

## gRPC Services

### `Join`

Users can join a chatroom by sending a `JoinRequest` to the `Join` method, which takes the user's ID and the room ID. The server adds the user to the room and starts listening for messages. When the user disconnects, they are automatically removed from the room.

### `SendMessage`

Users in a chatroom can send messages to other users in the same room. The server handles broadcasting the message to all users connected to that chatroom. If a chatroom does not exist, an error is returned.

## Usage

### Clone the Repository

```bash
git clone https://github.com/YuanziX/hello.git
cd hello
```

### Build and Run

Ensure you have [Go](https://golang.org/doc/install) installed.

```bash
go mod tidy
go run main.go
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
