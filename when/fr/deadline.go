package fr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olebedev/when/rules"
)

func Deadline(s rules.Strategy) rules.Rule {
	overwrite := s == rules.Override

	return &rules.F{
		RegExp: regexp.MustCompile(
			"(?i)(?:\\W|^)" +
				"(dans)\\s*" +
				"(" + INTEGER_WORDS_PATTERN + "|[0-9]+)\\s*" +
				// "(?:demi(?:-|\\s+))" +
				"(sec(?:onde)?s?|min(?:ute)?s?|heures?|jours?|semaines?|mois|ans?)\\s*" +
				"(et\\s+demi)?" +
				"(?:\\W|$)"),
		Applier: func(m *rules.Match, c *rules.Context, o *rules.Options, ref time.Time) (bool, error) {
			numStr := strings.TrimSpace(m.Captures[1])

			var num int
			var err error

			if n, ok := INTEGER_WORDS[numStr]; ok {
				num = n
			} else {
				num, err = strconv.Atoi(numStr)
				if err != nil {
					return false, fmt.Errorf("could not convert `%s` to int", numStr)
				}
			}

			exponent := strings.TrimSpace(m.Captures[2])

			switch {
			case strings.Contains(exponent, "sec"):
				if c.Duration == 0 || overwrite {
					c.Duration = time.Duration(num) * time.Second
				}
			case strings.Contains(exponent, "min"):
				if c.Duration == 0 || overwrite {
					c.Duration = time.Duration(num) * time.Minute
				}
			case strings.Contains(exponent, "heure"):
				if c.Duration == 0 || overwrite {
					c.Duration = time.Duration(num) * time.Hour
				}
			case strings.Contains(exponent, "jour"):
				if c.Duration == 0 || overwrite {
					c.Duration = time.Duration(num) * 24 * time.Hour
				}
			case strings.Contains(exponent, "semaine"):
				if c.Duration == 0 || overwrite {
					c.Duration = time.Duration(num) * 7 * 24 * time.Hour
				}
			case strings.Contains(exponent, "mois"):
				if c.Month == nil || overwrite {
					// c.Month = pointer.ToInt((int(ref.Month()) + num) % 12)
				}
			case strings.Contains(exponent, "an"):
				if c.Year == nil || overwrite {
					// c.Year = pointer.ToInt(ref.Year() + num)
				}
			}

			return true, nil
		},
	}
}
