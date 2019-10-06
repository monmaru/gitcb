package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
)

var s []prompt.Suggest

func completer(in prompt.Document) []prompt.Suggest {
	return prompt.FilterContains(s, in.GetWordBeforeCursor(), true)
}

func main() {
	out, err := runCommand(exec.Command("git", "branch", "-a"))
	exitIfError(err)

	lines := strings.Split(out, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Index(line, "*") == 0 {
			s = append(s, prompt.Suggest{Text: line[2:len(line)], Description: "*"})
		} else {
			s = append(s, prompt.Suggest{Text: line})
		}
	}

	in := prompt.Input("branch: ", completer,
		prompt.OptionTitle("git checkout"),
		prompt.OptionPrefixTextColor(prompt.Blue))
	r := regexp.MustCompile(`remotes/(.*)/(.*)`)
	result := r.FindAllStringSubmatch(in, -1)

	if len(result) == 0 {
		checkout(in, func() {
			out, err := runCommand(exec.Command("git", "checkout", in))
			exitIfError(err)
			fmt.Println(out)
		})
	} else {
		checkout(result[0][2], func() {
			out, err := runCommand(exec.Command("git", "checkout", "-b", result[0][2], result[0][1]+"/"+result[0][2]))
			exitIfError(err)
			fmt.Println(out)
		})
	}
}

func checkout(branch string, fn func()) {
	if branch != currentBranch() {
		isDirty := isWorkingTreeDirty()
		if !isDirty {
			exec.Command("git", "stash", "save")
		}
		fn()
		if !isDirty {
			exec.Command("git", "stash", "pop")
		}
	}
}

func currentBranch() string {
	out, _ := exec.Command("git", "branch").Output()
	lines := strings.Split(string(out), "\n")
	var s string = ""
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Index(line, "*") == 0 {
			s = line[2:len(line)]
		}
	}
	return s
}

func isWorkingTreeDirty() bool {
	out, _ := exec.Command("git", "status").Output()
	return strings.Index(string(out), "nothing to commit, working tree clean") == -1
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
