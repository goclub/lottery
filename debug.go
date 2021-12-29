package lottery

import "log"

var Debug = false
func debug(kind string, v interface{}) {
	if Debug {
		log.Printf(kind + ":%v", v)
	}
}
