package helpers

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/viant/toolbox"
)

func ParseDays(s string, allowInfinity bool) (int64, error) {
	r := regexp.MustCompile(`^(((\d+) days?)|(infinity))$`)
	matchedValue := r.FindStringSubmatch(s)
	if len(matchedValue) == 0 {
		return 0, fmt.Errorf("value `%s` is invalid, `X day`, `X days` or `infinity` expected", s)
	}
	var parsedFrom int64
	var err error
	if !allowInfinity && matchedValue[0] == "infinity" {
		return 0, fmt.Errorf("value `infinity` is invalid")
	} else if allowInfinity && matchedValue[0] == "infinity" {
		return math.MaxInt64, nil
	} else {
		parsedFrom, err = strconv.ParseInt(matchedValue[3], 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return parsedFrom, nil
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func ParseDate(name string, namingPattern string) (*time.Time, error) {
	dateLayout := toolbox.DateFormatToLayout(namingPattern)
	timeValue, err := time.Parse(dateLayout, name)
	if err != nil {
		return nil, err
	}
	return &timeValue, nil
}
