// name:    zombieDice
// author:  johnny
// version: 2018/10/09 19:30.  Complete rewrite
// stds:    space delimited
// notes:   actual gameplay still not done.
//          currently, only completed logic for rolls, die, color, face, and scoring
//          !!!!!	next step.. convert to prepopulated roll var to a slice
//          and remove an element when it becomes brain or shotgun
//          also.. adjust the min/max of the slice dynamically for the for/next to work properly
// urls:    https://github.com/logrusorgru/aurora
package main

import "fmt"
import "math/rand"                                                     // for random numbers
import "time"                                                          // for random seed
import . "github.com/logrusorgru/aurora"                               // for ANSI color output.. url: https://godoc.org/github.com/fatih/color

// VARS GLOBAL ========================================================
//   CONSTANTS --------------------------------------------------------
const brain int = 0                                                   // index for brain
const shotgun int = 1                                                 // index for shotgun
const runner int = 2                                                  // index for runner
const green int = 0                                                   // index for green    fyi - (6 green die)
const yellow int = 1                                                  // index for yellow   fyi - (4 yellow die)
const red int = 2                                                     // index for red      fyi - (3 red die)
const white int = 99                                                  // ANSI color white - arbitrary value
const cupMax int = 3                                                  // how many dice can the cup hold


func getDieColor(x int) Color {
  if x == green {
    return GreenFg              //|BlackBg
  } else if x == red {
    return RedFg                //|BlackBg
  } else if x == yellow {
    return BrownFg              //|BlackBg
  } else if x == white {
    return BlackFg               //|BlackBg
  } else {
    return BlackFg               //|BlackBg
  }
}



