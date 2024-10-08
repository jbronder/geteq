package logic

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrFlagMagOption = errors.New("--magnitude option invalid")
var ErrFlagTimeOption = errors.New("--time interval option invalid")
var ErrFlagFormatOption = errors.New("--output format option invalid")
var ErrEventIdInvalid = errors.New("eventid invalid")

const (
	RTENDPOINT   = "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary"
	FDSNENDPOINT = "https://earthquake.usgs.gov/fdsnws/event/1"
	ALPHABET     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NONALPHABET  = "\"&*?-.+-^%(){}[]_@:|\\!~`,"
	DIGITS       = "0123456789"
)

// ExtractRTParams parses the user flag values and returns a complete URL to
// send to the server.
func ExtractRTParams(formatFlag, magFlag, timeFlag string) (string, error) {
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
		return "", ErrFlagTimeOption
	}

	var fileSuffix string
	switch formatFlag {
	case "table":
		fallthrough
	case "json":
		fileSuffix = "geojson"
	case "csv":
		fileSuffix = "csv"
	default:
		return "", ErrFlagFormatOption
	}

	partial := fmt.Sprintf("%s_%s.%s", magRange, timeRange, fileSuffix)
	fullURL, err := url.JoinPath(RTENDPOINT, partial)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

// RequestContent performs the GET request for the resource.
func RequestContent(apiPath string) ([]byte, error) {
	response, err := http.Get(apiPath)
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

// ExtractFDSNParams resolves user input flag values and pairs them with the
// endpoint method to return a complete URL for a request.
func ExtractFDSNParams(endCmd, magFlag, formatFlag, dateTimeFlag string) (string, error) {
	v := url.Values{}

	switch formatFlag {
	case "table":
		fallthrough
	case "geojson":
		fallthrough
	case "json":
		v.Set("format", "geojson")
	case "text":
		v.Set("format", "text")
	case "csv":
		v.Set("format", "csv")
	default:
		return "", ErrFlagFormatOption
	}

	from, to, err := extractMagnitude(magFlag)
	if err != nil {
		return "", err
	}

	if len(from) != 0 {
		v.Set("minmagnitude", from)
	}

	if len(to) != 0 {
		v.Set("maxmagnitude", to)
	}

	startTime, endTime, err := extractTime(dateTimeFlag)
	if err != nil {
		return "", err
	}

	if len(startTime) != 0 {
		v.Set("starttime", startTime)
	}

	if len(endTime) != 0 {
		v.Set("endtime", endTime)
	}

	// Prepare URL Request
	fullURL, err := url.Parse(FDSNENDPOINT)
	if err != nil {
		return "", err
	}

	path, err := url.JoinPath(fullURL.Path, endCmd)
	if err != nil {
		return "", err
	}

	fullURL.Path = path
	fullURL.ForceQuery = true
	fullURL.RawQuery = v.Encode()
	return fullURL.String(), nil
}

/* >= <= 4.0 and ranges 4.45-6.0....*/
func extractMagnitude(mFlag string) (string, string, error) {

	if len(mFlag) == 0 {
		return "", "", nil
	}

	if strings.ContainsAny(mFlag, ALPHABET) {
		return "", "", ErrFlagMagOption
	}

	// process magnitude range
	if strings.Contains(mFlag, "-") && strings.Count(mFlag, "-") == 1 {
		fields := strings.Split(mFlag, "-")
		lower := strings.TrimSpace(fields[0])
		upper := strings.TrimSpace(fields[1])
		return lower, upper, nil
	}

	if strings.Contains(mFlag, ",") {
		fields := strings.Split(mFlag, ",")
		lower := strings.TrimSpace(fields[0])
		upper := strings.TrimSpace(fields[1])
		return lower, upper, nil
	}

	// minimum magnitude
	if strings.HasPrefix(mFlag, ">") {
		lower := strings.TrimPrefix(mFlag, ">")
		return lower, "", nil
	}

	// maximum magnitude
	if strings.HasPrefix(mFlag, "<") {
		upper := strings.TrimPrefix(mFlag, "<")
		return "", upper, nil
	}

	// exact magnitude
	if strings.ContainsAny(mFlag, DIGITS) {
		lower := strings.TrimSpace(mFlag)
		upper := strings.TrimSpace(mFlag)
		return lower, upper, nil
	}

	return "", "", ErrFlagMagOption
}

func extractTime(tFlag string) (string, string, error) {

	if len(tFlag) == 0 {
		return "", "", nil
	}

	if strings.Contains(tFlag, ",") {
		fields := strings.Split(tFlag, ",")
		begin := strings.TrimSpace(fields[0])
		end := strings.TrimSpace(fields[1])

		if len(begin) == 0 && len(end) == 0 {
			return "", "", nil
		}

		validBegin, err := parseTime(begin)
		if err != nil {
			return "", "", ErrFlagTimeOption
		}

		validEnd, err := parseTime(end)
		if err != nil {
			return "", "", ErrFlagTimeOption
		}

		return validBegin, validEnd, nil
	}

	return "", "", ErrFlagTimeOption
}

func parseTime(timeStr string) (string, error) {

	if strings.Count(timeStr, ":") == 2 && strings.Count(timeStr, "T") == 1 {
		timeStr = timeStr + "Z"
	}

	timeFormats := []string{time.RFC3339, time.DateOnly}
	for _, tf := range timeFormats {
		if t, err := time.Parse(tf, timeStr); err == nil {
			return t.Format(tf), nil
		}
	}
	return "", ErrFlagTimeOption
}

// ExtractId resolves an eventid flag value, a format flag and query method to
// return a complete URL for a request.
func ExtractId(endCmd, formatFlag, id string) (string, error) {

	v := url.Values{}

	switch formatFlag {
	case "table":
		fallthrough
	case "geojson":
		fallthrough
	case "json":
		v.Set("format", "geojson")
	case "text":
		v.Set("format", "text")
	case "csv":
		v.Set("format", "csv")
	default:
		return "", ErrFlagFormatOption
	}

	validId, err := validateId(id)
	if err != nil {
		return "", err
	}

	v.Set("eventid", validId)

	// Prepare URL Request
	fullURL, err := url.Parse(FDSNENDPOINT)
	if err != nil {
		return "", err
	}

	path, err := url.JoinPath(fullURL.Path, endCmd)
	if err != nil {
		return "", err
	}

	fullURL.Path = path
	fullURL.ForceQuery = true
	fullURL.RawQuery = v.Encode()
	return fullURL.String(), nil
}

func validateId(id string) (string, error) {

	if strings.ContainsAny(id, NONALPHABET) {
		return "", ErrEventIdInvalid
	}

	validId := strings.TrimSpace(id)

	// check if the eventid contains spaces between characters
	if strings.ContainsAny(validId, " ") {
		return "", ErrEventIdInvalid	
	}

	return validId, nil
}
