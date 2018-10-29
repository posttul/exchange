package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/posttul/exchange/storage"
)

func (s *Scraper) getDiarioOficial() (err error) {
	resp, err := http.Get("http://www.banxico.org.mx/tipcamb/tipCamMIAction.do")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return
	}
	dt := htmlquery.FindOne(doc, "/html/body/form/table/tbody/tr[5]/td/table/tbody/tr[1]/td/table/tbody/tr[1]/td[2]/table/tbody/tr[3]/td[4]")
	value := strings.Replace(htmlquery.InnerText(dt), "\n", "", -1)
	value = strings.Replace(value, " ", "", -1)
	valflo, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return
	}
	s.data.Rates["diario_oficial"] = storage.Rate{
		Value:      float64(valflo),
		LastUpdate: time.Now(),
	}
	return
}

func (s *Scraper) getFixer() error {
	res, err := http.Get("http://data.fixer.io/api/latest?access_key=f0c42e97b04a8ac7ef7b009c6ce07c54&base=EUR&symbols=USD")
	if err != nil {
		return err
	}
	fixerResponse := struct {
		Rates struct {
			USD float64
		} `json:"rates"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&fixerResponse)
	if err != nil {
		return err
	}
	s.data.Rates["fixer"] = storage.Rate{
		Value:      fixerResponse.Rates.USD,
		LastUpdate: time.Now(),
	}
	return err
}

func (s *Scraper) getBanxico() error {
	res, err := http.Get("https://www.banxico.org.mx/SieAPIRest/service/v1/series/SF43718/datos/oportuno?token=a7ae96bcfa73dfc242708c25681d37db9f1956dd73d4add684a8dd883a7d9677&mediaType=json")
	if err != nil {
		return err
	}
	banxicoResp := struct {
		BMX struct {
			Series []struct {
				Datos []struct {
					Dato string `json:"dato"`
				} `json:"datos"`
			} `json:"Series"`
		} `json:"bmx"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&banxicoResp)
	if err != nil {
		return err
	}
	if len(banxicoResp.BMX.Series) < 1 {
		return fmt.Errorf("Up's the serires is empty")
	}

	datos := banxicoResp.BMX.Series[0].Datos[0]
	value, err := strconv.ParseFloat(datos.Dato, 32)
	if err != nil {
		return err
	}
	s.data.Rates["banxico"] = storage.Rate{
		Value:      float64(value),
		LastUpdate: time.Now(),
	}
	return err
}

func NewScraper(d time.Duration) *Scraper {
	return &Scraper{
		sleep: d,
		data: storage.Response{
			Rates: make(map[string]storage.Rate),
		},
	}
}

// Scraper is gonna handle all
type Scraper struct {
	sleep time.Duration
	data  storage.Response
}

// GetData use to scrap web.
func (s *Scraper) GetData(w storage.Storage) {
	for {

		if err := s.getDiarioOficial(); err != nil {
			log.Printf("Could not get getDiarioOficial storage err -> %s", err.Error())
		}

		if err := s.getFixer(); err != nil {
			log.Printf("Could not get fixer storage err -> %s", err.Error())
		}

		if err := s.getBanxico(); err != nil {
			log.Printf("Could not get fixer storage err -> %s", err.Error())
		}

		s.writeToStorage(w)
		time.Sleep(s.sleep)
	}
}

func (s *Scraper) writeToStorage(w storage.Storage) error {
	bts, err := json.Marshal(s.data)
	if err != nil {
		log.Printf("Could not update storage err -> %s", err.Error())
	}
	err = w.Write(bts)
	if err != nil {
		log.Printf("Could not update storage err -> %s", err.Error())
	}
	log.Printf("Rates cache updated at %s", time.Now())
	return err
}
