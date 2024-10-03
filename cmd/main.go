package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	reel "github.com/pepperonirollz/slots/pkg"
)

func main() {
	args := os.Args[1:]

	// Check if there are any arguments
	if len(args) == 0 {
		fmt.Println("No arguments provided.")
		return
	}
	numLines, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Fatal("error parsing numLines arg")
	}
	creditsWagered, err := strconv.ParseInt(args[1], 10, 32)

	if err != nil {
		log.Fatal("error parsing numLines arg")
	}
	var numPlays int

	if len(args) > 2 {
		numPlays, err = strconv.Atoi(args[2]) // Correctly use args[1]

		if err != nil {
			log.Fatal("Error parsing numPlays argument:", err) // Proper error logging
		}
	} else {
		numPlays = 10_000_000 // Default value if no argument is provided
	}
	totalBet := numLines * creditsWagered
	game := reel.NewSlots()
	totals := 0
	i := 0
	for i < int(numPlays) {
		game.Spin()
		totals += game.Evaluate(int(numLines), int(creditsWagered))
		i++
	}
	bonusHitFreq := float64(game.BonusesHitRate) / (float64(numPlays * int(numLines)))
	lineHitFreq := float64(game.LineHitRate) / float64(numPlays*int(numLines))
	scatterHitFreq := float64(game.ScatterHitRate) / (float64(numPlays))
	fmt.Println("plays, lines, wagered ", numPlays, numLines, creditsWagered)
	fmt.Println("Credits won: ", totals, "/", totalBet*int64(numPlays))
	fmt.Println("Return to player: ", float64(totals)/(float64(totalBet*int64(numPlays))))
	fmt.Println("Bonuses Hit Freq: ", bonusHitFreq)
	fmt.Println("Bonuses Paid percentage: ", float64(game.BonusesHitPay)/(float64(totalBet)*float64(numPlays)))

	fmt.Println("Lines hit freq ", lineHitFreq)
	fmt.Println("Lines paid percentage: ", float64(game.LineHitPay)/(float64(totalBet)*float64(numPlays)))
	fmt.Println("Scatter hit freq", scatterHitFreq)
	fmt.Println("Scatter paid percentage: ", float64(game.ScatterHitPay)/(float64(totalBet)*float64(numPlays)))
	fmt.Println("Total hit freq", scatterHitFreq+lineHitFreq+bonusHitFreq)

}
