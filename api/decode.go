package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// decodes json into your provided struct. Using this to avoid making a massive all encompassing struct
// remember &dst, if you get json: Unmarshal(non-pointer) error
func decodeForm(req *http.Request, dst interface{}) {
	if err := json.NewDecoder(req.Body).Decode(dst); err != nil {
		fmt.Println("ERRRORORORORORO:", err)
	}
}
