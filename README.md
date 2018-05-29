# trig

Command-line tool to set associations between a dependent file and its triggering files. For example, if a file includes other files, it is depending on those files, and when those files are changed the dependent files is triggered, so when it handles a template it will "touch" all associated files (currently this is implemented by moving the file to a temporary directory, then moving it back after 2 seconds). Consequently a build tool (like [CodeKit](https://codekitapp.com)) can react and process those files. (Associations are stored in `~/.trig/associations.db`.)

## Usage

    trig <command> [<args>...]
    trig register -t <file-suffix> <commands>
    trig set <dependent-file> <triggering-file-1> <triggering-file-2>...
    trig list
    trig handle <triggering-file>
    trig unset <file>
    trig help

Commands:

* **register**: Registers a command for a file suffix (typically the file extension). The string `$1` in the command is replaced with the file name.
* **set:** Associates a dependent file with the triggering files. (If a currently associated triggering file is not mentioned, the association is removed.)
* **unset:** Removes all associations of a dependent or triggering file.
* **list:** List associations.
* **handle:** Executes the registered command on all dependent files that were associated with the provided triggering file.
* **help:** Prints the help screen.

Examples:

    $ trig register .go 'go run $1'

    $ trig set www/index.go templates/one.gohtml templates/two.gohtml

    $ trig list

    $ trig handle templates/two.gohtml

    $ trig unset templates/one.gohtml

    $ trig unset www/index.go

## Install

```bash
$ go get github.com/zengabor/trig
```

## Author

[Gabor Lenard](https://github.com/zengabor)
