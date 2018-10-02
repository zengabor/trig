# trig

This command-line tool extends [CodeKit](https://codekitapp.com). It can set associations between template files (triggering files) and dependent files (the ones using those template files). For example, if you update a template `footer.html` and call `trig` with the file path (`trig handle footer.html`) then `trig` will call CodeKit via AppleScript to process all the files that are including `footer.html`. This way you can be sure that everything is updated when you modify a template. (Associations between the files are stored in `~/.trig/associations.db`.)

## Usage

    trig <command> [<args>...]
    trig set <dependent-file> <triggering-file-1> <triggering-file-2>...
    trig set <dependent-file>
    trig list
    trig handle <triggering-file>
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

    $ trig list

    $ trig handle templates/two.gohtml

## Install

```bash
$ go get github.com/zengabor/trig
```

## Author

[Gabor Lenard](https://github.com/zengabor)
