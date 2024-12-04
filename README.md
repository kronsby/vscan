![build](https://github.com/kronsby/vscan/actions/workflows/go.yml/badge.svg)
## vscan - A Simple TCP Scanner

**vscan** is a simple TCP scanner written in Go. I wrote this while going through the Black Hat Go
book. I expanded a bit on it. Very fun stuff!

### Features

  * **Efficient Scanning:** Scans multiple ports concurrently for fast results.
  * **Port Range Support:** Specify individual ports, comma-separated lists, or ranges (e.g., 1-1024).

### Usage

```bash
vscan -address <target_address> -ports <ports_to_scan> 
```

**Example:**

```bash
vscan -address scanme.nmap.org -ports 21,80,443,8080
```

**Parameters:**

  * `-address`: The target IP address or hostname to scan.
  * `-ports`:  A comma-separated list of ports or port ranges to scan (e.g., 22,80,443,1-100).

**Note:** If the `-ports` option is not provided, all 65535 ports will be scanned.

### Building from Source

1.  Ensure you have Go installed on your system.
2.  Clone this repository: `git clone https://github.com/ice-monster/vscan.git`
3.  Build the executable: `go build`

### License

This project is licensed under the MIT License - see the [LICENSE](/LICENSE) file for details.
