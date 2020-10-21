package types

import (
	"github.com/yourbase/narwhal"
)

const (
	MANIFEST_FILE = ".yourbase.yml"
	DOCS_URL      = "https://docs.yourbase.io"
)

type BuildManifest struct {
	Dependencies DependencySet  `yaml:"dependencies"`
	Sandbox      bool           `yaml:"sandbox"`
	BuildTargets []*BuildTarget `yaml:"build_targets"`
	Build        *BuildTarget   `yaml:"build"`
	Exec         *ExecPhase     `yaml:"exec"`
	Package      *PackagePhase  `yaml:"package"`
	CI           *CIInfo        `yaml:"ci"`
}

type CIInfo struct {
	CIBuilds []*CIBuild `yaml:"builds"`
}

type CIBuild struct {
	Name        string `yaml:"name"`
	BuildTarget string `yaml:"build_target"`
	When        string `yaml:"when"`
}

type PackagePhase struct {
	Artifacts []string `yaml:"artifacts"`
}

type DependencySet struct {
	Build   []string `yaml:"build"`
	Runtime []string `yaml:"runtime"`
}

type ExecPhase struct {
	Name         string               `yaml:"name"`
	Dependencies ExecDependencies     `yaml:"dependencies"`
	Container    *ContainerDefinition `yaml:"container"`
	Commands     []string             `yaml:"commands"`
	Ports        []string             `yaml:"ports"`
	Environment  map[string][]string  `yaml:"environment"`
	LogFiles     []string             `yaml:"logfiles"`
	Sandbox      bool                 `yaml:"sandbox"`
	HostOnly     bool                 `yaml:"host_only"`
	BuildFirst   []string             `yaml:"build_first"`
}

type ContainerDefinition struct {
	Image         string        `yaml:"image"`
	Mounts        []string      `yaml:"mounts"`
	Ports         []string      `yaml:"ports"`
	Environment   []string      `yaml:"environment"`
	Command       string        `yaml:"command"`
	WorkDir       string        `yaml:"workdir"`
	PortWaitCheck PortWaitCheck `yaml:"port_check"`
	Label         string        `yaml:"label"`
}

func (def *ContainerDefinition) ToNarwhal() *narwhal.ContainerDefinition {
	image := "yourbase/yb_ubuntu:18.04"
	if def == nil {
		return &narwhal.ContainerDefinition{Image: image}
	}
	if def.Image != "" {
		image = def.Image
	}
	return &narwhal.ContainerDefinition{
		Image:         image,
		Mounts:        append([]string(nil), def.Mounts...),
		Ports:         append([]string(nil), def.Ports...),
		Environment:   append([]string(nil), def.Environment...),
		Command:       def.Command,
		WorkDir:       def.WorkDir,
		PortWaitCheck: *def.PortWaitCheck.ToNarwhal(),
		Label:         def.Label,
	}
}

type PortWaitCheck struct {
	Port         int `yaml:"port"`
	Timeout      int `yaml:"timeout"`
	LocalPortMap int
}

func (check *PortWaitCheck) ToNarwhal() *narwhal.PortWaitCheck {
	return &narwhal.PortWaitCheck{
		Port:         check.Port,
		Timeout:      check.Timeout,
		LocalPortMap: check.LocalPortMap,
	}
}

type BuildDependencies struct {
	Build      []string                        `yaml:"build"`
	Containers map[string]*ContainerDefinition `yaml:"containers"`
}

func (b BuildDependencies) ContainerList() []*narwhal.ContainerDefinition {
	containers := make([]*narwhal.ContainerDefinition, 0, len(b.Containers))
	for label, c := range b.Containers {
		c.Label = label
		containers = append(containers, c.ToNarwhal())
	}
	return containers
}

type ExecDependencies struct {
	Containers map[string]*ContainerDefinition `yaml:"containers"`
}

func (b ExecDependencies) ContainerList() []*narwhal.ContainerDefinition {
	containers := make([]*narwhal.ContainerDefinition, 0, len(b.Containers))
	for label, c := range b.Containers {
		c.Label = label
		containers = append(containers, c.ToNarwhal())
	}
	return containers
}

type BuildTarget struct {
	Name         string               `yaml:"name"`
	Container    *ContainerDefinition `yaml:"container"`
	Commands     []string             `yaml:"commands"`
	HostOnly     bool                 `yaml:"host_only"`
	Root         string               `yaml:"root"`
	Environment  []string             `yaml:"environment"`
	Tags         map[string]string    `yaml:"tags"`
	BuildAfter   []string             `yaml:"build_after"`
	Dependencies BuildDependencies    `yaml:"dependencies"`
}

// A Project is a YourBase project as returned by the API.
type Project struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
	OrgSlug     string `json:"organization_slug"`
}

type TokenResponse struct {
	TokenOK bool `json:"token_ok"`
}

type WorktreeSave struct {
	Hash    string
	Path    string
	Files   []string
	Enabled bool
}
