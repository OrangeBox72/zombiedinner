// name:    zombieDinner
// author:  johnny
// version:
//   2020-04-15 - cleanup
//   2020-03-29 - added issue#6 enhancement: (human input)
//   2020-03-28 - fixed issue#5 enhancement: (outOfPlay stats)
//   2020-03-26 - fixed issue#4 (almost had X brains - stat)
//   2020-03-24 - started using ANSI colors.. (happy birthday my-hanh)
//
// references:
//   global constants and slices.  https://qvault.io/2019/10/21/how-to-global-constant-maps-and-slices-in-go/
//   box drawing.  https://en.wikipedia.org/wiki/Box-drawing_character
//   input handling. https://www.socketloop.com/tutorials/golang-handling-yes-no-quit-query-input
//

package main

import (
  "fmt"
  "math/rand"                                                          // for random numbers
  "strconv"                                                            // string conversions
  "time"                                                               // for random seed
  "os"
  "github.com/fatih/color"
)

/*
 #include <stdio.h>
 #include <unistd.h>
 #include <termios.h>
 char getch(){
          char ch = 0;
       struct termios old = {0};
       fflush(stdout);
       if( tcgetattr(0, &old) < 0 ) perror("tcsetattr()");
       old.c_lflag &= ~ICANON;
       old.c_lflag &= ~ECHO;
       old.c_cc[VMIN] = 1;
       old.c_cc[VTIME] = 0;
       if( tcsetattr(0, TCSANOW, &old) < 0 ) perror("tcsetattr ICANON");
       if( read(0, &ch,1) < 0 ) perror("read()");
       old.c_lflag |= ICANON;
       old.c_lflag |= ECHO;
       if(tcsetattr(0, TCSADRAIN, &old) < 0) perror("tcsetattr ~ICANON");
       return ch;
   }
 */
 import "C"

import (
  "strings"
)


// GLOBAL CONSTANTS ----------------------------------------------------
const brain int = 0                                                    // index for brain
const shotgun int = 2                                                  // index for shotgun
const runner int = 1                                                   // index for runner
const green int = 0                                                    // index for green
const yellow int = 1                                                   // index for yellow
const red int = 2                                                      // index for red
const first int = 0                                                    // (position) first die in left hand
const second int = 1                                                   // (position) second die in left hand
const third int = 2                                                    // (position) third die in left hand
const totalNumberOfDice = 13                                           // total # of die
const totalNumberOfGreenDice = 6
const totalNumberOfYellowDice = 4
const totalNumberOfRedDice = 3

// GLOBAL VARS ---------------------------------------------------------
var rolld6 int
var x int                                                              // misc var
var y int                                                              // misc var
var dieFace [][]int                                                    // two dimensional array listing all possible sides for each colored dice
var icon [][]string
var myScore []int
var myCup []int                                                        // dice in cup
var myLeftHand []int                                                   // dice in hand (current roll)
var myRightHand []int                                                  // hand that temporarily holds dice that are not put out of play
var outOfPlay string                                                   // dice now out of play because of a brain/shotgun roll
var outOfPlayCounter int                                               // i need to count length so i can pad spaces. ie ANSI-colors mess up str-lengths
var gameMessage string
var gameOutcome bool
var roundIdx int
var gameState bool
var diePercentages [][]float32                                         // mathematical percentages for outcomes of (green,yellow,red) / (brain,runner,shotgun)
var handPercentages []int                                              // percentages of possible outcome for dice in left hand
var cupPercentages []int                                               // percentages of future possibilities for dice in cup (This is technically cheating, but i like it)
var spaces []rune=[]rune("             ")                              // some spaces because i'm too dumb to know how to do this a better way.
                                                                       // PLUS.. i like how 'rune' sounds.   RUNE..  RUNE..
// FUNCTIONS ===========================================================
func visualizeDice(dice []int) (visualOutput string){
  // NOTE: space-padding works here because input is integer-slice who's length can properly calculated.
  //       Trying to calculate length of an ANSI-COLOR string is impossible.
  var i, v int

  for i, v = range dice {
    visualOutput = visualOutput + icon[v][6]
  }
  visualOutput = visualOutput + string(spaces[0:(12-i)])               // pad spaces to ensure length of 13.

  return
}

