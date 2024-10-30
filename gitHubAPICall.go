package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// Structure de requète réutilisable
type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
	response   SearchResult
}

// Fonction pour créer un nouveau client GitHub
func newGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
		baseURL:    "https://api.github.com/search/repositories",
	}
}

// Fonction pour faire une requête à l'API GitHub
func (client *GitHubClient) getLastPublicGithubRepositories() error {
	//on s'assure que Done soit appeler

	gitHubUrl := "https://api.github.com/search/repositories"

	// Ajout les paramètres de requête
	params := url.Values{}
	params.Add("public", "q=is:public")
	params.Add("sort", "sort=created")
	params.Add("order", "order=desc")
	params.Add("per_page", "per_page=100")

	// construction de la requete
	urlWithParams := gitHubUrl + "?" + params.Encode()
	req, err := http.NewRequest("GET", urlWithParams, nil)
	if err != nil {
		return err
	}

	// Ajouter les en-têtes
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Envoyer la requête
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		return err
	}

	// Lire et décoder la réponse JSON
	var rawMessages []json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&rawMessages)
	if err != nil {
		fmt.Println("Erreur lors du parsing JSON :", err)
		return err
	}

	// Préparer les variables de synchronisation et un channel pour les résultats
	var wg sync.WaitGroup
	results := make(chan map[string]interface{})

	// Démarrer une goroutine pour chaque objet JSON
	// il faudrait tester si les perfs sont meilleurs en dinminuant le nombre de goroutines
	for _, raw := range rawMessages {
		wg.Add(1)
		go func(raw json.RawMessage) {
			defer wg.Done()

			// Décoder chaque `raw` en `SmallStruct`
			var item map[string]interface{}
			if err := json.Unmarshal(raw, &item); err != nil {
				fmt.Println("Erreur lors du décodage d'un élément :", err)
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
	}()

	// Collecter et afficher les résultats
	var itemList []map[string]interface{}
	for item := range results {
		itemList = append(itemList, item)
	}
	client.response.items = itemList
	return nil
}
