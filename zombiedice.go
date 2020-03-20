// name:    zombieDice
// author:  johnny
// version: 2020/03/12. WIP
// references:
//    global constants and slices.  https://qvault.io/2019/10/21/how-to-global-constant-maps-and-slices-in-go/
// notes: actual gameplay still not done.
//   * 2020-03-12 - started using ANSI colors.. but they mess up printf formatting, so i have to pad spaces on the colored variables.
package main

import "fmt"
import "math/rand"                                                     // for random numbers
import "time"                                                          // for random seed

import "os"

// CONSTANTS -----------------------------------------------------------
const brain int = 0                                                    // index for brain
const shotgun int = 2                                                  // index for shotgun
const runner int = 1                                                   // index for runner
const green int = 0                                                    // index for green
const yellow int = 1                                                   // index for yellow
const red int = 2                                                      // index for red
const totalNumberOfDice = 13                                           // total # of die

const greenBrain string     = "\033[38;5;15m\033[48;5;2mB\033[39;49m"
const greenRunner string    = "\033[38;5;15m\033[48;5;2mR\033[39;49m"
const greenShotgun string   = "\033[38;5;0m\033[48;5;2mS\033[39;49m"
const yellowBrain string    = "\033[38;5;15m\033[48;5;3mB\033[39;49m"
const yellowRunner string   = "\033[38;5;15m\033[48;5;3mR\033[39;49m"
const yellowShotgun string  = "\033[38;5;0m\033[48;5;3mS\033[39;49m"
const redBrain string       = "\033[38;5;15m\033[48;5;1mB\033[39;49m"
const redRunner string      = "\033[38;5;15m\033[48;5;1mR\033[39;49m"
const redShotgun string     = "\033[38;5;0m\033[48;5;1mS\033[39;49m"
//
const greenDie string       = "\033[38;5;15m\033[48;5;2mG\033[39;49m"
const yellowDie string      = "\033[38;5;15m\033[48;5;3mY\033[39;49m"
const redDie string         = "\033[38;5;15m\033[48;5;1mR\033[39;49m"

// NOTE: because of ANSI-Color escape sequences and printf..   1         2         3
//       you have to space out extra spaces (to 32).  12345678901234567890123456789012
const msgYouWereShotgunned string      = "\033[38;5;1mYou have been DESTROYED!                 \033[38;5;4m"
const msgYouSurvivedAnotherDay string  = "\033[38;5;2mYou're a really AWESOME Zombie!          \033[38;5;4m"


// VARS ----------------------------------------------------------------
var rolld6 int
var x int                                                              // misc var
var y int                                                              // misc var
var color int      // color picked from cup during the setup
var dieColors [][]int                                                   // two dimensional array listing all possible sides for each colored dice
var icon [][]string
var myScore []int
var myCup []int                                                        // dice in cup
var myLeftHand []int                                                       // dice in hand (current roll)
var myRightHand []int                                                  // hand that temporarily holds dice that are not put out of play
var outOfPlay []int                                                    // dice now out of play because of a shotgun roll
var gameMessage string
var gameOutcome bool
var brainTally int
var runnerTally int
var shotgunTally int
var rollCount int
var gameState bool
var percentages [][]float32
var percents []int

var (
  ansiReset  = "\033[39;49m"
  ansiGreen  = "\033[38;5;2m"
  ansiRed    = "\033[38;5;1m"
  ansiBlue   = "\033[38;5;4m"
//  ansiYellow = "\033[38;5;3m"
)


// FUNCTIONS ===========================================================
func debuga(msg string, i int, s []int) {
//  var v int
  var dieColors [3]string

  dieColors[0]=greenDie
  dieColors[1]=yellowDie
  dieColors[2]=redDie
//  fmt.Printf("dbg:%-10s  i:%d  len:[%02d] ",msg, i, len(s))
//  for _, v = range s {
//    fmt.Printf("%-2s ", dieColors[v])
//  }
//  fmt.Println("")
}

