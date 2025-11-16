package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:6380")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')

		conn.Write([]byte(cmd))

		reply := make([]byte, 1024)
		n, _ := conn.Read(reply)
		fmt.Println(string(reply[:n]))
	}
}
