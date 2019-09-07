package security

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"time"
)

const RootDir = "/tmp/camera"
const onFile = RootDir + "/on"
const freeFile = RootDir + "/freeSpace"

const GenRoot = "/tmp/generator"
const GenOnFile = GenRoot + "/on"
const GenPowerOut = GenRoot + "/power"

var ErrAlreadyStarted = errors.New("Camera is already on")
var ErrNoPower = errors.New("Not enough Power")

func init() {
	os.Mkdir(RootDir, 0777)
	go fillStorage()
}

func fillStorage() {
	space := 500
	tick := time.Tick(30 * time.Minute)
	f, err := os.Create(freeFile)
	f.Write([]byte(fmt.Sprintf("%d", space)))
	f.Close()
	for {
		<-tick
		space += 500
		//Seems this could be racey
		f, err = os.Create(freeFile)
		if err != nil {
			log.Fatal("could not tick:", err)
		}
		f.Write([]byte(fmt.Sprintf("%d", space)))
		f.Close()
	}
}

func haveEnoughPower() bool {
	if _, err := os.Stat(GenOnFile); err != nil {
		return false
	}
	f, err := os.Open(GenPowerOut)
	if err != nil {
		return false
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return false
	}
	pow, err := strconv.Atoi(string(b))
	if err != nil {
		return false
	}
	if pow < 10 {
		return false
	}
	return true
}

func Start() error {
	if !haveEnoughPower() {
		return ErrNoPower
	}
	if _, err := os.Stat(onFile); err == nil {
		return ErrAlreadyStarted
	}
	if f, err := os.Create(onFile); err != nil {
		return err
	} else {
		f.Close()
	}
	return nil
}

func Stop() error {
	if !haveEnoughPower() {
		return ErrNoPower
	}
	if err := os.Remove(onFile); err != nil {
		return err
	}
	if err := os.Remove(freeFile); err != nil {
		return err
	}
	return nil
}

func Status() string {
	isOn := "off"
	space := "COULD NOT READ"
	if _, err := os.Stat(onFile); err == nil {
		isOn = "on"
	}
	f, err := os.Open(freeFile)
	if isOn == "on" && err == nil {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		space = string(b)
	}
	return fmt.Sprintf("Camera State: %s, Used Space: %sGb\n", isOn, space)
}

func IsOK() bool {
	if _, err := os.Stat(onFile); err != nil {
		return false
	}

	f, err := os.Open(freeFile)
	if err != nil {
		return false
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return false
	}
	space, err := strconv.Atoi(string(b))
	if err != nil || space < 0 || space > 2*1000 {
		return false
	}
	return true
}
