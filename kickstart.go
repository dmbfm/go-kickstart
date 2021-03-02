package main

import (
	"fmt"
	"os"
	"os/exec"
)

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	var err error
	var cmd *exec.Cmd

	if len(os.Args) < 3 {
		fmt.Println("Not enough args!")
		os.Exit(1)
	}

	var template = os.Args[1]
	var name = os.Args[2]

	fmt.Printf("Creating project '%v' from template '%v'...\n", name, template)

	err = os.Mkdir(name, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Chdir(name)
	exitIfError(err)

	var gitPath = "git@github.com:dmbfm/" + template

	cmd = exec.Command("git", "clone", gitPath, ".")
	err = cmd.Run()
	exitIfError(err)

	err = os.RemoveAll("./.git")
	exitIfError(err)

	_, err = os.Open("package.json")
	// TODO: Maybe, instead of this, each project can have a "kickstart.go" file that we
	//       run here by calling exec.Command("go", "run", "kickstart.go")
	if err == nil {
		fmt.Println("Installing project dependencies...")
		cmd = exec.Command("npm", "install")
		err = cmd.Run()
		exitIfError(err)

		fmt.Println("Starting vscode...")
		code := exec.Command("code", ".")
		code.Start()

		fmt.Println("Running 'npm start'...")
		cmd = exec.Command("npm", "start")
		err = cmd.Run()
		exitIfError(err)
	}

	fmt.Println("Done!")
}
