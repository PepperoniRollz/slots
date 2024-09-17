package reel

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

type Reel struct {
	Symbols     []Symbol
	Size        int
	VirtualReel []Symbol
}

type Symbol struct {
	Name   string
	Weight int
}

func NewReel(symbols []Symbol) *Reel {
	return &Reel{
		Symbols: symbols,
	}
}
func (l *Slots) SymbolsToIcons() {
	symbolToUnicode := map[string]string{
		"WS": "🃏",  // Wild
		"LM": "🦞",  // Lobstermania
		"BU": "🎈",  // Buoy
		"BO": "⛵",  // Boat
		"LH": "🚨",  // Light House
		"TU": "🐟",  // Tuna
		"CL": "🦪",  // Clam
		"SG": "🕊️", // Sea Gull
		"SF": "🌟",  // Star Fish
		"LO": "🎰",  // Bonus
		"LT": "🎲",  // Scatter
	}
	unicodeSymbols := make([]string, len(l.SpinResult))

	for i, symbol := range l.SpinResult {
		if unicodeChar, ok := symbolToUnicode[symbol.Name]; ok {
			unicodeSymbols[i] = unicodeChar
		} else {
			unicodeSymbols[i] = "❓"
		}
	}
	l.PrettyReel = unicodeSymbols
}

func NewSymbol(name string, weight int) *Symbol {
	return &Symbol{Name: name,
		Weight: weight}
}
func NewSymbols(numSymbols int, names []string, weights []int) []Symbol {
	var symbols []Symbol
	for i := 0; i < numSymbols; i++ {
		symbols = append(symbols, *NewSymbol(names[i], weights[i]))
	}

	return symbols
}

type Payline struct {
	payline        string
	payout         int
	creditsWagered int
}

func NewPayline(symbols string, payout int, creditsWagered int) *Payline {
	return &Payline{
		payline:        symbols,
		payout:         payout,
		creditsWagered: creditsWagered,
	}
}

func buildReels() []*Reel {
	reels := []*Reel{
		NewReel([]Symbol{
			{Name: "WS", Weight: 2},
			{Name: "LM", Weight: 4},
			{Name: "BU", Weight: 4},
			{Name: "BO", Weight: 6},
			{Name: "LH", Weight: 5},
			{Name: "TU", Weight: 6},
			{Name: "CL", Weight: 6},
			{Name: "SG", Weight: 5},
			{Name: "SF", Weight: 5},
			{Name: "LO", Weight: 2},
			{Name: "LT", Weight: 2},
		}),
		NewReel([]Symbol{
			{Name: "WS", Weight: 2},
			{Name: "LM", Weight: 4},
			{Name: "BU", Weight: 4},
			{Name: "BO", Weight: 4},
			{Name: "LH", Weight: 4},
			{Name: "TU", Weight: 4},
			{Name: "CL", Weight: 6},
			{Name: "SG", Weight: 6},
			{Name: "SF", Weight: 5},
			{Name: "LO", Weight: 5},
			{Name: "LT", Weight: 2},
		}),
		NewReel([]Symbol{
			{Name: "WS", Weight: 1},
			{Name: "LM", Weight: 3},
			{Name: "BU", Weight: 5},
			{Name: "BO", Weight: 4},
			{Name: "LH", Weight: 6},
			{Name: "TU", Weight: 5},
			{Name: "CL", Weight: 5},
			{Name: "SG", Weight: 5},
			{Name: "SF", Weight: 6},
			{Name: "LO", Weight: 6},
			{Name: "LT", Weight: 2},
		}),
		NewReel([]Symbol{
			{Name: "WS", Weight: 4},
			{Name: "LM", Weight: 4},
			{Name: "BU", Weight: 4},
			{Name: "BO", Weight: 4},
			{Name: "LH", Weight: 6},
			{Name: "TU", Weight: 6},
			{Name: "CL", Weight: 6},
			{Name: "SG", Weight: 6},
			{Name: "SF", Weight: 8},
			{Name: "LO", Weight: 0},
			{Name: "LT", Weight: 2},
		}),

		NewReel([]Symbol{
			{Name: "WS", Weight: 2},
			{Name: "LM", Weight: 4},
			{Name: "BU", Weight: 5},
			{Name: "BO", Weight: 4},
			{Name: "LH", Weight: 7},
			{Name: "TU", Weight: 7},
			{Name: "CL", Weight: 6},
			{Name: "SG", Weight: 6},
			{Name: "SF", Weight: 7},
			{Name: "LO", Weight: 0},
			{Name: "LT", Weight: 2},
		}),
	}
	for _, reel := range reels {
		sum := 0
		for _, symbol := range reel.Symbols {
			sum += symbol.Weight
		}
		reel.Size = sum
	}

	for _, reel := range reels {
		var virtualReel []Symbol
		for len(virtualReel) < reel.Size {
			for _, symbol := range reel.Symbols {
				if count(virtualReel, symbol) < symbol.Weight && float64(count(virtualReel, symbol)) <= float64(symbol.Weight)/float64(reel.Size)*float64(len(virtualReel)) {
					virtualReel = append(virtualReel, symbol)
				}
			}
		}
		reel.VirtualReel = virtualReel
	}
	return reels
}

