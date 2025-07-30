package main

import (
	"fmt"

	"sumx/api"
	"sumx/llm"
	"sumx/db"
)

func main() {
	llmClient := llm.NewHFClient("Qwen/Qwen3-235B-A22B-Instruct-2507:together")
	db := db.Connect()
	server := api.NewServer(llmClient, db)

	fmt.Println("🚀 API server running on :8080")
	server.Router.Run(":8080")
}
