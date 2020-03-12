// name:    zombieDice
// author:  johnny
// version: 2018/10/09 19:30 complete rewrite
// notes: actual gameplay still not done.
//   currently, only completed logic for rolls, die, color, face, and scoring
//   !!!!! next step.. convert to prepopulated roll var to a slice
//     and remove an element when it becomes brain or shotgun
//     also.. adjust the min/max of the slice dynamically for the for/next to work properly
package main

import "fmt"
import "math/rand"                                                     // for random numbers
import "time"                                                          // for random seed

//import "os"

// CONSTANTS -----------------------------------------------------------
const brain int = 0                                                    // index for brain
const shotgun int = 1                                                  // index for shotgun
const runner int = 2                                                   // index for runner
const green int = 0                                                    // index for green
const yellow int = 1                                                   // index for yellow
const red int = 2                                                      // index for red
const totalNumberOfDice = 13                                           // total # of die

// VARS ----------------------------------------------------------------
var x int        // misc var
var dieColor int      // color picked from cup during the setup
var diceColorQuantity [3]int
var dieColorName []string                                              // die color description
var dieFace []string                                                   // die face description
// FUNCTIONS ===========================================================
func printSlice(z int, s []int) {
  fmt.Printf("z=%d :  len=%d    %#v      \n", z, len(s), s)
}

func randomizeDiceInCup(howManyDice int) (cup []int) {
  x=0
  for x < howManyDice {
    dieColor = rand.Intn(3)
    if diceColorQuantity[dieColor] > 0 {                                      // if there is still THIS-COLOR remaining
printSlice(x, cup)
      cup=append(cup,dieColor)                                                // prepoulate
      x+=1                                                                //
      diceColorQuantity[dieColor]-=1                                           // reduce THIS-COLOR by one
    }
  }
  return
}

// MAIN ================================================================
func main() {

// VARS ----------------------------------------------------------------
//   informational vars
//
  dieColorName=append(dieColorName,"Green")                            // 0  NOTE: These appends HAVE TO BE IN THIS ORDER
  dieColorName=append(dieColorName,"Yellow")                           // 1        so that they match up with the CONSTANTS above
  dieColorName=append(dieColorName,"Red")                              // 2        i.e. green=0 .. append(green) needs to be 0 index

  dieFace=append(dieFace,"Brain")                                      // 0  NOTE: These appends HAVE TO BE IN THIS ORDER
  dieFace=append(dieFace,"Runner")                                     // 1        so that they match up with the CONSTANTS above
  dieFace=append(dieFace,"SHOTGUN")                                    // 2        i.e. brain=0 .. append(brain) needs to be 0 index

 greenDieSides := []int{brain, brain, brain, runner, runner, shotgun}
 yellowDieSides := []int{brain, brain, runner, runner, shotgun, shotgun}
 redDieSides := []int{brain, runner, runner, shotgun, shotgun, shotgun}
 dieSides := [][]int{}                                                 // one array listing all possible sides for each colored dice
 dieSides = append(dieSides, greenDieSides)                            // 0. since dieSides is NULL/just initiated.. the first index will be 0.
 dieSides = append(dieSides, yellowDieSides)                           // 1
 dieSides = append(dieSides, redDieSides)                              // 2


//   game prep vars
//
 var numDieInCup int                                                   // how many dice in cup
 var dieOutOfPlay int                                                  // how many dice (shotguns/brains) that are now out of play
// dieInCup := make([]int, 13)                                           // the prepopulated random order of dice pulled from the bag
 dieInCup := []int{}                                           // the prepopulated random order of dice pulled from the bag
 
 diceColorQuantity[green] = 6     // there are initially six green die in the bag
 diceColorQuantity[yellow] = 4    // there are initially four yellow die in the bag
 diceColorQuantity[red] = 3     // there are initially three red die in the bag

//   normal vars
//
 var y int        // misc var
 var rolld6 int

 var score [3]int
 score[brain] = 0   // how many brains can you eat ;-)
 score[shotgun] = 0  // just for stats.. doesnt affect game outcome
 score[runner] = 0 // max shotgun hits you can receive


// INIT =========================================================
 numDieInCup = len(dieInCup)
 dieOutOfPlay = totalNumberOfDice - numDieInCup
 rand.Seed(time.Now().UnixNano())

 // titles
 fmt.Println("\n\nZombie Dice")
 fmt.Println("=============================================================")
  var myCup []int
  printSlice(98, myCup)
 // prepopulate the random dice order ie. the order that dice will be pulled from the cup
  myCup = randomizeDiceInCup(totalNumberOfDice)
  printSlice(99, myCup)




 // MAIN =========================================================
 y = 0
 for y < numDieInCup {
  rolld6 = rand.Intn(6)      // roll me a die
  score[dieSides[dieInCup[y]][rolld6]]+=1  // change the current-round score for this particular die roll
  fmt.Printf("die %2d - %-6s %-7s   score (brains:%2d  shotguns:%2d  runners:%2d    numDieInCup: %2d , outOfPlay: %2d)\n", y+1, dieColorName[dieInCup[y]], dieFace[dieSides[dieInCup[y]][rolld6]], score[brain], score[shotgun], score[runner], numDieInCup, dieOutOfPlay )
  y+=1
  if score[shotgun] >= 3 {
   fmt.Println("You have been DESTROYED!")
   fmt.Printf("\tyou had (brains:%2d  shotguns:%2d  runners:%2d)\n", score[brain], score[shotgun], score[runner] )
   y=99
  }
 }
 fmt.Println("\n\n")
}

// END  =========================================================

