package youtube

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	patternDuration      = regexp.MustCompile(`PT(?P<H>\d+H)?(?P<M>\d+M)?(?P<S>\d+S)?`)
	patternDurationNames = patternDuration.SubexpNames()
)

func ParseDuration(input string) (int64, error) {
	matches := patternDuration.FindStringSubmatch(input)
	results := map[string]int64{}
	for i, name := range patternDurationNames {
		if name == "" || matches[i] == "" {
			continue
		}

		trimmed := strings.Trim(matches[i], "HMS")
		value, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return 0, err
		}

		results[name] = value
	}

	var total int64 = 0

	if value, ok := results["H"]; ok {
		total += value * 60 * 60
	}

	if value, ok := results["M"]; ok {
		total += value * 60
	}

	if value, ok := results["S"]; ok {
		total += value
	}

	return total, nil
}
