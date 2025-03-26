package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"sync"
)

func main(){
	target := "scanme.nmap.org"
	dialer := net.Dialer{
		Timeout: 5 * time.Second,
	}
	maxRetries := 3
	port := 1
	maxPort := 90
	var wg sync.WaitGroup
	
	for i := port; i <= maxPort; i++ {
		portStr := strconv.Itoa(i)
		address := net.JoinHostPort(target, portStr)
		wg.Add(1)
		go func (address string){
			for j := 0; j < maxRetries; j++ {
				conn, err := dialer.Dial("tcp", address)
				if err == nil {
					defer conn.Close()
					fmt.Printf("Connection to %s was successful\n", address)
					break
				}
				backoff := time.Duration(1<<j) * time.Second
				fmt.Printf("Port %d Attempt %d failed. Waiting %v...\n", i, j+1, backoff)
				time.Sleep(backoff)
			}
			defer wg.Done()
		}(address)
	}
	wg.Wait()
}

	
	// if err != nil{
	// 	log.Fatalf("Unable to connect to %s: %v", address, err)
	// }
	// defer conn.Close()
	// fmt.Printf("Connection to %s was successful\n", address)
	
	// for port; port < 513; port++ {
	// 	for i := 0; i < maxRetries; i++ {