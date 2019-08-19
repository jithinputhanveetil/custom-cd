package customcd

type messageType int

const (
	// WARN is a warning
	WARN messageType = iota
	//ERROR is an error
	ERROR
)

type colorCode string

const (
	// RED is the color code for red
	RED colorCode = "\033[01;91m"
	// YELLOW is the color code for yellow
	YELLOW colorCode = "\033[01;33m"
	// WHITE is the color code for white
	WHITE colorCode = "\033[0m"
)
