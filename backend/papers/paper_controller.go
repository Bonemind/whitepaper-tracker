package papers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PaperController struct {
	Db *gorm.DB
}

func NewController(db *gorm.DB) *PaperController {
	return &PaperController{Db: db}
}

func (c *PaperController) ListWhitepapers(w http.ResponseWriter, r *http.Request) {
	var papers []Paper
	result := c.Db.Find(&papers)
	log.Println("GET /api/whitepapers")

	if result.Error != nil {
		log.Fatalf("Error reading items from db: %v", result.Error)
	}

	js, err := json.Marshal(papers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	w.Write(js)

}

type WhitepaperUpdateRequest struct {
	Id   string
	Read bool
}

func (c *PaperController) UpdateWhitepaper(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		log.Println("OPTIONS /api/whitepaper")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("POST /api/whitepaper")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var wpr WhitepaperUpdateRequest
	err = json.Unmarshal(body, &wpr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing body: %v", err), http.StatusBadRequest)
		return
	}

	var paper Paper
	result := c.Db.First(&paper, "id = ?", wpr.Id)

	if result.RowsAffected == 0 {
		http.Error(w, fmt.Sprintf("Paper not found: %v", result.Error), http.StatusNotFound)
		return
	}

	if wpr.Read {
		paper.DateRead = time.Now()
	} else {
		nilTime, _ := time.Parse(DATETIME_LAYOUT, "0001-01-01T00:00:00Z")
		paper.DateRead = nilTime
	}
	result = c.Db.Save(paper)

	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Error saving paper: %v", result.Error), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(paper)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	w.Write(js)
}
