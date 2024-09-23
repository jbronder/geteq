package logic

import (
	"encoding/json"
	"fmt"
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

// ExtractFeatures unmarshals a list of earthquake events into Features to
// prepare for formatting.
func ExtractFeatures(res []byte) (Features, error) {
	var usgsRes USGSResponse
	err := json.Unmarshal(res, &usgsRes)
	if err != nil {
		return nil, err
	}

	f := make(Features, 0, len(usgsRes.Features))
	for _, v := range usgsRes.Features {
		f = append(f, v)
	}
	return f, nil
}

// ExtractSingleFeature unmarshals one event into one Feature to prepare for
// formatting.
func ExtractSingleFeature(res []byte) (*Feature, error) {
	f := new(Feature)
	err := json.Unmarshal(res, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// StdoutFeatures outputs to standard output a list of Features that were
// unmarshaled from a response. As a snapshot, the main fields that are listed
// for each event are: eventid, time, place, latitude, and longitude.
func StdoutFeatures(features Features) {
	if len(features) == 0 {
		fmt.Fprintf(os.Stdout, "No records matched under the given criteria.\n")
		return
	}

	fmt.Printf("%-10s %s %-4s %-42s %-5s %s\n", "EventId", "Date-Time UTC+00:00", "Mag", "Place", "Lat", "Long")
	for _, f := range features {
		dateTimeVal := time.UnixMilli(f.Props.Time).UTC()
		dateTimeStr := dateTimeVal.Format(time.DateTime)
		fmt.Fprintf(os.Stdout, "%s %s %3.2f %-41s %6.2f %7.2f\n",
			f.Id, dateTimeStr, f.Props.Mag, f.Props.Place, f.Geo.Coordinates[1], f.Geo.Coordinates[0])
	}
}

// resolveMagType provides a mapping from the short-form 2 or 3 letter magnitude
// type into a more descriptive label of magnitude type. If a 2 or 3 letter
// magnitude type is not present in the hash table, return the original magType
// variable.
func resolveMagType(magType string) string {
	typeMapping := map[string]string{
		"md":  "Duration Magnitude (Md)",
		"ml":  "Richter Scale Magnitude (Ml)",
		"ms":  "20 sec Surface Wave Magnitude (Ms)",
		"mw":  "Moment Magnitude (Mw)",
		"mww": "Moment Magnitude, W-Phase (Mww)",
		"me":  "Energy Magnitude (Me)",
		"mi":  "Integrated P-Wave Magnitude (Mi)",
		"mb":  "Short-Period Body Wave Magnitude (Mb)",
		"mlg": "Short-Period Surface Wave Magnitude (Mlg)",
	}

	if fullMagString, hasKey := typeMapping[magType]; hasKey {
		return fullMagString
	} else {
		return magType
	}
}

// StdoutSingleEvent outputs detailed information about an earthquake event.
func StdoutSingleEvent(f *Feature) {
	if f == nil {
		fmt.Fprintf(os.Stdout, "No record matched under the given criteria.\n")
		return
	}

	datetime := time.UnixMilli(f.Props.Time).UTC()
	datetimeStr := datetime.Format(time.DateTime)

	updateDateTime := time.UnixMilli(f.Props.Updated).UTC()
	updateDateTimeStr := updateDateTime.Format(time.DateTime)

	fmt.Fprintf(os.Stdout, "Single Event Details\n--------------------\n")
	fmt.Fprintf(os.Stdout, "Event Id: %s\n", f.Id)
	fmt.Fprintf(os.Stdout, "Review Status: %s\n", f.Props.Status)
	fmt.Fprintf(os.Stdout, "Time (UTC+00:00): %s\n", datetimeStr)
	fmt.Fprintf(os.Stdout, "Updated Time (UTC+00:00): %s\n", updateDateTimeStr)
	fmt.Fprintf(os.Stdout, "Time Zone Offset: %d\n", f.Props.Tz)
	fmt.Fprintf(os.Stdout, "Place: %s\n", f.Props.Place)
	fmt.Fprintf(os.Stdout, "Magnitude: %3.2f\n", f.Props.Mag)
	fmt.Fprintf(os.Stdout, "Magnitude Type: %s\n", resolveMagType(f.Props.MagType))
	fmt.Fprintf(os.Stdout, "Depth: %.2f km\n", f.Geo.Coordinates[2])
	fmt.Fprintf(os.Stdout, "Latitude: %.2f\n", f.Geo.Coordinates[1])
	fmt.Fprintf(os.Stdout, "Longitude: %.2f\n", f.Geo.Coordinates[0])
	fmt.Fprintf(os.Stdout, "Horizontal distance (in deg) from epicenter to the nearest station: %f\n", f.Props.Dmin)
	fmt.Fprintf(os.Stdout, "Largest Azimuthal Gap between stations (deg): %.2f\n", f.Props.Gap)
	fmt.Fprintf(os.Stdout, "Root-Mean-Square (RMS) Travel Time Residual (sec): %.3f\n", f.Props.Rms)
	fmt.Fprintf(os.Stdout, "Seismic Event Type: %s\n", f.Props.Type)
	fmt.Fprintf(os.Stdout, "PAGER Alert Level: %s\n", f.Props.Alert)
	fmt.Fprintf(os.Stdout, "Number of Felt Reports of DYFI: %d\n", f.Props.Felt)
	fmt.Fprintf(os.Stdout, "Intensity Level: %.2f\n", f.Props.Cdi)
	fmt.Fprintf(os.Stdout, "Modified Mercalli Intensity (MMI): %.2f\n", f.Props.Mmi)
	fmt.Fprintf(os.Stdout, "Event Significance: %d\n", f.Props.Sig)
	fmt.Fprintf(os.Stdout, "Large Event in Oceanic Region: %d\n", f.Props.Tsunami)
	fmt.Fprintf(os.Stdout, "Number of Stations used to determine location: %d\n", f.Props.Nst)
	fmt.Fprintf(os.Stdout, "Associated Event Ids: %s\n", f.Props.Ids)
	fmt.Fprintf(os.Stdout, "Network Contributors: %s\n", f.Props.Sources)
	fmt.Fprintf(os.Stdout, "Preferred Contributor Id: %s\n", f.Props.Net)
	fmt.Fprintf(os.Stdout, "Event Id Code: %s\n", f.Props.Code)
}
