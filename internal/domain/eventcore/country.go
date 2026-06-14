package eventcore

import "strings"

type Country struct {
	ISO2       string
	ISO3       string
	IOCCode    string
	Name       string
	CodeNum    int
	Population int64
	AreaKm2    float64
	GDPUSD     float64
}

// EmojiFlag derives the Discord flag-emoji shortcode from the ISO2 code
// (e.g. "BR" -> ":flag_br:"), or "" when ISO2 is empty.
func (c Country) EmojiFlag() string {
	if c.ISO2 == "" {
		return ""
	}
	return ":flag_" + strings.ToLower(c.ISO2) + ":"
}

// IsThis matches value against name, IOC, ISO2 or ISO3 (case-insensitive).
func (c Country) IsThis(value string) bool {
	value = strings.ToLower(value)
	switch value {
	case strings.ToLower(c.Name),
		strings.ToLower(c.IOCCode),
		strings.ToLower(c.ISO2),
		strings.ToLower(c.ISO3):
		return value != ""
	}
	return false
}
