package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	var url string
	var secondsDuration int
	flag.StringVar(&url, "u", "default-url", "Specify url.")
	flag.IntVar(&secondsDuration, "s", 15, "Specify duration to capture a reading in seconds.")
	flag.Parse()
	if len(url) <= 0 {
		log.Fatalf("No url provided.")
	}
	perpetualReadings := map[string]float64{}
	dailyReadings := map[string]float64{}

	dailyAverages := map[string]float64{}

	programLaunchTime := time.Now().String()
	dayStart := time.Now()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n--- QUITTING ---\n Perpetual average since %s: %v\nAverage of current day starting at %s: %v\nAll daily averages: %v\n\n", programLaunchTime, calculateAverage(perpetualReadings), dayStart.String(), calculateAverage(dailyReadings), dailyAverages)
		os.Exit(1)
	}()

	for {
		time.Sleep(time.Duration(secondsDuration) * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making get request to: %s with error: %s\n", url, err.Error())
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %s\n", err.Error())
			continue
		}
		sb := string(body)
		re := regexp.MustCompile(`PS\[(\d* \d* \d* \d* \d* \d*)\]`)
		ps := re.FindString(sb)
		lastSpaceIndex := strings.LastIndex(ps, " ")
		if lastSpaceIndex < 0 {
			fmt.Printf("Proper regex not found, not calculating. String: %s\n", ps)
			continue
		}
		ps = ps[0:lastSpaceIndex]
		lastSpaceIndex = strings.LastIndex(ps, " ")
		if lastSpaceIndex < 0 {
			fmt.Printf("Proper regex not found in second space chopping, not calculating. String: %s\n", ps)
			continue
		}
		ps = ps[lastSpaceIndex+1:]
		reading, err := strconv.ParseFloat(ps, 32)
		if err != nil {
			fmt.Printf("Error converting int to string. Error: %s Original string: %s", err.Error(), ps)
			continue
		}
		nowTime := time.Now()
		perpetualReadings[nowTime.String()] = reading
		if nowTime.Sub(dayStart).Seconds() >= 86400 {
			yesterdaysAverage := calculateAverage(dailyReadings)
			dailyAverages[dayStart.String()] = yesterdaysAverage
			fmt.Printf("----- Starting a new day. Yesterdays average was: %v -----\n\n", yesterdaysAverage)
			dailyReadings = map[string]float64{}
			dayStart = time.Now()
		}
		dailyReadings[nowTime.String()] = reading
		fmt.Printf("Most recent reading (%s): %v\nTodays average (starting at %s): %v\nAverage since program start (%s): %v\n\n", nowTime.String(), reading, dayStart.String(), calculateAverage(dailyReadings), programLaunchTime, calculateAverage(perpetualReadings))

	}
}

func calculateAverage(readings map[string]float64) float64 {
	sumOfReadings := 0.0
	for _, reading := range readings {
		sumOfReadings += reading
	}
	return sumOfReadings / float64(len(readings))
}
