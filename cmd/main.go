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
		// Correctly parse args[1] (second argument) as an integer
		numPlays, err = strconv.Atoi(args[2]) // Correctly use args[1]

		if err != nil {
			numPlays = 1_000_000                               // Default value
			log.Fatal("Error parsing numPlays argument:", err) // Proper error logging
		}
	} else {
		numPlays = 1_000_000 // Default value if no argument is provided
	}
	totalBet := numLines * creditsWagered
	game := reel.NewSlots()
	totals := 0
	i := 0
	for i < int(numPlays) {
		game.Spin()
		totals += game.Evaluate(int(numLines))
		i++
	}
	// game.PrintReel()
	fmt.Println("Amount won!!!!! ", totals)
	fmt.Println("RTP!!!!! ", float64(totals)/(float64(totalBet*int64(numPlays))))
	fmt.Println("Bonuses Hit: ", float64(game.BonusesHit)/(float64(totalBet)*float64(numPlays)))
	fmt.Println("Lines Hit: ", float64(game.LinesPaidOut)/(float64(totalBet)*float64(numPlays)))

	fmt.Println("Scatter counter", game.ScatterPaidOut)
	fmt.Println("Scatter Hit: ", float64(game.ScatterPaidOut)/(float64(totalBet)*float64(numPlays)))

}