// MAIN	===============================================================
func main() {

// VARS  ==============================================================
// informational vars
  var dieColorName [3]string                                            // color names for the dice.  (human readable)
  dieColorName[green] = "Green"
  dieColorName[yellow] = "Yellow"
  dieColorName[red] = "Red"

  var dieFace [3]string                                                // face names for dice.  (human readable)
  dieFace[brain] =   "Brain   "
  dieFace[runner] =  "runner  "                                        // NOTE: Padded spaces to end in order to add ANSI later
  dieFace[shotgun] = "SHOTGUN "                                        //       couldnt combine Colorize with Sprintf %-6s

  var diceQuantity [3]int                                              // array of the different die quantities
  diceQuantity[green] = 6                                              // there are six green die
  diceQuantity[yellow] = 4                                             // there are four yellow die
  diceQuantity[red] = 3                                                // there are three red die

  //  dieSides := [][]int{}
  var dieSides [][]int                                                 // array of all the color-icons.  sorted by color, icon.
  greenDieSides := []int{brain, brain, brain, runner, runner, shotgun}
  yellowDieSides := []int{brain, brain, runner, runner, shotgun, shotgun}
  redDieSides := []int{brain, runner, runner, shotgun, shotgun, shotgun}
  dieSides = append(dieSides, greenDieSides)                           // populate greens first.   (NOTE: i used slice to allow utilize the append function / readability)
  dieSides = append(dieSides, yellowDieSides)                          //   then yellows
  dieSides = append(dieSides, redDieSides)                             //   finally reds

// game prep vars
  var dieAvailable int                                                 // how many dice in cup
  var dieOutOfPlay int                                                 // how many dice (shotguns/brains) that are now out of play
  var dieOrder []int                                                   // the prepopulated random order of dice pulled from the bag (For a slice, DONT put size) !!This is a slice
  var dieColor int                                                     // color picked from cup during the setup
  var dice [3]int                                                      // three dice in cup for rolling..
  var roll int                                                         // index for the three dice you are rolling  (0-2)
  var rollResult int                                                   // the result of a rolled die (ie side facing up)
  var runningScore int                                                 // the running score of this round
  var finalScore int                                                   // the total score of the game

//  normal vars
  var x int                                                            // misc var
  var nextDiceAvail int = 0                                            // index of diceAvail/slice position for drawing 3 dice into cup

  var score [3]int
  score[brain] = 0                                                     // how many brains can you eat ;-)
  score[shotgun] = 0                                                   // just for stats.. doesnt affect game outcome
  score[runner] = 0                                                    // max shotgun hits you can receive


//	INIT	=========================================================
//delthis. left here to remember about len and cap capabilities
//fmt.Println("len of dieOrder: ", len(dieOrder), ", ", cap(dieOrder))

  rand.Seed(time.Now().UnixNano())
  dieAvailable = 13
  dieOutOfPlay = 13 - dieAvailable

  // prepopulate the random dice order ie. the order that dice will be pulled from the available pool
  //   i.e. randomly picking the available colored dice and putting them in the diOrder slice/array.
  x=0
  for x < dieAvailable {
    dieColor = rand.Intn(3)
    if diceQuantity[dieColor] > 0 {                                    // if there is still THIS-COLOR remaining
      dieOrder=append(dieOrder, dieColor)                                             // prepoulate
      x+=1                                                             //
      diceQuantity[dieColor]-=1                                        // reduce THIS-COLOR by one
    }
  }


  // titles
  fmt.Println(Colorize("===================================================================================================", getDieColor(white)))
  fmt.Println(Colorize("== Zombie Dice                                                                                   ==", getDieColor(white)))
  fmt.Println(Colorize("===================================================================================================", getDieColor(white)))

  finalScore=0
  runningScore=0
  roll=0
  for roll < 3 {                                                       // (effectively) roll your 3 dice
    rollResult = dieSides[dieOrder[0]][rand.Intn(6)]                   //   randomly assign value to the roll
    score[rollResult]+=1                                               //   increase the tally for the type of die rolled
//probably need to reduce dieAvailable and update dieOutOfPlay  somewhere in here

    dice[roll] = dieOrder[0]                                           //   keep track of 3 dice

    copy(dieOrder[0:], dieOrder[1:])                                   //   remove first die from the dieOrder, hence eventually rolling die 2 and 3
    dieOrder[len(dieOrder)-1] = 0                                      //   erase last element (write to zero)
    dieOrder = dieOrder[:len(dieOrder)-1]                              //   truncate (end of) slice, effectively reducing size of slice and making position 0 the next die in line

    if rollResult == brain {
      runningScore+=1
    }
    if score[shotgun] == 3 {
//      fmt.Println("OUCH!  You received ", score[shotgun], " shotgun blasts!   At the time, you had ", runningScore, " points.")
//works kinda      fmt.Println(Gray("OUCH!  You received ").BgBlack(), Gray(score[shotgun]).BgBlack(), Gray(" shotgun blasts!   At the time, you had ").BgBlack(), Gray(runningScore).BgBlack(), Gray(" points.").BgBlack())
      fmt.Print(Colorize("OUCH!  You received THREE shotgun blasts!  ", getDieColor(red)))
      fmt.Print(Colorize("At the time, you had ", getDieColor(white)))
      fmt.Print(Colorize(runningScore, getDieColor(white)))
      fmt.Print(Colorize(" points.", getDieColor(white)))
      fmt.Println(Colorize("                          ", getDieColor(white)))
      runningScore=0
    }
    fmt.Print(Colorize("die:", getDieColor(white)))
    fmt.Print(Colorize(roll+1, getDieColor(white)))
    fmt.Print(Colorize("  ", getDieColor(white)))
    fmt.Print(Colorize(dieFace[rollResult],getDieColor(dice[roll])))       // print ANSI color of the name-color of the current (1of3) die roll
    fmt.Print(Colorize("stats(brains:", getDieColor(green)))
    fmt.Print(Colorize(score[brain], getDieColor(green)))
    fmt.Print(Colorize(" shotgun:", getDieColor(red)))
    fmt.Print(Colorize(score[shotgun], getDieColor(red)))
    fmt.Print(Colorize(" runner:", getDieColor(yellow)))
    fmt.Print(Colorize(score[runner], getDieColor(yellow)))
    fmt.Print(Colorize("  dieAvailable:", getDieColor(white)))
    fmt.Print(Colorize(dieAvailable - nextDiceAvail, getDieColor(white)))
    fmt.Print(Colorize(" outOfPlay:", getDieColor(white)))
    fmt.Print(Colorize(dieOutOfPlay, getDieColor(white)))
    fmt.Println(")")
    roll+=1
  }
  finalScore+=runningScore

  fmt.Print(Colorize("Running score: ", getDieColor(white)))           // dieColor(white) is standard GrayFg on BlackBg
  fmt.Print(Colorize(runningScore, getDieColor(white)))
  fmt.Print(Colorize(".  Final score: ", getDieColor(white)))             // dieColor(white) is standard GrayFg on BlackBg
  fmt.Print(Colorize(finalScore, getDieColor(white)))
  fmt.Println(Colorize("                                                                                     ", getDieColor(white)))
  fmt.Println(Colorize("---------------------------------------------------------------------------------------------------", getDieColor(white)))
  fmt.Println(Colorize("===================================================================================================", getDieColor(white)))
  fmt.Println("\n\n")
}
// END  =========================================================
