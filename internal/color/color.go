package color

const (
	Reset        = "\033[0m"
	Black        = "\033[30m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	DarkGreen    = "\033[38;5;22m"
	Orange       = "\033[38;5;208m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Magenta      = "\033[35m"
	Cyan         = "\033[36m"
	White        = "\033[37m"
	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"
	Bold         = "\033[1m"
	Dim          = "\033[2m"
)

func C(Text, Color string) string { return Color + Text + Reset }
func B(Text string) string        { return Bold + Text + Reset }
