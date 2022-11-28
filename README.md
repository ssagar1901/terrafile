# Terrafile

Terrafile is a binary written in Go to systematically manage external modules from Github for use in Terraform. This is a fork of the original version found here: https://github.com/coretech/terrafile

## How to install

### macOS

```sh
brew tap segmentio/packages
brew install segmentio/packages/terrafile
```

### Linux

Download your preferred flavor from the [releases](https://github.com/coretech/terrafile/releases/latest) page and install manually.

## How to use

Terrafile expects a file named `Terrafile` which will contain your terraform module dependencies in a yaml like format.

Terrafile config file in custom directory

```sh
$ terrafile -f config/Terrafile
```

Terraform modules exported to custom directory

```sh
$ terrafile -p /path/to/custom_directory
```

## Segment Terrafile Format

Segment Terrafile is an internal format that specifies sources with a list of its versions. Each version is vendored under `.terrafile/<user>/<repo>/<ref>`

> :warning: the major downside of using this format is that vendored modules has version in their path. So once you want to upgrade module version, you will need to modify both `Terrafile` and all module usages. Check out [Community Terrafile Format](#community-terrafile-format) which doesn't have such downside

An example Terrafile:

```yaml
git@github.com:segmentio/terracode-modules:
  - chamber_v2.0.0
  - iam_v1.0.0
  - master
```

```sh
$ terrafile
[*] Cloning   git@github.com:segmentio/terracode-modules
[*] Vendoring ref chamber_v2.0.0
[*] Vendoring ref iam_v1.0.0
[*] Vendoring ref master
```

```sh
$ ls -A1 .terrafile/segmentio/terracode-modules
chamber_v2.0.0
iam_v1.0.0
master
```

## Community Terrafile Format

Community Terrafile is a community supported format implemented in various languages (e.g. [github: coretech/terrafile](https://github.com/coretech/terrafile), [npm: terrafile](https://www.npmjs.com/package/terrafile)). Each version is vendored under `.terrafile/<alias>`

In comparison with [Segment Terrafile Format](#segment-terrafile-format), module versions are not included in the vendored path. So once you want to upgrade module version, only `Terrafile` needs to be updated and all module usages will be left untouched.

An example Terrafile:

```yaml
terracode-modules-chamber:
  source: "git@github.com:segmentio/terracode-modules"
  version: "chamber_v2.0.0"
terracode-modules-iam:
  source: "git@github.com:segmentio/terracode-modules"
  version: "iam_v1.0.0"
terracode-modules:
  source: "git@github.com:segmentio/terracode-modules"
  version: "master"
```

```sh
$ terrafile
[*] Cloning   git@github.com:segmentio/terracode-modules
[*] Vendoring ref chamber_v2.0.0
[*] Vendoring ref iam_v1.0.0
[*] Vendoring ref master
```

```sh
$ ls -A1 .terrafile
terracode-modules-chamber
terracode-modules-iam
terracode-modules
```
