package main

import (
	"bufio"
	"fmt"
	"os"
	"log"

	"github.com/nnaka2992/go_web_app/chapter4/thesaurus"
	"github.com/joho/godotenv"
)

// read env information
var _ = godotenv.Load("test.env")

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("Failed to query %q's synonyms: %v\n", word, err)
		}
		if len(syns) == 0 {
			log.Fatalf("%q's synonym not found\n")
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
