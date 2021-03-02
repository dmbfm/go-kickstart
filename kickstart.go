package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func randName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func genRandomProjectName() string {
	var name string
	for {
		name = randName(10)
		f, err := os.Open(name)
		if err != nil {
			f.Close()
			break
		}
	}

	return name
}

func main() {
	var err error
	var cmd *exec.Cmd
	var name string

	if len(os.Args) < 2 {
		fmt.Println("Not enough args!")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		name = genRandomProjectName()
	} else {
		name = os.Args[2]
	}

	var template = os.Args[1]

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
