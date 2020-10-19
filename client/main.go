package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
import "../common"

func parseParams(req *http.Request) (string, string, bool) {
	host, ok := common.GetParam(req, "host")

	if !ok {
		return "", "", false
	}

	uuid, ok := common.GetParam(req, "uuid")

	if !ok {
		return "", "", false
	}
	ip := "172.22.0." + host
	return ip, uuid, true

}
func startVM(w http.ResponseWriter, req *http.Request) {
	ip, uuid, ok := parseParams(req)
	if !ok {
		_, _ = fmt.Fprintln(w, "host and uuid required")
		return
	}
	resp, err := http.Get("http://" + ip + ":2333/start_vm?uuid=" + uuid)
	if err == nil {
		_, _ = fmt.Fprintln(w, resp)
	} else {
		_, _ = fmt.Fprintln(w, err)
	}
}

func stopVM(w http.ResponseWriter, req *http.Request) {
	ip, uuid, ok := parseParams(req)
	if !ok {
		_, _ = fmt.Fprintln(w, "host and uuid required")
		return
	}

	resp, err := http.Get("http://" + ip + ":2333/stop_vm?uuid=" + uuid)
	if err == nil {
		_, _ = fmt.Fprintln(w, resp)
	} else {
		_, _ = fmt.Fprintln(w, err)
	}
}

type VM struct {
	name    string
	uuid    string
	running bool
}

type indexPageData struct {
	vms []VM
}

func getVMs(ip string) []VM {
	r, _ := regexp.Compile("\"(.*)\" {([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})}")
	resp, err := http.Get("http://" + ip + ":2333/list_vm?running=1")

	if err != nil {
		fmt.Println(err)
		return []VM{}
	}
	s, _ := ioutil.ReadAll(resp.Body)

	lines := strings.Split(string(s), "\n")

	runningVmUuid := make(map[string]bool)

	for _, line := range lines {
		res := r.FindAllString(line, -1)
		if len(res) < 2 {
			continue
		}
		runningVmUuid[res[1]] = true
	}

	resp, err = http.Get("http://" + ip + ":2333/list_vm")
	if err != nil {
		fmt.Println(err)
		return []VM{}
	}
	s, _ = ioutil.ReadAll(resp.Body)

	lines = strings.Split(string(s), "\n")
	var data []VM

	for _, line := range lines {
		res := r.FindAllString(line, -1)
		if len(res) < 2 {
			continue
		}
		data = append(data, VM{
			name:    strings.Trim(res[0], " "),
			uuid:    res[1],
			running: runningVmUuid[res[1]],
		})
	}
	return data
}

func index(w http.ResponseWriter, req *http.Request) {
	host, ok := common.GetParam(req, "host")
	ip := "172.22.0." + host
	if !ok {
		host = "1"
	}
	tmpl, err := template.ParseFiles("client.html")
	if err != nil {
		_, _ = fmt.Fprintln(w, err)
		return
	}

	resp, err := http.Get("http://" + ip + ":2333/list_vm")
	if err != nil {
		_, _ = fmt.Fprintln(w, err)
		return
	}
	s, _ := ioutil.ReadAll(resp.Body)

	lines := strings.Split(string(s), "\n")
	var data []VM

	r, _ := regexp.Compile("\"(.*)\" {([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})}")
	for _, line := range lines {
		res := r.FindAllString(line, -1)
		if len(res) < 2 {
			continue
		}
		data = append(data, VM{
			name: strings.Trim(res[0], " "),
			uuid: res[1],
		})
	}
	d := indexPageData{
		vms: data,
	}
	_ = tmpl.Execute(w, d)
}

func main() {
	http.HandleFunc("/start_vm", startVM)
	http.HandleFunc("/stop_vm", stopVM)
	http.HandleFunc("/", index)
	_ = http.ListenAndServe(":2333", nil)
}
