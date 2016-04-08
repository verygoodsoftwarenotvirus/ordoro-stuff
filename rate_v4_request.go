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

type RateV4Request struct {
	UserID   string `xml:"USERID,attr"`
	Revision int
	Package  []Package
}

func (r *RateV4Request) validateRequest() error {
	// TODO: flesh out more.
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

func (r *RateV4Request) requestRate() (RateV4Response, error) {
	var response RateV4Response
	marshalledXML, err := xml.Marshal(r)
	if err != nil {
		log.Printf("Error encountered marshalling the request struct:\n\t%v", err)
		return response, err
	}
	requestXML := url.QueryEscape(string(marshalledXML))
	requestURI := fmt.Sprintf("http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML=%v", requestXML)

	// log.Printf("requestURI: %v", requestURI)
	res, err := http.Get(requestURI)
	if err != nil {
		log.Printf("Error encountered hitting the API:\n\t%v", err)
		return response, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response from server:\n\t%v", err)
		return response, err
	}

	err = xml.Unmarshal(contents, &response)
	if err != nil {
		log.Printf("Error unmarshalling the API's response:\n\t%v", err)
		return response, err
	}

	if err != nil {
		log.Printf("Error marshalling the response to JSON:\n\t%v", err)
		return response, err
	}
	return response, nil
}
