package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		email := scanner.Text()
		if isValidEmail(email) {
			fmt.Println("Email address is valid!")
			croppedEmail := cropEmail(email)
			checkDomain(croppedEmail)
		} else {
			fmt.Println("Email address is invalid!")
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}

}

func isValidEmail(email string) bool {
	// Regular expression pattern for validating email address
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the pattern into a regular expression object
	re := regexp.MustCompile(pattern)

	// Match the email address against the regular expression
	return re.MatchString(email)
}

func checkDomain(domain string) {

	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	maxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(maxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v, %v, %v, %v, %v, %v\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)

}

func cropEmail(email string) string {
	atIndex := strings.LastIndex(email, "@")
	return email[atIndex+1:]
}
