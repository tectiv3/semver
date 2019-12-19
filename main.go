package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if out, err := exec.Command("git", "pull", "--rebase").CombinedOutput(); err != nil {
		log.Fatal(string(out))
	}

	out, err := exec.Command("git", "describe", "--abbrev=0").CombinedOutput()
	if err != nil {
		log.Fatal(string(out))
	}
	describe := strings.Trim(string(out), "\n ")

	str := strings.ReplaceAll(describe, "v", "")
	old_version := strings.Split(str, ".")
	if len(old_version) < 2 {
		log.Fatalf("Invalid version: %q", old_version)
	}

	var major, minor, patch int
	if len(old_version) == 2 {
		major, minor = toInt(old_version[0]), toInt(old_version[1])
		patch = 0
	} else {
		major, minor, patch = toInt(old_version[0]), toInt(old_version[1]), toInt(old_version[2])
	}

	shortlog, err := exec.Command("git", "shortlog", "--no-merges", string(describe)+"..HEAD").CombinedOutput()
	if err != nil {
		log.Fatal(string(shortlog))
	}

	bump := "minor"
	if len(os.Args) > 1 {
		bump = os.Args[1]
	}

	new_version := ""

	switch bump {
	case "patch":
		patch += 1
	case "minor":
		minor += 1
		patch = 0
	case "major":
		major += 1
		minor = 0
		patch = 0
	default:
		new_version = bump
	}

	if len(new_version) == 0 {
		log.Println("Bumped " + bump + " version.")
		new_version = "v" + toStr(major) + "." + toStr(minor) + "." + toStr(patch)
	}

	if out, err := exec.Command("git", "flow", "release", "start", new_version).CombinedOutput(); err != nil {
		log.Fatal(string(out))
	}

	log.Println(new_version)

	tmpfile, err := ioutil.TempFile("", "shortlog")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(shortlog); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	if out, err := exec.Command("git", "flow", "release", "finish", "-f", tmpfile.Name()).CombinedOutput(); err != nil {
		log.Fatal(string(out))
	} else {
		log.Println(string(out))
	}
}

func toInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func toStr(i int) string {
	return strconv.Itoa(i)
}
