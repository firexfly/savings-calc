/*
savings-calc
A quick and dirty program to calculate interest on savings
*/
package main

import (
	"fmt"
	"os"
	"savings_calc/sc"
)

// remove an element from a slice
func remove(s []string, i int) []string {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

// check for -j (JSON) flag and return args array
func JSONFlagCheck(args []string) ([]string, bool) {
	for i, arg := range args {
		if arg == "-j" || arg == "--json" {
			args = remove(args, i)
			return args[1:], true
		}
	}
	return args[1:], false
}

// output in JSON
func printJSONSummary(s sc.InterestSummary) {
	fmt.Println(sc.JSONify(s))
}

// print a summary of savings
func printSavingsSummary(s sc.InterestSummary) {
	fmt.Printf("%s\n%s\t%s\n%s\t\t%s\n%s\t\t%s\n",
			   "--- Summary ---",
			   "Cumulative payments:",	s.Saved,
			   "Total interest:",		s.Interest,
			   "Savings:",				s.Total)
}

// choose mode and run the calculator
func main() {
	var on bool
	var scArgs []string
	var saved, earned sc.Cash
	outputFunc := printSavingsSummary
	if scArgs, on = JSONFlagCheck(os.Args); on {
		outputFunc = printJSONSummary
	}
	if len(scArgs) == 1 && (scArgs[0] == "-i" ||
							scArgs[0] == "--interactive") {
		saved, earned = sc.SavingsCalc(sc.InteractiveMode())
	} else if len(scArgs) <= 1 {
		sc.Help()
	} else {
		saved, earned = sc.SavingsCalc(sc.ParseArgs(scArgs))
	}
	outputFunc(sc.InterestSummary{saved, earned, saved+earned})
}
