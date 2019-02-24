# cwt

A **linux** morse code trainer. Close to being finished but still a beta version.

 The application still has some testing junk in some of it's output.

## The morse code key

My morse code key is wired to the contacts of the left button of a cheap usb mouse I ripped apart. So...

* the morse code keying is sensed as "mouseup" and "mousedown" events in the browser.
* the critical timing code is written with go compiled to wasm.

So this will be a test to see if go wasm can keep up with the keying.

## To install and build on linux

* The mainprocess/goalsa package is written in cgo which requires gcc.
* The mainprocess/goalsa package uses the alsa lib and requires libasound2-dev.

``` bash

sudo apt install gcc
sudo apt install libasound2-dev
go get -u github.com/josephbudd/cwt

```

**go get -u github.com/josephbudd/cwt** will import

* [the boltdb package.](https://github.com/boltdb/bolt)
* [the yaml package.](https://gopkg.in/yaml.v2)
* [the gorilla websocket package.](https://github.com/gorilla/websocket)

## The build

### The application code is physically and logically organized into 4 areas

1. The **domain/** folder contains domain ( shared ) logic.
1. The **mainprocess/** folder contains the main process logic.
1. The **renderer/** folder contains the renderer logic. The code in the **renderer/** folder is compiled into wasm at **site/app.wasm**.
1. The **site/** folder contains the wasm, templates, styles etc for the browser.

### The application has 2 processes

1. The **main process** is a web server running through whatever port you indicate in your application's http.yaml file. When you start the main process it opens a browser which loads and runs the renderer process from the **site/** folder.
1. The **renderer process** is all of the wasm, html, css, images, etc contained in the **site/** folder.

### The application has a 2 step build

So when you build the application, you build both the renderer process and the main process. The renderer process code is in the **renderer/** folder but it is built into the **site/** folder.

There is a shell script in the **renderer/** folder that builds the renderer process into the **site/** folder. It's **build.sh**.

``` bash

cd $GOPATH
cd src/github.com/josephbudd/cwt
cd renderer
./build.sh
cd ..
go build
./cwt

```

### Development issues

#### I can't key

The morse code key is like a musical instrument that plays only one note. Morse code is like a song that has

* a never changing beat,
* notes that must only last 1 or 3 beats,
* pauses that must only last 1 or 3 beats.

I can't do any of that properly so I added a metronome. I'm hoping it will help me learn to key. Right now the metronome is automatic but I will be making changes with the metronome.

### CGO issues

I think I fixed my CGO issue. It was revealed with my metronome. Time will tell.