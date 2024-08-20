package logic

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrFlagMagOption = errors.New("--magnitude option invalid")
var ErrFlagRangeOption = errors.New("--time interval option invalid")
var ErrFlagFormatOption = errors.New("--format option invalid")

const (
	ENDPOINT = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary"
)

// ExtractFileName parses the command line Options object and returns a string
// filling in the designated file name to send to the server.
func ExtractFileName(formatFlag, magFlag, timeFlag string) (string, error) {
	var magRange string
	switch magFlag {
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
	switch timeFlag {
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
	switch formatFlag {
	case "table":
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

func RequestContent(apiPath string) ([]byte, error) {
	response, err := http.Get(ENDPOINT + apiPath) 
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bContent, nil
}
