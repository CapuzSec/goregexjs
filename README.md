This is a command line tool for searching regular expressions in a list of JS URLs. It takes a list of URLs as input from standard input and looks for regular expressions in your JS content.

## Install
```
go build goregexjs.go
```

## Usage:

```
cat js-list.txt | goregexjs --regex-file regex.txt --show-chars 60 --threads 30
```

```
Usage of regexsearch:
  -regex-file string
        path to the file containing the regular expressions
  -show-chars int
        number of characters to show from the matched word (default 10)
  -threads int
        number of threads to use (default 10)
```


The tool expects a list of URLs to search in to be piped into it from the command line:

```
cat urls.txt | regexsearch -regex-file regexes.txt -show-chars 15 -threads 20
```




