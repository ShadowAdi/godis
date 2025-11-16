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
		decoder := resp.NewResp(conn)

		value, err := decoder.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("CLIENT SENT:", value)

		writer := resp.NewWriter(conn)

		writer.Write(resp.Value{Typ: "string", Str: "OK"})
	}
}
