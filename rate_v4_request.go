package main

import (
	"encoding/json"
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
	Revision int8
	Package  []Package
}

func (r *RateV4Request) validateRequest() error {
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

func (r *RateV4Request) requestRate() ([]byte, error) {
	marshalledXML, err := xml.Marshal(r)
	if err != nil {
		return nil, err
	}
	requestXML := url.QueryEscape(string(marshalledXML))
	requestURI := fmt.Sprintf("http://production.shippingapis.com/ShippingAPI.dll?API=RateV4&XML=%v", requestXML)

	// log.Printf("requestURI: %v", requestURI)
	res, err := http.Get(requestURI)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := RateV4Request{}
	err = xml.Unmarshal(contents, &response)
	if err != nil {
		return nil, err
	}

	output, err := json.MarshalIndent(response, "", "  ")
	return output, nil
}

type Package struct {
	// Package ID, arbitrarily defined by user
	ID string `xml:",attr"`
	//
	Service string
	// Possible FirstClassMailTypes are Letter, Flat, Parcel, Postcard, Package Service
	FirstClassMailType string `xml:",omitempty" json:",omitempty"`
	// Zip Code limitations: length=5, pattern=\d{5}.
	// I've chosen int32 because
	//      A) int16 was too small to store all possible zip codes.
	//      B) zip codes cannot be numerically negative
	ZipOrigination int32
	ZipDestination int32
	Pounds         int64
	Ounces         float64
	Container      string  `json:",omitempty"`
	Size           string  `xml:",omitempty" json:",omitempty"`
	Width          float64 `xml:",omitempty" json:",omitempty"`
	Length         float64 `xml:",omitempty" json:",omitempty"`
	Height         float64 `xml:",omitempty" json:",omitempty"`
	Girth          float64 `xml:",omitempty" json:",omitempty"`
	// Value and AmountToCollect are used to determine availability and cost of extra services
	Value             float64          `xml:",omitempty" json:",omitempty"`
	AmountToCollect   float64          `xml:",omitempty" json:",omitempty"`
	SpecialServices   []SpecialService `xml:",omitempty" json:",omitempty"`
	Content           *Content         `xml:",omitempty" json:",omitempty"`
	GroundOnly        bool             `xml:",omitempty" json:",omitempty"`
	SortBy            string           `xml:",omitempty" json:",omitempty"`
	Machinable        bool             `xml:",omitempty" json:",omitempty"`
	ReturnLocations   bool             `xml:",omitempty" json:",omitempty"`
	ReturnServiceInfo bool             `xml:",omitempty" json:",omitempty"`
	// when storing DropOffTime and ShipDate as time.Time fields, omitempty never triggers
	// so we are regrettably formatting them as strings :(
	DropOffTime string    `xml:",omitempty" json:",omitempty"`
	ShipDate    *ShipDate `xml:",omitempty" json:",omitempty"` // pattern=\d{2}-[a-zA-z]{3}-\d{4}
}

func (p Package) validatePackage() error {
	// TODO: ensure required fields are present and not empty.

	return nil
}

func (p Package) validateService() error {
	validServices := []string{"first class", "first class dommercial",
		"first class hfp commercial", "priority", "priority commercial",
		"priority cpp", "priority hfp commercial", "priority hfp cpp",
		"priority mail express", "priority mail express commercial",
		"priority mail express cpp", "priority mail express sh",
		"priority mail express sh commercial", "priority mail express hfp",
		"priority mail express hfp commercial", "hfp cpp", "standard post",
		"media", "library", "all", "online", "plus",
	}
	lowerService := strings.ToLower(p.Service)
	serviceIsValid := false
	for _, s := range validServices {
		if lowerService == s {
			serviceIsValid = true
		}
	}

	if !serviceIsValid {
		return fmt.Errorf("The assigned service for this Package is invalid: %v\n\nFor a list of valid service options, visit this URL: %v", p.Service, "https://www.usps.com/business/web-tools-apis/rate-calculator-api.htm#_Toc423593289")
	}
	return nil
}

func (p Package) validateMachinable() error {
	// Machinable is required when Service=('FIRST CLASS', 'STANDARD POST', 'ALL', or 'ONLINE') and (FirstClassMailType='LETTER' or FirstClassMailType='FLAT')
	machinableRequiredServices := map[string]string{
		"first class":   "first class",
		"standard post": "standard post",
		"all":           "all",
		"online":        "online",
	}
	machinableRequiredFirstClassMailTypes := map[string]string{
		"flat":   "flat",
		"letter": "letter",
	}
	_, serviceValid := machinableRequiredServices[strings.ToLower(p.Service)]
	_, mailTypeValid := machinableRequiredFirstClassMailTypes[strings.ToLower(p.FirstClassMailType)]

	if serviceValid && mailTypeValid {
		if &p.Machinable == nil {
			return fmt.Errorf("Machinable not set while Service and FirstClassMailType are set")
		}
	} else {
		if &p.Machinable != nil {
			return fmt.Errorf("Machinable set despite having invalid Service and FirstClassMailType values")
		}
	}
	return nil
}

func (p *Package) setSpecialServices(s []int16) {
	specialServices := map[int16]string{
		100: "Insurance",
		101: "Insurance – Priority Mail",
		102: "Return Receipt",
		103: "Collect on Delivery",
		104: "Certificate of Mailing (Form 3665)",
		105: "Certified Mail",
		106: "USPS Tracking",
		107: "Return Receipt for Merchandise",
		108: "Signature Confirmation",
		109: "Registered Mail",
		110: "Return Receipt Electronic",
		112: "Registered mail COD collection Charge",
		118: "Return Receipt – Priority Mail Express",
		119: "Adult Signature Required",
		120: "Adult Signature Restricted Delivery",
		125: "Insurance – Priority Mail Express",
		156: "Signature Confirmation Electronic",
		160: "Certificate of Mailing (Form 3817)",
		161: "Priority Mail Express 1030 AM Delivery",
		170: "Certified Mail Restricted Delivery",
		171: "Certified Mail Adult Signature Required",
		172: "Certified Mail Adult Signature Restricted Delivery",
		173: "Signature Confirm. Restrict. Delivery",
		174: "Signature Confirmation Electronic Restricted Delivery",
		175: "Collect on Delivery Restricted Delivery",
		176: "Registered Mail Restricted Delivery",
		177: "Insurance Restricted Delivery",
		178: "Insurance Restrict.  Delivery – Priority Mail",
		179: "Insurance Restrict. Delivery – Priority Mail Express",
		180: "Insurance Restrict. Delivery (Bulk Only)",
	}

	for _, x := range s {
		if _, ok := specialServices[x]; ok {
			p.SpecialServices = append(p.SpecialServices, SpecialService{SpecialService: x})
		} else {
			log.Printf("Erroneous input! %v is not a valid special service ID", x)
		}
	}
}

func (p Package) validateZipCodes() error {
	if p.ZipOrigination < 0 || p.ZipOrigination > 99999 {
		return fmt.Errorf("ZipOrigination is invalid")
	}
	if p.ZipDestination < 0 || p.ZipDestination > 99999 {
		return fmt.Errorf("ZipDestination is invalid")
	}
	return nil
}

type ShipDate struct {
	Date   string `xml:",innerxml"`
	Option string `xml:",attr,omitempty" json:",omitempty"`
}

func (s ShipDate) validateShipDateOption() error {
	if s.Option != "" {
		if s.Option != "HFP" && s.Option != "PEMSH" {
			return fmt.Errorf("Invalid ship date option provided: %v\n\nValid values are either 'PEMSH' or 'HFP'", s.Option)
		}
	}
	return nil
}

type SpecialService struct {
	SpecialService int16
}

type Content struct {
	ContentType        string `xml:",omitempty" json:",omitempty"`
	ContentDescription string `xml:",omitempty" json:",omitempty"`
}
