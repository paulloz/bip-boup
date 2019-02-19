package httpreq

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func HTTPGetAsXML(url string) (*xml.XmlDocument, error) {
	body, err := HTTPGet(url)
	if err != nil {
		return nil, err
	}

	return gokogiri.ParseXml(body)
}

func HTTPGetAsHTML(url string) (*html.HtmlDocument, error) {
	body, err := HTTPGet(url)
	if err != nil {
		return nil, err
	}

	return gokogiri.ParseHtml(body)
}

func HTTPPost(url string, params string) ([]byte, error) {
	body := strings.NewReader(params)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func HTTPPostAsHTML(url string, params string) (*html.HtmlDocument, error) {
	body, err := HTTPPost(url, params)
	if err != nil {
		return nil, err
	}

	return gokogiri.ParseHtml(body)
}
