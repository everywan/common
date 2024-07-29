package utils

import "log"

func HandleInitError(module string, err error) {
	if err != nil {
		log.Fatalf("init %s error. err:%s", module, err)
	}
}
