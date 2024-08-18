package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Options struct {
	format    string
	time      string
	magnitude string
}

var ErrFlagMagOption = errors.New("magnitude option invalid")
var ErrFlagRangeOption = errors.New("time interval option invalid")
var ErrFlagFormatOption = errors.New("format option invalid")

const (
	ENDPOINT = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary"
)

func main() {
	formatFlag := flag.String("format", "human", "string - {human, json, csv}")
	timeFlag := flag.String("time", "month", "string - {hour, day, week, month}")
	magFlag := flag.String("mag", "major", "string - {all, 1.0, 2.5, 4.5, major}")
	flag.Parse()

	opts := &Options{*formatFlag, *timeFlag, *magFlag}
	fileName, err := extractFileName(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(ENDPOINT + fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	bContent, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// Standard output format
	switch opts.format {
	case "human":
		if features, err := extractFeatures(bContent); err != nil {
			fmt.Fprint(os.Stderr, "%s\n", err)
			os.Exit(1)
		} else {
		  stdoutFeatures(features)
		}
	case "csv":
		fmt.Println(string(bContent))
	case "json":
		fmt.Println(string(bContent))
	}
}

// extractFileName parses the command line Options object and returns a string
// filling in the designated file name to send to the server.
func extractFileName(cmdOpt *Options) (string, error) {
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
		return "", ErrFlagMagOption
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
		return "", ErrFlagRangeOption
	}

	var fileSuffix string
	switch cmdOpt.format {
	case "human":
		fileSuffix = "geojson"
	case "csv":
		fileSuffix = "csv"
	case "json":
		fileSuffix = "geojson"
	default:
		return "", ErrFlagFormatOption
	}

	endpointFile := fmt.Sprintf("/%s_%s.%s", magRange, timeRange, fileSuffix)
	return endpointFile, nil
}

func stringifyDateTime(tVal time.Time) string {
	year, month, day := tVal.Date()
	hour, min, sec := tVal.Clock()
	result := fmt.Sprintf("%d %v %d - %d:%d:%d", year, month, day, hour, min, sec)
	return result
}
