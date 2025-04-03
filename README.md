<!-- README.md -->

<strong>Port Scanner</strong><br>
<strong>Joseph Koop</strong><br>
<strong>Systems Programming and Computer Organization</strong><br>
<strong>Test #1</strong><br>
<strong>March 30, 2024</strong><br>

<br>
<strong>Link to video:</strong>
https://youtu.be/LXhA6hgjZbo
________________________________________________________________________________________________________________________________________

Overview:<br>
<br>
This project is a basic port scanner. The scanner reads a list of targets and ports and attempts to connect to each target on every port provided.
The program uses channels, go routines, and waitgroups to scan the ports concurrently.
After the scan, the total number of ports, the number of open ports, and the time taken is output, along with banner information for open ports.
Command line flags such as target, ports, workers, timeout, and json make the program a lot more flexible for the user.
________________________________________________________________________________________________________________________________________

How to Use:<br>

You can run this program in your IDE using this base command:<br>
>    go run main.go <br>

If you want to alter the default options, add flags after the command. Change the provided examples with your requests:<br>
>    Target:<br>
>>       -target=example.com<br>
>>       -target "example.com,scanme.nmap.org"<br>
>    Ports:<br>
>>        -ports=80<br>
>>        -ports "20,22,80,443"<br>
>    Start Port:<br>
>>        -start=78<br>
>    End Port:<br>
>>        -end=82<br>
>    Workers:<br>
>>        -workers=200<br>
>    Timeout:<br>
>>        -timeout=1000ms<br>
>    JSON Output:<br>
>>        -json=true<br>

Alternatively, you can build and then run:<br>
>    go build main.go<br>
Then type ./main followed by any command line flags. You can do this as many times as desired without repeating go build main.go.
________________________________________________________________________________________________________________________________________

Example #1:<br>
<br>
go run main.go -target "scanme.nmap.org,example.com" -ports "20,22,80,443" -workers=200 -timeout=1000ms -json<br>
<br>
Scanning port scanme.nmap.org:20 out of 4 total ...<br>
Scanning port example.com:22 out of 4 total ...<br>
Scanning port scanme.nmap.org:22 out of 4 total ...<br>
Scanning port example.com:20 out of 4 total ...<br>
Scanning port scanme.nmap.org:80 out of 4 total ...<br>
Scanning port example.com:80 out of 4 total ...<br>
Scanning port example.com:443 out of 4 total ...<br>
Scanning port scanme.nmap.org:443 out of 4 total ...<br>
<br>
<br>
<br>
RESULTS<br>
{<br>
 "summary": {<br>
  "total": 8,<br>
  "ports": 4,<br>
  "time": 10.003285617<br>
 },<br>
 "results": [<br>
  {<br>
   "target": "example.com",<br>
   "port": 443,<br>
   "banner": "No banner received"<br>
  },<br>
  {<br>
   "target": "example.com",<br>
   "port": 80,<br>
   "banner": "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nETag: \"84238dfc8092e5d9c0dac8ef93371a07:1736799080.121134\"\r\nLast-Modified: Mon, 13 Jan 2025 20:11:20 GMT\r\nCache-Control: max-age=2343\r\nDate: Mon, 31 Mar 2025 05:49:48 GMT\r\nContent-Length: 1256\r\nConnection: close\r\nX-N: S"<br>
  },<br>
  {<br>
   "target": "scanme.nmap.org",<br>
   "port": 80,<br>
   "banner": "HTTP/1.1 200 OK\r\nDate: Mon, 31 Mar 2025 05:49:48 GMT\r\nServer: Apache/2.4.7 (Ubuntu)\r\nAccept-Ranges: bytes\r\nVary: Accept-Encoding\r\nConnection: close\r\nTransfer-Encoding: chunked\r\nContent-Type: text/html"<br>
  },<br>
  {<br>
   "target": "scanme.nmap.org",<br>
   "port": 22,<br>
   "banner": "SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13\r\n"<br>
  }<br>
 ]<br>
}<br>
________________________________________________________________________________________________________________________________________
<br>
Example #2<br>
<br>
go run main.go -target=scanme.nmap.org -start=78 end=82 -workers=200 -timeout=1000ms<br>
<br>
Scanning port scanme.nmap.org:78 out of 5 total ...<br>
Scanning port scanme.nmap.org:79 out of 5 total ...<br>
Scanning port scanme.nmap.org:81 out of 5 total ...<br>
Scanning port scanme.nmap.org:82 out of 5 total ...<br>
Scanning port scanme.nmap.org:80 out of 5 total ...<br>
<br>
<br>
<br>
OPEN PORTS<br>
<br>
Target: scanme.nmap.org<br>
Port: 80<br>
Banner: HTTP/1.1 200 OK<br>
Date: Mon, 31 Mar 2025 05:53:44 GMT<br>
Server: Apache/2.4.7 (Ubuntu)<br>
Accept-Ranges: bytes<br>
Vary: Accept-Encoding<br>
Connection: close<br>
Transfer-Encoding: chunked<br>
Content-Type: text/html<br>
<br>
<br>
SUMMARY<br>
<br>
Total ports scanned: 5 <br>
Number of open ports: 1 <br>
Time taken: 7.20 seconds<br>
