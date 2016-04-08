package main

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestExampleStructRequest(t *testing.T) {
	uspsProvidedXML := `<RateV4Request USERID="048NA0008090">
    <Revision>2</Revision>
    <Package ID="1ST">
        <Service>FIRST CLASS</Service>
        <FirstClassMailType>LETTER</FirstClassMailType>
        <ZipOrigination>44106</ZipOrigination>
        <ZipDestination>20770</ZipDestination>
        <Pounds>0</Pounds>
        <Ounces>3.5</Ounces>
        <Container/>
        <Size>REGULAR</Size>
        <Machinable>true</Machinable>
    </Package>
    <Package ID="2ND">
        <Service>PRIORITY</Service>
        <ZipOrigination>44106</ZipOrigination>
        <ZipDestination>20770</ZipDestination>
        <Pounds>1</Pounds>
        <Ounces>8</Ounces>
        <Container>NONRECTANGULAR</Container>
        <Size>LARGE</Size>
        <Width>15</Width>
        <Length>30</Length>
        <Height>15</Height>
        <Girth>55</Girth>
        <Value>1000</Value>
        <SpecialServices>
            <SpecialService>1</SpecialService>
        </SpecialServices>
    </Package>
    <Package ID="3RD">
        <Service>ALL</Service>
        <ZipOrigination>90210</ZipOrigination>
        <ZipDestination>96698</ZipDestination>
        <Pounds>8</Pounds>
        <Ounces>32</Ounces>
        <Container/>
        <Size>REGULAR</Size>
        <Machinable>true</Machinable>
        <DropOffTime>23:59</DropOffTime>
        <ShipDate>2016-04-11</ShipDate>
    </Package>
</RateV4Request>`

	uspsUserID := "048NA0008090"
	// from https://www.usps.com/business/web-tools-apis/rate-calculator-api.htm#_Toc423593290
	package1 := Package{
		ID:                 "1ST",
		Service:            "FIRST CLASS",
		FirstClassMailType: "LETTER",
		ZipOrigination:     44106,
		ZipDestination:     20770,
		Pounds:             0,
		Ounces:             3.5,
		Size:               "REGULAR",
		Machinable:         true,
	}

	package2 := Package{
		ID:             "2ND",
		Service:        "PRIORITY",
		ZipOrigination: 44106,
		ZipDestination: 20770,
		Pounds:         1,
		Ounces:         8,
		Container:      "NONRECTANGULAR",
		Size:           "LARGE",
		Width:          15,
		Length:         30,
		Height:         15,
		Girth:          55,
		Value:          1000,
		SpecialServices: []SpecialService{
			SpecialService{
				SpecialService: "1",
			},
		},
	}

	package3 := Package{
		ID:             "3RD",
		Service:        "ALL",
		ZipOrigination: 90210,
		ZipDestination: 96698,
		Pounds:         8,
		Ounces:         32,
		Size:           "REGULAR",
		Machinable:     true,
		DropOffTime:    "23:59",
		ShipDate: &ShipDate{
			Date: "2016-04-11",
		},
	}

	request := RateV4Request{
		UserID:   uspsUserID,
		Revision: 2,
		Package: []Package{
			package1,
			package2,
			package3,
		},
	}

	marshalledXML, err := xml.MarshalIndent(request, "", "    ")
	if err != nil {
		t.Errorf("Error returned while marshalling request struct: %v", err)
	}

	// Lame encoding hacks that we have to do
	// the raw string variable above contains carriage returns, but the generated xml doesn't for reasons
	input := []byte(strings.Replace(string(uspsProvidedXML), "\r", "", -1))
	// apparently the guy who wrote encoding/xml thought self closing tags were lame:
	// https://groups.google.com/d/msg/golang-nuts/guG6iOCRu08/VE-Cm_k528MJ
	// USPS apparently disagrees.
	output := []byte(strings.Replace(string(marshalledXML), "<Container></Container>", "<Container/>", -1))

	if !bytes.Equal(input, output) {
		t.Errorf("Produced XML is not equal to example XML provided by USPS.\n\ngiven:\n%v\ngenerated:\n%v", input, output)
	}
}
