package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"

	"github.com/Scalingo/go-utils/logger"
)
type DecoderSearchResult struct {
	TotalCount int                      `json:"total_count"`
	IncompleteResults	bool			`json:"incomplete_results"`
	Items      []json.RawMessage		`json:"items"`
}

type SearchResult struct {
	TotalCount int                      `json:"total_count"`
	IncompleteResults	bool			`json:"incomplete_results"`
	Items      []map[string]interface{} `json:"items"`
	reader		sync.WaitGroup
	writer		sync.WaitGroup
}

func (data *SearchResult) getAll(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data.writer.Wait()
	data.reader.Add(1)
	defer data.reader.Done()
	err := json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "Items": data.Items})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

func isType(a, b interface{}) bool {
    return reflect.TypeOf(a) == reflect.TypeOf(b)
}
func contains(element interface{}, value string) bool {
	if isType(element, value){
		return element.(string) == value
	}
	valueInt, err := strconv.Atoi(value)
	if err == nil && isType(element, valueInt) {

		return element == valueInt
	}
	return false
}

func (data *SearchResult) filterResults(filters url.Values) []map[string]interface{} {
	log := logger.Default()
	var filteredArray []map[string]interface{}
	log.Info("Filtered List : filtering list with those filters: -- ", filters)

	// Parcours de chaque élément dans le tableau source
	for _, element := range data.Items {
		matches := true

		// Vérification de chaque critère
		for criteriakey, criteriaValueList := range filters {
			for _ , criteriaValue := range criteriaValueList {
				// Si la clé existe et que la valeur correspond dans l'élément
				if value, ok := element[criteriakey]; !ok || !contains(value, criteriaValue) {
					matches = false
					break
				}
			}
		}

		// Si tous les critères correspondent, ajout de l'élément au tableau filtré
		if matches {
			filteredArray = append(filteredArray, element)
		}
	}

	return filteredArray
}

func (data *SearchResult) getFiltered(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	value, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.WithError(err).Error("Error parsing request")
		err = json.NewEncoder(w).Encode(map[string]string{"status": "ERROR"})
		if err != nil {
			log.WithError(err).Error("Fail to encode JSON")
		}
		return err
	}
	log.Info("filtering with those filters: ", value)

	data.writer.Wait()
	data.reader.Add(1)
	defer data.reader.Done()
	err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK","Items": data.filterResults(value)})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}