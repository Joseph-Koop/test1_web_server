// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go
// go run main.go -target=example.com -start=80 -end=90 -timeout=300ms


package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Target string `json:"target"`
	Port int `json:"port"`
	Banner string `json:"banner"`
}

type ScanSummary struct {
	TotalPorts int `json:"total"`
	OpenPorts int `json:"ports"`
	TimeTaken float64 `json:"time"`
}

type ScanReport struct {
	Summary ScanSummary `json:"summary"`
	Results []ScanResult `json:"results"`
}

var openPorts int = 0
var totalPorts int = 0
var mutex sync.Mutex		//implemented so that openPorts and total Ports update safely even though asynchronous functions are involved
var successfulPorts []ScanResult

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer, portCount int) {
	defer wg.Done()
	maxRetries := 3
	for addr := range tasks {
		success := false
		target, portStr, _ := net.SplitHostPort(addr)
		port, _ := strconv.Atoi(portStr)
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
				conn.Close()

				mutex.Lock()
				openPorts += 1
				totalPorts += 1
				successfulPorts = append(successfulPorts, ScanResult{Target: target, Port: port, Banner: banner})
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
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")

	flag.Parse()

	targetList := strings.Split(*target, ",")

	//portCount := *endPort - *startPort + 1
	var portsToScan []int
	if *ports != "" {
		portList := strings.Split(*ports, ",")
		for _, port := range portList {
			portInt, err := strconv.Atoi(port)
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

	if *jsonOutput {
		report := ScanReport{
			Results: successfulPorts,
			Summary: ScanSummary{
				TotalPorts: totalPorts,
				OpenPorts: openPorts,
				TimeTaken: time.Since(start).Seconds(),
			},
		}
		jsonData, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		fmt.Print("\n\n\n\nRESULTS\n")
		fmt.Println(string(jsonData))
	}else{
		fmt.Print("\n\n\n\nOPEN PORTS\n")
		for _, result := range successfulPorts {
			fmt.Printf("\nTarget: %s\nPort: %d\nBanner: %s\n", result.Target, result.Port, result.Banner)
		}
		fmt.Print("\n\nSUMMARY\n")
		fmt.Printf("\nTotal ports scanned: %d \nNumber of open ports: %d \nTime taken: %.2f seconds\n\n", totalPorts, openPorts, time.Since(start).Seconds())
	}
}
