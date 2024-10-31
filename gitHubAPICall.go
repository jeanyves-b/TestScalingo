package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Scalingo/go-utils/logger"
)

// Structure de requète réutilisable
type GitHubClient struct {
	httpClient	*http.Client
	request		*http.Request
	response	SearchResult
	timer		*time.Ticker
}

// Fonction pour créer un nouveau client GitHub
func newGitHubClient() (*GitHubClient, error) {
	timer := time.NewTicker(10*time.Second)
	log := logger.Default()
	gitHubUrl := "https://api.github.com/search/repositories?q=is:public"

	// Ajout les paramètres de requête
	params := url.Values{}
	params.Add("sort", "created")
	params.Add("order", "desc")
	params.Add("per_page", "100")

	// construction de la requete
	urlWithParams := gitHubUrl + "&" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, urlWithParams, nil)
	log.Info("Asking github api : " + urlWithParams)
	if err != nil {
		log.WithError(err).Error("Failed to create request")
		return nil, err
	}

	// Ajouter les en-têtes
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return &GitHubClient{
		httpClient: &http.Client{},
		request: req,
		timer: timer,
	}, nil
}

// Fonction pour faire une requête à l'API GitHub
func (client *GitHubClient) getLastPublicGithubRepositories() error {
	//initialisation du loggeur
	log := logger.Default()
	defer client.timer.Reset(10*time.Second)

	// Envoyer la requête
	resp, err := client.httpClient.Do(client.request)
	if err != nil {
		log.WithError(err).Error("Fail send request", "message: ", resp)
		return err
	}
	defer resp.Body.Close()

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		log.WithError(err).Error("Response not OK", "message: ", resp)
		return err
	}

	// Lire et décoder la réponse JSON
	var rawMessages DecoderSearchResult
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Error while reading the response")
		return err
	}
	err = json.Unmarshal(body, &rawMessages)
	if err != nil {
		log.WithError(err).Error("Fail to parse JSON")
		return err
	}
	client.response.TotalCount = rawMessages.TotalCount
	client.response.IncompleteResults = rawMessages.IncompleteResults
	log.Info("first parse done")

	// Préparer les variables de synchronisation et un channel pour les résultats
	var wg sync.WaitGroup
	results := make(chan map[string]interface{})

	log.Info("creating threads")
	// Démarrer une goroutine pour chaque objet JSON
	// il faudrait tester si les perfs sont meilleurs en dinminuant le nombre de goroutines
	for _, raw := range rawMessages.Items {
		wg.Add(1)
		go func(raw json.RawMessage) {
			defer wg.Done()
			log := logger.Default()

			// Décoder chaque `raw` en `SmallStruct`
			var item map[string]interface{}
			if err := json.Unmarshal(raw, &item); err != nil {
				log.WithError(err).Error("Fail to parse element")
				return
			}

			// Envoyer le résultat dans le channel
			results <- item
		}(raw)
	}

	// Fermer le channel après que toutes les goroutines aient terminées
	go func() {
		wg.Wait()
		close(results)
		log.Info("parsing done")
	}()

	// Collecter et afficher les résultats
	client.response.writer.Wait()
	client.response.writer.Add(1)
	defer client.response.writer.Done()
	client.response.reader.Wait()
	var itemList []map[string]interface{}
	for item := range results {
		itemList = append(itemList, item)
	}
	client.response.Items = itemList
	log.Info("List filled")

	return nil
}

func (client *GitHubClient) getAllUpdated(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	client.getLastPublicGithubRepositories()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	client.response.writer.Wait()
	client.response.reader.Add(1)
	defer client.response.reader.Done()
	err := json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "Items": client.response})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}

func (client *GitHubClient) getFilteredUpdated(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	client.getLastPublicGithubRepositories()

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
	
	client.response.writer.Wait()
	client.response.reader.Add(1)
	defer client.response.reader.Done()
	err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK","Items": client.response.filterResults(value)})
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
	}
	return nil
}