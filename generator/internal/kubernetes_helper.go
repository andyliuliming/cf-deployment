package internal

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/executor"
	"net/url"
	"net/http"
	"io/ioutil"
)

var proxyUrl = "http://cfproxy:8080/droplet"
var logAction = "kubernetes_helper"

func K8sRunContainer(logger lager.Logger, req *executor.RunRequest) {
	logger.Info(logAction, lager.Data{"request": req})
	if req.RunInfo.Setup == nil {
		logger.Info(logAction, lager.Data{"run for stage": "Skip STAGE"})
		return
	}
	downloadAction := req.RunInfo.Setup.SerialAction.Actions[0].DownloadAction
	runAction := req.RunInfo.Action.CodependentAction.Actions[0].RunAction
	logger.Info(logAction, lager.Data{"downloadAction": downloadAction})
	invokeCf2Kube(logger, downloadAction.CacheKey, downloadAction.ChecksumAlgorithm, downloadAction.ChecksumValue, downloadAction.To, downloadAction.From, runAction.Args[1])
}

func invokeCf2Kube(logger lager.Logger, cache_key, hash_method, hash_value, extract_to, download_url, start_cmd string) {
	params := url.Values{"cache_key": {cache_key}, "hash_method": {hash_method}, "hash_value": {hash_value}, "extract_to": {extract_to}, "download_url": {download_url}, "start_cmd": {start_cmd}}
	logger.Info(logAction, lager.Data{"invokeCf2Kube params": params})
	resp, err := http.PostForm(proxyUrl, params)
	if err != nil {
		logger.Info(logAction, lager.Data{"invokeCf2Kube invoke error": err})
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		logger.Info(logAction, lager.Data{"invokeCf2Kube response status": resp.StatusCode})
		if err != nil {
			logger.Info(logAction, lager.Data{"invokeCf2Kube response error": err})
		} else {
			logger.Info(logAction, lager.Data{"invokeCf2Kube response body": string(body)})
		}
	}
}
