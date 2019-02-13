package fr

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olebedev/when/rules"
)

func Hour(s rules.Strategy) rules.Rule {
	return &rules.F{
		RegExp: regexp.MustCompile("(?i)(?:\\W|^)" +
			"(" + INTEGER_WORDS_PATTERN + "|\\d{1,2})" +
			"(?:\\s*(h|heures?))" +
			"(?:\\W|$)"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			if c.Hour != nil && s != rules.Override {
				return false, nil
			}

			hourStr := strings.TrimSpace(m.Captures[0])
			var hour int
			var err error

			if n, ok := INTEGER_WORDS[hourStr]; ok {
				hour = n
			} else if hour, err = strconv.Atoi(hourStr); err != nil {
				return false, err
			}

			if hour > 24 {
				return false, nil
			}

			zero := 0

			c.Hour = &hour
			c.Minute = &zero

			return true, nil
		},
	}
}
