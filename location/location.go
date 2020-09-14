package location

//go get -v gopkg.in/webdeskltd/dadata.v2
import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"gopkg.in/webdeskltd/dadata.v2"
)

//Config -
type Config struct {
	APIKey    string
	SecretKey string
}

//Location -
type Location struct {
	Longitude float64
	Latitude  float64
}

//Metro -
type Metro struct {
	Name     string
	Line     string
	Distance float64
}

//Address - Адрес
type Address struct {
	Address     string //Адрес
	FullAddress string //адрес с индексом

	PostalCode     string
	Country        string
	Region         string //область
	Area           string //район
	City           string //город  нас пункт
	CityType       string //тип нас пункта
	Settlement     string
	SettlementType string
	Street         string
	House          string
	Location       Location
	Metro          []Metro
}

var config Config

//инициализация
func init() {
	file, err := os.Open("location/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}
}

//GetAddressByString -
func GetAddressByString(s string) (Address, int64) {
	geo := dadata.NewDaData(config.APIKey, config.SecretKey)
	result, err := geo.CleanAddresses(s)
	if err != nil {
		return Address{}, -1
	}
	d := result[0]
	lt, _ := strconv.ParseFloat(d.GeoLat, 64)
	ln, _ := strconv.ParseFloat(d.GeoLon, 64)
	return Address{
		Address:        d.Result,
		FullAddress:    d.Result,
		PostalCode:     d.PostalCode,
		Country:        d.Country,
		Region:         d.RegionWithType,
		Area:           d.AreaWithType,
		City:           d.City,
		CityType:       d.CityType,
		Settlement:     d.Settlement,
		SettlementType: d.SettlementType,
		Street:         d.Street,
		House:          d.House,
		Location: Location{
			Latitude:  lt,
			Longitude: ln,
		},
		Metro: func(m []dadata.Metro) []Metro {
			r := []Metro{}
			for _, v := range m {
				r = append(r, Metro{
					Name:     v.Name,
					Line:     v.Line,
					Distance: v.Distance,
				})
			}
			return r
		}(d.Metro),
	}, int64(d.QualityCode.(float64))
}

//GetAddressByLocation -
func GetAddressByLocation(l Location) Address {
	geo := dadata.NewDaData(config.APIKey, config.SecretKey)
	request := dadata.GeolocateRequest{Lat: float32(l.Latitude), Lon: float32(l.Longitude), Count: 1}
	result, err := geo.GeolocateAddress(request)
	if err != nil {
		return Address{}
	}
	d := result[0].Data
	lt, _ := strconv.ParseFloat(d.GeoLat, 64)
	ln, _ := strconv.ParseFloat(d.GeoLon, 64)
	return Address{
		Address:        result[0].Value,
		FullAddress:    result[0].UnrestrictedValue,
		PostalCode:     d.PostalCode,
		Country:        d.Country,
		Region:         d.RegionWithType,
		Area:           d.AreaWithType,
		City:           d.City,
		CityType:       d.CityType,
		Settlement:     d.Settlement,
		SettlementType: d.SettlementType,
		Street:         d.Street,
		House:          d.House,
		Location: Location{
			Latitude:  lt,
			Longitude: ln,
		},
		Metro: func(m []dadata.Metro) []Metro {
			r := []Metro{}
			for _, v := range m {
				r = append(r, Metro{
					Name:     v.Name,
					Line:     v.Line,
					Distance: v.Distance,
				})
			}
			return r
		}(d.Metro),
	}
}
