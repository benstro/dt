package typer

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed wordlist.txt
var wordlistRaw string

var wordList []string

func init() {
	wordList = strings.Fields(wordlistRaw)
}

// Sample returns n random words from the built-in word list.
func Sample(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = wordList[rand.Intn(len(wordList))]
	}
	return out
}
