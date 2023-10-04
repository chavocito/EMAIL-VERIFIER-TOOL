package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain, hasMx, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from standard input")
	}
}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecord, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecord) > 0 {
		hasMx = true
	}

	dnsTxtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dnsTxtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n\n", err)
	}

	for _, dmarcRec := range dmarcRecords {
		if strings.HasPrefix(dmarcRec, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = dmarcRec
			break
		}
	}

	fmt.Printf("%v %v %v %v %v %v", domain, hasMx, hasSPF, hasDMARC, spfRecord, dmarcRecord)
}
