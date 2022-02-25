package src

import (
	"encoding/json"
	"io/ioutil"
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
	responseBody, err := ioutil.ReadFile("./dummy_pokemon_list.json")
	if err != nil {
		return ListPokemonResponse{}, err
	}

	err = json.Unmarshal(responseBody, &pokemonResponse)
	if err != nil {
		return ListPokemonResponse{}, err
	}

	if len(pokemonResponse.Results) > pokemonResponse.Count-offset {
		//can get index out of bound [0:-2]
		count := (pokemonResponse.Count - offset)
		if count < 0 {
			count = 0
		}
		pokemonResponse.Results = pokemonResponse.Results[:count]
	}

	return pokemonResponse, nil
}
