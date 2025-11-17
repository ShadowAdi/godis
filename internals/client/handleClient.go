package client

import (
	"fmt"
	"godis/helper"
	"godis/internals/aof"
	"godis/internals/handlers"
	"godis/internals/resp"
	"net"
	"strings"
)

func HandleClient(conn net.Conn, aof *aof.Aof) {
	defer conn.Close()

	for {
		decoder := resp.NewResp(conn)
		value, err := decoder.Read()
		if err != nil {
			helper.LogConnection("Client disconnected")
			return
		}

		if value.Typ != "array" || len(value.Array) == 0 {
			helper.LogError("Invalid request")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handlers.Handlers[command]
		if !ok {
			helper.LogError(fmt.Sprintf("Unknown command: %s", command))
			resp.NewWriter(conn).Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		// AOF write
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		resp.NewWriter(conn).Write(result)
	}
}
