package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"charm.land/huh/v2"
)

const githubPath = "/home/cody/workspaces/github/CodyBense/"
const projectPath = "/home/cody/workspaces/projects/"

// createDir ...
func createDir(targetDir string) {
	err := os.Mkdir(targetDir, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// copyGitignore ...
func copyGitignore(path, projectName, language string) {
	if err := exec.Command(
		"cp",
		fmt.Sprintf("/home/cody/nix-config/dotfiles/git/%s.gitignore", language),
		fmt.Sprintf("%s%s/.gitignore", path, projectName)).Run(); err != nil {
		log.Fatalln(err)
		return
	}
}

// createDevenv ...
func createDevenv() {
	if err := exec.Command("devenv", "init").Run(); err != nil {
		log.Fatalln(err)
		return
	}
	if err := exec.Command("devenv", "allow").Run(); err != nil {
		log.Fatalln(err)
		return
	}
}

// createProject ...
func createProject(projectName, language, projectPath string) {
	var path string
	if projectPath == "github" {
		path = githubPath
	} else {
		path = projectPath
	}

	switch language {
	case "go":
		fmt.Println("Creating a Go project")
		createDir(path + projectName)
		os.Chdir(path + projectName)
		_, err := exec.Command("go", "mod", "init", fmt.Sprintf("github.com/CodyBense/%s", projectName)).Output()
		if err != nil {
			log.Fatalln(err)
			return
		}
		_, err = exec.Command("git", "init").Output()
		if err != nil {
			log.Fatalln(err)
			return
		}
		copyGitignore(path, projectName, language)
		createDevenv()
	case "rust":
		fmt.Println("Creating a Rust project")
		if err := exec.Command("cargo", fmt.Sprintf("%s%s", path, projectName)).Run(); err != nil {
			log.Fatalln(err)
			return
		}
		os.Chdir(path + projectName)
		createDevenv()
	case "zig":
		fmt.Println("Creating a Zig project")
		createDir(path + projectName)
		os.Chdir(path + projectName)
		if err := exec.Command("zig", "init"); err != nil {
			log.Fatalln(err)
			return
		}
		copyGitignore(path, projectName, language)
	case "python":
		fmt.Println("Creating a Python project")
		if err := exec.Command("uv", "init", fmt.Sprintf("%s%s", path, projectName)); err != nil {
			log.Fatalln(err)
			return
		}
	case "bash":
		fmt.Println("Creating a Bash proejct")
		createDir(path + projectName)
		os.Chdir(path + projectName)
		if err := exec.Command("git", "init"); err != nil {
			log.Fatalln(err)
			return
		}
		copyGitignore(path, projectName, language)
	case "c":
		fmt.Println("Creating a C project")
		createDir(path + projectName)
		os.Chdir(path + projectName)
		if err := exec.Command("git", "init"); err != nil {
			log.Fatalln(err)
			return
		}
		copyGitignore(path, projectName, language)
	}
}

func main() {
	var (
		projectName string
		language    string
		projectPath string
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of your project!").
				Prompt("> ").
				Value(&projectName),
		),

		huh.NewGroup(
			huh.NewSelect[string]().Title("Select a language").Options(
				huh.NewOption("go", "go"),
				huh.NewOption("rust", "rust"),
				huh.NewOption("zig", "zig"),
				huh.NewOption("python", "python"),
				huh.NewOption("bash", "bash"),
				huh.NewOption("c", "c"),
			).
				Value(&language),
		),

		huh.NewGroup(
			huh.NewSelect[string]().Title("Select a directory").
				Options(
					huh.NewOption("github", "github"),
					huh.NewOption("project", "project"),
				).
				Value(&projectPath),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}

	createProject(projectName, language, projectPath)
}
