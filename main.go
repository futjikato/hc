package main

import (
	"flag"
)

func main () {
	specCpu := flag.Bool("spec-cpu", false, "Spec for CPU heavy container. Primary health stat is load.")
	specIo := flag.Bool("spec-storage", false, "Spec for heavy disc usage. Primary health stat is IO.")
	specNet := flag.Bool("spec-network", false, "Spec for heavy network container. Primary health stat is traffic.")

	flag.Parse()

	var spec Spec
	if (*specIo) {
		spec = IoSpec{}
	}

	if (*specNet) {
		if (spec != nil) {
			panic("You must only provide one spec.")
		}
		spec = NetSpec{}
	}

	if (*specCpu) {
		if (spec != nil) {
			panic("You must only provide one spec.")
		}
		spec = CpuSpec{}
	}

	if (spec == nil) {
		panic("You must provide a spec.")
	}

	spec.getStats()
}