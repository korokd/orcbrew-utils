package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cespare/goclj/parse"
	"github.com/gopherjs/gopherjs/js"
)

var rawOutput = flag.Bool("raw", false, "Don't pretty-print JSON output")
var noSave = flag.Bool("nosave", false, "Don't save the JSON output")

// func main() {
// 	flag.Parse()

// 	args := flag.Args()
// 	if len(args) != 1 {
// 		printUsage()
// 		os.Exit(2)
// 	}

// 	filename := args[0]
// 	Convert(filename, false)
// }

func main() {
	js.Module.Get("exports").Set("Convert", map[string]interface{}{
		"Convert": Convert(),
	})
}

// Convert exposes the `convert()` function
// as a JavaScript Function
func Convert() *js.Object {
	return js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		filename := arguments[0].String()
		return convert(filename, true)
	})
}

// func MakeFunc(fn func(this *Object, arguments []*Object) interface{}) *Object

func convert(filename string, shouldJustReturn bool) string {
	contentsBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed when reading file %s: %s", filename, err)
		os.Exit(2)
	}

	// Remove BOM if its at the start of the file
	contentsBytes = bytes.TrimLeft(contentsBytes, "\xef\xbb\xbf")

	contents := string(contentsBytes)

	contents = regexp.MustCompile(`#:orcpub.dnd.e5{`).ReplaceAllString(contents, "{")
	contents = regexp.MustCompile(`#:orcpub.dnd.e5.character{`).ReplaceAllString(contents, "{")
	contents = regexp.MustCompile(`#:orcpub.dnd.e5.character`).ReplaceAllString(contents, "#")
	contents = regexp.MustCompile(`:orcpub.dnd.e5.character/`).ReplaceAllString(contents, "")
	contents = regexp.MustCompile(`orcpub.dnd.e5/`).ReplaceAllString(contents, "")
	contents = regexp.MustCompile(`orcpub.dnd.e5.[a-z-]*/`).ReplaceAllString(contents, "")

	buf := bytes.NewBufferString(contents)

	tt, err := parse.Reader(buf, "input.clj", 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s as Clojure: %s", filename, err)
		os.Exit(2)
	}

	jsonString := treeToJSON(tt)

	if shouldJustReturn == false {
		saveJSONFile(filename, jsonString)
	}
	fmt.Fprintf(os.Stdout, jsonString)
	return jsonString
}

func saveJSONFile(filename string, jsonString string) {
	if *noSave == false {
		fName := strings.TrimSuffix(filename, filepath.Ext(filename))
		fmt.Fprintf(os.Stdout, fmt.Sprintf("Saved to %s.json", fName))
		ioutil.WriteFile(fmt.Sprintf("%s.json", fName), []byte(jsonString), 0)
	} else {
		if *rawOutput {
			fmt.Fprintf(os.Stdout, jsonString)
		} else {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(jsonString), "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing JSON: %s\n%s", err, jsonString)

				os.Exit(2)
			}

			fmt.Fprint(os.Stdout, prettyJSON.String())
		}
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] inputFile\n", os.Args[0])
	flag.PrintDefaults()
}

func treeToJSON(tree *parse.Tree) string {
	return nodesToJSON(tree.Roots, 0)
}

func nodesToJSON(nodes []parse.Node, depth int) string {
	var buf bytes.Buffer
	for _, node := range nodes {
		buf.WriteString(strings.Repeat("  ", depth))
		buf.WriteString(nodeToJSON(node))
		buf.WriteString("\n")
	}

	return buf.String()
}

func nodeToJSON(node parse.Node) string {
	switch v := node.(type) {

	case *parse.KeywordNode:
		return fmt.Sprintf(`"%s"`, v.Val[1:])
	case *parse.StringNode:
		return fmt.Sprintf(`"%s"`, v.Val)
	case *parse.NumberNode:
		return v.Val
	case *parse.SymbolNode:
		return fmt.Sprintf(`"%s"`, v.Val)
	case *parse.SetNode:
		var vals []string
		for _, node := range filterNonValueNodes(v.Children()) {
			vals = append(vals, nodeToJSON(node))
		}

		return fmt.Sprintf("[%s]", strings.Join(vals, ","))
	case *parse.ListNode:
		var vals []string
		for _, node := range filterNonValueNodes(v.Children()) {
			vals = append(vals, nodeToJSON(node))
		}

		return fmt.Sprintf("[%s]", strings.Join(vals, ","))
	case *parse.VectorNode:
		var vals []string
		for _, node := range filterNonValueNodes(v.Children()) {
			vals = append(vals, nodeToJSON(node))
		}

		return fmt.Sprintf("[%s]", strings.Join(vals, ","))
	case *parse.MapNode:
		var keys []parse.Node
		var vals []parse.Node

		children := v.Children()
		children = filterNonValueNodes(children)

		if (len(children) % 2) != 0 {
			panic(v.String())
		}
		for idx, node := range children {
			if idx == 0 || (idx%2) == 0 {
				keys = append(keys, node)
			} else {
				// The value might be empty, remove the key in that case
				_, isNilValue := node.(*parse.NilNode)
				if isNilValue {
					keys = keys[0 : len(keys)-1]
				} else {
					vals = append(vals, node)
				}
			}
		}

		var entries []string
		for idx, key := range keys {
			var val = vals[idx]
			var keyString = nodeToJSON(key)
			if keyString[0] != '"' {
				keyString = fmt.Sprintf(`"%s"`, keyString)
			}
			entries = append(entries, fmt.Sprintf("%s: %s", keyString, nodeToJSON(val)))
		}

		return fmt.Sprintf("{%s}", strings.Join(entries, ","))
	case *parse.NewlineNode:
		return ""
	default:
		return node.String()
	}
}

func filterNonValueNodes(nodes []parse.Node) []parse.Node {
	var result []parse.Node

	for _, node := range nodes {
		switch v := node.(type) {
		case *parse.NewlineNode:
			break
		case *parse.CommentNode:
			break
		default:
			result = append(result, v)
		}
	}

	return result
}
