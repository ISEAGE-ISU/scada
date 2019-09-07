package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	os.Mkdir(RootDir, 0777)
	hardwareLine = make(chan struct{})
	go hardwareSecurity()
}

var ErrAlreadyStarted = errors.New("Generator is started")

const RootDir = "/tmp/generator"
const onFile = RootDir + "/on"
const powerOut = RootDir + "/power"

var hardwareLine chan struct{}

func hardwareSecurity() error {
	tick := time.Tick(10 * time.Second)
	hardwareLine <- struct{}{}
	for {
		<-tick
		hardwareLine <- struct{}{}
	}
}

func Start() error {
	if _, err := os.Stat(onFile); err == nil {
		return ErrAlreadyStarted
	}
	if f, err := os.Create(onFile); err != nil {
		return err
	} else {
		f.Close()
	}
	f, err := os.Create(powerOut)
	if err != nil {
		return err
	}
	defer f.Close()
	select {
	case <-hardwareLine:
		f.Write([]byte("50"))
	default:
		f.Write([]byte("-1"))
	}
	return nil
}

func Stop() error {
	if err := os.Remove(onFile); err != nil {
		return err
	}
	if err := os.Remove(powerOut); err != nil {
		return err
	}
	return nil
}

func Set(level string) error {
	f, err := os.Create(powerOut)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(level))
	return err
}

func Status() string {
	isOn := "off"
	power := "0"
	if _, err := os.Stat(onFile); err == nil {
		isOn = "on"
	}
	f, err := os.Open(powerOut)
	if isOn == "on" && err == nil {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		power = string(b)
	}
	if isOn == "off" && power != "0" {
		return "Inconsistant State"
	}
	return fmt.Sprintf("PowerState: %s, Power: %s%%\n", isOn, power)
}

func IsOK() bool {
	if _, err := os.Stat(onFile); err != nil {
		return false
	}

	f, err := os.Open(powerOut)
	if err != nil {
		return false
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return false
	}
	power, err := strconv.Atoi(string(b))
	if err != nil || power <= 0 {
		return false
	}
	return true
}
