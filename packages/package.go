package packages

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yourbase/yb/buildpacks"
	"github.com/yourbase/yb/internal/ybdata"
	"github.com/yourbase/yb/plumbing"
	"github.com/yourbase/yb/types"
	"gopkg.in/yaml.v2"
)

type Package struct {
	Name     string
	Path     string
	Manifest types.BuildManifest
}

func LoadPackage(name string, path string) (*Package, error) {
	manifest := types.BuildManifest{}
	buildYaml := filepath.Join(path, types.MANIFEST_FILE)
	if _, err := os.Stat(buildYaml); os.IsNotExist(err) {
		return nil, fmt.Errorf("Can't load %s: %v", types.MANIFEST_FILE, err)
	}

	buildyaml, _ := ioutil.ReadFile(buildYaml)
	err := yaml.Unmarshal([]byte(buildyaml), &manifest)
	if err != nil {
		return nil, fmt.Errorf("Error loading %s for %s: %v", types.MANIFEST_FILE, name, err)
	}
	err = mergeDeps(&manifest)
	if err != nil {
		return nil, fmt.Errorf("Error loading %s for %s: %v", types.MANIFEST_FILE, name, err)
	}

	return &Package{
		Path:     path,
		Name:     name,
		Manifest: manifest,
	}, nil
}

func (p Package) BuildRoot(dataDirs *ybdata.Dirs) (string, error) {
	// Are we a part of a workspace?
	workspaceDir, err := plumbing.FindWorkspaceRoot()
	if err != nil {
		// Nope, just ourselves...
		h := sha256.New()
		h.Write([]byte(p.Path))
		workspaceHash := fmt.Sprintf("%x", h.Sum(nil))
		workspaceDir = filepath.Join(dataDirs.Workspaces(), workspaceHash[0:12])
	}

	buildRoot := filepath.Join(workspaceDir, "build")
	if err := os.MkdirAll(buildRoot, 0777); err != nil {
		return "", fmt.Errorf("create workspace build directory: %w", err)
	}
	return buildRoot, nil
}

func (p Package) SetupBuildDependencies(ctx context.Context, dataDirs *ybdata.Dirs, target *types.BuildTarget) error {
	buildRoot, err := p.BuildRoot(dataDirs)
	if err != nil {
		return err
	}
	for _, dep := range target.Dependencies.Build {
		if err := buildpacks.Install(ctx, dataDirs, buildRoot, p.Path, dep); err != nil {
			return err
		}
	}
	return nil
}

func (p Package) SetupRuntimeDependencies(ctx context.Context, dataDirs *ybdata.Dirs) error {
	buildRoot, err := p.BuildRoot(dataDirs)
	if err != nil {
		return err
	}
	for _, dep := range p.Manifest.Dependencies.Runtime {
		if err := buildpacks.Install(ctx, dataDirs, buildRoot, p.Path, dep); err != nil {
			return err
		}
	}
	for _, dep := range p.Manifest.Dependencies.Build {
		if err := buildpacks.Install(ctx, dataDirs, buildRoot, p.Path, dep); err != nil {
			return err
		}
	}
	return nil
}
