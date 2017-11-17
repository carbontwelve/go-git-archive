# Go Git Archive [![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=Helpful%20command%20line%20tool%20for%20zipping%20up%20changes%20between%20two%20commits%20in%20git.&url=https://github.com/photogabble/go-git-archive&via=github&hashtags=golang,developers)

[![Go Report Card](https://goreportcard.com/badge/github.com/photogabble/go-git-archive)](https://goreportcard.com/report/github.com/photogabble/go-git-archive)

A small command line tool for zipping all files changed between two git commits.

## Usage
By default `-last` will be the current `HEAD` within your repository and therefore only the `-first` value is required.

```
Usage of git-archive.exe:
  -first string
        The git commit that we are to begin at.
  -last string
        The git commit that we are to end at. (default "...")
  -list
        List files rather than write to zip.
  -v    Toggle verbose output.
```

## Reasoning
I have been slowly learning Go as a hobby for the past few months, usually in blocks of an hour a month and have been working on small projects that I can easily get finished. This is one such project and it solves a problem I have at work where we often need to upload just the files changed between commit versions as a zip to an FTP endpoint.

## License

Distributed under MIT License, please see [license](LICENSE) file in code for more details.
