package main

import "github.com/zzztttkkk/h2tp"

func main() {
	router := h2tp.NewRouter()

	if err := h2tp.Run("127.0.0.1:8524", map[string]*h2tp.Router{"*": router}); err != nil {
		panic(err)
	}
}
