# Command Line Interfaces

All command line interfaces should be placed in here.

For example, if we need to test specific items we may create `/cmd/testingelement/name.go` and package name needs to
be `main`

From the root of the folder, running `make` should compile `cmd/main/*.go` to `bin/main` executable.

For different command line interfaces we may have to create different target labels.

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

## Links

## Notes

## Contributors

## Issues for Future Reference
