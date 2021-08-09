package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const baseUrl = "http://www.meb.gov.tr/baglantilar/okullar/index.php?ILKODU="
const okulJsonFilename = "okullar.json"

var turkceKeyz = unicode.TurkishCase

type Okul struct {
	Il   string `json:"il"`
	Ilce string `json:"ilce"`
	Okul string `json:"okul"`
}

func (o *Okul) Temizle() {
	o.Il = temizle(o.Il)
	o.Ilce = temizle(o.Ilce)
	o.Okul = temizle(o.Okul)
}

func temizle(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLowerSpecial(unicode.TurkishCase, str)
	strArr := strings.Fields(str)

	for i, s := range strArr {
		if len(s) == 0 {
			continue
		}

		tmp := []rune(s)
		tmp[0] = turkceKeyz.ToUpper(tmp[0])
		strArr[i] = string(tmp)
	}

	return strings.Join(strArr, " ")
}

func main() {
	var okullar []Okul
	for i := 1; ; i++ {
		iStr := strconv.Itoa(i)
		res, err := http.Get(baseUrl + iStr)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		sonuc := doc.Find("#icerik-listesi > tbody > tr").Each(func(i int, s *goquery.Selection) {
			ilILceOkulStr := s.First().Find("a").Text()
			ilILceOkulArr := strings.Split(ilILceOkulStr, "-")

			if len(ilILceOkulArr) < 3 {
				log.Fatal("hayırdır lan??")
			}
			okul := Okul{
				Il:   ilILceOkulArr[0],
				Ilce: ilILceOkulArr[1],
				Okul: ilILceOkulArr[2],
			}
			okul.Temizle()

			okullar = append(okullar, okul)
		})
		if sonuc.Nodes == nil {
			break
		}
	}

	file, err := json.MarshalIndent(okullar, "", "   ")
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(okulJsonFilename, file, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	// oku
	jsonFile, err := os.Open(okulJsonFilename)
	if err != nil {
		log.Fatalln(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &okullar); err != nil {
		log.Fatalln(err)
	}
	for _, okul := range okullar {
		log.Printf("Il: %s, Ilce: %s, Okul: %s\n", okul.Il, okul.Ilce, okul.Okul)
	}
}
