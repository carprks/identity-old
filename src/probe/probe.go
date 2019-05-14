package probe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTP responds with a healthy model
func HTTP(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)
	if len(buf) >= 1 {
		fmt.Println(fmt.Sprintf("probe request: %v", string(buf)))
	}

	// get probe response
	resp, _ := Probe()

	// send status
	j, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/health+json")
	w.Header().Set("Service", "Identity")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(j)
	if err != nil {
		fmt.Println(fmt.Errorf("write failed: %v", err))
	}

	return
}

// Probe responds with a healthy model
func Probe() (Healthy, error){
	return Healthy{
		Status: "pass",
	}, nil
}