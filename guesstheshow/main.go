package main

import (
	"fmt"
	"strings"

	"guesstheshow/data"
	"guesstheshow/game"
)

func main() {
	fmt.Println("Choose a genre:")
	for genre := range data.Shows {
		fmt.Println("-", genre)
	}

	var input string
	fmt.Print("Enter genre: ")
	fmt.Scanln(&input)

	// Normalize input
	normalizedInput := strings.ToLower(strings.ReplaceAll(input, " ", "-"))

	var selectedGenre string
	for genre := range data.Shows {
		if strings.ToLower(strings.ReplaceAll(genre, " ", "-")) == normalizedInput {
			selectedGenre = genre
			break
		}
	}

	if selectedGenre == "" {
		fmt.Println("Sorry, that genre isn't available.")
		return
	}

	// Start game
	game.Play(data.Shows[selectedGenre])
}
