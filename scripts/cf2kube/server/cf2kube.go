package main

import (
	"fmt"
	"net/http"
	"log"
	"os/exec"
	"os"
	"path"
)

// default runScript, set os.Env("CF2KUBE") will change this
var runScript = "/home/cfproxy/cf2kube/droplet2docker.sh"

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func dispatchToKube(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method Not Allowed")
		return
	}
	r.ParseForm()
	log.Println("form:", r.Form)
	log.Println("url:", r.URL)
	cacheKey := r.Form.Get("cache_key")
	hashMethod := r.Form.Get("hash_method")
	hashValue := r.Form.Get("hash_value")
	extractTo := r.Form.Get("extract_to")
	downloadUrl := r.Form.Get("download_url")
	startCmd := r.Form.Get("start_cmd")
	if len(cacheKey) == 0 || len(hashMethod) == 0 || len(hashValue) == 0 || len(extractTo) == 0 || len(downloadUrl) == 0 || len(startCmd) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "ERROR: params not satisfied.")
	} else if startCmd[0:2] == "sh" || startCmd[0:2] == "py" {	// static web & python web project supported
		cmd := exec.Command(runScript, cacheKey, hashMethod, hashValue, extractTo, downloadUrl, startCmd)
		log.Println("cmd", cmd)
		out, err := cmd.Output()
		if err != nil {
			log.Println("exec cmd error", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: exec cmd error", err)
		} else {
			log.Println("exec cmd success")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "Success: exec cmd success")
		}
		log.Println("exec cmd output", string(out))
		fmt.Fprint(w, "Output: exec cmd output", string(out))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "ERROR: app type not support(bad start_cmd).")
	}
}

func main() {
	CF2KUBE := os.Getenv("CF2KUBE")
	if CF2KUBE != "" {
		runScript = path.Join(CF2KUBE, "droplet2docker.sh")
	}
	if !PathExist(runScript) {
		log.Println("Error: CF2KUBE env is not define, and default script doesn't exist")
		return
	}
	http.HandleFunc("/droplet", dispatchToKube)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}