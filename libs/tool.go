/*
Package butterfly : The library for butterfly

Copyright(c) 2020 Star Inc. All Rights Reserved.
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/
package butterfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HTTPGet :
func HTTPGet(url string, recovery int) string {
	var output string
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	DeBug("GenRequest", err)
	req.Header.Set("User-Agent", Config.UserAgent)
	resp, err := client.Do(req)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		DeBug("ReadHTML", err)
		output = string(body)
	} else {
		if recovery == 0 {
			time.Sleep(time.Duration(2) * time.Second)
			HTTPGet(url, 1)
		} else {
			DeBug("GetHTTP", err)
		}
	}
	resp.Body.Close()
	return output
}

// DeBug : Print errors for debug and report
func DeBug(where string, err error) bool {
	if err != nil {
		fmt.Printf("Butterfly Error #%s\nReason:\n%s\n\n", where, err)
		return false
	}
	return true
}

// ReplaceSyntaxs : Remove space and syntax
func ReplaceSyntaxs(rawString string, filled string) string {
	var output bytes.Buffer
	rawString = strings.ReplaceAll(rawString, " ", "\x1e")
	rawString = strings.ReplaceAll(rawString, "\t", "\x1e")
	rawString = strings.ReplaceAll(rawString, "\n", "\x1e")
	stringSlice := strings.Split(rawString, "\x1e")
	for _, word := range stringSlice {
		if word != "" {
			output.WriteString(word + filled)
		}
	}
	return output.String()
}

// RemoveChildNode : Remove all child html node selected
func RemoveChildNode(rootNode *html.Node, removeNode *html.Node) {
	foundNode := false
	checkNode := make(map[int]*html.Node)

	i := 0
	for n := rootNode.FirstChild; n != nil; n = n.NextSibling {
		if n == removeNode {
			foundNode = true
			n.Parent.RemoveChild(n)
		}

		checkNode[i] = n
		i++
	}

	if !foundNode {
		for _, item := range checkNode {
			RemoveChildNode(item, removeNode)
		}
	}
}

// FindInSlice : Find out an item if exists in a slice
func FindInSlice(slice interface{}, value interface{}) (int, bool) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("SliceExists() given a non-slice type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == value {
			return i, true
		}
	}
	return -1, false
}

// NormalizeURI : Reformat a URI as the unique standard
func NormalizeURI(URI string) (string, *url.URL) {
	handleURI, _ := url.Parse(URI)

	if handleURI.Scheme == "" {
		handleURI.Scheme = "http"
	}

	return handleURI.String(), handleURI
}

// CallTextEditor : To call a text editor
func CallTextEditor(filePath string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = func() string {
			editors := []string{"vim", "vi", "emacs", "nano"}
			for _, editor := range editors {
				if r, _ := exec.LookPath(editor); r != "" {
					return editor
				}
			}
			return ""
		}()
		if editor == "" {
			fmt.Println("No text editor found.")
			return
		}
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
