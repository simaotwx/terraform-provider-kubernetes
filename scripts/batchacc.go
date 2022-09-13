package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"unicode"
)

func main() {
	cmd := exec.Command("./kubernetes.test", "-test.list", "^TestAccKubernetes")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	sanout := strings.TrimSuffix(string(out), "\n")
	tl := strings.Split(sanout, "\n")

	pt := PrefixTree{}
	for k := range tl {
		pt.addString(tl[k])
	}

	pfx := pt.prefixesToDepth(3)

	// print out json
	js := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(js)
	enc.Encode(pfx)
	fmt.Print(js.String())
}

func tokenizeCamelCase(s string) (out []string) {
	for len(s) > 0 {
		nx := 0
		for i := 1; i < len(s); i++ {
			nx = strings.IndexFunc(s[i:], func(r rune) bool {
				return unicode.IsUpper(r) || r == '_'
			})
			if nx == -1 {
				out = append(out, strings.TrimPrefix(s, "_"))
				return
			}
			if nx != 0 {
				nx = nx + i
				break
			}
			if i == len(s)-1 {
				nx = len(s)
			}
		}
		out = append(out, strings.TrimPrefix(s[:nx], "_"))
		s = s[nx:]
	}
	return
}

type PrefixTree map[string]PrefixTree

func (pt PrefixTree) addString(s string) {
	pt.addTokenized(tokenizeCamelCase(s))
}

func (pt PrefixTree) addTokenized(s []string) {
	if len(s) < 1 {
		return
	}
	var st PrefixTree
	st, ok := pt[s[0]]
	if !ok {
		st = make(PrefixTree)
	}
	st.addTokenized(s[1:])
	pt[s[0]] = st
}

func (pt PrefixTree) prefixesToDepth(d int) (px []string) {
	if d == 0 {
		for k := range pt {
			px = append(px, k)
		}
		return
	}
	for k := range pt {
		sfx := pt[k].prefixesToDepth(d - 1)
		if len(sfx) == 0 {
			px = append(px, k)
		}
		for _, s := range sfx {
			px = append(px, k+s)
		}
	}
	return
}
