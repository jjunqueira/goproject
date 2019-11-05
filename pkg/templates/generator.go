package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/jjunqueira/goproject/pkg/goproject"
)

// Project contains information regarding the project we are generating
type Project struct {
	GitPrefix string
	Name      string
	Tpl       *Template
}

// NewProject constructs a new Project struct with the provided settings
func NewProject(c *goproject.Config, gitPrefix string, tplName string, projectName string) (*Project, error) {
	p := new(Project)
	p.GitPrefix = gitPrefix
	p.Name = projectName
	tpl, err := Find(c, tplName)
	if err != nil {
		return nil, err
	}
	p.Tpl = tpl
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
		return fmt.Errorf("unable to get current working directory: %v", err)
	}

	fullPath := filepath.Join(cwd, p.Name)
	err = os.Mkdir(fullPath, 0777)
	if err != nil {
		return fmt.Errorf("unable to create project directory: %v", err)
	}

	err = gitInit(fullPath)
	if err != nil {
		return fmt.Errorf("unable to initialize git: %v", err)
	}

	err = initGoModule(fullPath, p.GitPrefix, p.Name)
	if err != nil {
		return fmt.Errorf("unable to initialize Go Modules: %v", err)
	}

	err = copyFiles(p.Tpl.path, fullPath)
	if err != nil {
		return fmt.Errorf("unable to copy template files: %v", err)
	}

	err = applyProjectToTemplates(p, fullPath)
	if err != nil {
		return fmt.Errorf("unable to execute templates: %v", err)
	}

	return nil
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q", dir)
	return cmd.Run()
}

func initGoModule(dir string, gitPrefix string, projectname string) error {
	var moduleName string
	if gitPrefix == "" {
		moduleName = projectname
	} else {
		moduleName = fmt.Sprintf("%s/%s", gitPrefix, projectname)
	}
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	return cmd.Run()
}

func copyFiles(src string, dest string) error {
	out, err := exec.Command("sh", "-c", fmt.Sprintf("cp -R %s %s", path.Join(src, "*"), dest)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to copy files %s %v", out, err)
	}
	return nil
}

func applyProjectToTemplates(p *Project, path string) error {
	filesToRemove := make([]string, 0, 512)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Processing file %s\n", info.Name())
		if !strings.Contains(info.Name(), "-tpl") {
			fmt.Printf("Skipping file %s\n", info.Name())
			return nil
		}

		outputFilename := strings.ReplaceAll(info.Name(), "-tpl", "")
		outputPath := strings.ReplaceAll(path, info.Name(), outputFilename)

		fmt.Printf("Reading file %s\n", info.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		fmt.Printf("Executing template on %s\n", info.Name())
		tpl, err := template.New(outputFilename).Parse(string(content))
		if err != nil {
			return err
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		fmt.Printf("Writing output file %s\n", info.Name())
		err = tpl.Execute(f, p)
		if err != nil {
			return err
		}

		filesToRemove = append(filesToRemove, path)
		return nil
	})

	for _, f := range filesToRemove {
		os.Remove(f)
	}

	return err
}
