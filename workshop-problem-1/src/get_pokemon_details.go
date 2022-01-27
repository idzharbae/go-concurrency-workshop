package src

import (
	"encoding/json"
	"os"
	"time"
)

func GetPokemonDetailsByName(name string) (PokemonDetails, error) {
	var pokemonResponse PokemonDetails

	// url := baseURL + fmt.Sprintf("/api/v2/pokemon-species/%s/", name)
	// resp, err := http.Get(url)

	// responseBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return PokemonDetails{}, err
	// }

	time.Sleep(time.Second / 10) // Simulate latency
	responseBody, err := os.ReadFile("./dummy_pokemon.json")
	if err != nil {
		return PokemonDetails{}, err
	}

	err = json.Unmarshal(responseBody, &pokemonResponse)
	if err != nil {
		return PokemonDetails{}, err
	}

	return pokemonResponse, nil
}