func count(reel []Symbol, symbol Symbol) int {
	count := 0
	for _, s := range reel {
		if s == symbol {
			count++
		}
	}
	return count
}

type Slots struct {
	Reels            []*Reel
	PayLines         [][]int
	Payouts          map[Payout]int
	CreditsWagered   int
	NumLinesSelected int
	SpinResult       []Symbol
	BonusesHit       int
	ScatterPaidOut   int
	LinesPaidOut     int
	PrettyReel       []string
}

func NewSlots() *Slots {
	return &Slots{
		Reels:    buildReels(),
		PayLines: buildPayLines(),
		Payouts:  buildPayoutMap(),
	}
}

func (l *Slots) Evaluate(linesPlayed int) int {

	totalPayout := 0
	for i := 0; i < linesPlayed; i++ {
		count := 1
		wildCounts := 0
		firstSymbolInPayline := l.SpinResult[l.PayLines[i][0]]
		if firstSymbolInPayline.Name == "WS" {
			wildCounts++
		}
		for j := 1; j < 5; j++ {

			currentSymbol := l.SpinResult[l.PayLines[i][j]].Name

			if currentSymbol == "WS" && wildCounts > 0 {
				wildCounts++
			}

			if currentSymbol == "LT" || firstSymbolInPayline.Name == "LT" {
				break
			}

			if firstSymbolInPayline.Name == "WS" && currentSymbol != "WS" && currentSymbol != "LO" {
				firstSymbolInPayline.Name = currentSymbol
			}
			if currentSymbol == firstSymbolInPayline.Name || currentSymbol == "WS" {
				if firstSymbolInPayline.Name == "LO" && currentSymbol == "WS" {
					break
				}
				count++

			} else {
				break
			}
		}

		if count >= 2 {
			linePay := int(math.Max(float64(l.Payouts[Payout{Symbol: firstSymbolInPayline.Name, Count: count}]), float64(l.Payouts[Payout{Symbol: "WS", Count: wildCounts}])))
			if linePay == 331 {
				l.BonusesHit += 331

			} else {
				l.LinesPaidOut += linePay
			}
			totalPayout += linePay
		}
	}

	//check scatter
	scatterCount := 0
	for _, symbol := range l.SpinResult {
		if symbol.Name == "LT" {
			scatterCount++
		}
	}
	if scatterCount >= 3 {
		scatterPay := l.Payouts[Payout{Symbol: "LT", Count: scatterCount}] * linesPlayed
		totalPayout += scatterPay
		l.ScatterPaidOut += scatterPay
	}

	return totalPayout
}

// func (l *LuckyLobster) bonusRound() {
// 	userSelection := 1
// 	numberOfBuoysWithPizes, _ := rand.Int(rand.Reader, big.NewInt(3))
//
// }
//
// func (l *LuckyLobster) buildBonusPrizes() {
// 	m := map[int]int{
// 		10:  10,
// 		14:  5,
// 		19:  6,
// 		24:  7,
// 		29:  8,
// 		39:  10,
// 		49:  12,
// 		59:  15,
// 		79:  20,
// 		99:  22,
// 		119: 25,
// 		139: 27,
// 		158: 30,
// 		180: 35,
// 		204: 45,
// 		223: 50,
// 		238: 55,
// 		253: 60,
// 		268: 65,
// 		283: 70,
// 		298: 75,
// 		308: 100,
// 		316: 150,
// 		321: 250,
// 	}
// 	upperBound := []int{
// 		10, 14, 19, 24, 29, 39, 49, 59, 79, 99, 119, 139, 158, 180, 204, 223, 238, 253, 268, 283, 298, 308, 316, 321,
// 	}
// 	bonus := []int{
// 		10, 5, 6, 7, 8, 10, 12, 15, 20, 22, 25, 27, 30, 35, 45, 50, 55, 60, 65, 70, 75, 100, 150, 250,
// 	}
// 	var bonuses []int
// 	for i := 0; i < upperBound[len(upperBound)-1]; i++ {
//
// 	}
//
// }

