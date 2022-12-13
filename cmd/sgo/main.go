// sgo used to switch go version
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

var (
	goroot                       string = "/usr/local/go"
	currentVersion, currentShort string = currentGo()

	supportVersions      map[string]struct{} = map[string]struct{}{}
	shortSupportVersions map[string]struct{} = map[string]struct{}{}
)

func currentGo() (string, string) {
	out, _ := exec.Command(goroot+"/bin/go", "version").CombinedOutput()
	// like: `go version go1.19 linux/amd64`
	version := strings.Fields(strings.TrimSpace(string(out)))[2]
	return version, version[2:]
}

func initSupportVersions() {
	supportVersions[currentVersion] = struct{}{}
	shortSupportVersions[currentVersion[2:]] = struct{}{}
	entry, err := os.ReadDir(filepath.Dir(goroot))
	if err != nil {
		log.Panic(err)
	}
	for _, e := range entry {
		if !e.IsDir() {
			continue
		}
		if idx := strings.Index(e.Name(), "go"); idx == -1 || len(e.Name()) == idx+2 {
			continue
		} else {
			supportVersions[e.Name()] = struct{}{}
			shortSupportVersions[e.Name()[2:]] = struct{}{}
		}
	}
}

func main() {
	initSupportVersions()

	if len(os.Args) == 1 {
		listSupportVersions()
		os.Exit(0)
	}
	switch str := os.Args[1]; str {
	case "list", "l":
		listSupportVersions()
		os.Exit(0)
	default:
		if _, exist := supportVersions[str]; exist {
			switchVersion(str)
			os.Exit(0)
		}
		if _, exist := shortSupportVersions[str]; exist {
			switchVersion("go" + str)
			os.Exit(0)
		}
		listSupportVersions()
		os.Exit(1)
	}
}

func listSupportVersions() {
	versions := make([]string, 0, len(shortSupportVersions))
	for version := range shortSupportVersions {
		versions = append(versions, version)
	}
	sort.Strings(versions)
	for _, version := range versions {
		tpl := " %s\n"
		if version == currentShort {
			tpl = "-%s\n"
		}
		fmt.Printf(tpl, version)
	}
}

func switchVersion(targetVersion string) {
	if targetVersion == currentVersion {
		return
	}

	if err := os.Rename(goroot, filepath.Join(filepath.Dir(goroot), currentVersion)); err != nil {
		log.Panic(err)
	}
	if err := os.Rename(filepath.Join(filepath.Dir(goroot), targetVersion), goroot); err != nil {
		log.Panic(err)
	}
}
