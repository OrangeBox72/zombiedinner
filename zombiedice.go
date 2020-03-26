// name:    zombieDice
// author:  johnny
// version: 2020/03/12. WIP
// references:
//    global constants and slices.  https://qvault.io/2019/10/21/how-to-global-constant-maps-and-slices-in-go/
// notes: actual gameplay still not done.
//   * 2020-03-24 - started using ANSI colors.. (happy birthday my-hanh)
package main

import (
  "fmt"
  "math/rand"                                                          // for random numbers
  "strconv"                                                            // string conversions
  "time"                                                               // for random seed
  "os"
  "github.com/fatih/color"
)

// GLOBAL CONSTANTS ---------------------------------------------------
const brain int = 0                                                    // index for brain
const shotgun int = 2                                                  // index for shotgun
const runner int = 1                                                   // index for runner
const green int = 0                                                    // index for green
const yellow int = 1                                                   // index for yellow
const red int = 2                                                      // index for red
const totalNumberOfDice = 13                                           // total # of die
const totalNumberOfGreenDice = 6
const totalNumberOfYellowDice = 4
const totalNumberOfRedDice = 3

const greenBrain string     = "\033[38;5;15m\033[48;5;2mB\033[39;49m"  // here are all of the icon-faces
const greenRunner string    = "\033[38;5;15m\033[48;5;2mR\033[39;49m"
const greenShotgun string   = "\033[38;5;0m\033[48;5;2mS\033[39;49m"
const yellowBrain string    = "\033[38;5;15m\033[48;5;3mB\033[39;49m"
const yellowRunner string   = "\033[38;5;15m\033[48;5;3mR\033[39;49m"
const yellowShotgun string  = "\033[38;5;0m\033[48;5;3mS\033[39;49m"
const redBrain string       = "\033[38;5;15m\033[48;5;1mB\033[39;49m"
const redRunner string      = "\033[38;5;15m\033[48;5;1mR\033[39;49m"
const redShotgun string     = "\033[38;5;0m\033[48;5;1mS\033[39;49m"
//
const greenDie string       = "\033[38;5;15m\033[48;5;2mg\033[39;49m"
const yellowDie string      = "\033[38;5;15m\033[48;5;3my\033[39;49m"
const redDie string         = "\033[38;5;15m\033[48;5;1mr\033[39;49m"

var msgYouSurvivedAnotherDay string

// VARS ----------------------------------------------------------------
var rolld6 int
var x int                                                              // misc var
var y int                                                              // misc var
var dieColors [][]int                                                  // two dimensional array listing all possible sides for each colored dice
var icon [][]string
var myScore []int
var tally [3]int                                                       // single roll tally (ie not accumlative)
var myCup []int                                                        // dice in cup
var myLeftHand []int                                                   // dice in hand (current roll)
var myRightHand []int                                                  // hand that temporarily holds dice that are not put out of play
var outOfPlay []int                                                    // dice now out of play because of a shotgun roll
var gameMessage string
var gameOutcome bool
var roundIdx int
var gameState bool
var diePercentages [][]float32
var handPercentages []int
var possiblePercentages []int
var cupVisual string
var spaces []rune=[]rune("             ")                              // some spaces because i'm too dumb to know how to do this a better way.
                                                                       // PLUS.. i like how 'rune' sounds.
// FUNCTIONS ===========================================================
func showCupContents(c []int) (cupContents string){
  var v int
  cupContents = ""
  for _, v = range c {
    cupContents = cupContents + icon[v][6]
  }
  return cupContents
}

func getCupPercentages(c []int, percentageType int) (cupPercs int){
  var v int
  var i int
  var percs float32

//fmt.Printf("\n")

  percs=0.0
  for i, v = range c {
//fmt.Printf("%2d - [%2d][%2d]=%2.3f\n",i,v,percentageType,diePercentages[v][percentageType])
    percs=percs+diePercentages[v][percentageType]
  }
//fmt.Printf("pf:%2.3f  ", ((percs/float32(i))+0.005)*100 )
//fmt.Printf("pi:%2d    ", int(((percs/float32(i))+0.005)*100)     )
  cupPercs = int(((percs/float32(i))+0.005)*100)
//  fmt.Printf("cp: %4d\n", cupPercs)
  return
}

