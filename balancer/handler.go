package balancer

import (
	"fmt"
	"io"
	"net"
)

func handleConnection(clientConn net.Conn, lb *LoadBalancer){
	defer clientConn.Close()

	//get the next server using round robin
	backend := lb.getNextServer()

	if backend == "" {
		fmt.Println("No running server found!!")
		send502Response(clientConn)
		return
	}
	
	fmt.Printf("Forwarding connection to %s\n", backend)

	backendConn, err := net.Dial("tcp", backend)
	if err != nil {
		fmt.Printf("Failed to connect to backend %s: %v\n", backend, err)
		send502Response(clientConn)
		return
	}

	defer backendConn.Close()

	//copy data bidirectionally
	//Go routing - client --> Backend
	go io.Copy(backendConn, clientConn)

	//backend --> client
	io.Copy(clientConn, backendConn)
}

func send502Response(conn net.Conn){
	response := "HTTP/1.1 502 Bad Gateway\r\n"
	response += "Content-Type: text/plain\r\n"
	response += "Content-Length: 21\r\n"
	response += "\r\n"
	response += "Backend Unavailable\n"
	conn.Write([]byte(response))
}