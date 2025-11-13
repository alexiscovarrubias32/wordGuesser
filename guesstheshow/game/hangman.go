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
	const maxIncorrect = 6
	guessedLetters := make(map[rune]bool)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nLet's start! You have 45 seconds per guess.")
	printDisplay(display)

	for incorrectGuesses < maxIncorrect && strings.Contains(string(display), "_") {
		// Concurrency for input with timeout
		inputChan := make(chan string)
		go func() {
			fmt.Print("Enter a letter or full word (or type 'steal' for other player to guess): ")
			input, _ := reader.ReadString('\n')
			inputChan <- input
		}()

		timer := time.NewTimer(45 * time.Second)
		var guess string

		select {
		case input := <-inputChan:
			timer.Stop()
			guess = strings.ToLower(strings.TrimSpace(input))
		case <-timer.C:
			fmt.Println("\nTime's up!")
			incorrectGuesses++
			fmt.Printf("That counts as a wrong guess. You have %d guesses left.\n", maxIncorrect-incorrectGuesses)
			fmt.Println() 
			printDisplay(display)
			fmt.Println() 
			continue
		}

		if guess == "steal" {
			// Steal attempt with its own timeout
			stealChan := make(chan string)
			go func() {
				fmt.Print("Second player, enter your full word guess: ")
				stealGuess, _ := reader.ReadString('\n')
				stealChan <- stealGuess
			}()

			stealTimer := time.NewTimer(45 * time.Second)
			var stealGuess string

			select {
			case input := <-stealChan:
				stealTimer.Stop()
				stealGuess = strings.ToLower(strings.TrimSpace(input))
			case <-stealTimer.C:
				fmt.Println("\nTime's up for the steal attempt! Back to the original player.")
				fmt.Println() 
				printDisplay(display)
				fmt.Println() 
				continue
			}

			if stealGuess == wordLower {
				fmt.Printf("Other player wins! The word was: %s\n", word)
				return
			} else {
				fmt.Println("Incorrect steal! Back to original player.")
				fmt.Println() 
				printDisplay(display)
				fmt.Println() 
				continue
			}
		}

		if len(guess) == 1 {
			// Single letter guess
			letter := rune(guess[0])
			if guessedLetters[letter] {
				fmt.Println("You already guessed that letter.")
				fmt.Println() 
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
		} else if len(guess) > 1 {
			// Full word guess
			normalizedGuess := strings.Join(strings.Fields(guess), " ")
			if normalizedGuess == wordLower {
				display = []rune(word)
				break
			} else {
				incorrectGuesses++
				fmt.Printf("Wrong! You have %d guesses left.\n", maxIncorrect-incorrectGuesses)
			}
		} else {
			// Handles empty input
			fmt.Println("Invalid input. Please enter a guess.")
		}

		printDisplay(display)
		fmt.Println() 
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
