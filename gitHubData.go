package main

import (
	"encoding/json"
	"net/http"

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
}

func (data *SearchResult) getAll(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data.Items)
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

func (data *SearchResult) filterResults(filters map[string]string) []map[string]interface{} {
	log := logger.Get(r.Context())
	var filteredArray []map[string]interface{}

	// Parcours de chaque élément dans le tableau source
	for _, element := range data.Items {
		matches := true

		// Vérification de chaque critère
		for criteriakey, criteriaValue := range filters {
			// Si la clé existe et que la valeur correspond dans l'élément
			if value, ok := element[criteriakey]; !ok || value != criteriaValue {
				matches = false
				log.Info("Filtered List : filtered element")
				break
			}
		}

		// Si tous les critères correspondent, ajout de l'élément au tableau filtré
		if matches {
			filteredArray = append(filteredArray, element)
			log.Info("Filtered List : added element")
		}
	}

	return filteredArray
}

func (data *SearchResult) getFiltered(w http.ResponseWriter, r *http.Request, filters map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Info("filtering with those filters: ", filters)

	err := json.NewEncoder(w).Encode(data.filterResults(filters))
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}