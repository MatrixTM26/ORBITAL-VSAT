package printer

import (
	"fmt"
	"strings"
)

func FmtBytes(B float64) string {
	switch {
	case B < 1024:
		return fmt.Sprintf("%.0fB", B)
	case B < 1048576:
		return fmt.Sprintf("%.1fKB", B/1024)
	case B < 1073741824:
		return fmt.Sprintf("%.1fMB", B/1048576)
	default:
		return fmt.Sprintf("%.2fGB", B/1073741824)
	}
}

func FmtRate(N float64) string {
	if N < 1000 {
		return fmt.Sprintf("%.1f/s", N)
	}
	return fmt.Sprintf("%.2fk/s", N/1000)
}

func ProgressBar(Done, Total, Width int) string {
	if Total == 0 {
		return strings.Repeat("░", Width)
	}
	Filled := int(float64(Done) / float64(Total) * float64(Width))
	if Filled > Width {
		Filled = Width
	}
	return strings.Repeat("█", Filled) + strings.Repeat("░", Width-Filled)
}

func PadRight(S string, N int) string {
	if len(S) >= N {
		return S
	}
	return S + strings.Repeat(" ", N-len(S))
}
