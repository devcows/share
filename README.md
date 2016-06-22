## Share
[![Build Status](https://travis-ci.org/devcows/share.svg?branch=master)](https://travis-ci.org/devcows/share)
[![Coverage Status](https://coveralls.io/repos/github/devcows/share/badge.svg?branch=master)](https://coveralls.io/github/devcows/share?branch=master)
[![Go Report](https://goreportcard.com/badge/github.com/devcows/share)](https://goreportcard.com/report/github.com/devcows/share)
[![License](http://img.shields.io/:license-Apache_v2-blue.svg)](https://raw.githubusercontent.com/devcows/share/master/LICENSE)

Share is a command line for sharing files by network. The goal is to share a file easy as **human-friendly** as possible. It provides an independent webserver for each file.

`TODO png or gif usage`

Share is written in GO.

### Installation

#### ‎Mac OS
```
$ brew tap devcows/apps
$ brew install share
```

#### Other

To install Share command line go to [relases](https://github.com/devcows/share/releases) and download package for your operating system.

### Usage

Launch server:
```
$ share server
```

Hello world CLI:
```
$ share add --file path_to_file
```

Synopsis:

```
$ share [command] [flags] FILE_PATH
```
See also `share --help`.
