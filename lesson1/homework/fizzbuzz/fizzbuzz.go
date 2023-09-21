package fizzbuzz

import "strconv"

func FizzBuzz(k int) string {
	if k%15 == 0 {
		return "FizzBuzz"
	} else if k%5 == 0 {
		return "Buzz"
	} else if k%3 == 0 {
		return "Fizz"
	} else {
		return strconv.Itoa(k)
	}
}
