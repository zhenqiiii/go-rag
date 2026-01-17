package main

import "go-rag/api"

func main() {
	// router
	r := api.SetupRouter()

	r.Run()
}
