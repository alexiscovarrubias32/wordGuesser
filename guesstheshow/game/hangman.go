package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Play starts a Hangman game with the given word list, showcasing several Go features:
// - Use of the standard library for randomization (`math/rand`), time (`time`), and input/output (`fmt`, `bufio`).
// - Basic control flow (`for` loop) and data structures (slices and maps).
// - Concurrency with goroutines, channels, and the `select` statement for the guess timeout feature.
func Play(wordList []string) {
	// GO FEATURE: Use of `math/rand` and `time` from the standard library.
	// We seed the random number generator to ensure a different word is chosen each time.
	rand.Seed(time.Now().UnixNano())
	word := wordList[rand.Intn(len(wordList))]
	wordLower := strings.ToLower(word)

	// --- Game State Initialization ---
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

	// --- Main Game Loop ---
	// The loop continues as long as the player has guesses left and the word has not been fully revealed.
	for incorrectGuesses < maxIncorrect && strings.Contains(string(display), "_") {

		// =====================================================================
		// GO FEATURE: Concurrency for Input with a Timeout
		// This section demonstrates Go's powerful and simple concurrency model.
		// =====================================================================

		// 1. GO FEATURE: Channels
		// A channel is created to safely pass the user's input from one goroutine to another.
		inputChan := make(chan string)

		// 2. GO FEATURE: Goroutines
		// An anonymous function is launched as a goroutine.
		// This allows the program to listen for user input without blocking the main game loop.
		go func() {
			fmt.Print("Enter a letter or full word (or type 'steal' for other player to guess): ")
			input, _ := reader.ReadString('\n')
			// Send the received input back to the main loop via the channel.
			inputChan <- input
		}()

		// A timer is started for the guess.
		timer := time.NewTimer(45 * time.Second)
		var guess string

		// 3. GO FEATURE: The `select` Statement
		// The `select` statement waits for one of multiple communication operations to complete.
		// Here, it creates a "race" between the user's input and the timer.
		select {
		case input := <-inputChan:
			// Case 1: Input was received from the user.
			timer.Stop() // Stop the timer because the user answered in time.
			guess = strings.ToLower(strings.TrimSpace(input))
		case <-timer.C:
			// Case 2: The timer finished before the user provided input.
			fmt.Println("\nTime's up!")
			incorrectGuesses++
			fmt.Printf("That counts as a wrong guess. You have %d guesses left.\n", maxIncorrect-incorrectGuesses)
			fmt.Println()
			printDisplay(display)
			fmt.Println()
			continue // Skip the rest of the loop and start the next turn.
		}
		// ================= End of Concurrency Section ======================

		// --- Guess Processing Logic ---

		if guess == "steal" {
			// The "steal" attempt also gets its own concurrent timeout logic.
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
				display = []rune(word) // Guessed correctly, reveal the word.
				break                 // Exit the loop.
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

	// --- Game Over Logic ---
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
