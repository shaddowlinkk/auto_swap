package Utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadConfig[C any](filename string) C {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read json api file err:%s\n", err)
	}
	var api C
	if err := json.Unmarshal(data, &api); err != nil {
		fmt.Println("Error in unmarshalling json data")
		fmt.Println(err)
	}
	return api
}
