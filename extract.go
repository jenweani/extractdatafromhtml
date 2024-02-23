package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func extractData(node *html.Node, class string, result *[]string) {
	if node.Type == html.ElementNode && node.Data == "span" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, class) {
				// Found the target span tag
				var text string
				// Extract text content from child nodes
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if child.Type == html.ElementNode && child.Data == "span" {
						text += child.FirstChild.Data
					}
				}
				*result = append(*result, text)
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractData(child, class, result)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run extract_data.go input.html outExtract.txt")
		os.Exit(1)
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]
	className := os.Args[3]

	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	var extractedData []string
	extractData(doc, className, &extractedData)

	// Save extracted data to the output file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for _, data := range extractedData {
		_, err := outputFile.WriteString(data + ",")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Extraction completed. Data saved to", outputFileName)
}