func randomizeDiceInCup(howManyDice int) (cup []int) {
  var diceColorsAvailable [3]int
  var randomColor int
  var z int

  diceColorsAvailable[green]=totalNumberOfGreenDice                    // there are initially six green die in the bag
  diceColorsAvailable[yellow]=totalNumberOfYellowDice                  // there are initially four yellow die in the bag
  diceColorsAvailable[red]=totalNumberOfRedDice                        // there are initially three red die in the bag
  z=0
  for z < howManyDice {
    randomColor = rand.Intn(3)
    if diceColorsAvailable[randomColor] > 0 {                          // if there is still THIS-COLOR remaining
      cup=append(cup,randomColor)                                      // prepoulate
      z+=1                                                             //
      diceColorsAvailable[randomColor]-=1                              // reduce THIS-COLOR by one
    }
  }
  return
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
  //    NOTE: in order to show blank dice (colors only).. ie before dice is rolled, i have added a 7th item to these 6-sided dies.
  //          i.e.. greenDie, yellowDie, and redDie.
  greenIcons = []string{greenBrain, greenBrain, greenBrain, greenRunner, greenRunner, greenShotgun, greenDie}
  yellowIcons = []string{yellowBrain, yellowBrain, yellowRunner, yellowRunner, yellowShotgun, yellowShotgun, yellowDie}
  redIcons = []string{redBrain, redRunner, redRunner, redShotgun, redShotgun, redShotgun, redDie}
  ds = append(ds, greenIcons)
  ds = append(ds, yellowIcons)
  ds = append(ds, redIcons)
  return ds
}

func rollResults() {
  var rolld6 int
  var v int                                                            // the type of die being utilized (GREEN, YELLOW, RED)
  var i int                                                            // index var:  current die being utilized
  var rolledVisual string                                              // visual die. ANSI colored with rolled value showing.
  var handVisual string                                                // visual die. ANSI colored.. but no face value. ie before roll.
  var rolledDieOnTable [3]int      // number of die to replenish after roll (ie how many taken out of play)

  roundIdx+=1
  fmt.Print(color.BlueString("┃ "))
  tally[brain]=0
  tally[runner]=0
  tally[shotgun]=0

  handPercentages = nil
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[0]][brain] + diePercentages[myLeftHand[1]][brain] + diePercentages[myLeftHand[2]][brain])/3)+0.005)*100))
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[0]][runner] + diePercentages[myLeftHand[1]][runner] + diePercentages[myLeftHand[2]][runner])/3)+0.005)*100))
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[0]][shotgun] + diePercentages[myLeftHand[1]][shotgun] + diePercentages[myLeftHand[2]][shotgun])/3)+0.005)*100))
  possiblePercentages = nil
  possiblePercentages = append(possiblePercentages, getCupPercentages(myCup, brain) )
  possiblePercentages = append(possiblePercentages, 33)
  possiblePercentages = append(possiblePercentages, getCupPercentages(myCup, shotgun) )
  for i, v = range myLeftHand {
    rolld6=rand.Intn(6)                                                // roll die (RANGOM NUMBER)
    cupVisual = showCupContents(myCup)
    cupVisual = cupVisual + string(spaces[0:(13-len(myCup))])
    // NOTE:  (I think) the SWITCH order needs to be SHOTGUN, BRAIN, RUNNER to ensure 3-SHOTGUNS will stop before BRAINS are added to score.
    switch dieColors[v][rolld6] {                                       // was the roll a BRAIN, RUNNER, or SHOTGUN
      // SHOTGUN ------------------------------------------------------
      case shotgun: {
        tally[shotgun]+=1
        myScore[shotgun]+=1
        rolledDieOnTable[i]=1
        if myScore[shotgun] > 2 {
          gameMessage=color.RedString("You have been Destroyed!") + "   (you almost had " + strconv.Itoa(tally[brain]) + " braaains.)"
//          gameMessage=color.RedString("You have been Destroyed!  (you almost had 999 braaains.)")
          myScore[brain]=0                                             // no BRAINS for you! You got blasted!
          gameOutcome=false
          gameState=false
        } //eoif 3-shotguns
        //  FIX THIS.. when cup is empty.. but you have 3 dice in hand.. you should get one more roll
        if len(myCup)==0 {gameState=false}
      } //eocase shotgun
      // BRAIN --------------------------------------------------------
      case brain:   {
        if gameState {
          tally[brain]+=1
          myScore[brain]+=1
          if tally[brain] > 6 {                                            //WINNING
            gameMessage=msgYouSurvivedAnotherDay                         //  FOR NOW.. i will say if brains are greater than 7 then quit turn.
            gameOutcome=true
            gameState=false
          }
          rolledDieOnTable[i]=1
          if len(myCup)==0 {gameState=false}
        } //eoifGameState
      } //eocase brain
      // RUNNER -------------------------------------------------------
      case runner: {
        tally[runner]+=1
      } //eocase runner
    } //eoswitch dieColors
    rolledVisual=rolledVisual+icon[v][rolld6]
    handVisual=handVisual+icon[v][6]
  } //eofor theRoll

  //now.. move (copy) all runners from leftHand to rightHand
  //AND   forget about brains and shotguns in left hand (they will go out of play and have been already tally'd)
  myRightHand=nil
  for i=0; i<3; i++ {
    switch rolledDieOnTable[i] {
      case 0:  { // was a runner
        myRightHand=append(myRightHand, myLeftHand[i])
      } //eocase0 rolledDieOnTable
      case 1:  { // was a brain or shotgun
        if len(myCup)>0 {                                              // Be sure there are dice left in cup to take
          myRightHand=append(myRightHand, myCup[len(myCup)-1])         // get another die from cup
          myCup=myCup[:len(myCup)-1]                                   //   therefore reducing the cup qty
        } else {
          gameMessage=msgYouSurvivedAnotherDay                         //MAYBE this means i won??? (since no more dice left in cup and havent received 3 shotguns)
          gameState=false
        } //eoif qtyDiceLeftInCup
      } //eocase1 rolledDieOnTable
    } //eoswitch rolledDieOnTable
  } //eofor rolledDieOnTable

  myLeftHand=myRightHand                                               // now take what is in right hand and put it back in left hand
  fmt.Printf("round %02d: (in hand: %-3s b:%2d%% s:%2d%%) (in play: b:%2d%% s:%2d%%  %-10s)    rolled: %-3s",
    roundIdx,
    handVisual, handPercentages[brain], handPercentages[shotgun],
    possiblePercentages[brain], possiblePercentages[shotgun], cupVisual,
    rolledVisual)
    fmt.Print(color.BlueString(" ┃\n"))

}

