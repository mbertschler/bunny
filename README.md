Bunny project management tool
=============================
![status: not yet usable](https://img.shields.io/badge/status-not%20yet%20usable-red.svg)

Bunny will be a simple and modern project management tool for individuals
and teams. It is not usable yet.

Why another project management tool?
------------------------------------
Bunny is an experimental project with the goal of creating a simple and
efficient project management tool.

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