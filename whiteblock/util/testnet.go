package util

func ReadTestnetStore(name string, out interface{}) error {
	return JsonRpcCallP("state::get", []interface{}{name}, out)
}

func WriteTestnetStore(name string, in interface{}) error {
	_, err := JsonRpcCall("set_extra", []interface{}{name, in})
	return err
}
