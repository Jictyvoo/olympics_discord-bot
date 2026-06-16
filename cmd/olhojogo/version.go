package main

import (
	"fmt"
	"runtime/debug"
)

// version is overridable via -ldflags "-X main.version=<v>"; when unset it falls
// back to the VCS revision from runtime/debug.
var version = ""

const gitShortShaLen = 7

func printVersion() {
	fmt.Printf("olhojogo %s\n", resolveVersion())
}

func resolveVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	var rev, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= gitShortShaLen {
				rev = s.Value[:gitShortShaLen]
			} else {
				rev = s.Value
			}
		case "vcs.modified":
			if s.Value == "true" {
				modified = "-dirty"
			}
		}
	}
	if rev == "" {
		return "dev"
	}
	return rev + modified
}
