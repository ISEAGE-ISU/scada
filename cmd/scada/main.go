package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ISEAGE-ISU/scada"
	"github.com/ISEAGE-ISU/scada/crane"
	"github.com/ISEAGE-ISU/scada/generator"
	"github.com/ISEAGE-ISU/scada/security"
)

func main() {
	s, found := os.LookupEnv("DEVICE")
	if found == false {
		log.Fatal("DEVICE not set")
	}

	if len(os.Args) > 1 {
		f, err := os.Create(os.Args[1])
		if err != nil {
			log.Fatal("could not create pid file:", err)
		}
		f.Write([]byte(fmt.Sprintf("%d", os.Getpid())))
	}

	var d *scada.Device
	switch s {
	case "gen":
		d = generator.Create()
	case "crane":
		d = crane.Create()
	case "camera":
		d = security.Create()
	default:
		log.Fatal("Device unknown: ", s)
	}

	go d.StartSCADA()
	log.Fatal(d.StartHTTP(":80"))
}
