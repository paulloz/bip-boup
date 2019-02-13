package fr

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olebedev/when/rules"
)

/*
	17h20 | 17 h 20 | 17H20 | 17 H 20
	17:20 | 17 : 20
	17 heures 20
*/

func HourMinute(s rules.Strategy) rules.Rule {
	return &rules.F{
		RegExp: regexp.MustCompile("(?i)(?:\\W|^)" +
			"(" + INTEGER_WORDS_PATTERN + "|\\d{1,2})" +
			"(?:\\s*(:|h|heures?))" +
			"\\s*(" + INTEGER_WORDS_PATTERN + "|\\d{1,2}|et\\s+demi|et\\s+quart|trois(\\s+|-)quarts?)" +
			"(?:\\W|$)"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			if (c.Hour != nil || c.Minute != nil) && s != rules.Override {
				return false, nil
			}

			var err error

			hourStr := strings.TrimSpace(m.Captures[0])
			var hour int

			if n, ok := INTEGER_WORDS[hourStr]; ok {
				hour = n
			} else if hour, err = strconv.Atoi(hourStr); err != nil {
				return false, err
			}

			minutesStr := strings.TrimSpace(m.Captures[2])
			var minutes int

			if n, ok := INTEGER_WORDS[minutesStr]; ok {
				minutes = n
			} else if minutes, err = strconv.Atoi(minutesStr); err != nil {
				cap := []byte(minutesStr)
				if ok, _ := regexp.Match("et\\s+demi", cap); ok {
					minutes = 30
				} else if ok, _ := regexp.Match("et\\s+quart", cap); ok {
					minutes = 15
				} else if ok, _ := regexp.Match("trois(\\s+|-)quarts?", cap); ok {
					minutes = 45
				} else {
					return false, err
				}
			}

			if hour > 23 || minutes > 59 {
				return false, nil
			}

			c.Hour = &hour
			c.Minute = &minutes

			return true, nil
		},
	}
}
