package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jjunqueira/goproject/pkg/goproject"
)

// Project contains information regarding the project we are generating
type Project struct {
	GitPrefix  string
	Name       string
	ModuleName string
	Tpl        *Template
}

// NewProject constructs a new Project struct with the provided settings
func NewProject(c *goproject.Config, gitPrefix string, tplName string, projectName string) (*Project, error) {
	p := new(Project)
	p.GitPrefix = gitPrefix
	p.Name = projectName

	if p.GitPrefix == "" {
		p.ModuleName = p.Name
	} else {
		p.ModuleName = path.Join(p.GitPrefix, p.Name)
	}

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
	tplPath := path.Join(c.TemplatesPath, tplName)

	_, err := os.Stat(tplPath)
	if err == nil {
		return &Template{name: tplName, path: tplPath}, err
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

	err = os.Mkdir(fullPath, 0750)
	if err != nil {
		return fmt.Errorf("unable to create project directory: %v", err)
	}

	err = copyFiles(p.Tpl.path, fullPath)
	if err != nil {
		return fmt.Errorf("unable to copy template files: %v", err)
	}

	err = applyProjectToTemplates(p, fullPath)
	if err != nil {
		return fmt.Errorf("unable to execute templates: %v", err)
	}

	err = fixCmdProjectFolderName(p, fullPath)
	if err != nil {
		return fmt.Errorf("unable to rename cmd project folder: %v", err)
	}

	err = gitCleanup(fullPath)
	if err != nil {
		return fmt.Errorf("unable to initialize git: %v", err)
	}

	return nil
}

func gitCleanup(dir string) error {
	gitInit := exec.Command("git", "init", "-q")
	gitInit.Dir = dir

	gitAdd := exec.Command("git", "add", "--all")
	gitAdd.Dir = dir

	gitCommit := exec.Command("git", "commit", "-am", "Initial commit")
	gitCommit.Dir = dir

	err := gitInit.Run()
	if err != nil {
		return err
	}

	err = gitAdd.Run()
	if err != nil {
		return err
	}

	err = gitCommit.Run()
	if err != nil {
		return err
	}

	return nil
}

func copyFiles(src string, dest string) error {
	allSources := path.Join(src, "*")
	copyCommand := fmt.Sprintf("cp -R %s %s", allSources, dest)

	out, err := exec.Command("sh", "-c", copyCommand).CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to copy files %s %v", out, err)
	}

	return nil
}

func applyProjectToTemplates(p *Project, path string) error {
	filesToRemove := make([]string, 0, 512)
	walkErr := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Don't keep walking if an error was encountered
		if err != nil {
			return err
		}

		if info.Name() == "gitignore" {
			return os.Rename(path, strings.ReplaceAll(path, "gitignore", ".gitignore"))
		}

		if info.Name() == "gitkeep" {
			return os.Rename(path, strings.ReplaceAll(path, "gitkeep", ".gitkeep"))
		}

		if info.Name() == "golangci.yml" {
			return os.Rename(path, strings.ReplaceAll(path, "golangci.yml", ".golangci.yml"))
		}

		if !strings.Contains(info.Name(), "-tpl") {
			return nil
		}

		outputFilename := strings.ReplaceAll(info.Name(), "-tpl", "")
		outputFilename = strings.ReplaceAll(outputFilename, strings.ToLower(p.Tpl.name), strings.ToLower(p.Name))
		outputPath := strings.ReplaceAll(path, info.Name(), outputFilename)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		tpl, err := template.New(outputFilename).Parse(string(content))
		if err != nil {
			return err
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}

		err = tpl.Execute(f, p)
		if err != nil {
			return err
		}

		filesToRemove = append(filesToRemove, path)
		return nil
	})

	for _, f := range filesToRemove {
		err := os.Remove(f)
		if err != nil {
			fmt.Printf("unable to remove file %s %v", f, err)
		}
	}

	return walkErr
}

func fixCmdProjectFolderName(p *Project, fullpath string) error {
	oldpath := path.Join(fullpath, "cmd", strings.ToLower(p.Tpl.name))
	newpath := path.Join(fullpath, "cmd", strings.ToLower(p.Name))

	return os.Rename(oldpath, newpath)
}
