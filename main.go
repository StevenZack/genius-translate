package main

import (
	"encoding/xml"
	"log"
	"strconv"
	"strings"

	. "github.com/gofaith/faithtop"
)

type Resources struct {
	XMLName xml.Name `xml:"resources"`
	String  []String `xml:"string"`
}
type String struct {
	XMLName xml.Name `xml:"string"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

var w IWindow
var resources Resources
var fullMode bool

func main() {
	app := NewApp()

	var in, enc, mid, dec IPlainTextArea
	Window().Assign(&w).DeferShow().Title("Genius Translate Tool").CenterWidget(VBox(
		Label2("Encode for <a href='https://translate.google.com/'>Google translate</a>").Interaction(TextBrowserInteraction).OpenExternalLinks(true),
		PlainTextArea().Assign(&in).MaxHeight(40),
		Button2("Encode").OnClick(func() {
			enc.PlainText(doEncode(in.GetText()))
		}),
		PlainTextArea().Assign(&enc).MaxHeight(40),
		Label2("Decode from Google translate"),
		PlainTextArea().Assign(&mid).MaxHeight(40),
		Button2("Decode").OnClick(func() {
			dec.PlainText(doDecode(mid.GetText()))
		}),
		PlainTextArea().Assign(&dec).MaxHeight(40),
	))

	app.Run()
}

func doEncode(s string) string {
	fullMode = strings.HasPrefix(s, "<resources>")
	if !fullMode {
		s = "<resources>" + s + "</resources>"
	}
	resources = Resources{}
	e := xml.Unmarshal([]byte(s), &resources)
	if e != nil {
		log.Println(e)
		return e.Error()
	}

	return resources.Encode()
}

func doDecode(s string) string {
	s = strings.TrimSuffix(s, "|")
	ss := strings.Split(s, "|")
	if len(ss) != len(resources.String) {
		MessageBox_Info(w, "Decode error", "string array length mismatch: "+strconv.Itoa(len(ss))+" != "+strconv.Itoa(len(resources.String)), StandardButton_Ok, StandardButton_Ok)
	}
	for i, str := range ss {
		if len(resources.String) > i {
			resources.String[i].Value = str
		}
	}

	if fullMode {
		b, e := xml.MarshalIndent(resources, "", "\t")
		if e != nil {
			return e.Error()
		}
		return string(b)
	}
	b, e := xml.MarshalIndent(resources.String, "", "\t")
	if e != nil {
		return e.Error()
	}
	return string(b)
}

func (r *Resources) Encode() string {
	builder := new(strings.Builder)
	for _, s := range r.String {
		builder.WriteString(s.Value)
		builder.WriteString(" | ")
	}
	return builder.String()
}
