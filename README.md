# cwt

## A linux Continuous Wave ( Morse Code ) Trainer

Sept 13, 2019

* Rebuilt with [kickwasm](https://github.com/josephbudd/kickwasm) version 8.2.3.
* Added metronome on or off during keying practice.

### Credit where credit is due

Thanks to **cocoonlife**. I downloaded the [goalsa](https://github.com/cocoonlife/goalsa) package into the **cwt/mainprocess/goalsa/** folder.

If you are interested in CGO, take a look at **cwt/mainprocess/goalsa/**. I only slightly modified the cocoonlife code in the package's file **alsa.go**. My mods allowed that package to compile and made playing sound a little friendlier.

## The application services

The application offers the user 3 services.

### Selections

The **Selections** service of the application allows the user to select the morse code characters which are to be simultaniously copied and keyed by the user. It also shows the user the keying and copying test scores for the characters at each of the various speeds.

### Copy

The **Copy** service of the application allows the user to practice or test copying morse code at the selected copy speed. The difference between practice and test is that test saves results.

### Key

The **Key** service of the application allows the user to practice or test keying morse code at the selected key speed. One difference between practice and test is that test involves more words and saves results.

#### Key practice

Key practice offers a very advanced feature. **The keying metronome.** The optional metronome during keying practice, allows the user who is keying, to hear the wpm rythem. Keying to the correct wpm rythem is required for keying that others will understand.

The **Key** service is designed to help me and you learn to key properly.

#### The current cwt keying principles

The morse code key is like a piano that only has one key.

Morse code is like a song that has

* a beat that never changes,
* short notes ( dits ) that must only last 1 beat,
* long notes ( dahs ) that must only last 3 beats,
* pauses between the dit and dahs in a character that must only last 1 beat,
* pauses between characters in a word that must only last 3 beats,
* pauses between words that must only last 7 beats.

My intention is that when I key with cwt, those priciples become a part of my concious thought process and then slip into my sub concious. Then I will key correctly without even thinking about it.

#### The current cwt implementation for keying practice

1. One randomly created 5 character word at a time. The word is created from the characters that the user has the worst test scores for.
1. An optional metronome to keep the beat.
1. A morse code cheat sheet with per beat keying instructions.
1. Tolerances.

## The morse code key

My morse code key is a straight key. It is wired to the contacts of the left button on a board that I ripped out of a cheap usb mouse. The board gets plugged into a usb port on the lap top. Pressing the key down causes a "mouse-down" event and letting the key up causes a "mouse-up" event.

## To run on ubuntu

The executable **cwt/cwt** is compiled for 64 bit ubuntu 18.04 on an amd64. Just download it and try it.

``` shell

$ go get github.com/josephbudd/cwt
$ cd ~/src/github.com/josephbudd/cwt
$ ./cwt

```

If it does not run try adding libasound2-dev but I think its already in the executable.

``` shell

$ sudo apt install libasound2-dev

```

## To build on linux

* The mainprocess/goalsa package is written in cgo which requires gcc.
* The mainprocess/goalsa package uses the alsa lib and requires libasound2-dev.

So install those with

``` shell

$ sudo apt install gcc
$ sudo apt install libasound2-dev

```

### Get cwt

``` shell

$ go get github.com/josephbudd/cwt

```

### Get cwt's other dependencies

* [the boltdb package.](https://github.com/boltdb/bolt)
* [the yaml package.](https://gopkg.in/yaml.v2)
* [the gorilla websocket package.](https://github.com/gorilla/websocket)

``` shell

$ go get github.com/boltdb/bolt/...
$ go get gopkg.in/yaml.v2
$ go get github.com/gorilla/websocket

```

### Finally download the kickpack tool

The new renderer build scripts which are **renderer/build.sh** and **renderer/buildPack.sh** use kickpack. So you will need to download, build and install kickpack.

``` shell

$ go get -u github.com/josephbudd/kickpack
$ cd ~/go/src/github.com/josephbudd/kickpack
$ go install

```
 
### The application code is physically and logically organized into 4 areas

1. The **cwt/domain/** folder contains domain ( shared ) logic.
1. The **cwt/mainprocess/** folder contains the main process logic.
1. The **cwt/renderer/** folder contains the renderer logic. The code in the **cwt/renderer/** folder is compiled into wasm at the file **cwt/site/app.wasm**.
1. The **cwt/site/** folder contains the wasm, templates, styles etc for the browser.

### The application has 2 processes

1. The **main process** is a web server running through whatever port you indicate in your application's http.yaml file. Port 0 allows any suitable open port to be selected. When you start the application, it runs the main process. The main process opens a browser and serves the renderer process from the **cwt/site/** folder to the browser.
1. The **renderer process** runs in the browser. It is all of the wasm, html, css, images, etc contained in the **cwt/site** folder.

### The application has a 2 step build

So when you build the application, you build both the renderer process and the main process.

#### Build Step 1: Build the renderer process

``` shell

$ cd $GOPATH
$ cd src/github.com/josephbudd/cwt/renderer
nil@NIL:~/go/src/github.com/josephbudd/cwt/renderer$ ./buildPack.sh 

STEP 1:
REMOVE YOUR PREVIOUS BUILD OF /home/nil/go/src/github.com/josephbudd/cwt/renderer/spawnpack
rm -r /home/nil/go/src/github.com/josephbudd/cwt/renderer/spawnpack

STEP 2:
WRITE THE SOURCE CODE FOR YOUR NEW spawnpack PACKAGE.
 * The spawnpack package is your renderer's spawn html templates.
cd /home/nil/go/src/github.com/josephbudd/cwt/site
kickpack -nu -o=/home/nil/go/src/github.com/josephbudd/cwt/renderer/spawnpack ./spawnTemplates
cd /home/nil/go/src/github.com/josephbudd/cwt/renderer
 * Success. The source code for your new spawnpack package is written.

STEP 3:
BUILD THE RENDERER GO CODE INTO WEB ASSEMBLY CODE AT ../site/app.wasm
GOARCH=wasm GOOS=js go build -o ../site/app.wasm main.go panels.go
 * Success. The renderer go code is compiled into web assembly code at ../site/app.wasm

STEP 4:
WRITE THE cwtsitepack PACKAGE SOURCE CODE.
 * The cwtsitepack package will be a file store
     containing all of the files in the /site folder.
 * So this process could take a while.
 * See func serveFileStore in Serve.go.
cd /home/nil/go/src/github.com/josephbudd/cwt
kickpack -nu -strict -o=/home/nil/go/src/github.com/josephbudd/cwtsitepack ./site ./Http.yaml
 * Success. The cwtsitepack package source code has been written.

STEP 5:
BUILD THE cwtsitepack PACKAGE.
 * This process could take use of all of your CPU cores.
 * This process will take a while.
cd /home/nil/go/src/github.com/josephbudd/cwtsitepack
go build
 * Success.
 * You have successfully compiled the cwtsitepack package object code.
 * In other words, you have compiled your entire renderer process into object code.

STEP 6:
BUILD THE ENTIRE APPLICATION INTO A SINGLE EXECUTABLE cwt.
   You will do so with the following 2 commands...
   cd ..
   go build

```

#### Build Step 2: Build the main process and the entire app into a single executable file ./cwt

``` shell

$ cd ..
$ go build
$ ./cwt

```

## Distribution

Sure you should build cwt but if you don't want to build cwt, that's ok. The ubuntu 18.04.2 linux executable **cwt/cwt** is included in the source. I assume that you can still distribute this app using the existing **cwt/cwt**.

### How to install

The entire application is contained in the executable file **cwt/cwt**. Put that file where ever you want it. You can start it however you want to start it.

### FYI

Cwt creates and stores your keying information at **~/.cwt_kwfw/boltdb/allstores.nosql**. If you delete the folder or file cwt will just create a new one the next time you run it.

### How to uninstall

1. Delete **~/.cwt_kwfw/**.

## Miscelanious

The **cwt/.kickwasm/** folder was created by kickwasm, the tool which created this application's framework. The folder contains information that rekickwasm needs. Rekickwasm is a refactoring tool for this application's framework.

## License

This application has an MIT License. You are free to take this source code and use it to make an even better application.

## The cwt video

This video is of an earlier version of cwt. I still have to make the new video for this current version of cwt.

[![Learning Morse Code with CWT.](https://i.vimeocdn.com/video/772644525.jpg)](https://vimeo.com/328175343)
