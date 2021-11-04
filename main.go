package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"whitepaper-tracker/fileserver"
	"whitepaper-tracker/papers"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	testMode := flag.Bool("test", false, "Whether to test if item fetch still works instead of starting the server")
	serverPort := flag.Int("port", 3000, "Port to listen on")
	flag.Parse()

	if *testMode {
		testItemCount := 2
		result, err := papers.TestItemFetch(testItemCount)
		if err != nil {
			log.Fatalf("Error fetching items: %v", err)
		}
		if result.Metadata.Count != testItemCount {
			log.Fatalf("Got unexpected number of items, expected: %d, actual: %d", testItemCount, result.Metadata.Count)
		}
		log.Println("Loaded expected number of items, integration works")
		return
	}

	db, err := gorm.Open(sqlite.Open("whitepapers.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}

	db.AutoMigrate(&papers.Paper{})

	// fmt.Println("Loading whitepapers")
	// err = papers.LoadItems(db)
	// if err != nil {
	// 	log.Fatalf("Error loading items: %v", err)
	// }

	// fmt.Println("Load done")

	controller := papers.NewController(db)

	http.HandleFunc("/api/whitepapers", controller.ListWhitepapers)
	http.HandleFunc("/api/whitepaper", controller.UpdateWhitepaper)
	http.HandleFunc("/", fileserver.ServeFrontend)
	log.Printf("Starting server on port %d\n", *serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)
}
