package util

import ()


func ReadTestnetStore(name string, outptr interface{}) error {
	return JsonRpcCallP("state::get", []interface{}{name}, outptr)
}

func WriteTestnetStore(name string, in interface{}) error {
	_, err := JsonRpcCall("set_extra", []interface{}{name, in})
	return err
}
