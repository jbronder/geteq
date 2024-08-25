package logic

import "testing"

type MagnitudeTest struct {
	in, outBegin, outEnd string
	err                  error
}

type TimeTest struct {
	in, outBegin, outEnd string
	err                  error
}

type ParseTimeTest struct {
	in, out string
	err     error
}

type EventIdTest struct {
	in, out string
	err     error
}

func TestExtractMagnitude(t *testing.T) {
	mTests := []MagnitudeTest{
		{"6.0", "6.0", "6.0", nil},
		{" 6.0", "6.0", "6.0", nil},
		{"6.0 ", "6.0", "6.0", nil},
		{"4.0,6.0", "4.0", "6.0", nil},
		{" 4.0-6.0", "4.0", "6.0", nil},
		{"4.50, 6.2", "4.50", "6.2", nil},
		{">4.0", "4.0", "", nil},
		{"<4.0", "", "4.0", nil},
		{"--", "", "", ErrFlagMagOption},
		{"", "", "", nil},
	}

	for _, test := range mTests {
		begin, end, err := extractMagnitude(test.in)
		if begin != test.outBegin || end != test.outEnd || err != test.err {
			t.Errorf("extractMagnitude(%q) = %v %v %v; want %v %v %v", test.in, begin, end, err, test.outBegin, test.outEnd, test.err)
		}
	}
}

func TestExtractTime(t *testing.T) {
	timeTests := []TimeTest{
		{"2024-12-01,2024-12-02", "2024-12-01", "2024-12-02", nil},
		{"2024-12-01T12:34:56,2024-12-02T12:34:56", "2024-12-01T12:34:56Z", "2024-12-02T12:34:56Z", nil},
		{"2024-12-01, 2024-12-02", "2024-12-01", "2024-12-02", nil},
		{" 2024-12-01, 2024-12-02 ", "2024-12-01", "2024-12-02", nil},
		{"--", "", "", ErrFlagTimeOption},
		{" , ", "", "", nil},
		{",", "", "", nil},
		{"2024-12-01-2024-12-02", "", "", ErrFlagTimeOption},
	}

	for _, test := range timeTests {
		begin, end, err := extractTime(test.in)
		if begin != test.outBegin || end != test.outEnd || err != test.err {
			t.Errorf("extractTime(%q) = %v %v %v; want %v %v %v", test.in, begin, end, err, test.outBegin, test.outEnd, test.err)
		}
	}
}

func TestParseTime(t *testing.T) {
	tTests := []ParseTimeTest{
		{"2024-01-01", "2024-01-01", nil},
		{"2024-12-01T12:34:56", "2024-12-01T12:34:56Z", nil},
		{"-12-01", "", ErrFlagTimeOption},
		{"20-12-01", "", ErrFlagTimeOption},
		{"2024-12-01T", "", ErrFlagTimeOption},
	}

	for _, test := range tTests {
		timeVal, err := parseTime(test.in)
		if timeVal != test.out || err != test.err {
			t.Errorf("parseTime(%q) = %v %v; want %v %v", test.in, timeVal, err, test.out, test.err)
		}
	}
}

func TestValidateId(t *testing.T) {
	idTests := []EventIdTest{
		{"ci40012345", "ci40012345", nil},
		{"^ci12345678", "", ErrEventIdInvalid},
	}

	for _, test := range idTests {
		id, err := validateId(test.in)
		if id != test.out || err != test.err {
			t.Errorf("validateId(%q) = %v %v; want %v %v", test.in, id, err, test.out, test.err)
		}
	}
}
