package cfg

import "flag"

var (
	verboseFlag    = flag.Bool("v", false, "Enables debug messages")
	quickStartFlag = flag.Bool("q", false, "Enabled quick start")
)

type Flags struct {
	Verbose bool
	Quick   bool
}

func ReadFlags() Flags {
	flag.Parse()

	return Flags{
		Verbose: *verboseFlag,
		Quick:   *quickStartFlag,
	}
}
