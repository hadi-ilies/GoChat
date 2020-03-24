package client

import (
	"bufio"
)

//getLine: get the line wrote by the client
func getLine(scanner *bufio.Scanner) string {
	// Scans a line from Stdin(Console)
	scanner.Scan()
	// Holds the string that scanned
	line := scanner.Text()
	return line
}
