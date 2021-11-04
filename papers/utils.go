package papers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const DATETIME_LAYOUT = "2006-01-02T15:04:05-0700"
const URL = "https://aws.amazon.com/api/dirs/items/search?item.directoryId=whitepapers&sort_by=item.additionalFields.sortDate&sort_order=desc&size=%d&item.locale=en_US&page=%d"
const ITEMS_PER_PAGE = 100

func getUrlForPage(page int) string {
	return fmt.Sprintf(URL, ITEMS_PER_PAGE, page)
}

func TestItemFetch(items int) (*Response, error) {
	var parsedResponse Response
	resp, err := http.Get(fmt.Sprintf(URL, items, 0))
	if err != nil {
		return nil, fmt.Errorf("Error fetching whitepapers: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing json response: %v", err)
	}
	return &parsedResponse, nil
}