func getCupPercentages(c []int, percentageType int) (cupPercs int){
  var v int
  var i int
  var pct float32                                                      // calculated percentage

  pct=0.0
  for i, v = range c {
    pct=pct+diePercentages[v][percentageType]
  }
  cupPercs = int(((pct/float32(i))+0.005)*100)
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

func prepDieFaces() [][]int {                                          // populate two dimensional array w/each die color and face
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

func prepIcons() [][]string {                                          // populate two dimensional array w/each die color and face
  var greenIcons, yellowIcons, redIcons []string
  var ds [][]string

  green  := color.New(color.FgWhite, color.BgGreen).SprintFunc()
  yellow := color.New(color.FgWhite, color.BgYellow).SprintFunc()
  red    := color.New(color.FgWhite, color.BgRed).SprintFunc()
  var greenBrain    string = green("B")
  var greenRunner   string = green("R")
  var greenShotgun  string = green("S")
  var yellowBrain   string = yellow("B")
  var yellowRunner  string = yellow("R")
  var yellowShotgun string = yellow("S")
  var redBrain      string = red("B")
  var redRunner     string = red("R")
  var redShotgun    string = red("S")
  var greenDie      string = green("g")
  var yellowDie     string = yellow("y")
  var redDie        string = red("r")

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

func continueOn() (theResult bool){
  char := C.getch()
  answer:=fmt.Sprintf("%c", char)
  answer=strings.ToLower(answer)
  theResult=false
  if answer=="y" {
    theResult=true
  }
  return
}

func howMuchBuckshot() {
  if myScore[shotgun] > 2 {
    gameMessage=color.RedString("You have been Destroyed!") + "   (you almost had " + strconv.Itoa(myScore[brain]) + " braaains.)"
    myScore[brain]=0                                                   // no BRAINS for you! You got blasted!
    gameOutcome=false
    gameState=false
  } //eoif 3-shotguns
}

func rollResults() {
  var rolld6 int
  var v int                                                            // the type of die being utilized (GREEN, YELLOW, RED)
  var i int                                                            // index var:  current die being utilized
  var rolledVisual string                                              // visual die. ANSI colored with rolled value showing.
  var handVisual string                                                // visual die. ANSI colored.. but no face value. ie before roll.
  var rolledDieOnTable [3]int                                          // number of die to replenish after roll (ie how many taken out of play)

  roundIdx+=1
  fmt.Print(color.BlueString("┃ "))
  handPercentages = nil                                                // populate percentages of dice in your left hand
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[first]][brain] + diePercentages[myLeftHand[second]][brain] + diePercentages[myLeftHand[third]][brain])/3)+0.005)*100))
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[first]][runner] + diePercentages[myLeftHand[second]][runner] + diePercentages[myLeftHand[third]][runner])/3)+0.005)*100))
  handPercentages = append(handPercentages, int((((diePercentages[myLeftHand[first]][shotgun] + diePercentages[myLeftHand[second]][shotgun] + diePercentages[myLeftHand[third]][shotgun])/3)+0.005)*100))
  cupPercentages = nil                                                 // populate percentages of diec in cup. Note 'Runner' percentages are always 33%.
  cupPercentages = append(cupPercentages, getCupPercentages(myCup, brain))
  cupPercentages = append(cupPercentages, getCupPercentages(myCup, runner))
  cupPercentages = append(cupPercentages, getCupPercentages(myCup, shotgun))
  for i, v = range myLeftHand {
    rolld6=rand.Intn(6)                                                // roll die (RANGOM NUMBER)
    // NOTE:  (I think) the SWITCH order needs to be SHOTGUN, BRAIN, RUNNER to ensure 3-SHOTGUNS will stop before BRAINS are added to score.
    switch dieFace[v][rolld6] {                                        // was the roll a BRAIN, RUNNER, or SHOTGUN
      // SHOTGUN -------------------------------------------------------
      case shotgun: {
        myScore[shotgun]+=1
        outOfPlay=outOfPlay+icon[v][rolld6]
        outOfPlayCounter+=1
        rolledDieOnTable[i]=1
        howMuchBuckshot()
        //  FIX THIS.. when cup is empty.. but you have 3 dice in hand.. you should get one more roll
        if len(myCup)==0 {gameState=false}
      } //eocase shotgun
      // BRAIN ---------------------------------------------------------
      case brain:   {
        if gameState {
          myScore[brain]+=1
          outOfPlay=outOfPlay+icon[v][rolld6]
          outOfPlayCounter+=1
          rolledDieOnTable[i]=1
          if len(myCup)==0 {gameState=false}
        } //eoifGameState
      } //eocase brain
      // RUNNER --------------------------------------------------------
      case runner: {
      } //eocase runner
    } //eoswitch dieFace
    rolledVisual=rolledVisual+icon[v][rolld6]
    handVisual=handVisual+icon[v][6]
  } //eofor theRoll

  fmt.Printf(" %02d   ┃ %-3s b:%02d%% s:%02d%% ┃ b:%02d%% s:%02d%%  %-13s┃  %-3s   ┃ %-13s",
    roundIdx,
    handVisual, handPercentages[brain], handPercentages[shotgun],
    cupPercentages[brain], cupPercentages[shotgun],
    visualizeDice(myCup),
    rolledVisual,
    (outOfPlay + string(spaces[0:(13-outOfPlayCounter)])) )

  fmt.Print(color.BlueString(" ┃\n"))
  if gameState {                                                       // if we havent already reached 3 shotguns..
    gameState=continueOn()                                             //   query user for next step
  }

  //now.. move (copy) all runners from leftHand to rightHand
  //AND   forget about brains and shotguns in left hand (they will go out of play and have been already scored)
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
          gameState=false
        } //eoif qtyDiceLeftInCup
      } //eocase1 rolledDieOnTable
    } //eoswitch rolledDieOnTable
  } //eofor rolledDieOnTable
  myLeftHand=myRightHand                                               // now take what is in right hand and put it back in left hand
}

