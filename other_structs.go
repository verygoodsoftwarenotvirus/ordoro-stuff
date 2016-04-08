package main

import (
	"fmt"
)

type PackageError struct {
	Number      int    `xml:",omitempty" json:",omitempty"`
	Source      string `xml:",omitempty" json:",omitempty"`
	Description string `xml:",omitempty" json:",omitempty"`
	HelpFile    string // reserved for future use, according to the API docs
	HelpContext string // reserved for future use, according to the API docs
}

type Restriction struct {
	Restrictions string `xml:",omitempty" json:",omitempty"`
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

type Content struct {
	ContentType        string `xml:",omitempty" json:",omitempty"`
	ContentDescription string `xml:",omitempty" json:",omitempty"`
}

type Postage struct {
	CLASSID         string           `xml:",attr"`
	MailService     string           `xml:",omitempty" json:",omitempty"`
	Rate            float32          `xml:",omitempty" json:",omitempty"`
	SpecialServices []SpecialService `xml:",omitempty" json:",omitempty"`
	CommitmentDate  string           `xml:",omitempty" json:",omitempty"`
	CommitmentName  string           `xml:",omitempty" json:",omitempty"`
}

type SpecialService struct {
	// SpecialService could (and maybe should) have been an int, but I was getting parseInt
	// errors on RateResponses, so it is a string ¯\_(ツ)_/¯
	SpecialService        string  `xml:",omitempty" json:",omitempty"`
	ServiceID             int     `xml:",omitempty" json:",omitempty"`
	ServiceName           string  `xml:",omitempty" json:",omitempty"`
	Available             bool    `xml:",omitempty" json:",omitempty"`
	AvailableOnline       bool    `xml:",omitempty" json:",omitempty"`
	Price                 float32 `xml:",omitempty" json:",omitempty"`
	PriceOnline           float32 `xml:",omitempty" json:",omitempty"`
	DeclaredValueRequired bool    `xml:",omitempty" json:",omitempty"`
	DueSenderRequired     bool    `xml:",omitempty" json:",omitempty"`
}
