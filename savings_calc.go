/*
savings-calc
A quick and dirty program to calculate interest on savings
TODO: Output as JSON or YAML
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Print help and exit
func help(){
	fmt.Print("usage: savings-calc [-m months] [-y years] [-b balance]\n",
			  "                    [-s save]   [-a AER]   [-i]    [-h]\n",
			  "\n",
			  "  -m --months      : number of months\n",
			  "  -y --years       : number of years\n",
			  "  -b --balance     : starting balance\n",
			  "  -s --save        : amount saved per month\n",
			  "  -a --aer         : interest rate (AER)\n",
			  "  -i --interactive : interactive mode\n",
			  "  -h --help        : help\n")
	os.Exit(0)
}

// prompt for a value to be used as option
func input_prompt(s string, input *bufio.Scanner) string {
	fmt.Print(s)
	input.Scan()
	value := strings.TrimSpace(input.Text())
	if value == "" { // allow blanks
		value = "0"
	}
	return value
}

// validate that an int or float conversion was successful
func validate_arg(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}

// prompt for an float using a string and pointer to a scanner
func get_float(s string, input *bufio.Scanner) float64 {
	f, err := strconv.ParseFloat(input_prompt(s, input), 64)
	validate_arg(err)
	return f
}

// prompt for an int using a string and pointer to a scanner
func get_int(s string, input *bufio.Scanner) int {
	i, err := strconv.Atoi(input_prompt(s, input))
	validate_arg(err)
	return i
}

// get values interactively and return them ready for the calculator
func interactive_mode() (float64, float64, float64, int) {
	input	:= bufio.NewScanner(os.Stdin)
	balance := get_float("What is your current balance? ",		input)
	save	:= get_float("How much will you save per month? ",	input)
	aer		:= get_float("What is your interest rate (AER)? ",	input)
	years	:= get_int("How many years will you save for? ",	input)
	months	:= get_int("How many months will you save for? ",	input)
	return balance, save, aer, months+years*12
}

// parse the command line args and return vlaues ready for the calculator
func parse_args(argv []string) (float64, float64, float64, int) {
	months, years, balance, save, aer := 0, 0, 0.0, 0.0, 0.0
	var err error
	for i := 0; i < len(argv); i+=2 {
		switch argv[i] {
		case "-m", "--months":
			months, err = strconv.Atoi(argv[i+1])
		case "-y", "--years":
			years, err = strconv.Atoi(argv[i+1])
		case "-b", "--balance":
			balance, err = strconv.ParseFloat(argv[i+1], 64)
		case "-s", "--save":
			save, err = strconv.ParseFloat(argv[i+1], 64)
		case "-a", "--aer":
			aer, err = strconv.ParseFloat(argv[i+1], 64)
		default:
			fmt.Println("Error: invalid command line argument.")
			help()
		}
		validate_arg(err)
	}
	return balance, save, aer, months+years*12
}

// return left hand side percent of right hand side value
func percent(p, x float64) float64 {
	return x / 100.0 * p
}

// Take balance, monthly payment and AER interest rate and length of time in
// months. Returns total interest and balance after the savings period.
func savings_calc(balance, save, aer float64, months int) (float64, float64) {
	interest := 0.0
	for i := 0; i < months; i++ {
		balance += save
		interest += percent(aer, balance) / 12.0
	}
	return balance, interest
}

// print a summary of savings
func print_savings_summary(b, i, t float64) {
	fmt.Printf("%s\n%s\t£%.2f\n%s\t\t£%.2f\n%s\t\t£%.2f\n",
			   "--- Summary ---",
			   "Cumulative payments:",	b, 
			   "Total interest:",		i, 
			   "Savings:",				t)
}

// choose mode and run the calculator
func main() {
	var saved, earned float64
	if len(os.Args) == 2 && (os.Args[1] == "-i" || 
							 os.Args[1] == "--interactive") {
		saved, earned = savings_calc(interactive_mode())
	} else if len(os.Args) <= 2 {
		help()
	} else {
		saved, earned = savings_calc(parse_args(os.Args[1:]))
	}
	print_savings_summary(saved, earned, saved+earned)
}
