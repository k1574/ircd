package format

import(
	"strings"
)


func
SplitTilSep(s, sep, endsep string) ([]string, string) {
	n := strings.LastIndex(s, endsep)
	var arg, str string
	if n != -1 {
		arg = s[:n]
		str = s[n+len(endsep):]
	} else {
		arg = s
		str = ""
	}

	return strings.Split(arg, sep), str
}


func
HasAnyOfPrefixes(s string, prefs []string) string {
	for _, v := range prefs {
		if strings.HasPrefix(s, v) {
			return v
		}
	}
	return ""
}
