# Go Git Archive

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
