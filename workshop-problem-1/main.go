package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/idzharbae/go-concurrency-workshop/workshop-problem-1/src"
)

func main() {
	debugFlag := flag.Bool("debug", false, "toggle debug log")
	flag.Parse()

	now := time.Now()

	pokemonListCh := make(chan src.PokemonResult)

	limit := 10
	offset := 0
	count := 0

	go func() {
		defer close(pokemonListCh)

		// Get all pokemons
		for {
			pokemonListResponse := getPokemonList(debugFlag, limit, offset)

			count += len(pokemonListResponse.Results)

			if *debugFlag {
				log.Printf("Fetched %d pokemons out of %d\n", count, pokemonListResponse.Count)
			}

			for _, pokemonListResult := range pokemonListResponse.Results {
				pokemonListCh <- pokemonListResult
			}

			offset += 10
			if offset > pokemonListResponse.Count {
				break
			}
		}
	}()

	// Get each pokemon details
	for pokemonFromList := range pokemonListCh {
		go getPokemonDetail(debugFlag, pokemonFromList)
	}

	log.Printf("Fetched %d pokemons in %v!\n", count, time.Since(now))
	PrintMemUsage()
}

func SavePokemonDummy(pokemon src.PokemonDetails) error {
	// Saving to DB
	time.Sleep(time.Second / 10)

	return nil
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getPokemonList(debugFlag *bool, limit, offset int) (pokemonList src.ListPokemonResponse) {
	pokemonListResponse, err := src.ListPokemon(limit, offset)
	if err != nil {
		log.Fatal(err)
	}

	return pokemonListResponse
}

func getPokemonDetail(debugFlag *bool, pokemonFromList src.PokemonResult) {
	if *debugFlag {
		log.Printf("Get detail pokemon %s\n", pokemonFromList.Name)
	}

	pokemonDetail, err := src.GetPokemonDetailsByName(pokemonFromList.Name)
	if err != nil {
		log.Fatal(err)
	}

	go setPokemonDetail(pokemonDetail)
}

func setPokemonDetail(pokemonDetail src.PokemonDetails) {
	err := SavePokemonDummy(pokemonDetail)
	if err != nil {
		log.Fatal(err)
	}
}
