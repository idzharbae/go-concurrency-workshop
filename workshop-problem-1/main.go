package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/idzharbae/go-concurrency-workshop/workshop-problem-1/src"
	"go.uber.org/atomic"
)

var debugFlag *bool
var numberOfPokemons atomic.Int32

func main() {
	debugFlag = flag.Bool("debug", false, "toggle debug log")
	flag.Parse()

	now := time.Now()

	maxWorkers := 80

	pokemonNameChan := make(chan string)
	pokemonDetailChan := make(chan src.PokemonDetails)

	// Get all pokemons
	limit := 10
	offset := 0

	// Get first page for max pokemon
	pokemonListResponse, err := src.ListPokemon(limit, offset)
	if err != nil {
		log.Fatal(err)
	}

	maxPokemon := pokemonListResponse.Count

	maxPokemonNameGenerators := maxWorkers / 10 // Similar runtime with maxGenerators = maxWorkers
	if maxPokemonNameGenerators == 0 {
		maxPokemonNameGenerators = 1
	}

	generatorLimiter := make(chan bool, maxPokemonNameGenerators)

	// Spawn pokemon name generators
	go func() {
		defer close(pokemonNameChan)

		for offset := 0; offset < maxPokemon; offset += limit {
			// Take one limiter slot
			generatorLimiter <- true
			go func(limit, offset int) {
				defer func() {
					<-generatorLimiter
				}()

				pokemonNameGenerator(pokemonNameChan, limit, offset)
			}(limit, offset)
		}

		// Wait until all go routines are done
		for i := 0; i < maxPokemonNameGenerators; i++ {
			generatorLimiter <- true
		}
	}()

	// Get each pokemon details
	var wg sync.WaitGroup

	wg.Add(maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wg.Done()
			pokemonDetailGenerator(pokemonNameChan, pokemonDetailChan)
		}()
	}

	// Close chan after all pokemons are fetched
	go func() {
		defer close(pokemonDetailChan)
		wg.Wait()
	}()

	var wgSave sync.WaitGroup
	wgSave.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func() {
			defer wgSave.Done()

			for pokemonDetail := range pokemonDetailChan {
				err := SavePokemonDummy(pokemonDetail)
				if err != nil {
					log.Fatal(err)
				}

				numberOfPokemons.Add(1)
				if *debugFlag {
					log.Printf("Saved %d out of %d pokemons.", numberOfPokemons, 898)
				}
			}
		}()
	}

	wgSave.Wait()

	log.Printf("Fetched %d pokemons in %v!\n", numberOfPokemons.Load(), time.Since(now))

	PrintMemUsage()
}

func pokemonDetailGenerator(pokemonNameChan <-chan string, pokemonDetailChan chan<- src.PokemonDetails) {
	for pokemonName := range pokemonNameChan {
		pokemonDetail, err := src.GetPokemonDetailsByName(pokemonName)
		if err != nil {
			log.Fatal(err)
		}

		pokemonDetailChan <- pokemonDetail
	}
}

func pokemonNameGenerator(pokemonNameChan chan<- string, limit, offset int) {
	pokemonListResponse, err := src.ListPokemon(limit, offset)
	if err != nil {
		log.Fatal(err)
	}

	for _, pokemon := range pokemonListResponse.Results {
		pokemonNameChan <- pokemon.Name
	}
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
