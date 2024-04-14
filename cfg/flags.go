package cfg

import "flag"

var (
	verboseFlag = flag.Bool("v", false, "Enables debug messages")
)

type Flags struct {
	Verbose bool
}

func ReadFlags() Flags {
	flag.Parse()

	return Flags{
		Verbose: *verboseFlag,
	}
}
