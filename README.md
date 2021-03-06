Bunny work management tool
==========================
[![Build Status](https://travis-ci.org/mbertschler/bunny.svg?branch=master)](https://travis-ci.org/mbertschler/bunny)
[![GoDoc](https://godoc.org/github.com/mbertschler/bunny?status.svg)](https://godoc.org/github.com/mbertschler/bunny)
![status: not yet usable](https://img.shields.io/badge/status-not%20yet%20usable-red.svg)
[![GoDoc](https://goreportcard.com/badge/github.com/mbertschler/bunny)](https://goreportcard.com/report/github.com/mbertschler/bunny)

Bunny will be a simple and modern work management tool for individuals
and teams. It is not usable yet.

Why another work management tool?
---------------------------------
Bunny is an experimental project with the goal of creating a simple and
efficient work management tool.

Building
--------

#### Requirements
- Go and dep
- Node.js and yarn

```bash
# get JS and CSS dependencies
cd js
yarn install

# build bunny
cd ..
dep ensure
go install github.com/mbertschler/bunny
```

### Build the Docker image

Yarn and dep have to be run on the host before starting the Docker
build process. 

```bash
docker build -t mbertschler/bunny:alpha-1 .
```

Running
-------

### Using Docker

```bash
docker run -p 3080:3080 mbertschler/bunny:alpha-1
```

License
-------
Bunny is released under the Apache 2.0 license. See [LICENSE](LICENSE).

--------------

Created by  Martin Bertschler `@mbertschler` in 2018.
