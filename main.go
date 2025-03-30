// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go
// go run main.go -target=example.com -start=80 -end=90 -timeout=300ms


package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var openPorts int = 0
var totalPorts int = 0
var mutex sync.Mutex		//implemented so that openPorts and total Ports update safely even though asynchronous functions are involved
var successfulPorts map[string][]string

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer, portCount int) {
	defer wg.Done()
	maxRetries := 3
	for addr := range tasks {
		var success bool
		fmt.Printf("\nScanning port %s out of %d total ...", addr, portCount)
		for i := range maxRetries {
			conn, err := dialer.Dial("tcp", addr)
			if err == nil {

				httpRequest := "GET / HTTP/1.1\r\nHost: " + addr + "\r\nConnection: close\r\n\r\n"		//a request was needed to get a response otherwise reading a successful port blocked
				_, err = conn.Write([]byte(httpRequest))

				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)

				var banner string				//saving success data so i can display it after everything else is cleared
				if err == nil {
					response := string(buffer[:n])
					headersEnd := strings.Index(response, "\r\n\r\n")		//I only wanted the basic information, not any of the html
					//fmt.Printf("\nConnection to %s was successful.", addr)		//This line should not get seperated from the main banner body
					if headersEnd != -1 {
						banner = response[:headersEnd]
						//fmt.Println("\nBanner: ", response[:headersEnd])
					} else {
						banner = response
						//fmt.Println("\nBanner: ", response)
					}

				} else {
					banner = "No banner received"
					//fmt.Printf("\nConnection to %s was successful.", addr)
					//fmt.Println("\nBanner Error: ", err)
				}
				fmt.Println()
				conn.Close()

				mutex.Lock()
				openPorts += 1
				totalPorts += 1
				successfulPorts[addr] = append(successfulPorts[addr], fmt.Sprintf("%s - %s", addr, banner))
				mutex.Unlock()

				success = true
				break
			}
			backoff := time.Duration(1<<i) * time.Second
			//fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
			time.Sleep(backoff)
		}
		if !success {
			//fmt.Printf("Failed to connect to %s after %d attempts.\n", addr, maxRetries)

			mutex.Lock()
			totalPorts += 1
			mutex.Unlock()
		}
	}
}

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	target := flag.String("target", "scanme.nmap.org", "Target hosts to scan")
	ports := flag.String("ports", "", "List of specific ports")
	startPort := flag.Int("start", 1, "Start of Port Range")
	endPort := flag.Int("end", 1024, "End of Port Range")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Duration("timeout", 500*time.Millisecond, "Connection Timeout")

	flag.Parse()

	targetList := strings.Split(*target, ",")
	successfulPorts = make(map[string][]string)

	//portCount := *endPort - *startPort + 1
	var portsToScan []int
	if *ports != "" {
		portList := strings.Split(*ports, ",")
		for _, port := range portList {
			portInt, err := strconv.Atoi(strings.TrimSpace(port))
			if err == nil {
				portsToScan = append(portsToScan, portInt)
			}
		}
	} else {
		for port := *startPort; port <= *endPort; port++ {
			portsToScan = append(portsToScan, port)
		}
	}

	portCount := len(portsToScan)

	dialer := net.Dialer{
		Timeout: *timeout,
	}

	start := time.Now()

	for i := 1; i <= *workers; i++ {
		wg.Add(1)
		go worker(&wg, tasks, dialer, portCount)
	}

	for _, indTarget := range targetList{
		for _, port := range portsToScan {
			address := net.JoinHostPort(indTarget, strconv.Itoa(port))
			tasks <- address
		}
	}
	close(tasks)
	wg.Wait()

	fmt.Print("\n\n\nOPEN PORTS\n")
	for target, ports := range successfulPorts {
		fmt.Printf("\nResults for: %s\n", target)
		for _, result := range ports{
			fmt.Println(result)
		}
	}
	fmt.Print("\n\n\nSUMMARY\n")
	fmt.Printf("\nTotal ports scanned: %d \nNumber of open ports: %d \nTime taken: %.2f seconds\n\n", totalPorts, openPorts, time.Since(start).Seconds())
}
