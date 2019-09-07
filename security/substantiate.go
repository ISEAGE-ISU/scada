package security

import (
	"log"

	"github.com/ISEAGE-ISU/scada"
)

func Create() *scada.Device {
	err := Start()
	if err != nil {
		log.Fatal(err)
	}
	return &scada.Device{Root, Status, IsOK, nil}
}
