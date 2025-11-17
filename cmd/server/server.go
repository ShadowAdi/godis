package main

import (
	"fmt"
	"godis/helper"
	"godis/internals/aof"
	"godis/internals/handlers"
	"godis/internals/resp"
	"net"
	"strings"
)

func main() {
	// Print banner and server info
	helper.PrintBanner()
	helper.PrintServerInfo(":6380")

	helper.LogInfo("Starting GODIS server...")

	l, err := net.Listen("tcp", ":6380")
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to start listener: %v", err))
		return
	}
	helper.LogInfo("TCP listener started successfully")

	aof, err := aof.NewAoF("database.aof")
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to initialize AOF: %v", err))
		return
	}
	defer aof.Close()
	helper.LogInfo("AOF initialized successfully")

	// Load data from AOF
	helper.LogInfo("Loading data from AOF file...")
	commandCount := 0
	aof.Read(func(value resp.Value) {
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]
		handler, ok := handlers.Handlers[command]
		if !ok {
			helper.LogError(fmt.Sprintf("Invalid command in AOF: %s", command))
			return
		}
		handler(args)
		commandCount++
	})
	helper.LogInfo(fmt.Sprintf("Loaded %d commands from AOF", commandCount))

	helper.LogInfo("Waiting for client connections...")

	conn, err := l.Accept()
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to accept connection: %v", err))
		return
	}
	defer conn.Close()

	helper.LogConnection(fmt.Sprintf("New client connected: %s", conn.RemoteAddr().String()))

	for {
		decoder := resp.NewResp(conn)

		value, err := decoder.Read()
		if err != nil {
			helper.LogError(fmt.Sprintf("Connection error: %v", err))
			helper.LogConnection("Client disconnected")
			return
		}

		if value.Typ != "array" {
			helper.LogError("Invalid request: expected array")
			continue
		}

		if len(value.Array) == 0 {
			helper.LogError("Invalid request: array length = 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		helper.LogCommand(command, len(args))

		writer := resp.NewWriter(conn)

		handler, ok := handlers.Handlers[command]
		if !ok {
			helper.LogError(fmt.Sprintf("Unknown command: %s", command))
			writer.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
