// File TestGitHubData.go
// A file used for small test utilities
package main

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
	"github.com/sirupsen/logrus"
)

func displayData(data map[string]interface{}, log logrus.FieldLogger) {
	for key, valueList := range data {
		log.Info("Printing data --- key = ", key, " - value = ", valueList)
	}
}

func (client *GitHubClient) printFirst(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	item := client.response.Items[0]
	displayData(item, log)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}
