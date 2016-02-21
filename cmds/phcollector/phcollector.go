package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/icphalanx/collector"
)

var (
	Port = flag.Int("port", 13890, "port for collector to listen on")
	Host = flag.String("host", "[::]", "host for collector to listen on")

	CertLocation    = flag.String("certLocation", "phagent.crt", "location to store certificate")
	PrivKeyLocation = flag.String("privKeyLocation", "phagent.key", "location to store private key")

	CALocation = flag.String("caLocation", "phalanx.crt", "location of CA certificate")
	CARemote   = flag.String("caRemote", "http://root.phalanx.lukegb.com:8888", "location of remote cfssl server")

	DBLocation = flag.String("dbLocation", "postgres://localhost/phalanx?sslmode=disable", "database connection URL")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *Host, *Port)

	dbConn, err := sql.Open("postgres", *DBLocation)
	if err != nil {
		log.Fatalln(err)
	}

	tlsC, err := collector.MakeTLSConfig(*CALocation, *CertLocation, *PrivKeyLocation)
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(collector.ListenAndServeAtAddr(addr, collector.Config{
		CARemote:  *CARemote,
		TLSConfig: tlsC,
		DBConn:    dbConn,
	}))
}
