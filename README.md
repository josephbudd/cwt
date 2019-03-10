# cwt

A **linux** morse code trainer.

## The framework for this application was built by kickwasm v3.0.0

This application will be finished when I'm convinced that I can use this application to learn to key. See **The key service** below.

## Credit where credit is due

I downloaded the awesome goalsa package at https://github.com/cocoonlife/goalsa into the **cwt/mainprocess/goalsa/** folder. The code is well written. I wish I would have taken the same approach with my own alsa code which I've given up on. Anyway, I slightly modified the code in the package's file **alsa.go**. My mods allowed that package to compile and made playing sound a little friendlier.

## The application services

### The reference service

The reference service of the application allows the user to select the morse code characters which are to be simultaniously copied and keyed by the user. It also shows the user the keying and copying test scores for the characters at each of the various speeds.

### The copy service

The copy service of the application allows the user to practice or test copying morse code at the selected copy speed. The difference between practice and test is that test saves results.

### The key service

The key service of the application allows the user to practice or test keying morse code at the selected key speed. The difference between practice and test is that test saves results.

That's all well and good for people who know how to key. But despite my attempts to teach myself to key years ago, I now realize that I didn't really understand what keying was and only taught myself a bunch of crap that I want to forget.

So the key service is designed to help me learn to key properly.

#### Here is my premise concerning keying

The morse code key is like a piano that only has one key.

Morse code is like a song that has

* a beat that never changes,
* short notes ( dits ) that must only last 1 beat,
* long notes ( dahs ) that must only last 3 beats,
* pauses between the dit and dahs in a character that must only last 1 beat,
* pauses between characters in a word that must only last 3 beats,
* pauses between words that must only last 7 beats.

#### My solution

1. So I have allowed for variances in user input.
1. I added a metronome to help me keep the beat.
1. I'm trying to come up with techiques for learning to key.

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

1. The **main process** is a web server running through whatever port you indicate in your application's http.yaml file. When you start the application, it runs the main process. The main process opens a browser which loads and runs the renderer process from the **cwt/site/** folder.
1. The **renderer process** is for the browser. It is all of the wasm, html, css, images, etc contained in the **cwt/site** folder.

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

If you don't want to or cant build cwt, that's ok. The ubuntu 18.04.2 linux binary **cwt** is included in the source so you can still distribute and run it on linux.

### How to

1. Make an empty distribution folder.
1. Copy the executable **cwt/cwt** file into your distribution folder.
1. Copy the **cwt/http.yaml** file into your distribution folder.
1. Copy the **cwt/site/** folder into your distribution folder.
1. Then put your distribution folder where you want it and run the executable in it.
1. You won't need the downloaded source code after that so you can delete it if you want.

Let me know if this technique works for you.

## License

This application has an MIT License. You are free to take this source code and use it to make a better application.