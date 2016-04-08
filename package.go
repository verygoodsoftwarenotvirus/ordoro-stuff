package main

import (
	"fmt"
	"strings"
)

type Package struct {
	// Package ID, arbitrarily defined by user
	ID                 string `xml:",attr"`
	Service            string
	FirstClassMailType string `xml:",omitempty" json:",omitempty"`
	ZipOrigination     string
	ZipDestination     string
	Pounds             int
	Ounces             float64
	Container          string           `json:",omitempty"`
	Size               string           `xml:",omitempty" json:",omitempty"`
	Width              float64          `xml:",omitempty" json:",omitempty"`
	Length             float64          `xml:",omitempty" json:",omitempty"`
	Height             float64          `xml:",omitempty" json:",omitempty"`
	Girth              float64          `xml:",omitempty" json:",omitempty"`
	Value              float64          `xml:",omitempty" json:",omitempty"`
	AmountToCollect    float64          `xml:",omitempty" json:",omitempty"`
	SpecialServices    []SpecialService `xml:",omitempty" json:",omitempty"`
	Content            *Content         `xml:",omitempty" json:",omitempty"`
	GroundOnly         bool             `xml:",omitempty" json:",omitempty"`
	SortBy             string           `xml:",omitempty" json:",omitempty"`
	Machinable         bool             `xml:",omitempty" json:",omitempty"`
	Zone               int              `xml:",omitempty" json:",omitempty"`
	ReturnLocations    bool             `xml:",omitempty" json:",omitempty"`
	ReturnServiceInfo  bool             `xml:",omitempty" json:",omitempty"`
	// when storing DropOffTime and ShipDate as time.Time fields, omitempty never triggers
	// so we are regrettably formatting them as strings :(
	DropOffTime string    `xml:",omitempty" json:",omitempty"`
	ShipDate    *ShipDate `xml:",omitempty" json:",omitempty"` // pattern=\d{2}-[a-zA-z]{3}-\d{4}
	// These fields are necessary for RateResponse Packages only
	Postage         []Postage     `xml:",omitempty" json:",omitempty"`
	Restriction     *Restriction  `xml:",omitempty" json:",omitempty"`
	RatePriceType   string        `xml:",omitempty" json:",omitempty"`
	RatePaymentType string        `xml:",omitempty" json:",omitempty"`
	Error           *PackageError `xml:",omitempty" json:",omitempty"`
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
