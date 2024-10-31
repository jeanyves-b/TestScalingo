package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := newConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration")
		os.Exit(1)
	}
	// Initialize web server and configure the following routes:
	// GET /repos
	client, err := newGitHubClient()
	if err != nil {
		log.WithError(err).Error("Fail to create http request")
		os.Exit(1)
	}

	
	// Lancer la thread qui va emettre des requète toutes les 5 secondes de façon à peupler puis mettre à jour la structure de données
	go func(client *GitHubClient) { 
		for {
			select {
				case <- client.timer.C:
					client.getLastPublicGithubRepositories()
			}
		}
	}(client)
	
	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", pongHandler)
	router.HandleFunc("/getAll", client.response.getAll)
	router.HandleFunc("/getFiltered", client.response.getFiltered)
	router.HandleFunc("/getAllUpdated", client.getAllUpdated)
	router.HandleFunc("/getFilteredUpdated", client.getFilteredUpdated)

	//méthodes de test
	router.HandleFunc("/PrintFirst", client.printFirst)
	
	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}

}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "pong"})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}