func showCupContents(s []int) {
  var v int
  var dieColors [3]string

  dieColors[0]=greenDie
  dieColors[1]=yellowDie
  dieColors[2]=redDie
  fmt.Print(ansiBlue, "┃ ",ansiReset, "random dice sequence: ")
  for _, v = range s {
    fmt.Printf("%-2s ", dieColors[v])
  }
  fmt.Print(ansiBlue, "             ┃", ansiReset, "\n")
}



func printSlice(z int, s []int) {
  fmt.Printf("z=%d :  len=%d    %v      \n", z, len(s), s)
}
func printSlice2(z int, s [][]int) {
  fmt.Printf("z=%d :  len=%d    %#v      \n", z, len(s), s)
}

func randomizeDiceInCup(howManyDice int) (cup []int) {
  var diceColorsAvailable [3]int
  var z int

  diceColorsAvailable[green]=6                                         // there are initially six green die in the bag
  diceColorsAvailable[yellow]=4                                        // there are initially four yellow die in the bag
  diceColorsAvailable[red]=3                                           // there are initially three red die in the bag
  z=0
  for z < howManyDice {
    color = rand.Intn(3)
    if diceColorsAvailable[color] > 0 {                                // if there is still THIS-COLOR remaining
      cup=append(cup,color)                                            // prepoulate
      z+=1                                                             //
      diceColorsAvailable[color]-=1                                    // reduce THIS-COLOR by one
    }
  }
  return
}

func resetScores() []int {
  return []int{0, 0, 0}
}

func updateScore(score int, theScores []int) {
  switch score {
    case brain, runner, shotgun: {
      theScores[score]++
    }
  }
}

func prepDieColors() [][]int {                                          // populate two dimensional array w/each die color and face
  var greenDieSides, yellowDieSides, redDieSides []int
  var ds [][]int

  greenDieSides = []int{brain, brain, brain, runner, runner, shotgun}    // prep the sides of a green die
  yellowDieSides = []int{brain, brain, runner, runner, shotgun, shotgun} // prep the sides of a yellow die
  redDieSides = []int{brain, runner, runner, shotgun, shotgun, shotgun}  // prep the sides of a red die
  ds = append(ds, greenDieSides)                                         // idx0: green face possibilities
  ds = append(ds, yellowDieSides)                                        // idx1: yellow face possibilities
  ds = append(ds, redDieSides)                                           // idx2: red face possibilities
  return ds
}

func prepIcons() [][]string {                                            // populate two dimensional array w/each die color and face
  var greenIcons, yellowIcons, redIcons []string
  var ds [][]string
  // -- prep sides for green, yellow, and red dies (icons) into a two dimensional array
  greenIcons = []string{greenBrain, greenBrain, greenBrain, greenRunner, greenRunner, greenShotgun}
  yellowIcons = []string{yellowBrain, yellowBrain, yellowRunner, yellowRunner, yellowShotgun, yellowShotgun}
  redIcons = []string{redBrain, redRunner, redRunner, redShotgun, redShotgun, redShotgun}
  ds = append(ds, greenIcons)
  ds = append(ds, yellowIcons)
  ds = append(ds, redIcons)
  return ds
}

func prepPercentages() [][]float32 {
  var myPercentages [][]float32
  var greenPercentages []float32
  var yellowPercentages []float32
  var redPercentages []float32

    greenPercentages =  []float32{0.500,0.333,0.166}
    yellowPercentages = []float32{0.333,0.333,0.333}
    redPercentages =    []float32{0.166,0.333,0.500}

    myPercentages = append(myPercentages,greenPercentages)
    myPercentages = append(myPercentages,yellowPercentages)
    myPercentages = append(myPercentages,redPercentages)
    return myPercentages
  }


