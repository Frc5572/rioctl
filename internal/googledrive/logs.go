package googledrive

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type LogDestination struct {
	Year  string
	Event string
}

var (
	matchLogPattern = regexp.MustCompile(
		`^akit_(\d{2})-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}_([A-Za-z0-9]+)_[pqe]\d+\.wpilog$`,
	)
	practiceLogPattern = regexp.MustCompile(
		`^akit_(\d{2})-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}\.wpilog$`,
	)
)

func ClassifyLog(filename string) (LogDestination, error) {
	filename = filepath.Base(filename)

	if matches := matchLogPattern.FindStringSubmatch(filename); matches != nil {
		return LogDestination{
			Year:  "20" + matches[1],
			Event: strings.ToUpper(matches[2]),
		}, nil
	}

	if matches := practiceLogPattern.FindStringSubmatch(filename); matches != nil {
		return LogDestination{
			Year:  "20" + matches[1],
			Event: "Test",
		}, nil
	}

	return LogDestination{}, fmt.Errorf(
		"unrecognized log filename: %s", filename,
	)
}
