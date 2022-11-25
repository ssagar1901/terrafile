package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/rendon/testcli"
	"github.com/stretchr/testify/assert"
)

var terrafileBinaryPath string
var workingDirectory string

func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	terrafileBinaryPath = workingDirectory + "/terrafile"
}

func TestTerrafile(t *testing.T) {
	tests := []*struct {
		Description      string
		TerrafileCreator terrafileCreator
		ExpectedModules  []string
	}{
		{
			Description:      "Segment Terrafile Format",
			TerrafileCreator: createSegmentTerrafile,
			ExpectedModules: []string{
				"terraform-aws-modules/terraform-aws-vpc/master",
				"terraform-aws-modules/terraform-aws-vpc/v1.46.0",
			},
		},
		{
			Description:      "Community Terrafile Format",
			TerrafileCreator: createCommunityTerrafile,
			ExpectedModules: []string{
				"terraform-aws-vpc",
				"terraform-aws-vpc-1.46.0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			folder, back := setup(t, test.TerrafileCreator)
			defer back()
			defer func() {
				assert.NoError(t, os.RemoveAll(path.Join(workingDirectory, "./.terrafile")))
			}()

			testcli.Run(terrafileBinaryPath, "-d", "-f", fmt.Sprint(folder, "/Terrafile"))

			if !testcli.Success() {
				t.Fatalf("Expected to succeed, but failed: %q \nStdout: %q \nStderr: %q", testcli.Error(), testcli.Stdout(), testcli.Stderr())
			}
			// Assert output
			for _, output := range []string{
				"Cloning   git@github.com:terraform-aws-modules/terraform-aws-vpc",
				"Vendoring ref master",
				"Vendoring ref v1.46.0",
			} {
				assert.Contains(t, testcli.Stdout(), output)
			}
			// Assert files exist
			for _, moduleName := range test.ExpectedModules {
				assert.DirExists(t, path.Join(workingDirectory, "./.terrafile", moduleName))
			}
		})
	}
}

func setup(t *testing.T, createTerrafile terrafileCreator) (current string, back func()) {
	folder, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	createTerrafile(t, folder)
	return folder, func() {
		assert.NoError(t, os.RemoveAll(folder))
	}
}

func createFile(t *testing.T, filename string, contents string) {
	assert.NoError(t, ioutil.WriteFile(filename, []byte(contents), 0644))
}

type terrafileCreator func(t *testing.T, folder string)

func createSegmentTerrafile(t *testing.T, folder string) {
	var yaml = `git@github.com:terraform-aws-modules/terraform-aws-vpc:
  - v1.46.0
  - master
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}

func createCommunityTerrafile(t *testing.T, folder string) {
	var yaml = `terraform-aws-vpc:
  source: "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: master
terraform-aws-vpc-1.46.0:
  source: "git@github.com:terraform-aws-modules/terraform-aws-vpc"
  version: v1.46.0
`
	createFile(t, path.Join(folder, "Terrafile"), yaml)
}
