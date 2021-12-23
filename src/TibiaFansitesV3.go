package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// TibiaFansitesV3 func
func TibiaFansitesV3() string {

	// Child of Fansite
	type ContentType struct {
		Statistics bool `json:"statistics"`
		Texts      bool `json:"texts"`
		Tools      bool `json:"tools"`
		Wiki       bool `json:"wiki"`
	}
	// Child of Fansite
	type SocialMedia struct {
		Discord   bool `json:"discord"`
		Facebook  bool `json:"facebook"`
		Instagram bool `json:"instagram"`
		Reddit    bool `json:"reddit"`
		Twitch    bool `json:"twitch"`
		Twitter   bool `json:"twitter"`
		Youtube   bool `json:"youtube"`
	}

	// Child of Fansites
	type Fansite struct {
		Name           string      `json:"name"`
		LogoURL        string      `json:"logo_url"`
		Homepage       string      `json:"homepage"`
		Contact        string      `json:"contact"`
		ContentType    ContentType `json:"content_type"`
		SocialMedia    SocialMedia `json:"social_media"`
		Languages      []string    `json:"languages"`
		Specials       []string    `json:"specials"`
		FansiteItem    bool        `json:"fansite_item"`
		FansiteItemURL string      `json:"fansite_item_url"`
	}

	// Child of JSONData
	type Fansites struct {
		PromotedFansites  []Fansite `json:"promoted"`
		SupportedFansites []Fansite `json:"supported"`
	}

	//
	// The base includes two levels: Fansites and Information
	type JSONData struct {
		Fansites    Fansites    `json:"fansites"`
		Information Information `json:"information"`
	}

	// Getting data with TibiadataHTMLDataCollectorV3
	BoxContentHTML := TibiadataHTMLDataCollectorV3("https://www.tibia.com/community/?subtopic=fansites")

	// Loading HTML data into ReaderHTML for goquery with NewReader
	ReaderHTML, err := goquery.NewDocumentFromReader(strings.NewReader(BoxContentHTML))
	if err != nil {
		log.Fatal(err)
	}

	// Creating empty PromotedFansitesData and SupportedFansitesData var
	var PromotedFansitesData []Fansite
	var SupportedFansitesData []Fansite

	// Running query over each tr in fansitesinnertable
	ReaderHTML.Find("#promotedfansitesinnertable tr").First().NextAll().Each(func(index int, s *goquery.Selection) {
		// #promotedfansitesinnertable
		// #supportedfansitesinnertable

		// Storing HTML into FansiteTrHTML
		FansiteTrHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Removing line breaks
		FansiteTrHTML = TibiadataHTMLRemoveLinebreaksV3(FansiteTrHTML)

		// Regex to get data for fansites
		regex1 := regexp.MustCompile(`<td><a href="(.*)" target.*img .*src="(.*)" alt="(.*)"\/><\/a>.*<a href=".*">(.*)<\/a><\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*<ul><li>(.*)<\/li><\/ul><\/td><td.*top;">(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(FansiteTrHTML, -1)

		if len(subma1) > 0 {

			// ContentType
			ContentTypeData := ContentType{}
			var imgRE1 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
			imgs1 := imgRE1.FindAllStringSubmatch(subma1[0][5], -1)
			out := make([]string, len(imgs1))
			for i := range out {
				if strings.Contains(imgs1[i][1], "Statistics") {
					ContentTypeData.Statistics = true
				} else if strings.Contains(imgs1[i][1], "ArticlesNews") {
					ContentTypeData.Texts = true
				} else if strings.Contains(imgs1[i][1], "Tools") {
					ContentTypeData.Tools = true
				} else if strings.Contains(imgs1[i][1], "Wiki") {
					ContentTypeData.Wiki = true
				}
			}

			// SocialMedia
			SocialMediaData := SocialMedia{}
			var imgRE2 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
			imgs2 := imgRE2.FindAllStringSubmatch(subma1[0][6], -1)
			out2 := make([]string, len(imgs2))
			for i := range out2 {
				if strings.Contains(imgs2[i][1], "Discord") {
					SocialMediaData.Discord = true
				} else if strings.Contains(imgs2[i][1], "Facebook") {
					SocialMediaData.Facebook = true
				} else if strings.Contains(imgs2[i][1], "Instagram") {
					SocialMediaData.Instagram = true
				} else if strings.Contains(imgs2[i][1], "Reddit") {
					SocialMediaData.Reddit = true
				} else if strings.Contains(imgs2[i][1], "Twitch") {
					SocialMediaData.Twitch = true
				} else if strings.Contains(imgs2[i][1], "Twitter") {
					SocialMediaData.Twitter = true
				} else if strings.Contains(imgs2[i][1], "Youtube") {
					SocialMediaData.Youtube = true
				}
			}

			// Languages
			re := regexp.MustCompile("iti__flag.iti__(..)")
			found := re.FindAllString(subma1[0][7], -1)
			FansiteLanguagesData := make([]string, len(found))
			for i := range FansiteLanguagesData {
				FansiteLanguagesData[i] = strings.ReplaceAll(found[i], "iti__flag iti__", "")
			}

			// Specials
			subma1[0][8] = html.UnescapeString(subma1[0][8])
			FansiteSpecialsData := strings.Split(subma1[0][8], "</li><li>")

			// FansiteItem & FansiteItemURL
			var FansiteItemData bool
			var FansiteItemURLData string
			regex2 := regexp.MustCompile(`.*src="(.*)" alt=".*`)
			subma1item := regex2.FindAllStringSubmatch(subma1[0][9], -1)
			if len(subma1item) > 0 {
				FansiteItemData = true
				FansiteItemURLData = subma1item[0][1]
			} else {
				FansiteItemData = false
				FansiteItemURLData = ""
			}

			PromotedFansitesData = append(PromotedFansitesData, Fansite{
				Name:           subma1[0][3],
				LogoURL:        subma1[0][2],
				Homepage:       subma1[0][1],
				Contact:        subma1[0][4],
				ContentType:    ContentTypeData,
				SocialMedia:    SocialMediaData,
				Languages:      FansiteLanguagesData,
				Specials:       FansiteSpecialsData,
				FansiteItem:    FansiteItemData,
				FansiteItemURL: FansiteItemURLData,
			})
		}

	})

	// Running query over each tr in fansitesinnertable
	ReaderHTML.Find("#supportedfansitesinnertable tr").First().NextAll().Each(func(index int, s *goquery.Selection) {
		// #promotedfansitesinnertable
		// #supportedfansitesinnertable

		// Storing HTML into FansiteTrHTML
		FansiteTrHTML, err := s.Html()
		if err != nil {
			log.Fatal(err)
		}

		// Removing line breaks
		FansiteTrHTML = TibiadataHTMLRemoveLinebreaksV3(FansiteTrHTML)

		// Regex to get data for fansites
		regex1 := regexp.MustCompile(`<td><a href="(.*)" target.*img .*src="(.*)" alt="(.*)"\/><\/a>.*<a href=".*">(.*)<\/a><\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*top;">(.*)<\/td><td.*<ul><li>(.*)<\/li><\/ul><\/td><td.*top;">(.*)<\/td>`)
		subma1 := regex1.FindAllStringSubmatch(FansiteTrHTML, -1)

		if len(subma1) > 0 {

			// ContentType
			ContentTypeData := ContentType{}
			var imgRE1 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
			imgs1 := imgRE1.FindAllStringSubmatch(subma1[0][5], -1)
			out := make([]string, len(imgs1))
			for i := range out {
				if strings.Contains(imgs1[i][1], "Statistics") {
					ContentTypeData.Statistics = true
				} else if strings.Contains(imgs1[i][1], "ArticlesNews") {
					ContentTypeData.Texts = true
				} else if strings.Contains(imgs1[i][1], "Tools") {
					ContentTypeData.Tools = true
				} else if strings.Contains(imgs1[i][1], "Wiki") {
					ContentTypeData.Wiki = true
				}
			}

			// SocialMedia
			SocialMediaData := SocialMedia{}
			var imgRE2 = regexp.MustCompile(`<img[^>]+\bsrc="([^"]+)"`)
			imgs2 := imgRE2.FindAllStringSubmatch(subma1[0][6], -1)
			out2 := make([]string, len(imgs2))
			for i := range out2 {
				if strings.Contains(imgs2[i][1], "Discord") {
					SocialMediaData.Discord = true
				} else if strings.Contains(imgs2[i][1], "Facebook") {
					SocialMediaData.Facebook = true
				} else if strings.Contains(imgs2[i][1], "Instagram") {
					SocialMediaData.Instagram = true
				} else if strings.Contains(imgs2[i][1], "Reddit") {
					SocialMediaData.Reddit = true
				} else if strings.Contains(imgs2[i][1], "Twitch") {
					SocialMediaData.Twitch = true
				} else if strings.Contains(imgs2[i][1], "Twitter") {
					SocialMediaData.Twitter = true
				} else if strings.Contains(imgs2[i][1], "Youtube") {
					SocialMediaData.Youtube = true
				}
			}

			// Languages
			re := regexp.MustCompile("iti__flag.iti__(..)")
			found := re.FindAllString(subma1[0][7], -1)
			FansiteLanguagesData := make([]string, len(found))
			for i := range FansiteLanguagesData {
				FansiteLanguagesData[i] = strings.ReplaceAll(found[i], "iti__flag iti__", "")
			}

			// Specials
			subma1[0][8] = html.UnescapeString(subma1[0][8])
			FansiteSpecialsData := strings.Split(subma1[0][8], "</li><li>")

			// FansiteItem & FansiteItemURL
			var FansiteItemData bool
			var FansiteItemURLData string
			regex2 := regexp.MustCompile(`.*src="(.*)" alt=".*`)
			subma1item := regex2.FindAllStringSubmatch(subma1[0][9], -1)
			if len(subma1item) > 0 {
				FansiteItemData = true
				FansiteItemURLData = subma1item[0][1]
			} else {
				FansiteItemData = false
				FansiteItemURLData = ""
			}

			SupportedFansitesData = append(SupportedFansitesData, Fansite{
				Name:           subma1[0][3],
				LogoURL:        subma1[0][2],
				Homepage:       subma1[0][1],
				Contact:        subma1[0][4],
				ContentType:    ContentTypeData,
				SocialMedia:    SocialMediaData,
				Languages:      FansiteLanguagesData,
				Specials:       FansiteSpecialsData,
				FansiteItem:    FansiteItemData,
				FansiteItemURL: FansiteItemURLData,
			})
		}

	})

	// Printing the PromotedFansitesData data to log
	// log.Println(PromotedFansitesData)
	// Printing the SupportedFansitesData data to log
	// log.Println(SupportedFansitesData)

	//
	// Build the data-blob
	jsonData := JSONData{
		Fansites{
			PromotedFansites:  PromotedFansitesData,
			SupportedFansites: SupportedFansitesData,
		},
		Information{
			APIVersion: TibiadataAPIversion,
			Timestamp:  TibiadataDatetimeV3(""),
		},
	}

	js, _ := json.Marshal(jsonData)
	if TibiadataDebug {
		fmt.Printf("%s\n", js)
	}
	return string(js)
}