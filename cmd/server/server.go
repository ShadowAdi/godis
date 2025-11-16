package main

import (
	"fmt"
	"godis/internals/resp"
	"net"
)

func main() {
	fmt.Println("Listening on port :6380")

	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := resp.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}
