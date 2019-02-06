# cwt

A linux morse code trainer.

I originally wrote this with sciter and it's tiscript worked well with the critical key timing. This version is written with a framework I generated with kickwasm. So the critical timing code in this version is written with go compiled to wasm. This will be a test to see if go wasm can keep up with the keying.

I was going to wait to put this up here at git hub but I just almost lost my hard drive so just to be safe it's up here.

So this uses alsa and requires alsa to build and maybe to run.

``` bash

sudo apt-get install alsa-base libasound2-dev

```