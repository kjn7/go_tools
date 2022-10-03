package main

import (
	"fmt"
	"net"
	"sort"
	"time"
	"flag"
	"strings"
	"os"
	"strconv"
)

func worker(host string, ports, results chan int) {
	for p := range ports {
		addr := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.DialTimeout("tcp", addr, 1 * time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func parse_ports(ports string) (int, int) {
	if !strings.Contains(ports,"-") {
		fmt.Printf("bad port specification\n")
		os.Exit(-1)
	}
	a := strings.Split(ports,"-")
	start, e1 := strconv.Atoi(a[0])
	end, e2 := strconv.Atoi(a[1])
	if e1 != nil || e2 != nil {
		fmt.Printf("Failed to parse ports\n")
		os.Exit(-1)
	}
	return start, end
}

func main() {
	var max_workers = flag.Int("w", 100, "max number of workers")
	var host = flag.String("t", "", "target host")
	var port_string = flag.String("p","1-1024", "ports to scan")

	flag.Parse()
	if len(*host) < 1 {
		fmt.Printf("Please specify host name\n")
		os.Exit(-1)
	}


	start_port, end_port := parse_ports(*port_string)

	var open []int

	ports := make(chan int, *max_workers)
	results := make(chan int)

	for i := 0; i < cap(ports); i++ {
		go worker(*host, ports, results)
	}

	go func() {
		for port := start_port; port <= end_port; port++ {
			ports <- port
		}
	}()

	for i := start_port; i <= end_port; i++ {
		port := <- results
		if port != 0 {
			open = append(open, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(open)
	for _, p := range open {
		fmt.Printf("port %d open\n", p)
	}
}
