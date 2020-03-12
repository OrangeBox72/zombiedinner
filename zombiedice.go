// name:    zombieDice
// author:  johnny
// version: 2018/10/09 19:30	complete rewrite
// notes:	actual gameplay still not done.
//			currently, only completed logic for rolls, die, color, face, and scoring
//			!!!!!	next step.. convert to prepopulated roll var to a slice
//					and remove an element when it becomes brain or shotgun
//					also.. adjust the min/max of the slice dynamically for the for/next to work properly
package main

import "fmt"
import "math/rand"								// for random numbers
import "time"									// for random seed


//	MAIN	=========================================================
func main() {
//	CONSTANTS =======================================================
	const brain int = 0							// index for brain
	const shotgun int = 1						// index for shotgun
	const runner int = 2						// index for runner
	const green int = 0							// index for green
	const yellow int = 1						// index for yellow
	const red int = 2							// index for red

//	VARS	=========================================================
//			informational vars
//

	dieColorName := make([]string, 3)	// color names for the dice		fyi 'make' creates a slice. ie this is a slice
	dieColorName[green] = "Green"
	dieColorName[yellow] = "Yellow"
	dieColorName[red] = "Red"

	dieFace := make([]string, 3)
	dieFace[brain] = "Brain"
	dieFace[runner] = "Runner"
	dieFace[shotgun] = "SHOTGUN"

	dieSides := [][]int{}
	greenDieSides := []int{brain, brain, brain, runner, runner, shotgun}
	yellowDieSides := []int{brain, brain, runner, runner, shotgun, shotgun}
	redDieSides := []int{brain, runner, runner, shotgun, shotgun, shotgun}
	dieSides = append(dieSides, greenDieSides)			// presume that first append is at index 0 which is green
	dieSides = append(dieSides, yellowDieSides)
	dieSides = append(dieSides, redDieSides)


//			game prep vars
//
	var dieInCup int						// how many dice in cup
	var dieOutOfPlay int					// how many dice (shotguns/brains) that are now out of play
	dieOrder := make([]int, 13)				// the prepopulated random order of dice pulled from the bag
	var dieColor int						// color picked from cup during the setup

	var diceQuantity [3]int
	diceQuantity[green] = 6					// there are initially six green die in the bag
	diceQuantity[yellow] = 4				// there are initially four yellow die in the bag
	diceQuantity[red] = 3					// there are initially three red die in the bag

//			normal vars
//
	var x int								// misc var
	var y int								// misc var
	var rolld6 int

	var score [3]int
	score[brain] = 0   // how many brains can you eat ;-)
	score[shotgun] = 0  // just for stats.. doesnt affect game outcome
	score[runner] = 0 // max shotgun hits you can receive


//	INIT	=========================================================
fmt.Println("len of dieOrder: ", len(dieOrder), ", ", cap(dieOrder))
	dieInCup = len(dieOrder)
	dieOutOfPlay = 13 - dieInCup
	rand.Seed(time.Now().UnixNano())

	// titles
	fmt.Println("\n\nZombie Dice")
	fmt.Println("=============================================================")
	//	prepopulate the random dice order ie. the order that dice will be pulled from the cup
	x=0
	for x < dieInCup {
		dieColor = rand.Intn(3)
		if diceQuantity[dieColor] > 0 {				// if there is still THIS-COLOR remaining
			dieOrder[x]=dieColor					// prepoulate
			x+=1									//
			diceQuantity[dieColor]-=1				// reduce THIS-COLOR by one
		}
	}

	//	MAIN	=========================================================
	y = 0
	for y < dieInCup {
		rolld6 = rand.Intn(6)						// roll me a die
		score[dieSides[dieOrder[y]][rolld6]]+=1		// change the current-round score for this particular die roll
		fmt.Printf("die %2d - %-6s %-7s   score (brains:%2d  shotguns:%2d  runners:%2d    dieInCup: %2d , outOfPlay: %2d)\n", y+1, dieColorName[dieOrder[y]], dieFace[dieSides[dieOrder[y]][rolld6]], score[brain], score[shotgun], score[runner], dieInCup, dieOutOfPlay )
		y+=1
		if score[shotgun] >= 3 {
			fmt.Println("You have been DESTROYED!")
			fmt.Printf("\tyou had (brains:%2d  shotguns:%2d  runners:%2d)\n", score[brain], score[shotgun], score[runner] )
			y=99
		}
	}
	fmt.Println("\n\n")
}
//	END		=========================================================
