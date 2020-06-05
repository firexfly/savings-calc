package sc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Cash float64
func (c Cash) String() string { return fmt.Sprintf("£%.2f", c) }

type InterestProps struct {
	Balance, Save, AER Cash
	Months, Years int
}

type InterestSummary struct {
	Saved, Interest, Total Cash
}

// Print help and exit
func Help(){
	fmt.Print("usage: savings-calc [-m months] [-y years] [-b balance]\n",
			  "                    [-s save]   [-a AER] [-j] [-i] [-h]\n",
			  "\n",
			  "  -m --months      : number of months\n",
			  "  -y --years       : number of years\n",
			  "  -b --balance     : starting balance\n",
			  "  -s --save        : amount saved per month\n",
			  "  -a --aer         : interest rate (AER)\n",
			  "  -j --json        : output as JSON\n",
			  "  -i --interactive : interactive mode\n",
			  "  -h --help        : help\n")
	os.Exit(0)
}

// prompt for a value to be used as option
func inputPrompt(s string, input *bufio.Scanner) string {
	fmt.Print(s)
	input.Scan()
	value := strings.TrimSpace(input.Text())
	if value == "" { // allow blanks
		value = "0"
	}
	return value
}

// validate that an int or float conversion was successful
func validateArg(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// prompt for an float using a string and pointer to a scanner
func getFloat(s string, input *bufio.Scanner) float64 {
	f, err := strconv.ParseFloat(inputPrompt(s, input), 64)
	validateArg(err)
	return f
}

// prompt for an int using a string and pointer to a scanner
func getInt(s string, input *bufio.Scanner) int {
	i, err := strconv.Atoi(inputPrompt(s, input))
	validateArg(err)
	return i
}

// get values interactively and return them ready for the calculator
func InteractiveMode() (*InterestProps) {
	p		 := InterestProps{}
	input	 := bufio.NewScanner(os.Stdin)
	p.Balance = Cash(getFloat("What is your current balance? ",		input))
	p.Save	  = Cash(getFloat("How much will you save per month? ",	input))
	p.AER	  = Cash(getFloat("What is your interest rate (AER)? ",	input))
	p.Years	  = getInt("How many years will you save for? ",		input)
	p.Months  = getInt("How many months will you save for? ",		input)
	return &p
}

// parse the command line args and return vlaues ready for the calculator
func ParseArgs(argv []string) (*InterestProps) {
	p := InterestProps{}
	var balance, save, aer float64 // temporary vars
	var err error
	for i := 0; i < len(argv); i+=2 {
		switch argv[i] {
		case "-m", "--months":
			p.Months, err = strconv.Atoi(argv[i+1])
		case "-y", "--years":
			p.Years, err = strconv.Atoi(argv[i+1])
		case "-b", "--balance":
			balance, err = strconv.ParseFloat(argv[i+1], 64)
			p.Balance = Cash(balance)
		case "-s", "--save":
			save, err = strconv.ParseFloat(argv[i+1], 64)
			p.Save = Cash(save)
		case "-a", "--aer":
			aer, err = strconv.ParseFloat(argv[i+1], 64)
			p.AER = Cash(aer)
		default:
			fmt.Println("Error: invalid command line argument.")
			Help()
		}
		validateArg(err)
	}
	return &p
}

// return left hand side percent of right hand side value
func percent(p, x float64) float64 {
	return x / 100.0 * p
}

// return output as JSON
func JSONify(data InterestSummary) string {
	json, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}

/*
Take balance, monthly payment, AER interest rate and length of time.
Returns total interest and balance after the savings period.
*/
func SavingsCalc(p *InterestProps) (Cash, Cash) {
	interest, balance := Cash(0.0), p.Balance
	for i := 0; i < p.Months+p.Years*12; i++ {
		balance += p.Save
		interest += Cash(percent(float64(p.AER), float64(balance)) / 12.0)
	}
	return balance, interest
}