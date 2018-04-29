# gohtml

Command-line tool to maintain associations between go files and templates, so when you invoke it with the path to a template it will `touch` all associated go files so that the build process can pick them up. (Associations are stored in the home directory, in `.gohtml/gohtml.db`.)

## Usage

    gohtml <command> <args>

Available commands:

* **set** Associates a go file with templates. Obsolete associations are removed.
* **run** Runs the go files that are associated with the template.
* **help** Prints the help screen.

Examples:

    gothml set path/index.go one.gohtml b/two.gohtml

    gohtml run b/two.gohtml

## Install

To install, use `go get`:

```bash
$ go get -d github.com/zengabor/gohtml
```

## Author

[Gabor Lenard](https://github.com/zengabor)
