package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
)

func main() {
	in := prompt.Input(
		"branch: ",
		makeBranchSelector(),
		prompt.OptionTitle("git checkout"),
		prompt.OptionPrefixTextColor(prompt.Blue))

	if len(in) == 0 {
		fmt.Println("Canceled")
		return
	}

	if strings.HasPrefix(in, "remotes/") {
		ss := strings.Split(in, "/")
		branch := strings.Join(ss[2:], "/")
		checkout(branch, exec.Command("git", "checkout", "-b", branch, in))
	} else {
		checkout(in, exec.Command("git", "checkout", in))
	}
}

func makeBranchSelector() func(in prompt.Document) []prompt.Suggest {
	out, err := runCommand(exec.Command("git", "branch", "-a"))
	exitIfError(err)

	var suggests []prompt.Suggest
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.Index(line, "*") == 0 {
			suggests = append(suggests, prompt.Suggest{Text: line[2:], Description: "*"})
		} else {
			suggests = append(suggests, prompt.Suggest{Text: line})
		}
	}

	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterContains(suggests, in.GetWordBeforeCursor(), true)
	}
}

func checkout(branch string, cmd *exec.Cmd) {
	if branch == currentBranch() {
		fmt.Printf("Already in %s\n", branch)
		return
	}

	out, err := runCommand(cmd)
	exitIfError(err)
	fmt.Println(out)
}

func currentBranch() string {
	out, _ := exec.Command("git", "branch").Output()
	lines := strings.Split(string(out), "\n")
	var branch string = ""
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.Index(line, "*") == 0 {
			branch = line[2:]
		}
	}
	return branch
}

func runCommand(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errStr := stderr.String()
		if len(errStr) == 0 {
			return "", err
		}
		return "", fmt.Errorf(errStr)
	}

	return stdout.String(), nil
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
