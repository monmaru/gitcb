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
	in := prompt.Input("branch: ", makeBranchSelector(),
		prompt.OptionTitle("git checkout"),
		prompt.OptionPrefixTextColor(prompt.Blue))

	if len(in) == 0 {
		fmt.Println("Canceled")
		return
	}

	if strings.HasPrefix(in, "remotes/") {
		strs := strings.Split(in, "/")
		branch := strings.Join(strs[2:], "/")
		checkout(branch, exec.Command("git", "checkout", "-b", branch, in))
	} else {
		checkout(in, exec.Command("git", "checkout", in))
	}
}

func makeBranchSelector() func(in prompt.Document) []prompt.Suggest {
	out, err := runCommand(exec.Command("git", "branch", "-a"))
	exitIfError(err)

	var s []prompt.Suggest
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.Index(line, "*") == 0 {
			s = append(s, prompt.Suggest{Text: line[2:len(line)], Description: "*"})
		} else {
			s = append(s, prompt.Suggest{Text: line})
		}
	}

	return func(in prompt.Document) []prompt.Suggest {
		return prompt.FilterContains(s, in.GetWordBeforeCursor(), true)
	}
}

func checkout(branch string, cmd *exec.Cmd) {
	if branch == currentBranch() {
		fmt.Printf("Already in %s\n", branch)
	} else {
		out, err := runCommand(cmd)
		exitIfError(err)
		fmt.Println(out)
	}
}

func currentBranch() string {
	out, _ := exec.Command("git", "branch").Output()
	lines := strings.Split(string(out), "\n")
	var s string = ""
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.Index(line, "*") == 0 {
			s = line[2:len(line)]
		}
	}
	return s
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