//func rollResults(theRoll []int) {
func rollResults() {
  var rolld6 int
  var v int                                                            // the type of die being utilized (GREEN, YELLOW, RED)
  var i int                                                            // index var:  current die being utilized
  var resultVisual string
  var tally [3]int                                                     // single roll tally (ie not accumlative)
  var rolledDieOnTable [3]int      // number of die to replenish after roll (ie how many taken out of play)

  rollCount+=1
  fmt.Print(ansiBlue, "┃ ",ansiReset)
// TEMP BELOW
//fmt.Print("\n")
//TEMP ABOVE
  tally[brain]=0
  tally[runner]=0
  tally[shotgun]=0



  for i, v = range myLeftHand {
//  for i, v = range theRoll {
    rolld6=rand.Intn(6)                                                // roll die (RANGOM NUMBER)
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// WILL BE CHANGING all these lines from  ....[v] to .....[i]
//     205, 220
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    percents = append(percents, int((((percentages[myLeftHand[0]][brain] + percentages[myLeftHand[1]][brain] + percentages[myLeftHand[2]][brain])/3)+0.005)*100))
    percents = append(percents, int((((percentages[myLeftHand[0]][runner] + percentages[myLeftHand[1]][runner] + percentages[myLeftHand[2]][runner])/3)+0.005)*100))
    percents = append(percents, int((((percentages[myLeftHand[0]][shotgun] + percentages[myLeftHand[1]][shotgun] + percentages[myLeftHand[2]][shotgun])/3)+0.005)*100))


// NOTE:  Maybe the SWITCH order needs to be SHOTGUN, BRAIN, RUNNER
    switch dieColors[v][rolld6] {                                       // was the roll a BRAIN, RUNNER, or SHOTGUN
      // SHOTGUN ------------------------------------------------------
      case shotgun: {
//        shotgunTally+=1
        tally[shotgun]+=1
        myScore[shotgun]+=1
//fmt.Printf("roll-s: (%d)  %1s\n", i, icon[v][rolld6])
        rolledDieOnTable[i]=1
        if myScore[shotgun] > 2 {
          gameMessage=msgYouWereShotgunned
          gameOutcome=false
          gameState=false
        } //eoif 3-shotguns
        //  FIX THIS.. when cup is empty.. but you have 3 dice in hand.. you should get one more roll
        if len(myCup)==0 {gameState=false}
      } //eocase shotgun
      // BRAIN --------------------------------------------------------
      case brain:   {
//        brainTally+=1
        tally[brain]+=1
        myScore[brain]+=1
//WINNING
//FOR NOW.. i will say if brains are greater than 7 then quit turn.
// LATER.. maybe do stats of how many REDs, YELLOWs, and GREENs are out of play and base WINNING on this.
        if brainTally > 6 {
          gameMessage=msgYouSurvivedAnotherDay
          gameOutcome=true
          gameState=false
        }
//fmt.Printf("roll-b: (%d)  %1s\n", i, icon[v][rolld6])
        rolledDieOnTable[i]=1
        if len(myCup)==0 {gameState=false}
      } //eocase brain
      // RUNNER -------------------------------------------------------
      case runner: {
//fmt.Printf("roll-r: (%d)  %1s\n", i, icon[v][rolld6])
  //      runnerTally+=1
        tally[runner]+=1
      } //eocase runner
    } //eoswitch dieColors
    resultVisual=resultVisual+icon[v][rolld6]
  } //eofor theRoll

  //now.. move (copy) all runners from leftHand to rightHand
  //AND   forget about brains and shotguns in left hand (they will go out of play and have been already tally'd)
  myRightHand=nil
  for i=0; i<3; i++ {
//fmt.Printf("i= %d\n",i)
    switch rolledDieOnTable[i] {
      case 0:  { // was a runner
//debuga("70-lhand",i,myLeftHand)
//debuga("70-rhand",i,myRightHand)
        myRightHand=append(myRightHand, myLeftHand[i])
//debuga("71-lhand",i,myLeftHand)
//debuga("71-rhand",i,myRightHand)
      } //eocase0 rolledDieOnTable
      case 1:  { // was a brain or shotgun
//debuga("72-cup",i,myCup)
//debuga("72-rhand",i,myRightHand)
        if len(myCup)>0 {                                              // Be sure there are dice left in cup to take
          myRightHand=append(myRightHand, myCup[len(myCup)-1])                     // get another die from cup
          myCup=myCup[:len(myCup)-1]                                     //   therefore reducing the cup qty
        } else {
//MAYBE this means i won??? (since no more dice left in cup and havent received 3 shotguns)
          gameMessage=msgYouSurvivedAnotherDay
          gameState=false
        }//eoif qtyDiceLeftInCup
//debuga("73-cup",i,myCup)
//debuga("73-rhand",i,myRightHand)
      } //eocase1 rolledDieOnTable
    } //eoswitch rolledDieOnTable
  } //eofor rolledDieOnTable

  // now take what is in right hand and put it back in left hand
  myLeftHand=myRightHand


//  // now remove (removable dice) from myLeftHand (left)
//  for i=0; i<3; i++ {
//    if rolledDieOnTable[i] == 1 {
//debuga("74-hand",i,myLeftHand)
//        myLeftHand=append(myLeftHand[:i], myLeftHand[i+1:]...)                     // remove BRAIN/SHOTGUN from hand
//debuga("75-hand",i,myLeftHand)
//    }
//  }
//  // now replenish (now removed dice) into myLeftHand (left)
//  for i=0; i<(3-len(myLeftHand)); i++ {
//debuga("78-hand",i,myLeftHand)
//    myLeftHand=append(myLeftHand[:i], myLeftHand[i+1:]...)                     // remove BRAIN/SHOTGUN from hand
//debuga("79-hand",i,myLeftHand)
//  }
//  fmt.Printf("roll %02d: %-3s   tally: brains:%02d    runners:%02d    shotguns:%02d",rollCount, resultVisual, brainTally, runnerTally, shotgunTally)
  fmt.Printf("roll %02d: %-3s  pre-roll%%: b:%2d%% s:%2d%%    tally: b:%02d r:%02d s:%02d",rollCount, resultVisual,
    percents[brain], percents[shotgun],
    tally[brain], tally[runner], tally[shotgun] )
  fmt.Print(ansiBlue, " ┃", ansiReset, "\n")
}

