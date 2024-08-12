package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type USGSResponse struct {
	Type     string    `json:"type"`
	Meta     Metadata  `json:"metadata"`
	Bbox     []float64 `json:"bbox"`
	Features `json:"features"`
}

type Metadata struct {
	Generated int64  `json:"generated"`
	Url       string `json:"url"`
	Title     string `json:"title"`
	Api       string `json:"api"`
	Count     int    `json:"count"`
	Status    int    `json:"status"`
}

type Feature struct {
	Type  string     `json:"type"`
	Props Properties `json:"properties"`
	Geo   Geometry   `json:"geometry"`
	Id    string     `json:"id"`
}

type Features []Feature

type Properties struct {
	Mag     float64 `json:"mag"`
	Place   string  `json:"place"`
	Time    int64   `json:"time"`
	Updated int64   `json:"updated"`
	Tz      int     `json:"tz"`
	Url     string  `json:"url"`
	Detail  string  `json:"detail"`
	Felt    int     `json:"felt"`
	Cdi     float64 `json:"cdi"`
	Mmi     float64 `json:"mmi"`
	Alert   string  `json:"alert"`
	Status  string  `json:"status"`
	Tsunami int     `json:"tsunami"`
	Sig     int     `json:"sig"`
	Net     string  `json:"net"`
	Code    string  `json:"code"`
	Ids     string  `json:"ids"`
	Sources string  `json:"sources"`
	Types   string  `json:"types"`
	Nst     int     `json:"nst"`
	Dmin    float64 `json:"dmin"`
	Rms     float64 `json:"rms"`
	Gap     float64 `json:"gap"`
	MagType string  `json:"magType"`
	Type    string  `json:"type"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func extractFeatures(res []byte) Features {
	var usgsRes USGSResponse
	err := json.Unmarshal(res, &usgsRes)
	if err != nil {
		log.Fatal(err)
	}

	f := make(Features, 0, len(usgsRes.Features))
	for _, v := range usgsRes.Features {
		f = append(f, v)
	}
	return f
}

func stdoutFeatures(features Features) {
	fmt.Printf("%s %18s %-40s %-4s %4s\n", "Date-Time", "Mag", "Place", "Lat", "Long")
	for _, f := range features {
		dateTimeVal := time.Unix(f.Props.Time, 0)
		dateTimeStr := stringifyDateTime(dateTimeVal)
		fmt.Fprintf(os.Stdout, "%s %3.2f %-40s %4.2f %4.2f\n",
			dateTimeStr, f.Props.Mag, f.Props.Place, f.Geo.Coordinates[1], f.Geo.Coordinates[0])
	}
}
