// Command worker prints the oxzoo greeting line to stdout every 10 seconds,
// forever. ox runs it as a systemd process with Restart=always and reads its
// journal to verify output; it opens no network socket.
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	tag := os.Getenv("GREETING_TAG")
	if tag == "" {
		fmt.Fprintln(os.Stderr, "GREETING_TAG must be set; start the worker with one, for example: GREETING_TAG=dev ./worker")
		os.Exit(1)
	}
	for {
		fmt.Printf("hello world oxzoo-worker-go_%s\n", tag)
		time.Sleep(10 * time.Second)
	}
}
