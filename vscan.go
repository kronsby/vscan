package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

func worker(addr string, ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func isValidAddress(address string) bool {
	ip := net.ParseIP(address)
	if ip != nil {
		return true
	}
	_, err := net.LookupHost(address)
	return err == nil
}

func expandRange(low int, high int) []int {
	var res []int
	for i := low; i <= high; i++ {
		res = append(res, i)
	}
	return res
}

func parsePorts(ports string) ([]int, error) {
	split_ports := strings.Split(ports, ",")
	var res []int
	for _, p := range split_ports {
		port_range := strings.Split(p, "-")
		// If no delimiter is found, Split returns a single port
		if len(port_range) > 2 {
			return nil, fmt.Errorf("invalid port range, must be a range between 2 numbers")
		} else if len(port_range) == 2 {
			low, err := strconv.Atoi(port_range[0])
			if err != nil {
				return nil, fmt.Errorf("not an int: %s", port_range[0])
			}
			high, err := strconv.Atoi(port_range[1])
			if err != nil {
				return nil, fmt.Errorf("not an int: %s", port_range[1])
			}
			if low >= high {
				return nil, fmt.Errorf("invalid port range, first number must be lower")
			}
			res = append(res, expandRange(low, high)...)
		} else if len(port_range) == 1 {
			port, err := strconv.Atoi(port_range[0])
			if err != nil {
				return nil, fmt.Errorf("not an int: %s", port_range[0])
			}
			res = append(res, port)
		}
	}
	return res, nil
}

func isValidPort(port int) bool {
	return (1 <= port && port <= 65535)
}

func cleanPorts(ports []int) []int {
	seen := make(map[int]bool)
	var res []int
	for _, p := range ports {
		if !seen[p] && isValidPort(p) {
			seen[p] = true
			res = append(res, p)
		}
	}
	return res
}

func displayKickoff(address string, ports []int) {
	fmt.Printf("Scanning target: %s\n", address)
	var p strings.Builder
	for i := 0; i < min(len(ports), 3); i++ {
		p.WriteString(strconv.Itoa(ports[i]))
		if i < len(ports)-1 && i < 2 {
			p.WriteString(", ")
		}
	}
	fmt.Printf("Ports to be scanned: %s ... (total: %d)\n", p.String(), len(ports))
}

func displayResults(openports []int) {
	fmt.Println("Scan completed successfully")
	for _, p := range openports {
		fmt.Printf("%d Open\n", p)
	}
}

func main() {
	address_in := flag.String("address", "", "The address to scan")
	ports_in := flag.String("ports", "", "A list of ports to scan. Valid inputs include comma separated list as well as port ranges (eg. 1-99)")

	flag.Parse()

	var address string

	if flag.NArg() > 0 {
		address = flag.Arg(0)
	} else if *address_in != "" {
		address = *address_in
	} else {
		flag.Usage()
		log.Fatal("Address is required")
	}

	if !isValidAddress(address) {
		flag.Usage()
		log.Fatalf("Invalid Address %s", address)
	}

	var ports []int
	if *ports_in != "" {
		parsedPorts, err := parsePorts(*ports_in)
		if err != nil {
			flag.Usage()
			log.Fatalf("Invalid port specified: %s", err)
		}
		ports = cleanPorts(parsedPorts)
	} else {
		ports = make([]int, 65535)
		for i := range ports {
			ports[i] = i + 1
		}
	}

	displayKickoff(address, ports)

	ports_chan := make(chan int, 10)
	results := make(chan int)
	var openports []int
	for i := 0; i < cap(ports_chan); i++ {
		go worker(address, ports_chan, results)
	}

	go func() {
		for _, p := range ports {
			ports_chan <- p
		}
	}()

	for range ports {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports_chan)
	close(results)
	sort.Slice(openports, func(i, j int) bool {
		return openports[i] < openports[j]
	})
	displayResults(openports)
}
