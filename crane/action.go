package crane

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const RootDir = "/tmp/crane"
const onFile = RootDir + "/on"
const heightFile = RootDir + "/heightFile"

const GenRoot = "/tmp/generator"
const GenOnFile = GenRoot + "/on"
const GenPowerOut = GenRoot + "/power"

func init() {
	os.Mkdir(RootDir, 0777)
}

var ErrAlreadyStarted = errors.New("Crane is already on")
var ErrNoPower = errors.New("Not enough Power")

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
	if pow < 60 {
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
	f, err := os.Create(heightFile)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte("10"))
	return nil
}

func Stop() error {
	if !haveEnoughPower() {
		return ErrNoPower
	}
	if err := os.Remove(onFile); err != nil {
		return err
	}
	if err := os.Remove(heightFile); err != nil {
		return err
	}
	return nil
}

func Set(level string) error {
	f, err := os.Create(heightFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(level))
	return err
}

func Status() string {
	isOn := "off"
	height := "0"
	if _, err := os.Stat(onFile); err == nil {
		isOn = "on"
	}
	f, err := os.Open(heightFile)
	if isOn == "on" && err == nil {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		height = string(b)
	}
	return fmt.Sprintf("Crane State: %s, height: %sft\n", isOn, height)
}

func IsOK() bool {
	if _, err := os.Stat(onFile); err != nil {
		return false
	}

	f, err := os.Open(heightFile)
	if err != nil {
		return false
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return false
	}
	height, err := strconv.Atoi(string(b))
	if err != nil || height < 0 || height > 100 {
		return false
	}
	return true
}
