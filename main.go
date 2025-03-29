// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var openPorts int = 0
var totalPorts int = 0
var mutex sync.Mutex

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer) {
	defer wg.Done()
	maxRetries := 3
    for addr := range tasks {
		var success bool
		for i := range maxRetries {      
			conn, err := dialer.Dial("tcp", addr)
			if err == nil {
				conn.Close()
				fmt.Printf("Connection to %s was successful.\n", addr)

				mutex.Lock()
				openPorts += 1
				totalPorts += 1
				mutex.Unlock()

				success = true
				break
			}
			backoff := time.Duration(1<<i) * time.Second
			fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
			time.Sleep(backoff)
	    }
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts.\n", addr, maxRetries)

			mutex.Lock()
			totalPorts += 1
			mutex.Unlock()
		}
	}
}

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	target := flag.String("target", "scanme.nmap.org", "Target host to scan")
	startPort := flag.Int("start", 1, "Start of Port Range")
	endPort := flag.Int("end", 1024, "End of Port Range")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Duration("timeout", 500 * time.Millisecond, "Connection Timeout")

	flag.Parse()

	dialer := net.Dialer {
		Timeout: *timeout,
	}

	start := time.Now()

    for i := 1; i <= *workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, dialer)
	}

	for p := *startPort; p <= *endPort; p++ {
		port := strconv.Itoa(p)
        address := net.JoinHostPort(*target, port)
		tasks <- address
	}
	close(tasks)

	wg.Wait()

	fmt.Printf("\nTotal ports scanned: %d \nNumber of open ports: %d \nTime taken: %.2f seconds\n\n", totalPorts, openPorts, time.Since(start).Seconds())
}