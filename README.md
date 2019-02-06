# cwt

A linux morse code trainer. Close to being finished but still a beta version.

## History

I originally wrote this with sciter and it's tiscript. The tiscript worked well with the critical key timing.

This version is written with a framework I generated with kickwasm. So the critical timing code in this version is written with go compiled to wasm. This will be a test to see if go wasm can keep up with the keying.

I was going to wait to put this up here at git hub but I just almost lost my hard drive so just to be safe it's up here.

The app is pretty much done but right now I'm working on the most critical part; recording the user's keying of morse code.

## The key

My key is wired to the contacts of the left button of a cheap usb mouse I ripped apart. So I record the times of the mouse up and down events in the renderer and let the main process copy those times into valid morse code and then into text.

## Build Requirements

### Alsa

So this uses alsa and requires alsa to build and maybe to run.

``` bash

sudo apt-get install alsa-base libasound2-dev

```

### Imports

* [the boltdb package.](https://github.com/boltdb/bolt)
* [the yaml package.](https://gopkg.in/yaml.v2)
* [the gorilla websocket package.](https://github.com/gorilla/websocket)

## The build

### The application code is physically and logically organized into 4 main levels

1. The **domain/** folder contains domain ( application ) level logic.
1. The **mainprocess/** folder contains the main process level logic.
1. The **renderer/** folder contains the renderer level logic.
1. The **site/** folder contains the wasm, templates, styles etc for the browser.

### The application has 2 processes

1. The **main process** is a web server running through whatever port you indicate in your application's http.yaml file. When you start the main process it opens a browser which loads and runs the renderer process from the **site/** folder.
1. The **renderer process** is all of the wasm, html, css, images, etc contained in the site/ folder.

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
