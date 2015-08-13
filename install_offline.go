package main

import (
	"os"
	"path/filepath"
)

func installOffline(args []string) error {
	allGoms, err := parseGomfile("Gomfile")
	if err != nil {
		return err
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	_, err = os.Stat(vendor)
	if err != nil {
		err = os.MkdirAll(vendor, 0755)
		if err != nil {
			return err
		}
	}
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}
	err = os.Setenv("GOBIN", filepath.Join(vendor, "bin"))
	if err != nil {
		return err
	}

	// 1. Filter goms to build
	goms := make([]Gom, 0)
	for _, gom := range allGoms {
		if group, ok := gom.options["group"]; ok {
			if !matchEnv(group) {
				continue
			}
		}
		if goos, ok := gom.options["goos"]; ok {
			if !matchOS(goos) {
				continue
			}
		}
		goms = append(goms, gom)
	}

	// 4. Build and install
	for _, gom := range goms {
		err = gom.Build(args)
		if err != nil {
			return err
		}
	}

	return nil
}
