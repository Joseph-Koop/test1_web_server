<!-- README.md -->

Port Scanner
Joseph Koop
Systems Programming and Computer Organization
Test #1
March 30, 2024
________________________________________________________________________________________________________________________________________

Overview:

This project is a basic port scanner. The scanner reads a list of targets and ports and attempts to connect to each target on every port provided.
The program uses channels, go routines, and waitgroups to scan the ports concurrently.
After the scan, the total number of ports, the number of open ports, and the time taken is output, along with banner information for open ports.
Command line flags such as target, ports, workers, timeout, and json make the program a lot more flexible for the user.
________________________________________________________________________________________________________________________________________

How to Use:

You can run this program in your IDE using this base command:
    go run main.go 

If you want to alter the default options, add flags after the command. Change the provided examples with your requests:
    Target:
        -target=example.com
        -target "example.com,scanme.nmap.org"
    Ports:
        -ports=80
        -ports "20,22,80,443"
    Start Port:
        -start=78
    End Port:
        -end=82
    Workers:
        -workers=200
    Timeout:
        -timeout=1000ms
    JSON Output:
        -json=true

Alternatively, you can build and then run:
    go build main.go
Then type ./main followed by any command line flags. You can do this as many times as desired without repeating go build main.go.
________________________________________________________________________________________________________________________________________

Example #1:

go run main.go -target "scanme.nmap.org,example.com" -ports "20,22,80,443" -workers=200 -timeout=1000ms -json

Scanning port scanme.nmap.org:20 out of 4 total ...
Scanning port example.com:22 out of 4 total ...
Scanning port scanme.nmap.org:22 out of 4 total ...
Scanning port example.com:20 out of 4 total ...
Scanning port scanme.nmap.org:80 out of 4 total ...
Scanning port example.com:80 out of 4 total ...
Scanning port example.com:443 out of 4 total ...
Scanning port scanme.nmap.org:443 out of 4 total ...



RESULTS
{
 "summary": {
  "total": 8,
  "ports": 4,
  "time": 10.003285617
 },
 "results": [
  {
   "target": "example.com",
   "port": 443,
   "banner": "No banner received"
  },
  {
   "target": "example.com",
   "port": 80,
   "banner": "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nETag: \"84238dfc8092e5d9c0dac8ef93371a07:1736799080.121134\"\r\nLast-Modified: Mon, 13 Jan 2025 20:11:20 GMT\r\nCache-Control: max-age=2343\r\nDate: Mon, 31 Mar 2025 05:49:48 GMT\r\nContent-Length: 1256\r\nConnection: close\r\nX-N: S"
  },
  {
   "target": "scanme.nmap.org",
   "port": 80,
   "banner": "HTTP/1.1 200 OK\r\nDate: Mon, 31 Mar 2025 05:49:48 GMT\r\nServer: Apache/2.4.7 (Ubuntu)\r\nAccept-Ranges: bytes\r\nVary: Accept-Encoding\r\nConnection: close\r\nTransfer-Encoding: chunked\r\nContent-Type: text/html"
  },
  {
   "target": "scanme.nmap.org",
   "port": 22,
   "banner": "SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13\r\n"
  }
 ]
}
________________________________________________________________________________________________________________________________________

Example #2

go run main.go -target=scanme.nmap.org -start=78 end=82 -workers=200 -timeout=1000ms

Scanning port scanme.nmap.org:78 out of 5 total ...
Scanning port scanme.nmap.org:79 out of 5 total ...
Scanning port scanme.nmap.org:81 out of 5 total ...
Scanning port scanme.nmap.org:82 out of 5 total ...
Scanning port scanme.nmap.org:80 out of 5 total ...



OPEN PORTS

Target: scanme.nmap.org
Port: 80
Banner: HTTP/1.1 200 OK
Date: Mon, 31 Mar 2025 05:53:44 GMT
Server: Apache/2.4.7 (Ubuntu)
Accept-Ranges: bytes
Vary: Accept-Encoding
Connection: close
Transfer-Encoding: chunked
Content-Type: text/html


SUMMARY

Total ports scanned: 5 
Number of open ports: 1 
Time taken: 7.20 seconds