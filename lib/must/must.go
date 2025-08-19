package must

import "log"

func Must[T any](v T, err error) T {
	if err != nil {
		log.Panic(err)
	}
	return v
}

func Run(err error) {
	if err != nil {
		log.Panic(err)
	}
}
