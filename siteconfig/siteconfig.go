package siteconfig

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

// Sites contains configuration for the various
// uggly sites that the server should run
type Sites struct {
	Sites []*struct {
		Name     string `yaml:"name"`
		DivBoxes []*struct {
			Name       string `yaml:"name"`
			Border     bool   `yaml:"border"`
			BorderW    int32    `yaml:"borderW"`
			BorderCharString string `yaml:"borderChar"`
			FillCharString   string `yaml:"fillChar"`
            BorderChar rune
            FillChar rune
			StartY     int32    `yaml:"startY"`
			StartX     int32    `yaml:"startX"`
			Width      int32    `yaml:"width"`
			Height     int32    `yaml:"height"`
            BorderSt   *style    `yaml:"borderSt"`
            FillSt   *style    `yaml:"fillSt"`
		} `yaml:"divBoxes"`
		Elements []*struct {
			TextBlobs []*struct {
				Content  string   `yaml:"content"`
				Wrap     bool     `yaml:"wrap"`
                Style  *style `yaml:"style"`
				DivNames []string `yaml:"divNames"`
			} `yaml:"textBlobs"`
		} `yaml:"elements"`
	}
}

// NewSiteConfig takes a yaml filename as input and
// attempts to parse it int32o a config object.
func NewSiteConfig(filename string) (*Sites, error) {
    var err error
    sc := Sites{}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return &sc, err
	}
	err = yaml.Unmarshal(yamlFile, &sc.Sites)
	if err != nil {
		return &sc, err
	}
    if len(sc.Sites) < 1 {
        err = errors.New("no sites parsed from config")
    }
    // validate some of the input
    for _, site := range sc.Sites {
        for _, divbox := range site.DivBoxes {
            bcr := []rune(divbox.BorderCharString)
            fcr := []rune(divbox.FillCharString)
            if len(fcr) > 1 || len(bcr) > 1 {
                err = errors.New(
                    "borderChar and fillChar must be string of"+
                    " length 1 so it can be parsed to rune")
                return &sc, err
            }
            divbox.BorderChar = bcr[0]
            divbox.FillChar = fcr[0]
        }
    }
	return &sc, err
}
