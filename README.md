# gohtml

Command-line tool to set associations between (go) files and templates, so when it handles a template it will "touch" all associated files (currently this is implemented by moving the file to a temporary directory, then moving it back after 2 seconds). Consequently a build tool (like [CodeKit](https://codekitapp.com)) can react and process those files. (Associations are stored in `~/.gohtml/gohtml.db`.)

## Usage

    gohtml <command> [<args>...]

Available commands:

* **set:** Associates a file with templates. (If a currently associated template is not mentioned, the association is removed.)
* **unset:** Removes all associations of a file or template.
* **list:** List associations.
* **handle:** "Touches" all files that are associated with the provided template.
* **help:** Prints the help screen.

Examples:

    $ gohtml set path/index.go one.gohtml b/two.gohtml

    $ gohtml list

    $ gohtml handle b/two.gohtml

    $ gohtml unset one.gohtml
 
    $ gohtml unset path/index.go

## Install

To install, use `go get`:

```bash
$ go get github.com/zengabor/gohtml
```

## Author

[Gabor Lenard](https://github.com/zengabor)
