package strings

import (
	"strings"
)

func Choose(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}

	return
}

func Every(ss []string, f func(string) string) (ret []string) {
	for _, s := range ss {
		ret = append(ret, f(s))
	}

	return
}

func Contains(ss []string, s string) bool {
	for _, _s := range ss {
		if _s == s {
			return true
		}
	}

	return false
}

func Capitalize(s string) string {
	if len(s) <= 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + string(s[1:])
}

func Reverse(ss []string) []string {
	for i := 0; i < (len(ss) / 2); i++ {
		j := len(ss) - i - 1
		ss[i], ss[j] = ss[j], ss[i]
	}
	return ss
}
