package GovRu

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicorn.dev.web-scrap/Tasks"
)

const (
	DamiaBaseUrl     string = "https://damia.ru/api-zakupki/"
	DamiaSearchPath  string = "zsearch"
	DamiaRNPPath     string = "rnp"
	DamiaZakupkaPath string = "zakupka"
)

const (
	KeywordParam  string = "q"
	FromDataParam string = "from_data"
	ToDateParam   string = "to_date"
	RegionParam   string = "region"
	OkpdParam     string = "okpd"
	StatusParam   string = "status"
	CustInnParam  string = "cust_inn"
	InnParam      string = "inn"
	RegnParam     string = "regn"
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

type DamiaConf struct {
	Key    string `json:"key"`
	Active bool
}

var damiaConf DamiaConf

func Configure(config DamiaConf) {
	damiaConf.Active = false

	if len(config.Key) > 0 {
		damiaConf.Key = config.Key
		damiaConf.Active = true
	}
}

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
	Page        int // unused
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

type rpnZakazchik struct {
	Ogrn      string `json:"ОГРН"`
	Inn       string `json:"ИНН"`
	NaimPoln  string `json:"НаимПолн"`
	NaimSokr  string `json:"НаимСокр"`
	AdresPoln string `json:"АдресПолн"`
	RukFIO    string `json:"РукФИО"`
	RukINNFL  string `json:"РукИННФЛ"`
}
type rpnZakupka struct {
	NomerIzveshcheniya string `json:"НомерИзвещения"`
	Lot                int    `json:"Лот"`
}
type rpnCost struct {
	Summa       float64 `json:"Сумма"`
	ValyutaKod  string  `json:"ВалютаКод"`
	ValyutaNaim string  `json:"ВалютаНаим"`
}
type rpnProdukt struct {
	Okpd     string `json:"ОКПД"`
	Nazvanie string `json:"Название"`
}
type rpn struct {
	Region       string       `json:"Регион"`
	Fz           int          `json:"AP"`
	DataPubl     string       `json:"ДатаПубл"`
	PrichinaVkl  string       `json:"ПричинаВкл"`
	OsnovanieVkl string       `json:"ОснованиеВкл"`
	DataVkl      string       `json:"ДатаВкл"`
	Status       string       `json:"Статус"`
	PrichinaIskl string       `json:"ПричинаИскл"`
	DataIskl     string       `json:"ДатаИскл"`
	Zakazchik    rpnZakazchik `json:"Заказчик"`
	Zakupka      rpnZakupka   `json:"Закупка"`
	Tsena        rpnCost      `json:"Цена"`
	Produkt      rpnProdukt   `json:"Продукт"`
}

type RpnRecord int

const (
	InactiveRpnRecord = iota
	ActiveRPNRecord   = 1
)

const (
	layoutISO = "2006-01-02"
)

type Etp struct {
	Kod          string `json:"Код"`
	Naimenovanie string `json:"Наименование"`
	URL          string `json:"Url"`
}
type Zakazchik struct {
	Ogrn      string `json:"ОГРН"`
	Inn       string `json:"ИНН"`
	NaimPoln  string `json:"НаимПолн"`
	NaimSokr  string `json:"НаимСокр"`
	AdresPoln string `json:"АдресПолн"`
	RukFIO    string `json:"РукФИО"`
	RukINNFL  string `json:"РукИННФЛ"`
	Telefon   string `json:"Телефон,omitempty"`
	Email     string `json:"Email,omitempty"`
}
type NachTsena struct {
	Summa       float64 `json:"Summa"`
	ValyutaKod  string  `json:"ВалютаКод"`
	ValyutaNaim string  `json:"ВалютаНаим"`
}
type ObespUchast struct {
	Summa      float64 `json:"Сумма"`
	Dolya      float64 `json:"Доля"`
	Bank       string  `json:"Банк"`
	Bik        string  `json:"БИК"`
	RaschSchet string  `json:"РасчСчет"`
	LitsSchet  string  `json:"ЛицСчет"`
}
type ObespIsp struct {
	Summa      float64 `json:"Сумма"`
	Dolya      float64 `json:"Доля"`
	Bank       string  `json:"Банк"`
	Bik        string  `json:"БИК"`
	RaschSchet string  `json:"РасчСчет"`
	LitsSchet  string  `json:"ЛицСчет"`
}
type ObespGarant struct {
}
type Produkt struct {
	Okpd      string        `json:"ОКПД"`
	Nazvanie  string        `json:"Название"`
	ObektyZak []interface{} `json:"ОбъектыЗак"`
}
type Usloviya struct {
}
type IP struct {
	Ogrnip  string `json:"ОГРНИП"`
	Innfl   string `json:"ИННФЛ"`
	Fio     string `json:"ФИО"`
	Telefon string `json:"Телефон"`
	Email   string `json:"Email"`
}
type YuL struct {
	Ogrn      string `json:"ОГРНИП"`
	Inn       string `json:"ИНН"`
	NaimPoln  string `json:"НаимПолн"`
	NaimSokr  string `json:"НаимСокр"`
	AdresPoln string `json:"АдресПолн"`
	RukFIO    string `json:"РукФИО"`
	RukINNFL  string `json:"РукИННФЛ"`
	Telefon   string `json:"Телефон"`
	Email     string `json:"Email,omitempty"`
}
type Zayavki struct {
	Nomer    string `json:"Номер"`
	IP       IP     `json:"ИП,omitempty"`
	Summa    int    `json:"Сумма"`
	Rezultat string `json:"Результат"`
	Prichina string `json:"Причина"`
	YuL      YuL    `json:"ЮЛ,omitempty"`
}
type Protokol struct {
	Tip     string    `json:"Тип"`
	Nomer   string    `json:"Номер"`
	Data    string    `json:"Дата"`
	Zayavki []Zayavki `json:"Заявки"`
	DopInfo string    `json:"ДопИнфо"`
	URL     string    `json:"Url"`
}
type Tsena struct {
	Summa       int    `json:"Сумма"`
	ValyutaKod  string `json:"ВалютаКод"`
	ValyutaNaim string `json:"ВалютаНаим"`
}
type Postavshchiki struct {
	YuL []interface{} `json:"ЮЛ"`
	IP  []IP          `json:"ИП"`
	Fl  []interface{} `json:"ФЛ"`
}
type Kontrakt struct {
	Nomer         string        `json:"Номер"`
	DataPodp      string        `json:"ДатаПодп"`
	Lot           int           `json:"Лот"`
	Tsena         Tsena         `json:"Цена"`
	Postavshchiki Postavshchiki `json:"Поставщики"`
}

type Status struct {
	Status   string `json:"Статус"`
	Prichina string `json:"Причина"`
	Data     string `json:"Дата"`
}

type Zakupka struct {
	Region       string              `json:"Регион"`
	Fz           int                 `json:"ФЗ"`
	DataPubl     string              `json:"ДатаПубл"`
	DataNach     string              `json:"ДатаНач"`
	VremyaNach   string              `json:"ВремяНач"`
	DataOkonch   string              `json:"ДатаОконч"`
	VremyaOkonch string              `json:"ВремяОконч"`
	DataRassm    string              `json:"ДатаРассм"`
	DataAukts    string              `json:"ДатаАукц"`
	VremyaAukts  string              `json:"ВремяАукц"`
	Etp          Etp                 `json:"ЭТП"`
	Zakazchik    Zakazchik           `json:"Заказчик"`
	Kontakty     []interface{}       `json:"Контакты"`
	SposobRazm   string              `json:"СпособРазм"`
	RazmRol      string              `json:"РазмРоль"`
	SMPiSONO     bool                `json:"СМПиСОНО"`
	NachTsena    NachTsena           `json:"НачЦена"`
	ObespUchast  ObespUchast         `json:"ОбеспУчаст"`
	ObespIsp     ObespIsp            `json:"ОбеспИсп"`
	ObespGarant  ObespGarant         `json:"ОбеспГарант"`
	Produkt      Produkt             `json:"Продукт"`
	Usloviya     Usloviya            `json:"Условия"`
	Protokol     Protokol            `json:"Протокол"`
	Kontrakty    map[string]Kontrakt `json:"Контракты"`
	Status       Status              `json:"Статус"`
}

func rpnIsActive(rpnData *rpn) RpnRecord {
	nowTime := time.Now()

	if len(rpnData.DataIskl) == 0 {
		return ActiveRPNRecord
	}

	excludeTime, err := time.Parse(layoutISO, rpnData.DataIskl)
	if err != nil {
		return ActiveRPNRecord
	}

	if excludeTime.Before(nowTime) {
		return InactiveRpnRecord
	}

	return ActiveRPNRecord
}

func CheckUnscrupulousOrganisation(inn string, result *Tasks.TaskResult) error {
	reputation := Tasks.TaskResultReputationGood
	defer func() {
		result.Reputation = reputation
	}()

	params := url.Values{}
	params.Add(KeyParam, damiaConf.Key)
	params.Add(InnParam, inn)

	response, err := http.Get(DamiaBaseUrl + DamiaRNPPath + "?" + params.Encode())
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(strings.NewReader(string(responseBody)))
	for {
		var responseJson map[string]map[string]rpn
		if err := dec.Decode(&responseJson); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		for innKey, rpns := range responseJson {
			if inn != innKey {
				continue
			}

			for _, rpnValue := range rpns {
				if rpnIsActive(&rpnValue) == ActiveRPNRecord {
					reputation = Tasks.TaskResultReputationBad
					return nil
				} else {
					reputation = Tasks.TaskResultReputationMed
					return nil
				}
			}

			return nil
		}
	}

	return nil
}

type distributor struct {
	Name  string
	Email string
	Phone string
	Inn   string
	Cost  int
}

func getDistributors(regn string) ([]distributor, error) {
	params := url.Values{}
	params.Add(KeyParam, damiaConf.Key)
	params.Add(RegnParam, regn)

	requestUrl := DamiaBaseUrl + DamiaZakupkaPath + "?" + params.Encode()
	response, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(strings.NewReader(string(responseBody)))

	distributors := make([]distributor, 0)

	for {
		var responseJson map[string]Zakupka
		if err := dec.Decode(&responseJson); err == io.EOF {
			break
		} else if err != nil && len(responseJson) == 0 {
			return nil, err
		}

		for _, zakupka := range responseJson {
			for _, zayavka := range zakupka.Protokol.Zayavki {
				var d distributor
				if len(zayavka.IP.Innfl) > 0 {
					d = distributor{
						Name:  zayavka.IP.Fio,
						Email: zayavka.IP.Email,
						Phone: zayavka.IP.Telefon,
						Inn:   zayavka.IP.Innfl,
						Cost:  zayavka.Summa,
					}
				} else if len(zayavka.YuL.Inn) > 0 {
					d = distributor{
						Name:  zayavka.YuL.RukFIO,
						Email: zayavka.YuL.Email,
						Phone: zayavka.YuL.Telefon,
						Inn:   zayavka.YuL.Inn,
						Cost:  zayavka.Summa,
					}
				} else {
					continue
				}
				distributors = append(distributors, d)
			}
		}
	}

	return distributors, nil
}

type zakazch struct {
	Inn      string `json:"ИНН"`
	NaimPoln string `json:"НаимПолн"`
}

type regRec struct {
	DataPubl     string  `json:"ДатаПубл"`
	DataOkonch   string  `json:"ДатаОконч"`
	VremyaOkonch string  `json:"ВремяОконч"`
	Produkt      string  `json:"Продукт"`
	Zakazchik    zakazch `json:"Заказчик"`
	Summa        float64 `json:"Сумма"`
	Status       string  `json:"Статус"`
}

func Search(query SearchQuery, task *Tasks.Task) {
	tasKStatus := Tasks.TaskStatusDone
	taskProgress := "Готово"
	defer func() {
		task.Status = tasKStatus
		task.Progress = taskProgress
	}()
	task.Status = Tasks.TaskStatusActive

	if !damiaConf.Active {
		tasKStatus = Tasks.TaskStatusError
		taskProgress = "Не поддерживается"
	}

	params := url.Values{}
	params.Add(KeyParam, damiaConf.Key)

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

	task.Progress = "Поиск поставок..."
	request := DamiaBaseUrl + DamiaSearchPath + "?" + params.Encode()
	response, err := http.Get(request)
	if err != nil {
		taskProgress = "Не удалось сделать запрос"
		tasKStatus = Tasks.TaskStatusError
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		taskProgress = "Плохой ответ от сервера"
		tasKStatus = Tasks.TaskStatusError
		return
	}

	dec := json.NewDecoder(strings.NewReader(string(responseBody)))

	task.Progress = "Собираем данные..."
	distributors := make([]distributor, 0)
	for {
		var responseJson map[string]map[string]regRec
		if err := dec.Decode(&responseJson); err == io.EOF {
			break
		} else if err != nil && len(responseJson) == 0 {
			taskProgress = "Ошибка парсинга"
			tasKStatus = Tasks.TaskStatusError
			return
		}

		task.Result = make([]Tasks.TaskResult, 0)
		for _, regKeys := range responseJson {
			for zakup_regn, _ := range regKeys {
				result := Tasks.TaskResult{}
				log.Print("Retrieving personal information, zakup_regn: ", zakup_regn)
				distrs, err := getDistributors(zakup_regn)
				if err != nil {
					fmt.Errorf("Unable to get distributors of " + zakup_regn)
					continue
				}

				for _, distrib := range distrs {
					distributors = append(distributors, distrib)
				}
				task.Result = append(task.Result, result)
			}
		}
	}

	task.Progress = "Ищем лучших поставщиков..."
	for _, distr := range distributors {
		result := Tasks.TaskResult{}
		log.Print("Checking for unscrupulous, inn: ", distr.Inn)
		err := CheckUnscrupulousOrganisation(distr.Inn, &result)
		if err != nil {
			log.Print("Cannot check unscrupulous, inn: ", distr.Inn)
		}

		result.CompanyName = distr.Name
		result.ContactPersons = append(result.ContactPersons, distr.Name)
		result.Emails = append(result.Emails, distr.Email)
		result.Phones = append(result.Phones, distr.Phone)
		result.AverageCapitalization = strconv.Itoa(distr.Cost)

		task.Result = append(task.Result, result)
	}

}
