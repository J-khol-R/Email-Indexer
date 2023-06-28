package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/J-khol-R/Email-Indexer/models"
)

func RequestZincsearch(key string, inicio, fin int) (models.ResponseZinc, error) {
	url := "http://localhost:4080/api/enron_mails/_search"

	query := `{
	    "search_type": "match",
	    "query":
	    {
	        "term": "` + key + `"
	    },
	    "from": ` + fmt.Sprint(inicio) + `,
	    "max_results": ` + fmt.Sprint(fin) + `,
	    "_source": []
	}`

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(query))

	if err != nil {
		fmt.Println(err)
		return models.ResponseZinc{}, err
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return models.ResponseZinc{}, err
	}
	defer res.Body.Close()

	var response models.ResponseZinc
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return models.ResponseZinc{}, err
	}

	return response, nil
}
