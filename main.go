package main

import (
	"fmt"
	"log"
	"net"

	"github.com/josephbudd/cwt/domain/data/settings"
	"github.com/josephbudd/cwt/domain/store"
	"github.com/josephbudd/cwt/mainprocess/lpc"
	"github.com/josephbudd/cwt/mainprocess/lpc/dispatch"
)

/*
	DO NOT EDIT THIS FILE.

	kicklpc and kickstore will alter this file.

*/

/*

	Data Storage:
	 * /domain/store/storer is the storer interfaces.
	 * /domain/store/storing is the bolt implementations of the storer interfaces.
	 * /domain/store/record is the record definitions.

*/

func main() {
	var err error
	// Build the application's data store APIs.
	var stores *store.Stores
	if stores, err = buildStores(); err != nil {
		log.Println(err)
		return
	}
	// Open the stores.
	if err = stores.Open(); err != nil {
		log.Println(err)
		return
	}
	// Close the stores later.
	defer stores.Close()

	// get the application's host and port and then setup the listener.
	var appSettings *settings.ApplicationSettings
	if appSettings, err = settings.NewApplicationSettings(); err != nil {
		log.Println(err)
		return
	}

	// initialize and start the listener.
	// the listener may have reset the address if "localhost:0".
	// use the listener's address.
	location := fmt.Sprintf("%s:%d", appSettings.Host, appSettings.Port)
	var listener net.Listener
	if listener, err = net.Listen("tcp", location); err != nil {
		log.Println(err)
		return
	}
	// get the channels
	sendChan, receiveChan, eojChan := lpc.Channels()
	quitChan := make(chan struct{}, 1)
	// process incoming lpcs.
	go func() {

		defer func() {
			log.Println("eoj for processing incoming lpcs.")
		}()

		// wait for the server to end and then stop lpc go funcs
		log.Println("listening for receiveChan")
		for {
			select {
			case cargo := <-receiveChan:
				log.Println("main: got cargo := <-receiveChan")
				dispatch.Do(cargo, sendChan, eojChan, stores)
			case <-quitChan:
				log.Println("main: got <-quitChan")
				eojChan.Signal()
				return
			}
		}
	}()
	// make the call server.
	server := lpc.NewServer(listener, quitChan, receiveChan, sendChan)
	server.Run(serve)
}
