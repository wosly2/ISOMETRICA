package stailoh

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func buildRandomSymbol(doSpacesAndNewlines bool) Symbol {
	rand.Seed(time.Now().UnixNano())
	var symbol Symbol

	if doSpacesAndNewlines {
		symbol.IsSpace = rand.Intn(3) == 1
		if !symbol.IsSpace {
			symbol.IsNewline = rand.Intn(7) == 1
		}
	}

	if !symbol.IsSpace && !symbol.IsNewline {
		symbol.Consonant = consonantSounds[rand.Intn(len(consonantSounds))]
		symbol.Vowel = vowelSounds[rand.Intn(len(vowelSounds))]
		symbol.Reversed = rand.Intn(2) == 1
	}

	return symbol
}

func TestSymbolToString(t *testing.T) {
	var testCount int = 10

	for i := 0; i < testCount; i++ {
		var symbol Symbol = buildRandomSymbol(false)

		fmt.Printf("Consonant: %s\nVowel: %s\nReversed: %t\n", symbol.Consonant, symbol.Vowel, symbol.Reversed)
		fmt.Printf("String: %s\n\n", symbol.ToString())
	}
}

func TestSymbolSetToString(t *testing.T) {
	var testCount int = 10
	var symbolCount int = 5000

	for i := 0; i < testCount; i++ {
		// build a symbol set
		var symbolSet SymbolSet

		for i := 0; i < symbolCount; i++ {
			symbolSet.Symbols = append(symbolSet.Symbols, buildRandomSymbol(true))
		}

		fmt.Printf("Symbols: %v\n", symbolSet.Symbols)
		fmt.Printf("String: \n%s\n\n", symbolSet.ToString())
	}
}

func TestDrawSymbol(t *testing.T) {
	// build a new ebiten loop

}
