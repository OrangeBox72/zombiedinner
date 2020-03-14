// name:    zombieDice
// author:  johnny
// version: 2020/03/12. WIP
// references:
//    global constants and slices.  https://qvault.io/2019/10/21/how-to-global-constant-maps-and-slices-in-go/
// DELTHIS//    ansi colors.  https://gist.github.com/ik5/d8ecde700972d4378d87
// notes: actual gameplay still not done.
//   * 2020-03-12 - started using ANSI colors.. but they mess up printf formatting, so i have to pad spaces on the colored variables.
//   * old notes
//     currently, only completed logic for rolls, die, color, face, and scoring
//     !!!!! next step.. convert to prepopulated roll var to a slice
//       and remove an element when it becomes brain or shotgun
//       also.. adjust the min/max of the slice dynamically for the for/next to work properly
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

//const greenDie string =  "\033[38;5;0m\033[48;5;2mg\033[39;49m"
//const yellowDie string = "\033[38;5;0m\033[48;5;3my\033[39;49m"
//const redDie string =    "\033[38;5;0m\033[48;5;1mr\033[39;49m"

//const brainResult string   = "\033[38;5;0m\033[48;5;2mB\033[39;49m"
//const runnerResult string  = "\033[38;5;0m\033[48;5;3mR\033[39;49m"
//const shotgunResult string = "\033[38;5;0m\033[48;5;1mS\033[39;49m"

const greenBrain string     = "\033[38;5;15m\033[48;5;2mB\033[39;49m"
const greenRunner string    = "\033[38;5;15m\033[48;5;2mR\033[39;49m"
const greenShotgun string   = "\033[38;5;0m\033[48;5;2mS\033[39;49m"
const yellowBrain string    = "\033[38;5;15m\033[48;5;3mB\033[39;49m"
const yellowRunner string   = "\033[38;5;15m\033[48;5;3mR\033[39;49m"
const yellowShotgun string  = "\033[38;5;0m\033[48;5;3mS\033[39;49m"
const redBrain string       = "\033[38;5;15m\033[48;5;1mB\033[39;49m"
const redRunner string      = "\033[38;5;15m\033[48;5;1mR\033[39;49m"
const redShotgun string     = "\033[38;5;0m\033[48;5;1mS\033[39;49m"


var rolld6 int

// VARS ----------------------------------------------------------------
var x int                                                              // misc var
var y int                                                              // misc var
var color int      // color picked from cup during the setup
//var colorName []string                                                 // die color description
//var faceName []string                                                  // die face description
var dieSides [][]int                                                   // two dimensional array listing all possible sides for each colored dice
var iconSides [][]string
var myScore []int
var myCup []int                                                        // dice in cup
var myHand []int                                                       // dice in hand (current roll)
var outOfPlay []int                                                    // dice now out of play because of a shotgun roll
var gameOutcome string
var brainTally int
var runnerTally int
var shotgunTally int
var rollCount int

//var (
//  Black   = Color("\033[38;5;0m%s\033[39;49m")
//  Red     = Color("\033[38;5;1m%s\033[39;49m")
//  Green   = Color("\033[38;5;2m%s\033[39;49m")
//  Yellow  = Color("\033[38;5;3m%s\033[39;49m")
//  Blue    = Color("\033[38;5;4m%s\033[39;49m")
//  Magenta = Color("\033[38;5;5m%s\033[39;49m")
//  Teal    = Color("\033[38/5;6m%s\033[39;49m")
//  White   = Color("\033[38;5;7m%s\033[39;49m")
//)
//var BorderColor = Blue

var (
  ansiReset  = "\033[39;49m"
//  ansiGreen  = "\033[38;5;2m"
//  ansiYellow = "\033[38;5;3m"
//  ansiRed    = "\033[38;5;1m"
  ansiBlue   = "\033[38;5;4m"
)
//var (
//  ansiGreenBG   = "\033[38;5;0m\033[48;5;2m"
//  ansiYellowBG  = "\033[38;5;0m\033[48;5;3m"
//  ansiRedBG     = "\033[38;5;0m\033[48;5;1m"
//)


// FUNCTIONS ===========================================================
func printSlice(z int, s []int) {
  fmt.Printf("z=%d :  len=%d    %#v      \n", z, len(s), s)
}
func printSlice2(z int, s [][]int) {
  fmt.Printf("z=%d :  len=%d    %#v      \n", z, len(s), s)
}

//func Color(colorString string) func(...interface{}) string {             // printing ansi colors (possibly not used)
//  sprint := func(args ...interface{}) string {
//    return fmt.Sprintf(colorString,
//      fmt.Sprint(args...))
//  }
//  return sprint
//}

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

//func getColorNames() []string {
//  return []string{ansiGreenBG+"Green "+ansiReset, ansiYellowBG+"Yellow"+ansiReset, ansiRedBG+" Red  "+ansiReset}
//}

//func getFaceNames() []string {
//  return []string{ansiGreen+"Brain  "+ansiReset, ansiYellow+"Runner "+ansiReset, ansiRed+"SHOTGUN"+ansiReset}
//}

