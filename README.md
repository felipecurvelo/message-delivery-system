# message-delivery-system

This is a simplified message delivery system built in Go that runs over a command line terminal (Tested with bash/zsh). 
The main goal of this project was to exercise network comunication using Go.

## Starting a server
Run the binary passing the port as argument:

```
go run ./cmd/server 1234
```

## Connecting a client to a server
Run the binary passing the ip and port as arguments:

```
go run ./cmd/client 127.0.0.1 1234
```

## Listing connected users
From the client command line, send the `list` command and the server should answer with the list of connected users.

## Getting client identity
From the client command line, send the `whoami` command and the server should answer with the connected user id.

## Relaying messages to other users
From the client command line, send the `ids|msg` command using the `|` separator as shown, where ids could be a single id or a comma separated list of destination ids and the message should be a string.

## Required Improvements List

* Improve message contract introducing json format.
* Test communication over local network.
* Test with windows based terminal (MS-DOS, Powershell).