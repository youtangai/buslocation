package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/PuerkitoBio/goquery"
)

const (
	originHost = "http://www.hakobus.jp"
	originPath = "/result.php"
)

//Info is hogehoge
type Info struct {
	time    string
	via     string
	landing string
	dest    string
	status  string
}

var (
	start = flag.String("s", "sisho", "乗車場所 school or sisho or weather or fun or airport")
	end   = flag.String("e", "fun", "降車場所 school or sisho or weather or fun or airport")
	m     = make(map[int]Info)
	info  = Info{}
)

//SjisToUTF8 is hogehoge
func SjisToUTF8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

// SetTableInfo is hogehoge
func SetTableInfo(start, end string) {
	doc, err := goquery.NewDocument(originHost + originPath + "?in=" + start + "&out=" + end)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("table tr td").Each(func(i int, s *goquery.Selection) {
		text, err := SjisToUTF8(s.Text())
		if err != nil {
			log.Fatal(err)
		}

		index := (i - 11) % 9
		switch index {
		case 0:
			info.time = text
			break
		case 1:
			info.via = text
			break
		case 2:
			info.landing = text
			break
		case 3:
			info.dest = text
			break
		case 6:
			info.status = text
			key := (i - 11) / 9
			m[key] = info
			info = Info{}
			break
		default:
			break
		}
	})
}

//PrintMap is hogehoge
func PrintMap() {
	keys := make([]int, len(m))

	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("%d\n停留所時刻:%s\n系統:%s乗り場:%s\n行き先:%s\n運行状況:%s\n\n",
			key,
			m[key].time,
			m[key].via,
			m[key].landing,
			m[key].dest,
			m[key].status,
		)
	}
}

func main() {
	flag.Parse()
	dic := map[string]string{
		"school":  "453",
		"sisho":   "155",
		"fun":     "165",
		"weather": "156",
		"airport": "506",
	}
	s := dic[*start]
	e := dic[*end]
	SetTableInfo(s, e)
	PrintMap()
}
