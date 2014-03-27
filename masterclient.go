// masterclient.go
package main

import (
	"net"

	"github.com/coreos/go-log/log"
)

var ctx_buffer = make(chan e2e_client_ctx, 8)

func enfreshen_scb() {
	log.Alert("Enfreshen!")
	for i := 0; i < 8; i++ {
	retry:
		log.Alert("Enfreshen!")
		thing, err := build_subcircuit()
		if err != nil {
			goto retry
		}
		ctx_buffer <- make_e2e_client_ctx(thing.wire)
	}
}

func run_client_loop() {
	enfreshen_scb()
	// Round robin, basically
	get_ctx := func() e2e_client_ctx {
		toret := <-ctx_buffer
		ctx_buffer <- toret
		return toret
	}
	// Main loop
	listener, err := net.Listen("tcp", MasterConfig.General.ORAddr)
	if err != nil {
		panic(err)
	}
	for {
		nconn, err := listener.Accept()
		if err != nil {
			log.Debug("Error while accepting: ", err.Error())
			continue
		}
		go get_ctx().AttachClient(nconn)
	}
}
