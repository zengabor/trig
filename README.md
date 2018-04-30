# gohtml

Command-line tool to set associations between (go) files and templates, so when you invoke it to handle a template it will `touch` all associated files. Consequently the build process will process those (go) files. (Associations are stored in the home directory, in `.gohtml/gohtml.db`.)

## Usage

    gohtml <command> <args>

Available commands:

* **set:** Associates a file with templates. Obsolete associations are removed.
* **handle:** Touches the all files that are associated with the provided template.
* **help:** Prints the help screen.

Examples:

    gothml set path/index.go one.gohtml b/two.gohtml

    gohtml handle b/two.gohtml

## Install

To install, use `go get`:

```bash
$ go get -d github.com/zengabor/gohtml
```

## Author

[Gabor Lenard](https://github.com/zengabor)
