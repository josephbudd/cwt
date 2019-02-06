# cwt

A linux morse code trainer. Close to being finished but still a beta version.

## History

I originally wrote this with sciter and it's tiscript worked well with the critical key timing. This version is written with a framework I generated with kickwasm. So the critical timing code in this version is written with go compiled to wasm. This will be a test to see if go wasm can keep up with the keying.

I was going to wait to put this up here at git hub but I just almost lost my hard drive so just to be safe it's up here.

I'ts pretty much done but right now I'm working on the most critical part; recording the user's keying of morse code.

## The key

My key is wired to the contacts of the left button on a usb mouse. So I record the time it takes for mouse up and down in the renderer and let the main process copy those times into valid morse code and then into text. Not rocket science. It's all pretty simple stuff.

## Alsa

So this uses alsa and requires alsa to build and maybe to run.

``` bash

sudo apt-get install alsa-base libasound2-dev

```