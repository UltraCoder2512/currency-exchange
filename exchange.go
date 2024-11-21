package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattevans/dinero"
)

func main() {
	start := time.Now()
	if len(os.Args) != 4 {
		log.Fatal("Correct usage: exchange <base currency ISO code> <required currency ISO code> <Amount>")
	}

	rate := GetExchangeRatesForOtherBase(os.Args[1], os.Args[2])

	amt, err := strconv.ParseFloat(strings.TrimSpace(os.Args[3]), 64)

	if err != nil {
		log.Fatal("Failed to parse amount")
	}

	fmt.Printf("%v %s is %v %s\n in %v", FormatDecimal(amt, 2), os.Args[1], FormatDecimal(rate*amt, 2), os.Args[2], time.Since(start))
}

func GetExchangeRates(TOISOCode string, data chan<- *float64) {
	client := dinero.NewClient(os.Getenv("CURRENCY_API_KEY"), "USD", time.Duration(time.Minute*5))
	rate, err := client.Rates.Get(TOISOCode)
	data <- rate
	if err != nil {
		log.Fatal("Fatal error: ", err)
	}
}

func GetExchangeRatesForOtherBase(BASEISOCode string, TOISOCode string) float64 {
	data1 := make(chan *float64)
	go GetExchangeRates(BASEISOCode, data1)
	rate1, ok := <-data1
	if !ok {
		log.Fatal("Failed to get exchange rates")
	}
	defer close(data1)

	data2 := make(chan *float64)
	go GetExchangeRates(TOISOCode, data2)
	rate2, ok := <-data2
	if !ok {
		log.Fatal("Failed to get exchange rates")
	}
	defer close(data2)

	rate := *rate2 / *rate1
	return rate
}

func FormatDecimal(val float64, prec int) float64 {
	v := math.Pow10(prec)
	k := math.Round(val * v)
	return k / v
}
