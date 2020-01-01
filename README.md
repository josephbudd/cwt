# cwt

## A linux Continuous Wave ( Morse Code ) Trainer

January 1, 2020

* Rebuilt with [kickwasm](https://github.com/josephbudd/kickwasm) version 15.0.0. Fixed a few issues from the crappy rebuild a week or so ago.


* I have not tested with a input of keying's timing because internet is 1.5 hours from home and it's not here with me. I'm assuming that there are no issues as I have not touched that code. I've been keying with the mouse and that works so the keying with the key should work.

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

### Install libasound2-dev.

``` shell

$ sudo apt install libasound2-dev

```

### Download kickwasm and it's tools and dependencies

``` shell
~$ go get -u github.com/josephbudd/kickwasm
~$ cd ~/go/src/github.com/josephbudd/kickwasm
~/go/src/github.com/josephbudd/kickwasm$ make install
~/go/src/github.com/josephbudd/kickwasm$ make test
~/go/src/github.com/josephbudd/kickwasm$ make dependencies
~/go/src/github.com/josephbudd/kickwasm$ make proofs
```

### Build and run the CWT

``` shell
$ cd ~/go/src/github.com/josephbudd/cwt
$ kickbuild -rp -mp -run
```

### How to uninstall

1. Delete **~/.cwt_kwfw/** which is where the database is stored.
1. Delete **~/go/src/github.com/josephbudd/crud**.

## License

This application has an MIT License. You are free to take this source code and use it to make an even better application.

## The cwt video

This video is of an earlier version of cwt. I still have to make the new video for this current version of cwt.

[![Learning Morse Code with CWT.](https://i.vimeocdn.com/video/772644525.jpg)](https://vimeo.com/328175343)
