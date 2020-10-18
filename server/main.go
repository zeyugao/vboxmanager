package main

import (
	"../common"
	"fmt"
	"net/http"
	"os/exec"
)

const vboxmanage = "vboxmanage.exe"

func listVM(w http.ResponseWriter, req *http.Request) {
	var out []byte
	var err error

	running, ok := common.GetParam(req, "running")

	var subcommand string
	if ok && (running != "1") {
		subcommand = "vms"
	} else {
		subcommand = "runningvms"
	}
	out, err = exec.Command(vboxmanage, "list", subcommand).Output()
	if err == nil {
		_, _ = fmt.Fprint(w, out)
	} else {
		_, _ = fmt.Fprintln(w, err)
	}
}

func startVM(w http.ResponseWriter, req *http.Request) {
	uuid, ok := common.GetParam(req, "uuid")

	if !ok {
		_, _ = fmt.Fprintln(w, "uuid required")
		return
	}
	if out, err := exec.Command(vboxmanage, "snapshot", "restorecurrent").Output(); err == nil {
		_, _ = fmt.Fprintln(w, out)
	} else {
		_, _ = fmt.Fprint(w, err)
		return
	}

	if out, err := exec.Command(vboxmanage, "startvm", "--headless", uuid).Output(); err == nil {
		_, _ = fmt.Fprintln(w, out)
	} else {
		_, _ = fmt.Fprint(w, err)
	}
}

func stopVM(w http.ResponseWriter, req *http.Request) {
	uuid, ok := common.GetParam(req, "uuid")

	if !ok {
		_, _ = fmt.Fprintln(w, "uuid required")
		return
	}

	if out, err := exec.Command(vboxmanage, "controlvm", uuid, "poweroff").Output(); err == nil {
		_, _ = fmt.Fprintln(w, out)
	} else {
		_, _ = fmt.Fprint(w, err)
	}
}

func main() {
	http.HandleFunc("/list_vm", listVM)
	http.HandleFunc("/start_vm", startVM)
	http.HandleFunc("/stop_vm", stopVM)
	_ = http.ListenAndServe(":2333", nil)
}