// MAIN ===============================================================
func main() {
  var d1,d2,d3 int

  // INIT -------------------------------------------------------------
msgYouSurvivedAnotherDay=color.GreenString("You're a really AWESOME Zombie!")

  gameState=true
  rand.Seed(time.Now().UnixNano())
  myScore=[]int{0, 0, 0}                                               // reset 3 scores (brains, runners, shotguns) to zeroes
  dieColors=prepDieColors()
  icon=prepIcons()
  diePercentages = append(diePercentages, []float32{0.500,0.333,0.167})  //GREEN:  brain, runner, shotgun - percentages
  diePercentages = append(diePercentages, []float32{0.333,0.333,0.333})  //YELLOW: brain, runner, shotgun - percentages
  diePercentages = append(diePercentages, []float32{0.167,0.333,0.500})  //RED:    brain, runner, shotgun - percentages

  myCup = randomizeDiceInCup(totalNumberOfDice)                        // prepopulate the random dice order ie. the order that dice will be pulled from the cup

  // Title ------------------------------------------------------------
  color.Blue("┏━━━━━━━━━━━━━━━━┓")
  color.Blue("┃  Zombie Dice   ┃")
  color.Blue("┣━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")


  d1=myCup[len(myCup)-1]                                               // get first 3 (already randomized) dice from cup
  myCup=myCup[:len(myCup)-1]
  d2=myCup[len(myCup)-1]
  myCup=myCup[:len(myCup)-1]
  d3=myCup[len(myCup)-1]
  myCup=myCup[:len(myCup)-1]
  myLeftHand=append(myLeftHand, d1,d2,d3)                                      //   and place in hand for first roll
  for {
    rollResults()
    if !gameState {
      break
    }
  }

  if (len(myCup)==0 && tally[shotgun]<3) {
    gameMessage=msgYouSurvivedAnotherDay
  }

  color.Blue("┣━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫")
  color.Blue("┃  Stats         ┃ %-80s  ┃\n", gameMessage)
  color.Blue("┣━━━━━━━━━━━━━━━━┫                                                                          ┃")
  color.Blue("┃ Rolls:    %02d   ┃                                                                          ┃", roundIdx)
  color.Blue("┣━━━━━━━━━━━━━━━━┫                                                                          ┃")
  color.Blue("┃ Braaains: %02d   ┃                                                                          ┃", myScore[brain])
  color.Blue("┃ Shotguns: %02d   ┃                                                                  madRobot┃", myScore[shotgun])
  color.Blue("┗━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")

  if gameOutcome {
    os.Exit(0)
  } else {
    os.Exit(1)
  }

}
// END  =========================================================
