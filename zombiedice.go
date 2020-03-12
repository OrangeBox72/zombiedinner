package main

import "fmt"
import "math/rand"

func main() {
	var dieColor int // temp die.  contains die color
	var roll int     //
	//	var rollSide int // the side 1-6 of the die that was rolled
	var diceInHand int = 0
	var rollResults [3]int // holds  3 results of rolled die

	//	var dice int = 0
	var diceInCup int
	//	var rollOne int = 0
	//	var rollTwo int = 0
	//	var rollThree int = 0

	var brain int = 0   // index for brain
	var shotgun int = 1 // index for shotgun
	var runner int = 2  // index for runner.  note: this is last (for void dice reduction)
	//	var greenDieQty int = 6
	//	var yellowDieQty int = 4
	//	var redDieQty int = 3
	//	var greenDie int = 0
	//	var yellowDie int = 1
	//	var redDie int = 2
	//      [0]=num of green  [1]=num of yellow  [2]=num of red
	var differentDiceColors [3]int
	differentDiceColors[brain] = 6
	differentDiceColors[shotgun] = 3
	differentDiceColors[runner] = 4

	var resultTotals [3]int
	resultTotals[brain] = 0   // how many brains can you eat ;-)
	resultTotals[shotgun] = 0 // max shotgun hits you can receive
	resultTotals[runner] = 0  // just for stats.. doesnt affect game outcome

	//	var cupOfDice [13]int

	//	var brainsEaten int = 0
	//	var shotgunBlasts int = 0
	//	var shotgunBlastsMax int = 2 // max = 3  or  0-2

	dieFace := make([]string, 4)
	dieFace[brain] = "Brain"
	dieFace[runner] = "Runner"
	dieFace[shotgun] = "SHOTGUN"
	dieFace[3] = "JOHNNYTEST"

	//rows green=0, yellow=1, red=2
	dieSides := [][]int{}
	greenDieSides := []int{brain, brain, brain, runner, runner, shotgun}
	yellowDieSides := []int{brain, brain, runner, runner, shotgun, shotgun}
	redDieSides := []int{brain, runner, runner, shotgun, shotgun, shotgun}
	dieSides = append(dieSides, greenDieSides)
	dieSides = append(dieSides, yellowDieSides)
	dieSides = append(dieSides, redDieSides)

	diceInCup = 13
	diceInHand = 0
	for diceInCup > 0 && resultTotals[shotgun] < 3 {
		if diceInHand < 3 {
			dieColor = rand.Intn(3)                // is the die green,yellow, or red
			if differentDiceColors[dieColor] > 0 { //  we can only continue if there was that color of die still avail
				diceInCup--                     // take a die out of the cup ONLY if that color still exists in cup
				diceInHand++                    // now you have x of 3
				differentDiceColors[dieColor]-- //   yes.. color was still in cup, so reduce num of that color avail
				//		roll=rand.Intn(6)			// roll it and get a side.
				roll = dieSides[dieColor][rand.Intn(6)] // did the color/side result in brain,run,shot
				resultTotals[roll]++                    // increment the stats
				rollResults[diceInHand] = roll          // keep 3 results..
				if roll < runner {
					//  need something here
				}
				fmt.Println("diceincup: ", diceInCup, " diceinhand: ", diceInHand, ", rolled: ", dieFace[roll], " .")
			}
		}
	}
}
