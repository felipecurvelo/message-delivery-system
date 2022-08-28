# message-delivery-system

This is a simplified message delivery system.

## Starting a server

```go
import "github.com/felipecurvelo/message-delivery-system/internal/server"

port := "1234"

s := server.NewServer()
s.Start(port)

defer s.Close()
```

## Connecting a client

```go
import "github.com/felipecurvelo/message-delivery-system/internal/client"

client := client.NewClient()
err := client.Connect("127.0.0.1:1234")
if err != nil {
	//Handle the error
}
defer client.Close()

clientOutChan := make(chan string)
go client.HandleMessages(clientOutChan)

//Listen to channel for incoming messages
msg := <-clientOutChan
```

## Listing connected users ids
Calling the List() method will make the server respond with a list of connected users

```go
client.List()
```

## Getting client identity
Calling the WhoAmI() method will make the server respond with the client id

```go
client.WhoAmI()
```

## Relaying messages to other users
Calling the SendMessage() method will send a message to a list of recipients. 

Destination ids should be a list of one or more ids separated by commas. The list supports up to 255 destination ids.

Message should be a byte array with a size up to 1024Kb.

```go
client.SendMessage("1,2,3", []byte("message body"))
```

## Required Improvements List

* Add tests covering a higher number of concurrent clients
* Improve input validation
* Implement client and server command-line interface