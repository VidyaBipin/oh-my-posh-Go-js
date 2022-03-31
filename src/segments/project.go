package segments

import (
	"encoding/json"
	"oh-my-posh/environment"
	"oh-my-posh/properties"

	"github.com/BurntSushi/toml"
)

type ProjectItem struct {
	Name    string
	File    string
	Fetcher func(item ProjectItem) (string, string)
}

//Node.js Package
type NodePackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Rust Cargo package
type CargoTOML struct {
	Package CargoPackageTOML
}

type CargoPackageTOML struct {
	Name    string
	Version string
}

// Python Poetrey package
type PyprojectTOML struct {
	Tool PyprojectToolTOML
}

type PyprojectToolTOML struct {
	Poetry PyprojectPoetryTOML
}

type PyprojectPoetryTOML struct {
	Name    string
	Version string
}

type Project struct {
	props properties.Properties
	env   environment.Environment

	projects []*ProjectItem
	Version  string
	Name     string
	Error    string
}

func (n *Project) Enabled() bool {
	var enabled = false
	for _, item := range n.projects {
		if !enabled {
			enabled = n.hasProjectFile(item)
		}
	}

	return enabled
}

func (n *Project) Template() string {
	return " {{ if .Error }}{{ .Error }}{{ else }}{{ if .Version }}\uf487 {{.Version}}{{ end }} {{ if .Name }}{{ .Name }}{{ end }}{{ end }} "
}

func (n *Project) Init(props properties.Properties, env environment.Environment) {
	n.props = props
	n.env = env

	n.projects = []*ProjectItem{
		{
			Name:    "node",
			File:    "package.json",
			Fetcher: n.getNodePackage,
		},
		{
			Name:    "cargo",
			File:    "Cargo.toml",
			Fetcher: n.getCargoPackage,
		},
		{
			Name:    "poetry",
			File:    "pyproject.toml",
			Fetcher: n.getPoetryPackage,
		},
	}

	n.Version = ""
	n.Name = ""
	for _, item := range n.projects {
		if n.hasProjectFile(item) {
			n.Version, n.Name = item.Fetcher(*item)
			break
		}
	}
}

func (n *Project) hasProjectFile(p *ProjectItem) bool {
	return n.env.HasFiles(p.File)
}

func (n *Project) getNodePackage(item ProjectItem) (string, string) {
	content := n.env.FileContent(item.File)

	var data NodePackageJSON
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		n.Error = err.Error()
		return "", ""
	}

	return data.Version, data.Name
}

func (n *Project) getCargoPackage(item ProjectItem) (string, string) {
	content := n.env.FileContent(item.File)

	var data CargoTOML
	_, err := toml.Decode(content, &data)
	if err != nil {
		n.Error = err.Error()
		return "", ""
	}

	return data.Package.Version, data.Package.Name
}

func (n *Project) getPoetryPackage(item ProjectItem) (string, string) {
	content := n.env.FileContent(item.File)

	var data PyprojectTOML
	_, err := toml.Decode(content, &data)
	if err != nil {
		n.Error = err.Error()
		return "", ""
	}

	return data.Tool.Poetry.Version, data.Tool.Poetry.Name
}
