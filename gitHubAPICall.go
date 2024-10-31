package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/Scalingo/go-utils/logger"
)

// Structure de requète réutilisable
type GitHubClient struct {
	httpClient *http.Client
	response   SearchResult
}

// Fonction pour créer un nouveau client GitHub
func newGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
	}
}

// Fonction pour faire une requête à l'API GitHub
func (client *GitHubClient) getLastPublicGithubRepositories() error {
	//initialisation du loggeur
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
		return err
	}

	// Ajouter les en-têtes
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Envoyer la requête
	resp, err := client.httpClient.Do(req)
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

	// Fermer le channel après que toutes les goroutines ont terminé
	go func() {
		wg.Wait()
		close(results)
		log.Info("parsing done")
	}()

	// Collecter et afficher les résultats
	var itemList []map[string]interface{}
	for item := range results {
		itemList = append(itemList, item)
	}
	client.response.Items = itemList
	log.Info("List filled")

	return nil
}
