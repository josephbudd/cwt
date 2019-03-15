# cwt

## A linux Continuous Wave ( Morse Code ) Trainer

According to [wikipedia](https://en.wikipedia.org/wiki/Continuous_wave)

> Continuous wave is also the name given to an early method of radio transmission, in which a sinusoidal carrier wave is switched on and off. Information is carried in the varying duration of the on and off periods of the signal, for example by Morse code in early radio.

## Credit where credit is due

Thanks to **cocoonlife**. I downloaded the goalsa package at https://github.com/cocoonlife/goalsa into the **cwt/mainprocess/goalsa/** folder. I wish I would have taken the same approach with my own alsa code which I've given up on. Take a look at it if you are interested in CGO. Anyway, I slightly modified the code in the package's file **alsa.go**. My mods allowed that package to compile and made playing sound a little friendlier.

## The application services

### Selections

The **Selections** service of the application allows the user to select the morse code characters which are to be simultaniously copied and keyed by the user. It also shows the user the keying and copying test scores for the characters at each of the various speeds.

### Copy

The **Copy** service of the application allows the user to practice or test copying morse code at the selected copy speed. The difference between practice and test is that test saves results.

### Key

The **Key** service of the application allows the user to practice or test keying morse code at the selected key speed. One difference between practice and test is that test saves results.

That's all well and good for people who know how to key. But despite my attempts to teach myself to key years ago, I now realize that I didn't really understand what keying was and only taught myself a bunch of crap that I want to forget.

The **Key** service is designed to help me and you learn to key properly.

#### The cwt keying rational

The morse code key is like a piano that only has one key.

Morse code is like a song that has

* a beat that never changes,
* short notes ( dits ) that must only last 1 beat,
* long notes ( dahs ) that must only last 3 beats,
* pauses between the dit and dahs in a character that must only last 1 beat,
* pauses between characters in a word that must only last 3 beats,
* pauses between words that must only last 7 beats.

Did you store those principles into your memory? Don't worry, you don't have to. Cwt will passively lead you to follow those principles. Without realizing it, those priciples will become a part of your concious thought process and then slip into your sub concious. Then you will key correctly without even thinking about it.

#### The cwt solution for practicing keying

1. One randomly created five character word at a time. The word is created from the characters that the user has the worst test scores for.
1. A metronome to keep the beat.
1. A morse code cheat sheet with per beat keying instructions.
1. Tolerances.

## The morse code key

My morse code key is a straight key. It is wired to the contacts of the left button on a board that I ripped out of a cheap usb mouse. The board gets plugged into a usb port on the lap top. Pressing the key down causes a "mouse-down" event and letting the key up causes a "mouse-up" event.

## To install and build on linux

* The mainprocess/goalsa package is written in cgo which requires gcc.
* The mainprocess/goalsa package uses the alsa lib and requires libasound2-dev.

So install those with

``` bash

sudo apt install gcc
sudo apt install libasound2-dev

```

Then get cwt with

``` bash

go get -u github.com/josephbudd/cwt

```

Doing so will also import the following packages for cwt

* [the boltdb package.](https://github.com/boltdb/bolt)
* [the yaml package.](https://gopkg.in/yaml.v2)
* [the gorilla websocket package.](https://github.com/gorilla/websocket)

## The build

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

* The main process source code is for your operating system ( linux ). It is compiled into the executable file at **cwt/cwt**.
* The renderer process source code is in the **cwt/renderer/** folder. It is compiled into web assembly code at **cwt/site/app.wasm** file.

There is a shell script at **cwt/renderer/build.sh** which builds the renderer process into the **cwt/site/** folder.

``` bash

cd $GOPATH
cd src/github.com/josephbudd/cwt
go build
cd renderer
./build.sh
cd ..
./cwt

```

## Distribution

Sure you should build cwt but if you don't want to build cwt, that's ok. The ubuntu 18.04.2 linux executable **cwt/cwt** is included in the source. I assume that you can still distribute this app using the existing **cwt/cwt**.

### How to install

1. Make an empty distribution folder.
1. Copy the executable **cwt/cwt** file into your distribution folder.
1. Copy the **cwt/http.yaml** file into your distribution folder.
1. Copy the **cwt/site/** folder into your distribution folder.
1. Then put your distribution folder where you want it and run the executable in it.
1. You won't need the downloaded source code after that so you can delete it if you want.

### FYI

1. To run cwt just double click on it with the mouse or run it in a terminal.
1. Cwt creates and stores your keying information at **~/.cwt_kwfw/boltdb/allstores.nosql**. If you delete the folder or file cwt will just create a new one the next time you run it.

### How to uninstall

1. Delete the distribution folder.
1. Delete **~/.cwt_kwfw/**.

## Miscelanious

The **cwt/.kickwasm/** folder was created by kickwasm, the tool which created this application's framework. The folder contains information that rekickwasm needs. Rekickwasm is a refactoring tool for this application's framework.

## License

This application has an MIT License. You are free to take this source code and use it to make an even better application.
