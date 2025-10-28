package balancer

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type LoadBalancer struct {
	servers 		[]string
	currentIndex 	int
	mu 				sync.Mutex

	healthy			map[string]bool
	healthyMu		sync.RWMutex
}

func NewLoadBalancer(servers []string) *LoadBalancer{

	healthy := make(map[string]bool)

	for _, server := range servers {
		healthy[server] = true
	}

	return &LoadBalancer{
		servers: 		servers,
		currentIndex: 	0,
		healthy: 		healthy,
	}
}

func (lb *LoadBalancer) isHealthy(server string) bool{
	lb.healthyMu.RLock()
	defer lb.healthyMu.RUnlock()
	return lb.healthy[server]
}

func (lb *LoadBalancer) setHealthy(server string, status bool){
	lb.healthyMu.Lock()
	defer lb.healthyMu.Unlock()
	lb.healthy[server] = status
}

func (lb *LoadBalancer) checkHealth(server string){
	//try to connect with the server
	conn, err := net.DialTimeout("tcp", server,2*time.Second)

	if err != nil {
		//failed to connect - server is unhealthy
		if lb.isHealthy(server) {
			//log unhealthy only if it's status changed
			fmt.Printf("Server %s marked as UNHEALTHY: %v\n", server, err)
		}
		lb.setHealthy(server, false)
		return
	}

	//Successfully connected now close the connection
	conn.Close()

	if !lb.isHealthy(server) {
		//log only when status changed
		fmt.Printf("Server %s marked as HEALTHY\n", server)
	}

	lb.setHealthy(server, true)
}

func (lb *LoadBalancer) startHealthChecker() {
	ticker := time.NewTicker(10 * time.Second)

	fmt.Println("Health checker started (checking every second)")

	for range ticker.C {
		fmt.Println("Running health checks...")

		for _, server := range lb.servers {
			lb.checkHealth(server)
		}
	}
}

func (lb *LoadBalancer) getNextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// startIndex := lb.currentIndex
	attempts := 0

	for attempts < len(lb.servers){
		server := lb.servers[lb.currentIndex]
		lb.currentIndex = (lb.currentIndex + 1) % len(lb.servers)

		if lb.isHealthy(server) {
			return server
		}
		
		attempts++
	}
	
	return ""
}

func (lb *LoadBalancer) Start(address string) error {
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Load Balancer Listening on %s\n", address)
	fmt.Printf("Forwarding to backends: %v\n", lb.servers)

	//start health checker in background
	go lb.startHealthChecker()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, lb)
	}
}