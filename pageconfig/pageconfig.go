package pageconfig

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type style struct {
	Fg   string
	Bg   string
	Attr string
}

// Page contains all of the fields that can be parsed
// from the pageconfig YAML for each page
type Page struct {
	Name     string `yaml:"name"`
	Links   []*struct {
		KeyStroke  string `yaml:"keyStroke"`
		PageName   string `yaml:"pageName"`
		Server     string `yaml:"server"`
		Port       string `yaml:"port"`
	} `yaml:"links"`
	DivBoxes []*struct {
		Name             string `yaml:"name"`
		Border           bool   `yaml:"border"`
		BorderW          int32  `yaml:"borderW"`
		BorderCharString string `yaml:"borderChar"`
		FillCharString   string `yaml:"fillChar"`
		BorderChar       rune
		FillChar         rune
		StartY           int32  `yaml:"startY"`
		StartX           int32  `yaml:"startX"`
		Width            int32  `yaml:"width"`
		Height           int32  `yaml:"height"`
		BorderSt         *style `yaml:"borderSt"`
		FillSt           *style `yaml:"fillSt"` } `yaml:"divBoxes"`
	Elements []*struct {
		TextBlobs []*struct {
			Content  string   `yaml:"content"`
			Wrap     bool     `yaml:"wrap"`
			Style    *style   `yaml:"style"`
			DivNames []string `yaml:"divNames"`
		} `yaml:"textBlobs"`
	} `yaml:"elements"`
}

// Pages contains configuration for the various
// uggly pages that the server should run
type Pages struct {
	Pages []*Page
}

// NewPageConfig takes a yaml filename as input and
// attempts to parse it into a config object.
func NewPageConfig(filename string) (*Pages, error) {
	var err error
	pc := Pages{}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return &pc, err
	}
	err = yaml.Unmarshal(yamlFile, &pc.Pages)
	if err != nil {
		return &pc, err
	}
	if len(pc.Pages) < 1 {
		err = errors.New("no pages parsed from config")
	}
	// validate some of the input
	for _, page := range pc.Pages {
		for _, divbox := range page.DivBoxes {
			bcr := []rune(divbox.BorderCharString)
			fcr := []rune(divbox.FillCharString)
			if len(fcr) > 1 || len(bcr) > 1 {
				err = errors.New(
					"borderChar and fillChar must be string of" +
						" length 1 so it can be parsed to rune")
				return &pc, err
			}
			if len(bcr) > 0 {
				divbox.BorderChar = bcr[0]
			}
			divbox.FillChar = fcr[0]
		}
	}
	return &pc, err
}
