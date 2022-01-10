package contextobj

var (
	obj = make(map[string]interface{})
)

func Set(key string, val interface{}) {
	obj[key] = val
}

func Get() map[string]interface{} {
	return obj
}
