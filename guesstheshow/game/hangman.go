package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Play starts a Hangman game with the given word list
func Play(wordList []string) {
	rand.Seed(time.Now().UnixNano())
	word := wordList[rand.Intn(len(wordList))]
	wordLower := strings.ToLower(word)

	// Initialize display with underscores, preserve spaces
	display := make([]rune, len(word))
	for i := range display {
		if word[i] == ' ' {
			display[i] = ' '
		} else {
			display[i] = '_'
		}
	}

	incorrectGuesses := 0
	const maxIncorrect = 3
	guessedLetters := make(map[rune]bool)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nLet's start! Guess the word or letters:")
	printDisplay(display)

	for incorrectGuesses < maxIncorrect && strings.Contains(string(display), "_") {
		fmt.Print("Enter a letter or full word (or type 'steal' for other player to guess): ")
		input, _ := reader.ReadString('\n')
		guess := strings.ToLower(strings.TrimSpace(input))

		if guess == "steal" {
			// Steal attempt
			fmt.Print("Other player, enter your full word guess: ")
			stealGuess, _ := reader.ReadString('\n')
			stealGuess = strings.ToLower(strings.TrimSpace(stealGuess))

			if stealGuess == wordLower {
				fmt.Printf("Other player wins! The word was: %s\n", word)
				return
			} else {
				fmt.Println("Incorrect steal! Back to original player.")
				printDisplay(display)
				continue
			}
		}

		if len(guess) == 1 {
			// Single letter guess
			letter := rune(guess[0])
			if guessedLetters[letter] {
				fmt.Println("You already guessed that letter.")
				continue
			}

			guessedLetters[letter] = true

			if strings.ContainsRune(wordLower, letter) {
				for i, c := range wordLower {
					if c == letter {
						display[i] = rune(word[i])
					}
				}
				fmt.Println("Correct!")
			} else {
				incorrectGuesses++
				fmt.Printf("Wrong! You have %d guesses left.\n", maxIncorrect-incorrectGuesses)
			}
		} else {
			// Full word guess
			normalizedGuess := strings.Join(strings.Fields(guess), " ")
			if normalizedGuess == wordLower {
				display = []rune(word)
				break
			} else {
				incorrectGuesses++
				fmt.Printf("Wrong! You have %d guesses left.\n", maxIncorrect-incorrectGuesses)
			}
		}

		printDisplay(display)
	}

	if strings.Contains(string(display), "_") {
		fmt.Printf("You lose! The word was: %s\n", word)
	} else {
		fmt.Println("Congratulations! You guessed the word!")
	}
}

// Helper function to print the current state
func printDisplay(display []rune) {
	for _, c := range display {
		fmt.Printf("%c ", c)
	}
	fmt.Println()
}
