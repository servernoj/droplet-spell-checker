package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

func lineCounter(r io.Reader) (int, error) {
	// credit: https://stackoverflow.com/a/24563853/2989718
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

type Ordered interface {
	~string
}

func binarySearch[T Ordered](a []T, x T) int {
	// credit: https://stackoverflow.com/a/72532012/2989718
	start, mid, end := 0, 0, len(a)-1
	for start <= end {
		mid = (start + end) >> 1
		switch {
		case a[mid] > x:
			end = mid - 1
		case a[mid] < x:
			start = mid + 1
		default:
			return mid
		}
	}
	return -mid
}

func readDict(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	numberOfLines, err := lineCounter(f)
	if err != nil {
		return nil, err
	}
	buf := make([]string, numberOfLines)
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	idx := 0
	for scanner.Scan() {
		word := scanner.Text()
		buf[idx] = word
		idx++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !slices.IsSortedFunc(buf, strings.Compare) {
		return nil, errors.New("invalid dictionary, not sorted")
	}
	return buf, nil
}

type ResultItem struct {
	Suggestions []string
	Line        int
	Column      int
	Word        string
	Context     string
}

func analyzeText(fileName string, dict []string) ([]ResultItem, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lineIdx := 0
	result := []ResultItem{}
	for scanner.Scan() {
		line := scanner.Text()
		words := regexp.MustCompile(`(\s|[()])+`).Split(line, -1)
		for wordIdx, word := range words {
			isEmpty := len(word) == 0
			isCapitalized := regexp.MustCompile(`^[A-Z]+`).MatchString(word)
			// TODO: account on punctuation from the previous line
			// TODO: add handling for proper nouns after punctuation
			isAfterPunctuation := wordIdx > 0 && regexp.MustCompile(`[.,!?]$`).MatchString(words[wordIdx-1])
			if isEmpty || (isCapitalized && !isAfterPunctuation) {
				continue
			}
			word := strings.ToLower(
				strings.Trim(
					word,
					",.!?-",
				),
			)
			searchResult := binarySearch(dict, word)
			if searchResult > 0 {
				// word found in the dictionary, continue analyzis...
				continue
			}
			// identify suggestions as neighbours to the closest match in dict
			// TODO: the size of "left" and "right" wings should be adjusted based on whether in dict the close match is located
			suggestionCenter := -searchResult + 1
			suggestionWingSize := 2
			suggestionLeft := max(0, suggestionCenter-suggestionWingSize)
			suggestionRight := min(len(dict), suggestionCenter+suggestionWingSize)
			suggestions := dict[suggestionLeft:suggestionRight]
			// context
			contextWingSize := 2
			contextLeft := max(0, wordIdx-contextWingSize)
			contextRight := min(len(words), wordIdx+contextWingSize)
			contextWords := words[contextLeft:contextRight]
			contextRegex := strings.Join(contextWords, ".+")
			context := regexp.MustCompile(contextRegex).FindString(line)
			wordIndexInContext := strings.Index(context, word)
			contextIndexInLine := strings.Index(line, context)
			wordIndexInLine := contextIndexInLine + wordIndexInContext
			if wordIndexInLine > 0 {
				wordIndexInLine++
			}
			// result
			result = append(result, ResultItem{
				Word:        word,
				Line:        lineIdx + 1,
				Column:      wordIndexInLine,
				Suggestions: suggestions,
				Context:     context,
			})
		}
		lineIdx++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, err
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <dict-file> <file-to-check>", os.Args[0])
	}
	dict, err := readDict(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	result, err := analyzeText(os.Args[2], dict)
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
