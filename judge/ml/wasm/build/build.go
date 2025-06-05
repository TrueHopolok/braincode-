package main

//go:generate go run . -dir ../../../../frontend/static/wasm/ -src ../src

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func try[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func try1(err error) {
	if err != nil {
		panic(err)
	}
}

func globalize(path string) string {
	return try(filepath.Abs(path))
}

func main() {
	dst := flag.String("dir", ".", "dir where all files will be placed")
	src := flag.String("src", ".", "wasm module source location")
	flag.Parse()

	*dst = globalize(*dst)
	*src = globalize(*src)

	goroot := strings.TrimSpace(string(try(exec.Command("go", "env", "GOROOT").Output())))

	wasmOut := globalize(filepath.Join(*dst, "markleft.wasm"))

	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", wasmOut, globalize(*src))
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm", "GOTMPDIR="+os.TempDir(), "GOCACHE="+filepath.Join(os.TempDir(), "go-build"))

	res, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(res)
	}
	try1(err)

	data := try(os.ReadFile(filepath.Join(goroot, "lib/wasm/wasm_exec.js")))

	try1(os.WriteFile(globalize(filepath.Join(*dst, "wasm_exec.js")), data, 0666))
}
