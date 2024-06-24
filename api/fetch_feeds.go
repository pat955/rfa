package api

import (
	"encoding/xml"
	"net/http"
)

type Amount struct {
	Amount int32 `json:"amount"`
}

func GetNextFeedsToFetch(w http.ResponseWriter, r *http.Request) {
	var n Amount
	decodeForm(r, &n)
	db := connect().DB

	feedsToUpdate, err := db.GetNextFeedsToFetch(r.Context(), n.Amount)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}
	respondWithJSON(w, 200, dbFeedsToFeeds(feedsToUpdate))
}

// add errors back
func FetchFeed(endpoint string) RSS {
	response, err := http.Get(endpoint)
	if err != nil {
		return RSS{}
	}

	defer response.Body.Close()
	var rss RSS
	decodeXMLForm(response, &rss)
	return rss
}

type RSS struct {
	XMLName  xml.Name `xml:"rss"`
	CharData string   `xml:",chardata"`
	Atom     string   `xml:"atom,attr"`
	Version  string   `xml:"version,attr"`
	Channel  struct {
		CharData string `xml:",chardata"`
		Title    string `xml:"title"`
		Link     struct {
			CharData string `xml:",chardata"`
			Href     string `xml:"href,attr"`
			Rel      string `xml:"rel,attr"`
			Type     string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Items         []struct {
			CharData    string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}
