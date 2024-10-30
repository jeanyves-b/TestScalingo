package main

type SearchResult struct {
	TotalCount int                      `json:"total_count"`
	items      []map[string]interface{} `json:"items"`
}

func (data *SearchResult) filterResults(filters map[string]interface{}) []map[string]interface{} {
	var filteredArray []map[string]interface{}

	// Parcours de chaque élément dans le tableau source
	for _, element := range data.items {
		matches := true

		// Vérification de chaque critère
		for criteriakey, criteriaValue := range filters {
			// Si la clé existe et que la valeur correspond dans l'élément
			if value, ok := element[criteriakey]; !ok || value != criteriaValue {
				matches = false
				break
			}
		}

		// Si tous les critères correspondent, ajout de l'élément au tableau filtré
		if matches {
			filteredArray = append(filteredArray, element)
		}
	}

	return filteredArray
}
