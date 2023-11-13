package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

func getWords(data string) int {
	words := strings.Fields(data)
	return len(words)
}

func getTotalLines(data string) int {
	scanner := bufio.NewScanner(strings.NewReader(data))
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount
}

func getTotalCounts(files []string) (int, int, int, int) {
	var totalBytes, totalLines, totalWords, totalChars int

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			continue
		}

		totalBytes += len(data)
		totalLines += getTotalLines(string(data))
		totalWords += getWords(string(data))
		totalChars += utf8.RuneCountInString(strings.ReplaceAll(string(data), "\n", ""))
	}

	return totalBytes, totalLines, totalWords, totalChars
}

func processFiles(files []string, countBytes, countLines, countWords, countChars bool) {
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
			continue
		}

		// Print counts based on flags
		printCounts(string(data), countBytes, countLines, countWords, countChars)
		fmt.Printf("%s\n", file)
	}

	// If more than one file, print total counts
	if len(files) > 1 {
		totalBytes, totalLines, totalWords, totalChars := getTotalCounts(files)
		if countBytes {
			fmt.Printf("\t%d ", totalBytes)
		}
		if countLines {
			fmt.Printf("\t%d ", totalLines)
		}
		if countWords {
			fmt.Printf("\t%d ", totalWords)
		}
		if countChars {
			fmt.Printf("\t%d ", totalChars)
		}

		fmt.Print("total\n")
	}
}

func printCounts(data string, countBytes, countLines, countWords, countChars bool) {
	bytesCount := len(data)
	linesCount := getTotalLines(string(data))
	wordsCount := getWords(data)
	charsCount := utf8.RuneCountInString(data)

	// Print counts based on flags
	if countBytes {
		fmt.Printf("\t%d ", bytesCount)
	}
	if countLines {
		fmt.Printf("\t%d ", linesCount)
	}
	if countWords {
		fmt.Printf("\t%d ", wordsCount)
	}
	if countChars {
		fmt.Printf("\t%d ", charsCount)
	}
}

func main() {

	// Defining flags
	countBytes := flag.Bool("c", false, "Count bytes")
	countLines := flag.Bool("l", false, "Count lines")
	countWords := flag.Bool("w", false, "Count words")
	countChars := flag.Bool("m", false, "Count characters")

	// Parsing command line arguments
	flag.Parse()

	// Get the list of files from positional arguments
	files := flag.Args()

	// If no flags are provided, default to counting bytes, lines, and words
	if !(*countBytes || *countLines || *countWords || *countChars) {
		*countBytes, *countLines, *countWords, *countChars = true, true, true, false
	}

	// Process each file or read from standard input
	if len(files) > 0 {
		// Process files provided as arguments
		processFiles(files, *countBytes, *countLines, *countWords, *countChars)
	} else {
		// Read from standard input
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("usage: ggwc [-clmw] [file ...]")
			return
		}

		// Print counts based on flags
		printCounts(string(data), *countBytes, *countLines, *countWords, *countChars)
		fmt.Print("\n")
	}
}
