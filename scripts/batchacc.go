package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("./kubernetes.test", "-test.list", "^TestAccKubernetes")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	sanout := strings.TrimSuffix(string(out), "\n")
	tl := strings.Split(sanout, "\n")

	// print out json
	js := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(js)
	enc.Encode(tl)
	fmt.Print(js.String())
}
