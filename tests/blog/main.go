package main

import "github.com/zzztttkkk/ink"

func main() {
	router := ink.NewRouter()

	if err := ink.Run("127.0.0.1:8524", map[string]*ink.Router{"*": router}); err != nil {
		panic(err)
	}
}
