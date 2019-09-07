package main

import (
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
	log.Fatal(d.StartHTTP(":8080"))
}
