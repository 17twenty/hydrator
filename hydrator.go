package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"time"
)

// The usage of variable length args is to prevent
// missing fields from breaking the templates
func formatAsDollars(valueInCents ...int) (string, error) {
	if len(valueInCents) != 1 {
		return "$0", nil
	}
	dollars := valueInCents[0] / 100
	cents := valueInCents[0] % 100
	return fmt.Sprintf("$%d.%2d", dollars, cents), nil
}

func formatAsDate(t ...string) string {
	if len(t) != 1 {
		return ""
	}
	// Mon Jan 2 15:04:05 MST 2006
	d, _ := time.Parse("2/1/2006", t[0])
	// log.Println(err)
	year, month, day := d.Date()
	return fmt.Sprintf("%d/%d/%d", day, month, year)
}

func titleCase(input ...string) string {
	if len(input) != 1 {
		return ""
	}
	return strings.Title(strings.ToLower(input[0]))
}

// UnstructuredJSON is a simple cover for map[string]interface{}
type UnstructuredJSON map[string]interface{}

var (
	tpl = flag.String("t", "./transform.tmpl", "the template to use to transform the data")
	src = flag.String("s", "./source.json", "source JSON")
	out = flag.String("w", "", "Specifies an output file if present")
)

func main() {
	flag.Parse()

	filebytes, err := ioutil.ReadFile(*src)
	if err != nil {
		log.Fatalln("Couldn't open file '", *src, "' -", err)
		return
	}

	jsonBlob := UnstructuredJSON{}
	json.Unmarshal(filebytes, &jsonBlob)

	// Wire up our FuncMap
	// TODO: Dynamically load these using plugins!
	fmap := template.FuncMap{
		"formatAsDollars": formatAsDollars,
		"formatAsDate":    formatAsDate,
		"titleCase":       titleCase,
	}
	t := template.Must(template.New(path.Base(*tpl)).Funcs(fmap).ParseFiles(*tpl))

	// Exexute the template, on our source.json and output to a holding buffer
	outputBuffer := bytes.NewBuffer(nil)
	err = t.Execute(outputBuffer, jsonBlob)
	if err != nil {
		log.Fatalln("Execute:", err)
	}

	// Cleanup the JSON
	tmp := UnstructuredJSON{}
	err = json.Unmarshal(outputBuffer.Bytes(), &tmp)
	if err != nil {
		log.Fatalln("Unmarshal:", err, "-", outputBuffer.String())
	}

	indentedBytes, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		log.Fatalln("MarshalIndent:", err)
	}
	// It's lumpy to marshal and unmarshal but it does make it pretty
	if *out != "" {
		if err = ioutil.WriteFile(*out, indentedBytes, 0644); err != nil {
			log.Fatalln("WriteFile:", err)
		}
	} else {
		fmt.Printf("%s\n", string(indentedBytes))
	}
}
