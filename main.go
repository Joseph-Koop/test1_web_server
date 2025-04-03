// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go
// go run main.go -target=example.com -start=80 -end=90 -timeout=300ms

package main

import ( //Required packages for the program
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanReport struct { //These structs are used to store the results in JSON format
	Summary ScanSummary  `json:"summary"`
	Results []ScanResult `json:"results"`
}

type ScanResult struct { //Holds results of successful ports
	Target string `json:"target"`
	Port   int    `json:"port"`
	Banner string `json:"banner"`
}

type ScanSummary struct { //Holds summary of ports scanned
	TotalPorts int     `json:"total"`
	OpenPorts  int     `json:"ports"`
	TimeTaken  float64 `json:"time"`
}

var openPorts int = 0
var totalPorts int = 0
var mutex sync.Mutex             //Challenge: Implemented so that openPorts and total Ports update correctly even though asynchronous functions are involved
var successfulPorts []ScanResult //Needed to display ports after scanning is completed

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer, portCount int) {
	defer wg.Done()
	maxRetries := 3           //The number of times the function will try to connect to a port
	for addr := range tasks { //Each routine keeps going until the tasks channel is depleted
		success := false                              //Variable to store if connection succeeded
		target, portStr, _ := net.SplitHostPort(addr) //Unmerges target and port
		port, _ := strconv.Atoi(portStr)
		fmt.Printf("\nScanning port %s out of %d total ...", addr, portCount) //Progress Indicator
		for i := range maxRetries {                                           //Makes as many attempts as specified
			conn, err := dialer.Dial("tcp", addr) //Tries to make TCP connection
			if err == nil {

				httpRequest := "GET / HTTP/1.1\r\nHost: " + addr + "\r\nConnection: close\r\n\r\n" //Challenge: a request was needed to get a response otherwise reading a successful port blocked
				_, err = conn.Write([]byte(httpRequest))

				buffer := make([]byte, 1024) //Captures only hte first 1024 bytes of the response
				n, err := conn.Read(buffer)

				var banner string //Saving success data so i can display it after everything else is cleared
				if err == nil {
					response := string(buffer[:n])
					headersEnd := strings.Index(response, "\r\n\r\n") //Challenge: I only wanted the basic information, not any of the html
					if headersEnd != -1 {
						banner = response[:headersEnd]
					} else {
						banner = response
					}

				} else {
					banner = "No banner received"
				}
				conn.Close() //Closes connection

				mutex.Lock() //Locks so that openPorts, totalPorts, and successfulPorts can be updated correctly
				openPorts += 1
				totalPorts += 1
				successfulPorts = append(successfulPorts, ScanResult{Target: target, Port: port, Banner: banner})
				mutex.Unlock()

				success = true
				break
			}
			backoff := time.Duration(1<<i) * time.Second //Backoff increases by the power of 2: left shift
			time.Sleep(backoff)                          //Routine waites until backoff period is complete
		}
		if !success {

			mutex.Lock() //TotalPorts still gets updated
			totalPorts += 1
			mutex.Unlock()
		}
	}
}

func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100) //This channel will hold the port addresses the workers use

	target := flag.String("target", "scanme.nmap.org", "Target hosts to scan")      //Which sites the scanner will try to connect to
	ports := flag.String("ports", "", "List of specific ports")                     //An optional flag for giving the scanner specific ports
	startPort := flag.Int("start", 1, "Start of Port Range")                        //The first port if scanning a range of ports
	endPort := flag.Int("end", 1024, "End of Port Range")                           //The last port if scanning a range of ports
	workers := flag.Int("workers", 100, "Number of concurrent workers")             //The number of go routines that will be employed at a time
	timeout := flag.Duration("timeout", 500*time.Millisecond, "Connection Timeout") //The length of time before an attempt to connect will be considered failed
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")         //The option to output results in JSON format

	flag.Parse()

	targetList := strings.Split(*target, ",") //Splits provided targets into a list

	var portsToScan []int //Creates array of ports to scan
	if *ports != "" {     //If specific ports were provided, it converts them into an array and uses them only
		portList := strings.Split(*ports, ",")
		for _, port := range portList {
			portInt, err := strconv.Atoi(port)
			if err == nil {
				portsToScan = append(portsToScan, portInt)
			}
		}
	} else {
		for port := *startPort; port <= *endPort; port++ { //If no specific ports were provided, it uses the startPort and endPort to determine the range
			portsToScan = append(portsToScan, port)
		}
	}

	portCount := len(portsToScan) //Number of ports

	dialer := net.Dialer{ //Specific tool in go for dialing
		Timeout: *timeout, //Sets the timeout equal to the flag
	}

	start := time.Now() //Starts a time tracker used in summary

	for i := 1; i <= *workers; i++ { //Calls the number of worker functions equal to the workers flag
		wg.Add(1)
		go worker(&wg, tasks, dialer, portCount)
	}

	for _, indTarget := range targetList { //Writes every port of every target into the tasks channel, making them available to the go routines
		for _, port := range portsToScan {
			address := net.JoinHostPort(indTarget, strconv.Itoa(port)) //Merges target and port before sending them to the channel
			tasks <- address
		}
	}
	close(tasks)
	wg.Wait()

	if *jsonOutput { //If JSON output was requested, results are printed based on the structs' structure
		report := ScanReport{
			Results: successfulPorts,
			Summary: ScanSummary{
				TotalPorts: totalPorts,
				OpenPorts:  openPorts,
				TimeTaken:  time.Since(start).Seconds(),
			},
		}
		jsonData, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		fmt.Print("\n\n\n\nRESULTS\n")
		fmt.Println(string(jsonData))
	} else { //If JSON output was not requested, results are printed in basic format
		fmt.Print("\n\n\n\nOPEN PORTS\n")
		for _, result := range successfulPorts {
			fmt.Printf("\nTarget: %s\nPort: %d\nBanner: %s\n", result.Target, result.Port, result.Banner)
		}
		fmt.Print("\n\nSUMMARY\n")
		fmt.Printf("\nTotal ports scanned: %d \nNumber of open ports: %d \nTime taken: %.2f seconds\n\n", totalPorts, openPorts, time.Since(start).Seconds())
	}
}
