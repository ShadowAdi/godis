package helper

import (
	"fmt"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBold   = "\033[1m"
)

// PrintBanner prints the GODIS ASCII art banner
func PrintBanner() {
	banner := `
   ██████╗  ██████╗ ██████╗ ██╗███████╗
  ██╔════╝ ██╔═══██╗██╔══██╗██║██╔════╝
  ██║  ███╗██║   ██║██║  ██║██║███████╗
  ██║   ██║██║   ██║██║  ██║██║╚════██║
  ╚██████╔╝╚██████╔╝██████╔╝██║███████║
   ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝╚══════╝
`
	fmt.Println(ColorBlue + ColorBold + banner + ColorReset)
}

// PrintServerInfo prints server startup information
func PrintServerInfo(port string) {
	fmt.Println(ColorCyan + "════════════════════════════════════════════" + ColorReset)
	fmt.Printf(ColorGreen + "  ✓ Server Status: " + ColorBold + "RUNNING" + ColorReset + "\n")
	fmt.Printf(ColorGreen+"  ✓ Port: "+ColorBold+"%s"+ColorReset+"\n", port)
	fmt.Printf(ColorGreen+"  ✓ Time: "+ColorBold+"%s"+ColorReset+"\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(ColorGreen + "  ✓ AOF: " + ColorBold + "Enabled (database.aof)" + ColorReset + "\n")
	fmt.Println(ColorCyan + "════════════════════════════════════════════" + ColorReset)
	fmt.Println()
}

// LogInfo prints an info log message
func LogInfo(message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(ColorCyan+"[%s]"+ColorReset+" "+ColorGreen+"INFO:"+ColorReset+" %s\n", timestamp, message)
}

// LogError prints an error log message
func LogError(message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(ColorCyan+"[%s]"+ColorReset+" "+"\033[31m"+"ERROR:"+ColorReset+" %s\n", timestamp, message)
}

// LogCommand prints a command execution log
func LogCommand(command string, args int) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(ColorCyan+"[%s]"+ColorReset+" "+ColorYellow+"CMD:"+ColorReset+" %s (args: %d)\n", timestamp, command, args)
}

// LogConnection prints connection events
func LogConnection(message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf(ColorCyan+"[%s]"+ColorReset+" "+ColorBlue+"CONN:"+ColorReset+" %s\n", timestamp, message)
}
