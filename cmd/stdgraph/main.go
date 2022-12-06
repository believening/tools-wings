package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

const initIndex int = -1

var (
	// dot 格式文件
	graphTemplate = `digraph {{.name}} {
node [shape=box];
{{ range $pkg, $id := .ids -}}
{{ $id }} [label="{{ $pkg }}"];
{{ end -}}
{{- range $modId, $deps := .deps -}}
{{- range $depId, $_ := $deps -}}
{{ $modId }} -> {{ $depId }};
{{  end -}}
{{- end -}}
}
`

	// IDS and DEPS holds the topology relationships as shown below:
	// package A depends on package B and package C
	//   pkgA ---> [pkgB pkgC ...]
	// we use the hash value of pkgA as the key value of the
	// map and store package info under the corresponding key.
	IDS    map[string]int
	DEPS   map[int]map[int]struct{}
	serial int

	path string
	name string
)

type Package struct {
	ImportPath string   `json:",omitempty"` // import path of package in dir
	Imports    []string `json:",omitempty"` // import paths used by this package
}

func putPkg(pkg Package) {
	id := getOrPutID(pkg.ImportPath)
	_, ok := DEPS[id]
	if !ok {
		DEPS[id] = make(map[int]struct{})
	}
	for _, p := range pkg.Imports {
		pid := getOrPutID(p)
		DEPS[id][pid] = struct{}{}
	}
}

func getOrPutID(pkg string) int {
	id, ok := IDS[pkg]
	if !ok {
		id = serial
		IDS[pkg] = id
		serial++
	}
	return id
}

func logF(v ...interface{}) {
	_, f, l, _ := runtime.Caller(1)
	log.Fatalln(fmt.Sprintf("[%s:%d]", f, l), v)
}

func cmdGoList() {
	cmd := exec.Command("go", []string{"list", "-json", "./..."}...)
	cmd.Dir = filepath.Clean(path)
	name = filepath.Base(cmd.Dir)

	var out, eOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &eOut

	err := cmd.Run()
	if err != nil {
		logF(err)
	}
	rawPkgs := out.Bytes()
	if len(rawPkgs) == 0 {
		logF("len of stdout is zero,", eOut.String())
	}

	var pkgs []Package
	beg, end, pairCnt := initIndex, initIndex, 0
	for i := 0; i < len(rawPkgs); i++ {
		switch rawPkgs[i] {
		case '{':
			if pairCnt&1 == 0 {
				beg = i
				if end != initIndex {
					logF(fmt.Sprintf("beg %d and end %d not match", beg, end))
				}
			}
			pairCnt++
		case '}':
			if pairCnt&1 == 1 {
				end = i
				if beg == initIndex {
					logF(fmt.Sprintf("beg %d and end %d not match", beg, end))
				}
				var tempPkg Package
				_ = json.Unmarshal(rawPkgs[beg:end+1], &tempPkg)
				pkgs = append(pkgs, tempPkg)
				beg, end = initIndex, initIndex
			}
			pairCnt++
		default:
			continue
		}
	}

	for _, p := range pkgs {
		putPkg(p)
	}
}

func render() {
	tpl, err := template.New(name).Parse(graphTemplate)
	if err != nil {
		logF(err)
	}
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, map[string]interface{}{
		"name": name,
		"ids":  IDS,
		"deps": DEPS,
	})
	fmt.Println(buf.String())
}

func init() {
	IDS = make(map[string]int)
	DEPS = make(map[int]map[int]struct{})
	serial = 1

	flag.StringVar(&path, "p", filepath.Join(runtime.GOROOT(), "src"), "root path of the go codes")
	flag.Parse()
}

func main() {
	cmdGoList()
	render()
}
