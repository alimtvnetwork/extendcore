<!-- ![Use Package logo](UseLogo) -->

# `ExtendCore` Intro

`extendcore` comes after core, enum, errorwrapper. It will core package but with errorwrapper included.

## Git Clone

`git clone https://gitlab.com/evatix-go/extendcore.git`

### 2FA enabled, for linux

`git clone https://[YourGitLabUserName]:[YourGitlabAcessTokenGenerateFromGitlabsTokens]@gitlab.com/evatix-go/extendcore.git`

### Prerequisites

- Update git to latest 2.29
- Update or install the latest of Go 1.15.2
- Either add your ssh key to your gitlab account
- Or, use your access token to clone it.

## Installation

`go get gitlab.com/evatix-go/extendcore`

### Go get issue for private package

- Update git to 2.29
- Enable go modules. (Windows : `go env -w GO111MODULE=on`, Unix : `export GO111MODULE=on`)
- Add `gitlab.com/evatix-go` to go env private

To set for Windows:

`go env -w GOPRIVATE=[AddExistingOnes;]gitlab.com/evatix-go`

To set for Unix:

`expoort GOPRIVATE=[AddExistingOnes;]gitlab.com/evatix-go`

To go get as root:

- [Linux, Go Get Fix with SSH GitLabs](https://gitlab.com/evatix-go/os-manuals/-/issues/43)

## Why extendcore?

Provides and reduces code lines for dealing with common DRYness of the codebase.


## Folder Structure

```
.
├── cmd
│   ├── main
│       └──── main.go
│   ├── server // optional - server cmd cli
│       └──── main.go
│   ├── client // optional - client cmd cli
│       └──── main.go
├── assets
│   ├── ...contains assets
├── configs
│   ├── ... contains common config
├── .gitattributes
├── .gitignore
└── Readme.md
```

## Examples

`Code Smaples`

## Acknowledgement

Any other packages used

## Links

- [Directory Structure Sampling](https://stackoverflow.com/questions/19699059/representing-directory-file-structure-in-markdown-syntax)

## Issues

- [Create your issues](https://gitlab.com/evatix-go/extendcore/-/issues)

## Notes

## Contributors

## License

[Evatix MIT License](/LICENSE)
