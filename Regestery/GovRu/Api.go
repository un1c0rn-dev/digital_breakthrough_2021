package GovRu

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"unicorn.dev.web-scrap/Tasks"
)

const (
	DamiaBaseUrl    string = "https://damia.ru/api-zakupki/"
	DamiaSearchPath string = "zsearch"
)

const (
	KeywordParam  string = "q"
	FromDataParam string = "from_data"
	ToDateParam   string = "to_date"
	RegionParam   string = "region"
	OkpdParam     string = "okpd"
	StatusParam   string = "status"
	InnParam      string = "cust_inn"
	PlacingParam  string = "placing"
	ETPParam      string = "etp"
	MaxPriceParam string = "max_price"
	MinPriceParam string = "min_price"
	SMPParam      string = "smp"
	FzParam       string = "fz"
	PageParam     string = "page"
	KeyParam      string = "key"
)

type SearchStatusCode int
type SearchFZ int
type SearchETP int
type SearchPlacing int

const (
	AnyStatus  SearchStatusCode = -1
	Accepted                    = 1
	Commission                  = 2
	Done                        = 3
	Discarded                   = 4
)

const (
	AnyFZ SearchFZ = -1
	FZ44           = 44
	FZ223          = 223
)

/*
	Способ отбора победителя:
	1 – Электронный аукцион
	2 – Запрос котировок
	3 – Конкурс
	4 – Закупка у единственного поставщика
	5 – Запрос предложений
	99 – Другие способы
*/

const (
	AnyPlacing     SearchPlacing = -1
	EAuction                     = 1
	KQuery                       = 2
	Contest                      = 3
	SingleProvider               = 4
	OfferQuery                   = 5
	OtherPlacing                 = 99
)

/*
	Площадка проведения торгов:
	1 – Сбербанк-АСТ
	2 – РТС-тендер
	3 – ЕЭТП
	4 – ZakazRF (АГЗРТ)
	5 – ЭТП НЭП (ММВБ)
	6 – РАД (Lot-Online)
	7 – B2B-Center
	8 – Фабрикант
	9 – ЭТП ГПБ
	10 – OTC.RU
	11 – ТЭК-Торг
	12 – ЭТПРФ
	13 – Газнефтеторг
	14 – Тендер.Про
	15 – Аукционный Конкурсный Дом
	16 – ПолюсГолд
	99 – Другие площадки
*/

const (
	AnyETP       SearchETP = -1
	Sberbank               = 1
	RTSTender              = 2
	EETP                   = 3
	ZakazRF                = 4
	ETPNEP                 = 5
	RAD                    = 6
	B2BCenter              = 7
	Fabrikant              = 8
	ETPGPB                 = 9
	OTCRU                  = 10
	TEKTorg                = 11
	ETPRF                  = 12
	Gazneftetorg           = 13
	TenderPro              = 14
	AKD                    = 15
	PolusGold              = 16
	Other                  = 99
)

/* https://damia.ru/apizakupki#zsearch */

type SearchQuery struct {
	Keywords    []string
	FromDateYMD [3]int
	ToDateYMD   [3]int
	Region      []int
	Okpd        string
	Status      SearchStatusCode
	Placing     []int
	Etp         []int
	MinPrice    int
	MaxPrice    int
	Fz          SearchFZ
	Page        int
}

func makeParamFromStringList(list []string) string {
	p := ""
	for i, keyword := range list {
		p += keyword
		if i < len(list)-1 {
			p += ","
		}
	}
	return p
}

func makeParamFromIntList(list []int) string {
	p := ""
	for i, keyword := range list {
		p += strconv.Itoa(keyword)
		if i < len(list)-1 {
			p += ","
		}
	}
	return p
}

func NewSearchQuery() SearchQuery {
	return SearchQuery{
		Keywords:    make([]string, 0),
		FromDateYMD: [3]int{0, 0, 0},
		ToDateYMD:   [3]int{0, 0, 0},
		Region:      make([]int, 0),
		Okpd:        "",
		Status:      AnyStatus,
		Placing:     make([]int, 0),
		Etp:         make([]int, 0),
		MinPrice:    0,
		MaxPrice:    0,
		Fz:          AnyFZ,
		Page:        0,
	}
}

func Search(key string, query SearchQuery, task *Tasks.Task) {
	tasKStatus := Tasks.TaskStatusDone
	taskProgress := "Готово"
	defer func() {
		task.Status = tasKStatus
		task.Progress = taskProgress
	}()
	task.Status = Tasks.TaskStatusActive

	params := url.Values{}
	params.Add(KeyParam, key)

	if query.Keywords != nil && len(query.Keywords) > 0 {
		params.Add(KeywordParam, makeParamFromStringList(query.Keywords))
	}

	if query.Region != nil && len(query.Region) > 0 {
		params.Add(RegionParam, makeParamFromIntList(query.Region))
	}

	if len(query.Okpd) > 0 {
		params.Add(OkpdParam, query.Okpd)
	}

	if query.Placing != nil && len(query.Placing) > 0 {
		params.Add(PlacingParam, makeParamFromIntList(query.Placing))
	}

	if query.Etp != nil && len(query.Etp) > 0 {
		params.Add(ETPParam, makeParamFromIntList(query.Etp))
	}

	if query.Fz != AnyFZ {
		params.Add(FzParam, strconv.Itoa(int(query.Fz)))
	}

	if query.MaxPrice != 0 {
		params.Add(MaxPriceParam, strconv.Itoa(query.MaxPrice))
	}

	if query.MinPrice != 0 {
		params.Add(MinPriceParam, strconv.Itoa(query.MinPrice))
	}

	if query.Status != AnyStatus {
		params.Add(StatusParam, strconv.Itoa(int(query.Status)))
	}

	request := DamiaBaseUrl + DamiaSearchPath + "?" + params.Encode()
	response, err := http.Get(request)
	if err != nil {
		taskProgress = "Не удалось сделать запрос"
		tasKStatus = Tasks.TaskStatusError
		return
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		taskProgress = "Плохой ответ от сервера"
		tasKStatus = Tasks.TaskStatusError
		return
	}

	fmt.Println(string(bodyBytes))
}
