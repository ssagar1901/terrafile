package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/nritholtz/stdemuxerhook"
	log "github.com/sirupsen/logrus"
)

var opts struct {
	ModulePath string `short:"p" long:"module_path" default:"./.terrafile" description:"File path to install generated terraform modules"`

	TerrafilePath string `short:"f" long:"terrafile_file" default:"./Terrafile" description:"File path to the Terrafile file"`
	Debug         bool   `short:"d" long:"debug"`
}

// To be set by goreleaser on build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var tempDir string

func init() {
	// Needed to redirect logrus to proper stream STDOUT vs STDERR
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))
	var err error
	tempDir, err = ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}
}

func gitClone(repositoryPath string) string {
	pathParts := strings.Split(repositoryPath, ":")
	repositoryName := pathParts[1]

	repoPath := fmt.Sprintf("%s/%s", tempDir, repositoryName)
	cmd := exec.Command("git", "clone", repositoryPath, repoPath)
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return repoPath
}

func gitCheckoutRef(repositoryPath string, ref string, destinationDir string) {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = repositoryPath
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	destWithSlash := fmt.Sprintf("%s/", destinationDir)
	cmd = exec.Command("git", "checkout-index", "--prefix", destWithSlash, "-a")
	cmd.Dir = repositoryPath
	if opts.Debug {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func main() {
	// fmt.Printf("Terrafile: version %v, commit %v, built at %v \n", version, commit, date)
	_, err := flags.Parse(&opts)

	// Invalid choice
	if err != nil {
		panic("invalid arguments")
	}

	// Read File
	yamlFile, err := ioutil.ReadFile(opts.TerrafilePath)
	if err != nil {
		panic(err)
	}

	// Parse Terrafile
	sourceDependenciesMap, err := parseTerrafile(yamlFile)
	if err != nil {
		panic(err)
	}

	// Cleanup module path
	if err := os.RemoveAll(opts.ModulePath); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(opts.ModulePath, os.ModePerm); err != nil {
		panic(err)
	}

	for source, dependencies := range sourceDependenciesMap {
		fmt.Printf("[*] Cloning   %s\n", source)
		repo := gitClone(source)
		for _, dependency := range dependencies {
			fmt.Printf("[*] Vendoring ref %s\n", dependency.Version)
			targetPath, err := dependency.GetTargetPath(opts.ModulePath)
			if err != nil {
				panic(err)
			}

			gitCheckoutRef(repo, dependency.Version, targetPath)
		}
	}
}
