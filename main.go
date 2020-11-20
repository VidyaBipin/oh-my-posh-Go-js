package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
)

var currentVersion = "development"

type args struct {
	ErrorCode   *int
	PrintConfig *bool
	PrintShell  *bool
	Config      *string
	Shell       *string
	PWD         *string
	Version     *bool
	Debug       *bool
}

func main() {
	args := &args{
		ErrorCode: flag.Int(
			"error",
			0,
			"Error code of previously executed command"),
		PrintConfig: flag.Bool(
			"print-config",
			false,
			"Print the current config in json format"),
		PrintShell: flag.Bool(
			"print-shell",
			false,
			"Print the current shell name"),
		Config: flag.String(
			"config",
			"",
			"Add the path to a configuration you wish to load"),
		Shell: flag.String(
			"shell",
			"",
			"Override the shell you are working in"),
		PWD: flag.String(
			"pwd",
			"",
			"the path you are working in"),
		Version: flag.Bool(
			"version",
			false,
			"Print the current version of the binary"),
		Debug: flag.Bool(
			"debug",
			false,
			"Print debug information"),
	}
	flag.Parse()
	env := &environment{
		args: args,
	}
	settings := GetSettings(env)
	if *args.PrintConfig {
		theme, _ := json.MarshalIndent(settings, "", "    ")
		fmt.Println(string(theme))
		return
	}
	if *args.PrintShell {
		fmt.Println(env.getShellName())
		return
	}
	if *args.Version {
		fmt.Println(currentVersion)
		return
	}
	colorWriter := &Renderer{
		Buffer: new(bytes.Buffer),
	}
	shell := env.getShellName()
	if *args.Shell != "" {
		shell = *args.Shell
	}
	colorWriter.init(shell)
	engine := &engine{
		settings: settings,
		env:      env,
		renderer: colorWriter,
	}
	engine.render()
}
