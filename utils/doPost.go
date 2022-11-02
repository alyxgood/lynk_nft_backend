package utils

import (
	"alyx_nft_backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func QueryNFTInfo(address, data, url string) (model *models.JsonRPCModel, err error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "eth_call",
		"params": [
			{
				"to": "%s",
				"data": "%s"
			},
			"latest"
   	 ]
	}`, address, data)
	return jsonInfoRPC(body, url)
}

func jsonInfoRPC(body, url string) (data *models.JsonRPCModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DoPost(url string, contentType string, body string) (res []byte, err error) {
	resp, err := http.Post(url, contentType, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
