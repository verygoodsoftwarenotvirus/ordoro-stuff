package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type RateV4Response struct {
	UserID  string
	Package []Package
}

func (r *RateV4Response) validateRequest() error {
	errors := []string{}
	if len(r.Package) > 25 {
		errors = append(errors, fmt.Sprintf("Too many Packages included in payload. Max is 25, but %v are present", len(r.Package)))
	}

	if r.UserID == "" {
		errors = append(errors, "No User ID provided to request")
	}

	if len(errors) > 0 {
		return fmt.Errorf("The following errors occurred validating your request:\n\n%v", strings.Join(errors, "\n"))
	}
	return nil
}

func (r *RateV4Response) requestRate() ([]byte, error) {
	marshalledXML, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}
	requestXML := url.QueryEscape(string(marshalledXML))
	requestURI := fmt.Sprintf("http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML=%v", requestXML)

	log.Printf("requestURI: %v", requestURI)
	res, err := http.Get(requestURI)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
