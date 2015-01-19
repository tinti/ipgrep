package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	networkAddressParameter := flag.String("n", "0.0.0.0", "Network address with network mask or multiple entries. Example: 0.0.0.0/0 or '10.0.0.0/8 192.168.0.0/24'.")
	numberOfThreadsParameter := flag.Int("t", 1, "Number of concurrent threads parsing the file. Default 2.")
	fieldNumberParameter := flag.Int("c", 1, "Field that should contain the ip address. Note multiple spaces are considered as one.")
	compressedParameter := flag.Bool("z", false, "Compress all io with gzip.")
	inputFileNameParameter := flag.String("i", "", "Input file.")
	outputFileNameParameter := flag.String("o", "", "Output file.")

	flag.Parse()

	numberOfThreads := *numberOfThreadsParameter
	fieldNumber := *fieldNumberParameter
	compressed := *compressedParameter
	inputFileName := *inputFileNameParameter
	outputFileName := *outputFileNameParameter

	// Log to stderr
	log.SetOutput(os.Stderr)

	// Load addresses in array
	networkAddressesParameterArray := strings.Split(*networkAddressParameter, " ")

	var networkAddresses []*net.IPNet

	for _, networkAddressParameterEntry := range networkAddressesParameterArray {
		// Initialize filter
		_, networkAddress, err := net.ParseCIDR(networkAddressParameterEntry)

		if networkAddress == nil || err != nil {
			log.Fatalf("Unable to parse suplied networks '%v'\n", *networkAddressParameter)
			os.Exit(1)
		} else {
			networkAddresses = append(networkAddresses, networkAddress)
		}
	}

	// IO
	inputStream := io.Reader(os.Stdin)
	outputStream := io.Writer(os.Stdout)

	if inputFileName != "" {
		inputFile, _ := os.Open(inputFileName)

		defer func() {
			inputFile.Close()
		}()

		inputStream = io.Reader(inputFile)
	}

	if outputFileName != "" {
		outputFile, _ := os.Create(outputFileName)

		defer func() {
			outputFile.Close()
		}()

		outputStream = io.Writer(outputFile)
	}

	if compressed {
		gzipInputStream, _ := gzip.NewReader(inputStream)
		gzipOutputStream := gzip.NewWriter(outputStream)

		defer func() {
			gzipInputStream.Close()
			gzipOutputStream.Close()
		}()

		inputStream = io.Reader(gzipInputStream)
		outputStream = io.Writer(gzipOutputStream)
	}

	// Create input scanner
	scanner := bufio.NewScanner(inputStream)

	var plains []chan string
	var cleans []chan string

	for i := 0; i < numberOfThreads; i += 1 {
		plains = append(plains, make(chan string, 1))
		cleans = append(cleans, make(chan string, 1))

		go func(in <-chan string, out chan<- string) {
			for {
				text := <-in
				textFiltered := ""
				ipText := ""
				tokens := strings.Fields(text)

				if len(tokens) >= fieldNumber {
					ipText = tokens[fieldNumber-1]
				}

				ipAddress := net.ParseIP(ipText)

				if ipAddress == nil {
					log.Printf("Unable to parse extract ip from entry '%v'\n", text)
				} else {
					for _, networkAddress := range networkAddresses {
						if networkAddress.Contains(ipAddress) {
							textFiltered = text
							break
						}
					}
				}

				out <- textFiltered
			}
		}(plains[i], cleans[i])
	}

	// Scan and filter
	for running := true; running; {
		i := 0

		// Enqueue
		for ; i < numberOfThreads; i += 1 {
			if running = scanner.Scan(); running {
				plains[i] <- scanner.Text()
			} else {
				break
			}
		}

		// Check for errors
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading standard input:", err)
			// TODO(tinti) return 1 in main
			break

		}

		// Dequeue
		for j := 0; j < i; j += 1 {
			if text := <-cleans[j]; text != "" {
				fmt.Fprintf(outputStream, "%v\n", text)
			}
		}
	}
}
