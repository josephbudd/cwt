# cwt

A **linux** morse code trainer. This is far from being finished.

## Services

### The reference service

The reference service of the application allows the user to select the morse code characters which are to be both copied and keyed by the user. It also shows the user the keying and copying test scores for the characters at each of the various speeds.

### The copy service

The copy service of the application allows the user to practice or test copying morse code at the selected copy speed.

### The key service

The key service of the application allows the user to practice or test keying morse code at the selected key speed. This is the part I am working on at the moment. There are changes in this service as I try to figure out how to teach myself to key.

## The morse code key

My morse code key is wired to the contacts of the left button on a board that I ripped out of a cheap usb mouse. The board gets plugged into a usb port on the lap top. Pressing the key down causes a "mouse-down" event and letting the key up causes a "mouse-up" event. The go wasm code does a fine job of handling those events by appending the time of the event to a slice of time.Times.

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

## Development issues

### I can't key

The morse code key is like a musical instrument that can only play notes at one pitch.

Morse code is like a song that has

* a never changing beat,
* short notes ( dits ) that must only last 1 beat,
* long notes ( dahs ) that must only last 3 beats,
* pauses between the dit and dahs in a character that must only last 1 beat,
* pauses between characters in a word that must only last 3 beats,
* pauses between words that must only last 7 beats.

I can't do any of that properly so I have allowed for variances in user input and I added a metronome. I'm hoping it will help me learn to key.

### CGO issues

My big issue here is package mainprocess/goalsa. I'm getting fatals in my CGO every once in a while with the new metronome. So I need to look at that and fix it.