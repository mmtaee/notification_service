package phone_parser

import (
	"errors"
	"github.com/nyaruka/phonenumbers"
	"strings"
)

type Parser struct {
	Number         string
	Code           int32
	National       string
	International  string
	NationalNumber uint64
	Masked         string
}

type ParsedNumberInterface interface {
	Parse(number string) (*Parser, error)
	IsIranNumber() bool
	IranMasked() *Parser
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(number string) (*Parser, error) {
	num, err := phonenumbers.Parse(number, "")
	if err != nil {
		return nil, err
	}
	if !phonenumbers.IsValidNumber(num) {
		return nil, errors.New("invalid number")
	}
	numObj := &Parser{
		Number:         num.String(),
		Code:           num.GetCountryCode(),
		National:       phonenumbers.Format(num, phonenumbers.NATIONAL),
		International:  phonenumbers.Format(num, phonenumbers.INTERNATIONAL),
		NationalNumber: num.GetNationalNumber(),
		Masked:         num.String(),
	}
	return numObj, nil
}

func (p *Parser) IsIranNumber() bool {
	return p.Code == 98
}

func (p *Parser) IranMasked() *Parser {
	p.Masked = strings.Replace(p.Number, "+98", "0", -1)
	return p
}
