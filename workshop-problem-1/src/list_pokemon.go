package src

import (
	"encoding/json"
	"os"
	"time"
)

func ListPokemon(limit, offset int) (ListPokemonResponse, error) {
	var pokemonResponse ListPokemonResponse

	// url := baseURL + fmt.Sprintf("/api/v2/pokemon-species/?limit=%d&offset=%d", limit, offset)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return ListPokemonResponse{}, err
	// }

	// responseBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return ListPokemonResponse{}, err
	// }

	time.Sleep(time.Second / 10) // Simulate latency
	responseBody, err := os.ReadFile("./dummy_pokemon_list.json")
	if err != nil {
		return ListPokemonResponse{}, err
	}

	err = json.Unmarshal(responseBody, &pokemonResponse)
	if err != nil {
		return ListPokemonResponse{}, err
	}

	if len(pokemonResponse.Results) > pokemonResponse.Count-offset {
		pokemonResponse.Results = pokemonResponse.Results[:(pokemonResponse.Count - offset)]
	}

	return pokemonResponse, nil
}