// MAIN ===============================================================
func main() {
  var d1,d2,d3 int

  // INIT -------------------------------------------------------------
  gameState=true
  rand.Seed(time.Now().UnixNano())
  myScore=resetScores()                                                // reset 3 scores (brains, runners, shotguns) to zeroes
  dieColors=prepDieColors()
  icon=prepIcons()
  percentages=prepPercentages()
  myCup = randomizeDiceInCup(totalNumberOfDice)                        // prepopulate the random dice order ie. the order that dice will be pulled from the cup
  // Title ------------------------------------------------------------
  fmt.Print(ansiBlue)
  fmt.Print("┏━━━━━━━━━━━━━━━━┓\n")
  fmt.Print("┃  Zombie Dice   ┃\n")
  fmt.Print("┣━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
  fmt.Print(ansiReset, "\n")
  showCupContents(myCup)
  fmt.Print( "┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫\n")
  d1=myCup[len(myCup)-1]                                               // get first 3 (already randomized) dice from cup
  myCup=myCup[:len(myCup)-1]
  d2=myCup[len(myCup)-1]
  myCup=myCup[:len(myCup)-1]
  d3=myCup[len(myCup)-1]
  myCup=myCup[:len(myCup)-1]
  myLeftHand=append(myLeftHand, d1,d2,d3)                                      //   and place in hand for first roll
  for {
//    rollResults(myLeftHand)
    rollResults()
    if !gameState {
      break
    }
  }


  if (len(myCup)==0 && shotgunTally<3) {
    gameMessage=msgYouSurvivedAnotherDay
  }

//  gameMessage="johnny was here"
  fmt.Print(ansiBlue)
  fmt.Print( "┣━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫\n")
  fmt.Printf("┃  Stats         ┃ %-42s   ┃\n",gameMessage)
  fmt.Print( "┣━━━━━━━━━━━━━━━━┫                                             ┃\n")
  fmt.Printf("┃ Rolls:    %02d   ┃                                             ┃\n", rollCount)
  fmt.Print( "┣━━━━━━━━━━━━━━━━┫                                             ┃\n")
  fmt.Printf("┃ Braaains: %02d   ┃                                             ┃\n", myScore[brain])
//  fmt.Printf("┃ Runners:  %02d   ┃                                             ┃\n", myScore[runner])
  fmt.Printf("┃ Shotguns: %02d   ┃                                     madRobot┃\n", myScore[shotgun])
  fmt.Print( "┗━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
  fmt.Print(ansiReset, "\n")

  if gameOutcome {
    os.Exit(0)
  } else {
    os.Exit(1)
  }

}
// END  =========================================================