func (l *Slots) Spin() {
	resultGrid := make([][]Symbol, 3)

	results := make([]Symbol, 15, 15)
	for i := range resultGrid {
		resultGrid[i] = make([]Symbol, len(l.Reels))
	}

	for reelIndex, reel := range l.Reels {
		startIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(reel.Size)))
		startIndexInt := (int(startIndex.Int64()-1)%reel.Size + reel.Size) % reel.Size

		for row := 0; row < 3; row++ {
			index := (startIndexInt + row) % reel.Size

			resultGrid[row][reelIndex] = reel.VirtualReel[index]

			flatIndex := row*len(l.Reels) + reelIndex
			results[flatIndex] = reel.VirtualReel[index]
		}
	}
	l.SpinResult = results
}

func (l *Slots) PrintReel() {
	rows := 3
	cols := 5
	fmt.Print("\033[H\033[2J")
	fmt.Print("= = = = = = = = = = =\n")

	l.AnimateSpin()
	for i := 0; i < rows; i++ {
		row := ""
		for j := 0; j < cols; j++ {
			index := i*cols + j
			row += l.PrettyReel[index] + " "
		}
		fmt.Println(row)
	}

	fmt.Print("= = = = = = = = = = =\n")
}
func (l *Slots) AnimateSpin() {

	endTime := time.Now().Add(3 * time.Second)

	updateInterval := 100 * time.Millisecond
	var symbols = l.PrettyReel

	fmt.Print("\033[s")
	for time.Now().Before(endTime) {
		rows := 3
		cols := 5
		for i := 0; i < rows; i++ {
			row := ""
			for j := 0; j < cols; j++ {
				r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
				row += l.PrettyReel[r.Int64()] + " "
			}
			fmt.Printf("%s ", row)
			fmt.Println()
		}
		time.Sleep(updateInterval)

		fmt.Print("\033[u")
	}
}
func buildPayLines() [][]int {

	//slot machine
	//[
	//[0, 1, 2, 3, 4],
	//[5, 6, 7, 8, 9],
	//[10,11,12,13,14]
	//]
	return [][]int{
		{5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4},
		{10, 11, 12, 13, 14},
		{0, 6, 12, 8, 4},
		{10, 6, 2, 8, 14}, //5,
		{10, 11, 7, 3, 4},
		{0, 1, 7, 13, 14},
		{5, 11, 7, 3, 9},
		{5, 1, 7, 13, 9}, //9,
		{10, 6, 7, 8, 4},
		{0, 6, 7, 8, 14},
		{5, 11, 12, 8, 4},
		{5, 1, 2, 8, 14},
		{5, 6, 12, 8, 4},
		{5, 6, 2, 8, 14}, //15,
	}
}

type Payout struct {
	Symbol string
	Count  int
}

func buildPayoutMap() map[Payout]int {
	return map[Payout]int{

		{Symbol: "WS", Count: 2}: 5,
		{Symbol: "WS", Count: 3}: 100,
		{Symbol: "WS", Count: 4}: 500,
		{Symbol: "WS", Count: 5}: 10_000,

		{Symbol: "LM", Count: 2}: 2,
		{Symbol: "LM", Count: 3}: 40,
		{Symbol: "LM", Count: 4}: 200,
		{Symbol: "LM", Count: 5}: 1_000,

		{Symbol: "BU", Count: 3}: 25,
		{Symbol: "BU", Count: 4}: 100,
		{Symbol: "BU", Count: 5}: 500,

		{Symbol: "BO", Count: 3}: 25,
		{Symbol: "BO", Count: 4}: 100,
		{Symbol: "BO", Count: 5}: 500,

		{Symbol: "LH", Count: 3}: 10,
		{Symbol: "LH", Count: 4}: 50,
		{Symbol: "LH", Count: 5}: 500,

		{Symbol: "TU", Count: 3}: 10,
		{Symbol: "TU", Count: 4}: 50,
		{Symbol: "TU", Count: 5}: 250,

		{Symbol: "CL", Count: 3}: 5,
		{Symbol: "CL", Count: 4}: 30,
		{Symbol: "CL", Count: 5}: 200,

		{Symbol: "SG", Count: 3}: 5,
		{Symbol: "SG", Count: 4}: 30,
		{Symbol: "SG", Count: 5}: 200,

		{Symbol: "SF", Count: 3}: 5,
		{Symbol: "SF", Count: 4}: 30,
		{Symbol: "SF", Count: 5}: 150,

		{Symbol: "LO", Count: 3}: 331,

		{Symbol: "LT", Count: 3}: 5,
		{Symbol: "LT", Count: 4}: 25,
		{Symbol: "LT", Count: 5}: 200,
	}
}