// MAIN ================================================================
func main() {
  // INIT --------------------------------------------------------------
  gameMessage=color.GreenString("You're a really AWESOME Zombie!")

  gameState=true
  rand.Seed(time.Now().UnixNano())
  myScore=[]int{0, 0, 0}                                               // reset 3 scores (brains, runners, shotguns) to zeroes
  dieFace=prepDieFaces()
  icon=prepIcons()
  diePercentages = append(diePercentages, []float32{0.500,0.333,0.167})  //GREEN:  brain, runner, shotgun - percentages
  diePercentages = append(diePercentages, []float32{0.333,0.333,0.333})  //YELLOW: brain, runner, shotgun - percentages
  diePercentages = append(diePercentages, []float32{0.167,0.333,0.500})  //RED:    brain, runner, shotgun - percentages

  myCup = randomizeDiceInCup(totalNumberOfDice)                        // prepopulate the random dice order ie. the order that dice will be pulled from the cup

  myLeftHand=append(myLeftHand, myCup[len(myCup)-1])                   // its the FIRST roll.. get three dice, and put them in your left hand
  myCup=myCup[:len(myCup)-1]
  myLeftHand=append(myLeftHand, myCup[len(myCup)-1])
  myCup=myCup[:len(myCup)-1]
  myLeftHand=append(myLeftHand, myCup[len(myCup)-1])
  myCup=myCup[:len(myCup)-1]

  // Title -------------------------------------------------------------
  color.Blue("┏━━━━━━━━━━━━━━━━┓")
  color.Blue("┃ Zombie Dinner  ┃")
  color.Blue("┣━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
  color.Blue("┃ 'y'- to continue to roll. any other key ends round.                          ┃")
  color.Blue("┣━━━━━━━┳━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━┫")
  color.Blue("┃ round ┃     in hand     ┃          in cup           ┃ rolled ┃  out of play  ┃")
  for {
    rollResults()
    if !gameState {
      break
    }
  }
  color.Blue("┣━━━━━━━┻━━━━━━━━┳━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━━┫")
  color.Blue("┃  Stats         ┃ %-67s  ┃\n", gameMessage)
  color.Blue("┣━━━━━━━━━━━━━━━━┫                                                             ┃")
  color.Blue("┃ Rolls:    %02d   ┃                                                             ┃", roundIdx)
  color.Blue("┃ Braaains: %02d   ┃                                                             ┃", myScore[brain])
  color.Blue("┃ Shotguns: %02d   ┃                                                     madRobot┃", myScore[shotgun])
  color.Blue("┗━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛")

  if gameOutcome {
    os.Exit(0)
  } else {
    os.Exit(1)
  }

}
// END  ================================================================
