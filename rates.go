package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
)

var specialServices map[int16]string

func init() {
	specialServices = map[int16]string{
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
}

type RateV4Request struct {
	USERID   string `xml:",attr"`
	Revision int8
	Package  []Package
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
	// Zip Code the package in question starts at
	ZipOrigination int32
	// Zip Code the package in question ends up in
	ZipDestination int32
	// Parcel dimensions
	Pounds    int64
	Ounces    float64
	Container string  `json:",omitempty"`
	Size      string  `xml:",omitempty" json:",omitempty"`
	Width     float64 `xml:",omitempty" json:",omitempty"`
	Length    float64 `xml:",omitempty" json:",omitempty"`
	Height    float64 `xml:",omitempty" json:",omitempty"`
	Girth     float64 `xml:",omitempty" json:",omitempty"`
	// Value and AmountToCollect are used to determine availability and cost of extra services
	// Package value
	Value float64 `xml:",omitempty" json:",omitempty"`
	// Collect on delivery amount
	AmountToCollect float64          `xml:",omitempty" json:",omitempty"`
	SpecialServices []SpecialService `xml:",omitempty" json:",omitempty"`
	Content         *Content         `xml:",omitempty" json:",omitempty"`
	GroundOnly      bool             `xml:",omitempty" json:",omitempty"`
	SortBy          string           `xml:",omitempty" json:",omitempty"`
	// Machinable is required when Service=('FIRST CLASS', 'STANDARD POST', 'ALL', or 'ONLINE') and (FirstClassMailType='LETTER' or FirstClassMailType='FLAT')
	Machinable bool `xml:",omitempty" json:",omitempty"`
	// Include Dropoff Locations in Response if available
	ReturnLocations bool `xml:",omitempty" json:",omitempty"`
	// Include mail service specific information in Response if available
	ReturnServiceInfo bool `xml:",omitempty" json:",omitempty"`
	// when storing DropOffTime and ShipDate as time.Time fields, omitempty doesn't work properly
	// so we are regrettably formatting them as strings :(
	DropOffTime string `xml:",omitempty" json:",omitempty"`
	ShipDate    string `xml:",omitempty" json:",omitempty"` // pattern=\d{2}-[a-zA-z]{3}-\d{4}
	// The value of this attribute specifies how the RateV4Response will structure the Priority Express Mail Commitment data elements.
	// default=PEMSH, other option is HFP
	ShipDateOption string `xml:",omitempty" json:",omitempty"`
}

type SpecialService struct {
	SpecialService int16
}

type Content struct {
	ContentType        string `xml:",omitempty" json:",omitempty"`
	ContentDescription string `xml:",omitempty" json:",omitempty"`
}

func (p *Package) setSpecialServices(s []int16) {
	for _, x := range s {
		if _, ok := specialServices[x]; ok {
			p.SpecialServices = append(p.SpecialServices, SpecialService{SpecialService: x})
		} else {
			log.Fatalf("Error encountered! %v is not a valid special service ID", x)
		}
	}
}

func (p *Package) setZipOrigin(zc int32) {
	if zc < 0 || zc > 99999 {
		log.Fatal("Zip code is invalid")
	} else {
		p.ZipOrigination = zc
	}
}

func (p *Package) setZipDestination(zc int32) {
	if zc < 0 || zc > 99999 {
		log.Fatal("Zip code is invalid")
	} else {
		p.ZipDestination = zc
	}
}

func main() {
	package1 := Package{
		ID:                 "1ST",
		Service:            "First Class",
		FirstClassMailType: "Letter",
		ZipOrigination:     44106,
		ZipDestination:     20770,
		Pounds:             0,
		Ounces:             3.5,
		Size:               "Regular",
		Machinable:         true,
	}

	package2 := Package{
		ID:             "2ND",
		Service:        "Priority",
		ZipOrigination: 44106,
		ZipDestination: 20770,
		Pounds:         1,
		Ounces:         8,
		Container:      "Nonrectangular",
		Size:           "Large",
		Width:          15,
		Length:         30,
		Height:         15,
		Girth:          55,
		Value:          1000,
		SpecialServices: []SpecialService{
			SpecialService{SpecialService: 1},
		},
	}

	package3 := Package{
		ID:             "3RD",
		Service:        "All",
		ZipOrigination: 90210,
		ZipDestination: 96698,
		Pounds:         8,
		Ounces:         32,
		Size:           "Regular",
		Machinable:     true,
		DropOffTime:    "23:59",
		ShipDate:       "2013-07-28",
	}

	request := RateV4Request{
		USERID:   "048NA0008090",
		Revision: 2,
		Package: []Package{
			package1,
			package2,
			package3,
		},
	}

	outputType := "xml"

	var output []byte
	var err error

	if outputType == "xml" {
		output, err = xml.MarshalIndent(request, "", "  ")
	} else if outputType == "json" {
		output, err = json.MarshalIndent(request, "", "  ")
	}

	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write([]byte("\n"))
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
}
