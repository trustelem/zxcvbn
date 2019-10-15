// +build no_embedded_dict

package frequency

import (
	"encoding/json"
	"os"
)

var FrequencyLists map[string][]string

func init() {
	file := os.Getenv("ZXCVBN_DEFAULT_DICTIONARIES_JSON")
	if file == "" {
		panic("missing ZXCVBN_DEFAULT_DICTIONARIES_JSON env")
	}
	fd, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	err = json.NewDecoder(fd).Decode(&FrequencyLists)
	if err != nil {
		panic(err)
	}
}
