/*
 */
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

	// fmt.Println("extracting", ENDPOINT+fileName)

	response, err := http.Get(ENDPOINT + fileName)
	if err != nil {
		log.Fatal(err)
	}

	bContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(bContent))
	// using for parsing to human readable format
	// bContent, err := os.ReadFile("eq-out.csv")
	if err != nil {
		log.Fatal(err)
	}

	// flag switches from command line args to
	switch opts.format {
	case "human":
		csvRecords := createCSVRecords(bContent)
		humanize(csvRecords)
	case "csv":
		fmt.Println(string(bContent))
	case "json":
		fmt.Println("json todo")
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

	csvfile := fmt.Sprintf("/%s_%s.csv", magRange, timeRange)
	return csvfile
}

// humanize outputs a human readable table that fits within a terminal.
// It outputs fields that a person curious about earthquake occurrences may be
// commonly interested in such as: Date, time, magnitude, where it occurred, and
// the latitude and longitudinal coordinates of the location.
func humanize(recs []*CSVRecord) {
	//Date       Time     Mag   Place                            Lat    Long
	fmt.Printf("%s %18s %-40s %-4s %4s\n", "Date-Time", "Mag", "Place", "Lat", "Long")
	for i := range recs {
    dateTimeVal := processTime(recs[i].time)
    dateTimeStr := stringifyDateTime(dateTimeVal)
		fmt.Fprintf(os.Stdout, "%s %3.2f %-40s %4.2f %4.2f\n",
			dateTimeStr, recs[i].mag, recs[i].place, recs[i].latitude, recs[i].longitude)
	}
}

// createCSVRecords processes the CSV data received from the server and creates
// individual records out of each line CSV data
func createCSVRecords(resContent []byte) []*CSVRecord {
	strContent := string(resContent)
	csvReader := csv.NewReader(strings.NewReader(strContent))

	var records []*CSVRecord

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		if !strings.Contains(record[0], "time") {
			csvRec := NewRecord(record)
			records = append(records, csvRec)
		}
	}

	return records
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
