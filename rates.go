package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
)

func main() {
	// from https://www.usps.com/business/web-tools-apis/rate-calculator-api.htm#_Toc423593290
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
		ShipDate: &ShipDate{
			Date: "2016-04-08",
		},
	}

	animal := Package{
		ID:             "0",
		Service:        "Priority",
		ZipOrigination: 22201,
		ZipDestination: 26301,
		Pounds:         8,
		Ounces:         2,
		Container:      "Flat Rate Envelope",
		Size:           "Regular",
		Content: &Content{
			ContentType:        "Lives",
			ContentDescription: "Other",
		},
		Machinable: true,
	}

	requestType := "Revision Tag"

	var request RateV4Request
	if requestType == "Revision Tag" {
		request = RateV4Request{
			UserID:   "048NA0008090",
			Revision: 2,
			Package: []Package{
				package1,
				package2,
				package3,
			},
		}
	} else if requestType == "Live Animal Sample" {
		request = RateV4Request{
			UserID:   "048NA0008090",
			Revision: 2,
			Package: []Package{
				animal,
			},
		}
	}

	outputType := "xml"

	var output []byte
	var err error

	if outputType == "xml" {
		// output, err = request.requestRate()
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
