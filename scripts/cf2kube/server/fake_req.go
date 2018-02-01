package main

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
)

func main() {
	proxyUrl := "http://cfproxy:8080/droplet"
	params := url.Values{"_cache_key": {"_cache_key"}, "_hash_method": {"_hash_method"}, "_hash_value": {"_hash_value"},
		"_extract_to": {"_extract_to"}, "_download_url": {"_download_url"}, "_start_cmd": {"_start_cmd"}}
	resp, err := http.PostForm(proxyUrl, params)
	if err != nil {
		fmt.Println("invokeCf2Kube invoke error", err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("invokeCf2Kube response error", err)
		} else {
			fmt.Println("invokeCf2Kube response body", string(body))
		}
	}
}