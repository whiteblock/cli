package util

import(
    "os"
    "log"
    "io/ioutil"
)

var storeDirectory string 

func init(){
    home := os.Getenv("HOME")
    storeDirectory = home+"/.config/whiteblock/cli/"
    err := os.MkdirAll(storeDirectory, 0755)
    if err != nil {
        log.Fatalf("Could not create directory: %s", err)
    }
}

func ReadStore(name string) ([]byte,error) {
   return ioutil.ReadFile(storeDirectory+name)
}

func WriteStore(name string,data []byte) error {
    return Write(storeDirectory+name,data)
}