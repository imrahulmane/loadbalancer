package main

import (
	"fmt"
	"loadbalancer/balancer"
)

func main()  {
	servers := []string{
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	}

	lb := balancer.NewLoadBalancer(servers)

	fmt.Println("Starting New Loadbalancer...")
	err := lb.Start(":8090")

	if err != nil {
		fmt.Println("Error starting load balancer:", err)
	}

}