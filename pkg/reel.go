package reel

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

type Reel struct {
	Symbols     []WeightedSymbol
	Size        int
	VirtualReel []Symbol
}

type WeightedSymbol struct {
	Symbol Symbol
	Weight int
}

func NewReel(symbols []WeightedSymbol) *Reel {
	return &Reel{
		Symbols: symbols,
	}
}
func (l *Slots) SymbolsToIcons() {
	symbolToUnicode := map[Symbol]string{
		WILD:         "🃏",  // Wild
		LOBSTERMANIA: "🦞",  // Lobstermania
		BUOY:         "🎈",  // Buoy
		BOAT:         "⛵",  // Boat
		LIGHTHOUSE:   "🚨",  // Light House
		TUNA:         "🐟",  // Tuna
		CLAM:         "🦪",  // Clam
		SEAGULL:      "🕊️", // Sea Gull
		STARFISH:     "🌟",  // Star Fish
		BONUS:        "🎰",  // Bonus
		SCATTER:      "🎲",  // Scatter
	}
	unicodeSymbols := make([]string, len(l.SpinResult))

	for i, symbol := range l.SpinResult {
		if unicodeChar, ok := symbolToUnicode[symbol]; ok {
			unicodeSymbols[i] = unicodeChar
		} else {
			unicodeSymbols[i] = "❓"
		}
	}
	l.PrettyReel = unicodeSymbols
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
		NewReel([]WeightedSymbol{
			{Symbol: WILD, Weight: 2},
			{Symbol: LOBSTERMANIA, Weight: 4},
			{Symbol: BUOY, Weight: 4},
			{Symbol: BOAT, Weight: 6},
			{Symbol: LIGHTHOUSE, Weight: 5},
			{Symbol: TUNA, Weight: 6},
			{Symbol: CLAM, Weight: 6},
			{Symbol: SEAGULL, Weight: 5},
			{Symbol: STARFISH, Weight: 5},
			{Symbol: BONUS, Weight: 2},
			{Symbol: SCATTER, Weight: 2},
		}),
		NewReel([]WeightedSymbol{
			{Symbol: WILD, Weight: 2},
			{Symbol: LOBSTERMANIA, Weight: 4},
			{Symbol: BUOY, Weight: 4},
			{Symbol: BOAT, Weight: 4},
			{Symbol: LIGHTHOUSE, Weight: 4},
			{Symbol: TUNA, Weight: 4},
			{Symbol: CLAM, Weight: 6},
			{Symbol: SEAGULL, Weight: 6},
			{Symbol: STARFISH, Weight: 5},
			{Symbol: BONUS, Weight: 5},
			{Symbol: SCATTER, Weight: 2},
		}),
		NewReel([]WeightedSymbol{
			{Symbol: WILD, Weight: 1},
			{Symbol: LOBSTERMANIA, Weight: 3},
			{Symbol: BUOY, Weight: 5},
			{Symbol: BOAT, Weight: 4},
			{Symbol: LIGHTHOUSE, Weight: 6},
			{Symbol: TUNA, Weight: 5},
			{Symbol: CLAM, Weight: 5},
			{Symbol: SEAGULL, Weight: 5},
			{Symbol: STARFISH, Weight: 6},
			{Symbol: BONUS, Weight: 6},
			{Symbol: SCATTER, Weight: 2},
		}),
		NewReel([]WeightedSymbol{
			{Symbol: WILD, Weight: 4},
			{Symbol: LOBSTERMANIA, Weight: 4},
			{Symbol: BUOY, Weight: 4},
			{Symbol: BOAT, Weight: 4},
			{Symbol: LIGHTHOUSE, Weight: 6},
			{Symbol: TUNA, Weight: 6},
			{Symbol: CLAM, Weight: 6},
			{Symbol: SEAGULL, Weight: 6},
			{Symbol: STARFISH, Weight: 8},
			{Symbol: BONUS, Weight: 0},
			{Symbol: SCATTER, Weight: 2},
		}),

		NewReel([]WeightedSymbol{
			{Symbol: WILD, Weight: 2},
			{Symbol: LOBSTERMANIA, Weight: 4},
			{Symbol: BUOY, Weight: 5},
			{Symbol: BOAT, Weight: 4},
			{Symbol: LIGHTHOUSE, Weight: 7},
			{Symbol: TUNA, Weight: 7},
			{Symbol: CLAM, Weight: 6},
			{Symbol: SEAGULL, Weight: 6},
			{Symbol: STARFISH, Weight: 7},
			{Symbol: BONUS, Weight: 0},
			{Symbol: SCATTER, Weight: 2},
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
				if count(virtualReel, symbol.Symbol) < symbol.Weight && float64(count(virtualReel, symbol.Symbol)) <= float64(symbol.Weight)/float64(reel.Size)*float64(len(virtualReel)) {
					virtualReel = append(virtualReel, symbol.Symbol)
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
	Reels          []*Reel
	PayLines       [][]int
	Payouts        map[Payout]int
	CreditsWagered int
	NumLinesPlayed int
	SpinResult     []Symbol
	BonusesHit     float64
	LineHitRate    float64
	BonusesHitRate float64
	BonusesHitPay  float64
	ScatterHitRate float64
	LineHitPay     float64
	ScatterHitPay  float64
	PrettyReel     []string
}

func NewSlots() *Slots {
	return &Slots{
		Reels:    buildReels(),
		PayLines: buildPayLines(),
		Payouts:  buildPayoutMap(),
	}
}

func (l *Slots) Evaluate(linesPlayed, wagerPerLine int) int {

	totalPayout := 0
	for i := 0; i < linesPlayed; i++ {
		count := 1
		wildCounts := 0
		firstSymbolInPayline := l.SpinResult[l.PayLines[i][0]]
		if firstSymbolInPayline == WILD {
			wildCounts++
		}
		for j := 1; j < 5; j++ {

			currentSymbol := l.SpinResult[l.PayLines[i][j]]
			currentIsWild := currentSymbol == WILD

			if currentIsWild && firstSymbolInPayline == WILD {
				wildCounts++
			}

			if firstSymbolInPayline == BONUS && currentSymbol == WILD || firstSymbolInPayline == WILD && currentSymbol == BONUS {
				break
			}
			if currentSymbol == SCATTER || firstSymbolInPayline == SCATTER {
				break
			}
			if firstSymbolInPayline == WILD && !currentIsWild {
				firstSymbolInPayline = currentSymbol
			}
			if currentSymbol == firstSymbolInPayline || currentIsWild {
				count++
			} else {
				break
			}
		}

		if count >= 2 {
			linePay := int(math.Max(float64(l.Payouts[Payout{Symbol: firstSymbolInPayline, Count: count}]), float64(l.Payouts[Payout{
				Symbol: WILD,
				Count:  wildCounts,
			}])))
			if linePay == 10000 {
				fmt.Println(linePay, count, wildCounts)
				l.PrintReel()
			}
			if linePay == 331 {
				l.BonusesHitPay += 331 * float64(wagerPerLine)
				l.BonusesHitRate++

			} else if linePay > 0 {
				l.LineHitPay += float64(linePay) * float64(wagerPerLine)
				l.LineHitRate++
			}
			totalPayout += int(linePay) * wagerPerLine
		}
	}

	//check scatter
	scatterCount := 0
	for _, symbol := range l.SpinResult {
		if symbol == SCATTER {
			scatterCount++
		}
	}
	if scatterCount >= 3 {
		scatterPay := l.Payouts[Payout{Symbol: SCATTER, Count: scatterCount}] * linesPlayed * wagerPerLine
		totalPayout += scatterPay
		l.ScatterHitPay += float64(scatterPay)
		l.ScatterHitRate++
	}

	return int(totalPayout)
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
	// fmt.Print("\033[H\033[2J")
	fmt.Print("= = = = = = = = = = =\n")

	// l.AnimateSpin()
	l.SymbolsToIcons()
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
	Symbol Symbol
	Count  int
}

type Symbol uint8

const (
	WILD Symbol = iota
	LOBSTERMANIA
	BUOY
	BOAT
	LIGHTHOUSE
	TUNA
	CLAM
	SEAGULL
	STARFISH
	BONUS
	SCATTER
)

func buildPayoutMap() map[Payout]int {
	return map[Payout]int{

		{Symbol: WILD, Count: 2}: 5,
		{Symbol: WILD, Count: 3}: 100,
		{Symbol: WILD, Count: 4}: 500,
		{Symbol: WILD, Count: 5}: 10_000,

		{Symbol: LOBSTERMANIA, Count: 2}: 2,
		{Symbol: LOBSTERMANIA, Count: 3}: 40,
		{Symbol: LOBSTERMANIA, Count: 4}: 200,
		{Symbol: LOBSTERMANIA, Count: 5}: 1_000,

		{Symbol: BUOY, Count: 3}: 25,
		{Symbol: BUOY, Count: 4}: 100,
		{Symbol: BUOY, Count: 5}: 500,

		{Symbol: BOAT, Count: 3}: 25,
		{Symbol: BOAT, Count: 4}: 100,
		{Symbol: BOAT, Count: 5}: 500,

		{Symbol: LIGHTHOUSE, Count: 3}: 10,
		{Symbol: LIGHTHOUSE, Count: 4}: 50,
		{Symbol: LIGHTHOUSE, Count: 5}: 500,

		{Symbol: TUNA, Count: 3}: 10,
		{Symbol: TUNA, Count: 4}: 50,
		{Symbol: TUNA, Count: 5}: 250,

		{Symbol: CLAM, Count: 3}: 5,
		{Symbol: CLAM, Count: 4}: 30,
		{Symbol: CLAM, Count: 5}: 200,

		{Symbol: SEAGULL, Count: 3}: 5,
		{Symbol: SEAGULL, Count: 4}: 30,
		{Symbol: SEAGULL, Count: 5}: 200,

		{Symbol: STARFISH, Count: 3}: 5,
		{Symbol: STARFISH, Count: 4}: 30,
		{Symbol: STARFISH, Count: 5}: 150,

		{Symbol: BONUS, Count: 3}: 331,

		{Symbol: SCATTER, Count: 3}: 5,
		{Symbol: SCATTER, Count: 4}: 25,
		{Symbol: SCATTER, Count: 5}: 200,
	}
}
