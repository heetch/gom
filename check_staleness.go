package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrStaledDependencies = fmt.Errorf("Dependencies staled. Run `gom install` to fix the issue")
)

func getVcsCommand(vendor string, path string) (*vcsCmd, string, error) {

	for {
		if isDir(filepath.Join(path, ".git")) {
			return git, path, nil
		} else if isDir(filepath.Join(path, ".hg")) {
			return hg, path, nil
		} else if isDir(filepath.Join(path, ".bzr")) {
			return bzr, path, nil
		}

		path = filepath.Clean(filepath.Join(path, ".."))
		if path == vendor {
			break
		}
	}
	return nil, "", fmt.Errorf("Unable to get the VCS")
}

func locateGomfile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		file := filepath.Join(dir, "Gomfile")
		if isFile(file) {
			return file, nil
		}
		next := filepath.Clean(filepath.Join(dir, ".."))
		if next == dir {
			break
		}
		dir = next
	}

	return "", errors.New("Can't locate Gomfile")
}

func checkStaleness() error {
	gomfile, err := locateGomfile()
	if err != nil {
		return err
	}

	var allGoms []Gom
	allGoms, err = parseGomfile(gomfile)
	if err != nil {
		return err
	}

	for _, g := range allGoms {
		rawCommit := g.options["commit"]
		commit, ok := rawCommit.(string)
		if !ok || commit == "" {
			fmt.Printf("[%s] No commit set. Please set a revion with :commit => 'SHA1'\n", g.name)
			continue
		}

		vendor, err := filepath.Abs(filepath.Join(filepath.Dir(gomfile), vendorFolder))
		if err != nil {
			return err
		}

		p := filepath.Join(vendor, "src", g.name)
		if !isDir(p) {
			return ErrStaledDependencies
		}

		vcs, path, err := getVcsCommand(vendor, p)
		if err != nil {
			return err
		}

		if vcs == nil {
			fmt.Printf("[%s] Unable to check the revision. Reason: unknown VCS\n", g.name)
			continue
		}

		revision, err := vcs.Revision(path)
		if err != nil {
			return err
		}

		if commit != revision {
			return ErrStaledDependencies
		}
	}

	return nil
}
