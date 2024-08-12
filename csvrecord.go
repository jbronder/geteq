package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type CSVRecord struct {
	time            string
	latitude        float64
	longitude       float64
	depth           float64
	mag             float64
	magType         string
	nst             int
	gap             float64
	dmin            float64
	rms             float64
	net             string
	id              string
	updated         string
	place           string
	typeEvent       string
	horizontalError float64
	depthError      float64
	magError        float64
	magNst          int
	status          string
	locationSource  string
	magSource       string
}

const (
	NUM_CSV_FIELDS = 22
)

/*
Assumes all records have 22 fields. Each record string[] needs to be processed
equentially
*/
func NewRecord(record []string) *CSVRecord {
	if len(record) != NUM_CSV_FIELDS {
		fmt.Fprintf(os.Stderr, "CSV Records found: %d\n", len(record))
	}

	csvR := new(CSVRecord)

	csvR.time = processString(record[0])
	csvR.latitude = processFloat(record[1])
	csvR.longitude = processFloat(record[2])
	csvR.depth = processFloat(record[3])
	csvR.mag = processFloat(record[4])
	csvR.magType = processString(record[5])
	csvR.nst = processInt(record[6])
	csvR.gap = processFloat(record[7])
	csvR.dmin = processFloat(record[8])
	csvR.rms = processFloat(record[9])
	csvR.net = processString(record[10])
	csvR.id = processString(record[11])
	csvR.updated = processString(record[12])
	csvR.place = processString(record[13])
	csvR.typeEvent = processString(record[14])
	csvR.horizontalError = processFloat(record[15])
	csvR.depthError = processFloat(record[16])
	csvR.magError = processFloat(record[17])
	csvR.magNst = processInt(record[18])
	csvR.status = processString(record[19])
	csvR.locationSource = processString(record[20])
	csvR.magSource = processString(record[21])

	return csvR
}

func processString(rec string) string {
	if len(rec) == 0 {
		return ""
	}
	return strings.TrimSpace(rec)
}

func processFloat(rec string) float64 {
	if len(rec) == 0 {
		return 0
	}
	if floatField, err := strconv.ParseFloat(rec, 64); err != nil {
		fmt.Fprint(os.Stderr, "%v\n", err)
		return -1
	} else {
		return floatField
	}
}

func processInt(rec string) int {
	if len(rec) == 0 {
		return 0
	}
	if intField, err := strconv.ParseInt(rec, 10, 0); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return -1
	} else {
		return int(intField)
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

// humanize outputs a human readable table that fits within a terminal.
// It outputs fields that a person curious about earthquake occurrences may be
// commonly interested in such as: Date, time, magnitude, where it occurred, and
// the latitude and longitudinal coordinates of the location.
func humanizeCSV(recs []*CSVRecord) {
	//Date       Time     Mag   Place                            Lat    Long
	fmt.Printf("%s %18s %-40s %-4s %4s\n", "Date-Time", "Mag", "Place", "Lat", "Long")
	for i := range recs {
		dateTimeVal := processTime(recs[i].time)
		dateTimeStr := stringifyDateTime(dateTimeVal)
		fmt.Fprintf(os.Stdout, "%s %3.2f %-40s %4.2f %4.2f\n",
			dateTimeStr, recs[i].mag, recs[i].place, recs[i].latitude, recs[i].longitude)
	}
}
