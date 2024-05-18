package external

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const URL = "https://development.kpi-drive.ru/_api/"

type Fact struct {
	PeriodStart         string
	PeriodEnd           string
	PeriodKey           string
	IndicatorToMoID     string
	IndicatorToMoFactID int
}

// SaveFact - сохраняет факт в базу данных
func SaveFact(fact map[string]string) (Fact, error) {
	op := "SaveFact"
	endPoint := URL + "facts/save_fact"
	req, err := genReqAuth(endPoint, fact)
	if err != nil {
		return Fact{}, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Fact{}, err
	}
	defer req.Body.Close()
	var factResponse SaveFactResponse
	err = json.NewDecoder(resp.Body).Decode(&factResponse)
	if err != nil {
		log.Println(op, err)
		return Fact{}, err
	}
	if factResponse.Status != "OK" {
		// создаем новую ошибку с ошибками, полученными в ответе
		return Fact{}, errors.New(factResponse.Messages.Error[0].(string))
	}
	log.Println(op, "fact saved succesfully ID:", factResponse.Data.IndicatorToMoFactID)

	return Fact{
		PeriodStart:         fact["period_start"],
		PeriodEnd:           fact["period_end"],
		PeriodKey:           fact["period_key"],
		IndicatorToMoID:     fact["indicator_to_mo_id"],
		IndicatorToMoFactID: factResponse.Data.IndicatorToMoFactID,
	}, nil
}

type SaveFactResponse struct {
	Data struct {
		IndicatorToMoFactID int `json:"indicator_to_mo_fact_id"`
	} `json:"DATA"`
	Status   string   `json:"STATUS"`
	Messages Messages `json:"MESSAGES"`
}

type Messages struct {
	Error []interface{} `json:"error"`
	// Warning []interface{} `json:"warning"` // we can parse it but is not used, so is not necessary
	// Info    []interface{} `json:"info"` 		// ...
}

// IsFactSaved - проверяет, сохранен ли факт
func IsFactSaved(fact *Fact) (bool, error) {
	op := "IsFactSaved"
	endPoint := URL + "indicators/get_facts"

	params := make(map[string]string)
	params["period_start"] = fact.PeriodStart
	params["period_end"] = fact.PeriodEnd
	params["period_key"] = fact.PeriodKey
	params["indicator_to_mo_id"] = fact.IndicatorToMoID

	req, err := genReqAuth(endPoint, params)
	if err != nil {
		log.Println(op, err)
		return false, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(op, err, resp.StatusCode)
		return false, err
	}

	defer req.Body.Close()

	var isFactSavedResponse IsFactSavedResponse
	err = json.NewDecoder(resp.Body).Decode(&isFactSavedResponse)
	if err != nil {
		log.Println(op, err)
		return false, err
	}
	if len(isFactSavedResponse.Messages.Error) != 0 {
		log.Println(isFactSavedResponse.Messages.Error...)
		return false, errors.New("cannot check if fact is saved")
	}
	for _, v := range isFactSavedResponse.Data.Rows {
		if v.IndicatorToMoFactID == fact.IndicatorToMoFactID {
			return true, nil
		}
	}
	return false, err
}

type IsFactSavedResponse struct {
	Messages Messages `json:"MESSAGES"`
	Data     struct {
		// Page int `json:"page"` 							// ...
		// PagesCount int `json:"pages_count"`  // ...
		RowsCount int `json:"rows_count"`
		Rows      []struct {
			IndicatorToMoFactID int `json:"indicator_to_mo_fact_id"` // the other fiels are not necessary
			//...
		} `json:"rows"`
	} `json:"DATA"`
}
