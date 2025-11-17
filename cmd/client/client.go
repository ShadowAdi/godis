package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func toRESP(input string) string {
	parts := strings.Fields(input)
	resp := fmt.Sprintf("*%d\r\n", len(parts))

	for _, p := range parts {
		resp += fmt.Sprintf("$%d\r\n%s\r\n", len(p), p)
	}

	return resp
}

func main() {
	conn, _ := net.Dial("tcp", "localhost:6380")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		resp := toRESP(cmd)
		conn.Write([]byte(resp))

		reply := make([]byte, 1024)
		n, _ := conn.Read(reply)
		fmt.Println(string(reply[:n]))
	}
}
