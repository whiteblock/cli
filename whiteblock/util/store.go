package util

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func init() {
	err := os.MkdirAll(conf.StoreDirectory, 0755)
	if err != nil {
		log.Fatalf("Could not create directory: %s", err)
	}
}

func ReadStore(name string) ([]byte, error) {
	return ioutil.ReadFile(conf.StoreDirectory + name)
}

func WriteStore(name string, data []byte) error {
	return ioutil.WriteFile(conf.StoreDirectory+name, data, 0664)
}

func DeleteStore(name string) error {
	return os.Remove(conf.StoreDirectory + name)
}

func StoreExists(name string) bool {
	file := conf.StoreDirectory + name
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func ReadTestnetStore(name string, outptr interface{}) error {
	return JsonRpcCallP("state::get", []interface{}{name}, outptr)
}

func WriteTestnetStore(name string, in interface{}) error {
	_, err := JsonRpcCall("set_extra", []interface{}{name, in})
	return err
}