func prepDieSides() [][]int {                                          // populate two dimensional array w/each die color and face
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

func rollResults(theRoll []int) {
  var rolld6 int
  var v int
//  var i int
  var resultVisual string

  rollCount+=1
  fmt.Print(ansiBlue, "┃ ",ansiReset)
  for _, v = range theRoll {
    rolld6=rand.Intn(6)
    switch dieSides[v][rolld6] {
      case brain:   brainTally+=1
      case runner:  runnerTally+=1
      case shotgun: shotgunTally+=1
    } //eoswitch
    resultVisual=resultVisual+iconSides[v][rolld6]
  } //eofor
  fmt.Printf("roll %02d: %-3s   tally: brains:%02d    runners:%02d    shotguns:%02d",rollCount, resultVisual, brainTally, runnerTally, shotgunTally)
  fmt.Print(ansiBlue, " ┃", ansiReset, "\n")

//  tally=append(tally,brainTally)
//  tally=append(tally,shotgunTally)


//  fmt.Println(i) // dummy line to keep compiler from squawking that 'i' isnt used anywhere
//  fmt.Printf("roll: %2d.\n",rolld6)


//  switch score {
//    case brain, runner, shotgun: {
//      theScores[score]++
//    }
//  }



// - %-6s %-7s   tally (brains:%2d  shotguns:%2d  runners:%2d    die in cup: %2d , outOfPlay: %2d)\n", y+1, colorName[myCup[y]], faceName[dieSides[myCup[y]][rolld6]], myScore[brain], myScore[shotgun], myScore[runner], numDieInCup, dieOutOfPlay )

}

// MAIN ===============================================================
func main() {

// var numDieInCup int                                                   // how many dice in cup
// var dieOutOfPlay int                                                  // how many dice (shotguns/brains) that are now out of play

  // Title ------------------------------------------------------------
  fmt.Print(ansiBlue)
  fmt.Print("┏━━━━━━━━━━━━━━━━┓\n")
  fmt.Print("┃  Zombie Dice   ┃\n")
  fmt.Print("┣━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓")
  fmt.Print(ansiReset, "\n")

  // INIT -------------------------------------------------------------
  rand.Seed(time.Now().UnixNano())
  myScore=resetScores()                                                // reset 3 scores (brains, runners, shotguns) to zeroes
//  colorName=getColorNames()
//  faceName=getFaceNames()
  dieSides=prepDieSides()
  iconSides=prepIcons()
  // prepopulate the random dice order ie. the order that dice will be pulled from the cup
  myCup = randomizeDiceInCup(totalNumberOfDice)


//numDieInCup = len(myCup)
//dieOutOfPlay = totalNumberOfDice - numDieInCup
printSlice(22,myCup)
printSlice(23,myHand)
var a,b,c int
a=myCup[len(myCup)-1]
myCup=myCup[:len(myCup)-1]
b=myCup[len(myCup)-1]
myCup=myCup[:len(myCup)-1]
c=myCup[len(myCup)-1]
myCup=myCup[:len(myCup)-1]
myHand=append(myHand,  a,b,c  )   // remove first die from cup and into hand
printSlice(24,myCup)
printSlice(25,myHand)

rollResults(myHand)




//os.Exit(22)

  // MAIN -------------------------------------------------------------
// y = 0
// for y < numDieInCup {
//  rolld6 = rand.Intn(6)      // roll me a die
//  myScore[dieSides[myCup[y]][rolld6]]+=1  // change the current-round score for this particular die roll
//  fmt.Printf(BorderColor("┃ "))
//  fmt.Printf("die %2d - %-6s %-7s   tally (brains:%2d  shotguns:%2d  runners:%2d    die in cup: %2d , outOfPlay: %2d)\n", y+1, colorName[myCup[y]], faceName[dieSides[myCup[y]][rolld6]], myScore[brain], myScore[shotgun], myScore[runner], numDieInCup, dieOutOfPlay )
//  y+=1
//  if myScore[shotgun] >= 3 {
//    gameOutcome=Red("You have been DESTROYED!")
//   y=99
//  }
// }
// if y==totalNumberOfDice {
//   gameOutcome=Green("You're a really AWESOME Zombie!")
// }

  gameOutcome="johnny was here"
  fmt.Print(ansiBlue)
  fmt.Print( "┣━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫\n")
  fmt.Printf("┃  Stats         ┃ %-32s            ┃\n",gameOutcome)
  fmt.Print( "┣━━━━━━━━━━━━━━━━┫                                             ┃\n")
  fmt.Printf("┃ Braaains: %02d   ┃                                             ┃\n", myScore[brain])
  fmt.Printf("┃ Runners:  %02d   ┃                                             ┃\n", myScore[runner])
  fmt.Printf("┃ Shotguns: %02d   ┃                                     madRobot┃\n", myScore[shotgun])
  fmt.Print( "┗━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
  fmt.Print(ansiReset, "\n")

  os.Exit(0)
}
// END  =========================================================
