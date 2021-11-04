package main

import (
	"fmt"
	"log"
	"net/http"
	"whitepaper-tracker/fileserver"
	"whitepaper-tracker/papers"

	"github.com/jasonlvhit/gocron"
	"github.com/namsral/flag"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	testMode := flag.Bool("test_fetch", false, "Whether to test if item fetch still works instead of starting the server")
	skipInitialLoad := flag.Bool("skipload", false, "Whether to skip the initial whitepaper load, useful for testing")
	dbLocation := flag.String("db_location", "papers.db", "The location of the sqlite db")
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

	db, err := gorm.Open(sqlite.Open(*dbLocation), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}

	db.AutoMigrate(&papers.Paper{})

	if *skipInitialLoad {
		log.Println("Skipping initial load")
	} else {
		log.Println("Loading whitepapers")
		err = papers.LoadPapers(db)
		if err != nil {
			log.Fatalf("Error loading items: %v", err)
		}

		log.Println("Load done")
	}

	// Schedule paper renewal on a daily basis
	gocron.Every(1).Minute().Do(papers.LoadPapers, db)

	log.Println("Starting gocron worker...")
	gocron.Start()
	log.Println("Gocron worker started")

	paperController := papers.NewController(db)

	http.HandleFunc("/api/whitepapers", paperController.ListWhitepapers)
	http.HandleFunc("/api/whitepaper", paperController.UpdateWhitepaper)
	http.HandleFunc("/", fileserver.ServeFrontend)
	log.Printf("Starting server on port %d\n", *serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil)
}
