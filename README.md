# trig

This command-line tool extends [CodeKit](https://codekitapp.com) by remembering associations between triggering files (template files) and dependent files (which are including those templates). After a file was modified, you call trig, e.g., `trig handle footer.html` (this call typically happens in a CodeKit hook). Since `trig` knows exatcly which files include `footer.html` (this should be automated) it asks CodeKit via AppleScript to process each of those dependent files. (Associations are stored in `~/.trig/associations.db`.)

## Usage

    trig <command> [<args>...]
    trig set <dependent-file> <triggering-file-1> <triggering-file-2>...
    trig set <dependent-file>
    trig handle <triggering-file>
    trig list
    trig help

Commands:

* **set:** Associates a dependent file with triggering filess. If a currently associated triggering file is not mentioned, that association is removed. If no triggering files are provided then all are removed.
* **list:** List associations.
* **handle:** Executes the registered command on all dependent files that were associated with the provided triggering file.
* **help:** Prints the help screen.

Examples:

    $ trig register .go 'go run $1'

    $ trig set www/index.go templates/one.gohtml templates/two.gohtml

    $ trig set www/index.go

    $ trig handle templates/two.gohtml

    $ trig list

## Install

```bash
$ go get github.com/zengabor/trig
```

## Author

[Gabor Lenard](https://github.com/zengabor)
