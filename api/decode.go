package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// decodes json into your provided struct. Using this to avoid making a massive all encompassing struct
// remember &dst, if you get json: Unmarshal(non-pointer) error
func decodeForm(req *http.Request, dst interface{}) {
	if err := json.NewDecoder(req.Body).Decode(dst); err != nil {
		fmt.Println(err)
	}
}

func decodeXMLForm(response *http.Response, dst interface{}) {
	if err := xml.NewDecoder(response.Body).Decode(dst); err != nil {
		fmt.Println(err)
	}
}
