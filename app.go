package main

import (
	"bufio"
	"bytes"
	"os"
	"os/signal"
	"syscall"

	"github.com/linkedin/goavro"
	"github.com/nsqio/go-nsq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	topic := kingpin.Flag("topic", "topic to listen to").Required().String()
	lookupdHTTPAddr := kingpin.Flag("lookupd-http-addr", "NSQlookupd address with port, e.g. 127.0.0.1:4161").Required().String()
	maxInFlight := kingpin.Flag("max-in-flight", "NSQ consumer max-in-flight number").Default("100").Int()
	kingpin.Parse()

	channel := "nsqtail#ephemeral" // use #ephemeral for temporary channel

	conf := nsq.NewConfig()
	conf.MaxInFlight = *maxInFlight

	cons, err := nsq.NewConsumer(*topic, channel, conf)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	cons.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		m.Finish()

		// NOTE(sonny): let's show only if avro parsing succeeds
		// TODO(sonny): need be able to handle avro / plain text message
		r := bytes.NewReader(m.Body)
		br := bufio.NewReader(r)

		ocrf, err := goavro.NewOCFReader(br)
		if err != nil {
			// leave it here - next will handle text message too
			// log.Infof("err (maybe json message): %v", err)
			return nil
		}

		var datum interface{}
		for ocrf.Scan() {
			datum, err = ocrf.Read()
			if err != nil {
				log.Errorf("failed to read ocrf: %v", err)
				return err
			}
		}

		log.Infoln(datum)
		return nil
	}))

	err = cons.ConnectToNSQLookupd(*lookupdHTTPAddr)
	if err != nil {
		log.Fatalf("failed to connect to lookupd: %v", err)
	}

	sch := make(chan os.Signal)
	go func(ch <-chan os.Signal) {
		for s := range ch {
			log.Infof("receiving signal %s - exiting gracefully", s.String())
			os.Exit(0)
		}
	}(sch)
	signal.Notify(sch, os.Interrupt, syscall.SIGTERM)

	// block forever
	select {}
}
