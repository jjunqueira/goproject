package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/jjunqueira/goproject/pkg/goproject"
	"github.com/mitchellh/go-homedir"
)

// Project contains information regarding the project we are generating
type Project struct {
	gitPrefix   string
	projectName string
	tpl         *Template
}

// NewProject constructs a new Project struct with the provided settings
func NewProject(c *goproject.Config, gitPrefix string, tplName string, projectName string) (*Project, error) {
	p := new(Project)
	p.gitPrefix = gitPrefix
	p.projectName = projectName
	tpl, err := Find(c, tplName)
	if err != nil {
		return nil, err
	}
	p.tpl = tpl
	return p, nil
}

// Template represents a template that exists on disk
type Template struct {
	name string
	path string
}

// Find attempts to load a template by name and the given search paths
func Find(c *goproject.Config, tplName string) (*Template, error) {
	// search for custom templates first, they will take precedence over base templates
	for _, t := range c.CustomTemplates {
		if t.Name == tplName {
			return &Template{name: t.Name, path: t.Path}, nil
		}
	}

	// No custom template was found, try to find a default one that matches the name
	var tpl *Template
	err := filepath.Walk(c.TemplatesPath, func(path string, info os.FileInfo, err error) error {
		if info.Name() == tplName {
			tpl = &Template{name: info.Name(), path: path}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if tpl != nil {
		return tpl, nil
	}
	return nil, fmt.Errorf("unable to find template '%s' in default or custom template paths", tplName)
}

// Generate generates a new project based on a template
func Generate(c *goproject.Config, p *Project) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(cwd, p.projectName)
	err = os.Mkdir(fullPath, 0777)
	if err != nil {
		return err
	}

	err = gitInit(fullPath)
	if err != nil {
		return err
	}

	err = initGoModule(fullPath, p.gitPrefix, p.projectName)
	if err != nil {
		return err
	}

	err = copyTemplate("empty", fullPath)
	return err
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q", dir)
	return cmd.Run()
}

func initGoModule(dir string, gitPrefix string, projectname string) error {
	moduleName := fmt.Sprintf("%s/%s", gitPrefix, projectname)
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	return cmd.Run()
}

func copyTemplate(tplName string, dest string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	tplPath := path.Join(home, ".config", "goproject", "templates", tplName)
	if err != nil {
		return err
	}
	err = copyFiles(tplPath, dest)
	return err
}

func copyFiles(src string, dest string) error {
	cmd := exec.Command("cp", "-r", path.Join(src, "*"), path.Join(dest, "*"))
	return cmd.Run()
}
