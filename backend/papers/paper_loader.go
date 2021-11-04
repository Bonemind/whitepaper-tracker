package papers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Response struct {
	Metadata struct {
		Count     int `json:"count"`
		TotalHits int `json:"totalHits"`
	} `json:"metadata"`
	FieldTypes struct {
		UpdateDate        string `json:"updateDate"`
		ImageSrcURL       string `json:"imageSrcUrl"`
		FeatureFlag       string `json:"featureFlag"`
		Description       string `json:"description"`
		SortDate          string `json:"sortDate"`
		DocTitle          string `json:"docTitle"`
		PrimaryURL        string `json:"primaryURL"`
		DatePublished     string `json:"datePublished"`
		PublishedText     string `json:"publishedText"`
		FooterInfoSubtext string `json:"footerInfoSubtext"`
		SubHeadline       string `json:"subHeadline"`
		EnableShare       string `json:"enableShare"`
		Category          string `json:"category"`
		ContentType       string `json:"contentType"`
	} `json:"fieldTypes"`
	Items []struct {
		Tags []struct {
			TagNamespaceID string `json:"tagNamespaceId"`
			CreatedBy      string `json:"createdBy"`
			Name           string `json:"name"`
			DateUpdated    string `json:"dateUpdated"`
			Locale         string `json:"locale"`
			LastUpdatedBy  string `json:"lastUpdatedBy"`
			DateCreated    string `json:"dateCreated"`
			Description    string `json:"description"`
			ID             string `json:"id"`
		} `json:"tags"`
		Item struct {
			CreatedBy        string  `json:"createdBy"`
			Locale           string  `json:"locale"`
			Author           string  `json:"author"`
			DateUpdated      string  `json:"dateUpdated"`
			Score            float64 `json:"score"`
			Name             string  `json:"name"`
			NumImpressions   int     `json:"numImpressions"`
			DateCreated      string  `json:"dateCreated"`
			AdditionalFields struct {
				DatePublished string `json:"datePublished"`
				FeatureFlag   string `json:"featureFlag"`
				PublishedText string `json:"publishedText"`
				Description   string `json:"description"`
				DocTitle      string `json:"docTitle"`
				SortDate      string `json:"sortDate"`
				EnableShare   string `json:"enableShare"`
				ContentType   string `json:"contentType"`
				PrimaryURL    string `json:"primaryURL"`
			} `json:"additionalFields"`
			ID            string `json:"id"`
			DirectoryID   string `json:"directoryId"`
			LastUpdatedBy string `json:"lastUpdatedBy"`
		} `json:"item"`
	} `json:"items"`
}

func LoadPapers(db *gorm.DB) error {
	totalCount := ITEMS_PER_PAGE * 2
	fetchedItems := 0
	log.Println("Started paper loading")

	for i := 0; fetchedItems < totalCount; i++ {
		var parsedResponse Response
		resp, err := http.Get(getUrlForPage(i))
		if err != nil {
			return fmt.Errorf("Error fetching whitepapers: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading response body: %v", err)
		}

		err = json.Unmarshal(body, &parsedResponse)
		if err != nil {
			return fmt.Errorf("Error parsing json response: %v", err)
		}

		totalCount = parsedResponse.Metadata.TotalHits
		fetchedItems += parsedResponse.Metadata.Count

		for _, item := range parsedResponse.Items {
			dateCreated, err := time.Parse(DATETIME_LAYOUT, item.Item.DateCreated)
			if err != nil {
				return fmt.Errorf("Error parsing creation date, got: %s, error: %v", item.Item.DateCreated, err)
			}

			dateUpdated, err := time.Parse(DATETIME_LAYOUT, item.Item.DateUpdated)
			if err != nil {
				return fmt.Errorf("Error parsing update date, got: %s, error: %v", item.Item.DateUpdated, err)
			}

			paper := &Paper{
				Id:            item.Item.ID,
				Name:          item.Item.Name,
				DatePublished: dateCreated,
				DateUpdated:   dateUpdated,
				Url:           item.Item.AdditionalFields.PrimaryURL,
				Type:          item.Item.AdditionalFields.ContentType,
				Title:         item.Item.AdditionalFields.DocTitle,
			}

			// If a paper already exists, we want to update to a new url, and want to update the update date
			db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"date_updated", "url"}),
			}).Create(paper)
		}
	}
	log.Printf("Done loading papers, fetched %d papers\n", fetchedItems)

	return nil
}
