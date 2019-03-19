package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"os"
)

func main() {
	var natsAddr string

	flag.StringVar(&natsAddr, "nats-addr", "localhost:4222", "host:port of NATS server")

	flag.Parse()

	u := uuid.Must(uuid.NewV4())
	subj := fmt.Sprintf("$$.nats-health-check.%s", u)

	nc, err := nats.Connect(natsAddr)
	if err != nil {
		fmt.Println("Failed to connect to NATS server.")
		os.Exit(1)
	}

	sub, _ := nc.SubscribeSync(subj)

	_ = nc.Publish(subj, nil)

	_, err = sub.NextMsg(1e+10) // 10sec in nanoseconds
	if err != nil {
		nc.Close()
		fmt.Println("NATS health-check ERROR.")
		os.Exit(1)
	}

	fmt.Println("NATS health-check OK.")

	nc.Close()
}
