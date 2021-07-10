package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"log"

	"github.com/nnaka2992/go_web_app/chapter4/thesaurus"
	"github.com/joho/godotenv"
)

func main() {
	// read env information
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	envPath := filepath.Dir(ex) + "/test.env"
	var _ = godotenv.Load(envPath)
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		fmt.Println(word)
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
		fmt.Println(err)
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
