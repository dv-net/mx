package testutils

import (
	"encoding/json"
	"log"
)

func PrintJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ") //nolint:errchkjson
	if len(b) != 0 {
		log.Println(string(b))
	}
}
