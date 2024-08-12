/*
 */
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Options struct {
	format    string
	time      string
	magnitude string
}

const (
	ENDPOINT = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary"
)

func main() {
	formatFlag := flag.String("format", "human", "string - {human, json, csv}")
	timeFlag := flag.String("time", "hour", "string - {hour, day, week, month}")
	magFlag := flag.String("mag", "major", "string - {all, 1.0, 2.5, 4.5, major}")
	flag.Parse()

	opts := &Options{*formatFlag, *timeFlag, *magFlag}
	fileName := extractFileName(opts)

	response, err := http.Get(ENDPOINT + fileName)
	if err != nil {
		log.Fatal(err)
	}

	bContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// flag switches from command line args to
	switch opts.format {
	case "human":
		features := extractFeatures(bContent)
		stdoutFeatures(features)
	case "csv":
		fmt.Println(string(bContent))
	case "json":
		fmt.Println(string(bContent))
	default:
		fmt.Println("todo")
	}
}

// extractFileName parses the command line Options object and returns a string
// filling in the designated file name to send to the server.
func extractFileName(cmdOpt *Options) string {
	var magRange string
	switch cmdOpt.magnitude {
	case "major":
		magRange = "significant"
	case "4.5":
		magRange = "4.5"
	case "2.5":
		magRange = "2.5"
	case "1.0":
		magRange = "1.0"
	case "all":
		magRange = "all"
	default:
		fmt.Println("todo") // is this needed?
	}

	var timeRange string
	switch cmdOpt.time {
	case "hour":
		timeRange = "hour"
	case "day":
		timeRange = "day"
	case "week":
		timeRange = "week"
	case "month":
		timeRange = "month"
	default:
		fmt.Println("todo") // is this needed?
	}

	var fileSuffix string
	switch cmdOpt.format {
	case "human":
		fileSuffix = "geojson"
	case "csv":
		fileSuffix = "csv"
	case "json":
		fileSuffix = "geojson"
	}

	endpointFile := fmt.Sprintf("/%s_%s.%s", magRange, timeRange, fileSuffix)
	return endpointFile
}

func processTime(timeStr string) time.Time {
	tVal, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		log.Fatal(err)
	}
	return tVal
}

func stringifyDateTime(tVal time.Time) string {
	year, month, day := tVal.Date()
	hour, min, sec := tVal.Clock()
	result := fmt.Sprintf("%d %v %d - %d:%d:%d", year, month, day, hour, min, sec)
	return result
}
