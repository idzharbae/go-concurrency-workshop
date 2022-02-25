package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/idzharbae/go-concurrency-workshop/workshop-problem-1/src"
)

func main() {
	debugFlag := flag.Bool("debug", false, "toggle debug log")
	flag.Parse()

	now := time.Now()

	var pokemonList []src.PokemonResult

	limit := 10
	offset := 0
	respCount := 0
	chanPokemonList := make(chan []src.PokemonResult)

	var wg sync.WaitGroup

	//get count first
	pokemonListResponse, err := src.ListPokemon(limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	// respCount = 898
	respCount = pokemonListResponse.Count

	// Get all pokemons
	for {
		if offset > respCount {
			offset = respCount
		}

		go func() {
			wg.Add(1)
			defer wg.Done()

			pokemonListResponse, err := src.ListPokemon(limit, offset)
			if err != nil {
				log.Fatal(err)
			}

			chanPokemonList <- pokemonListResponse.Results

			if *debugFlag {
				log.Printf("Fetched %d pokemons out of %d\n", len(pokemonList), pokemonListResponse.Count)
			}
		}()

		if offset >= respCount {
			break
		}
		offset += 10
	}

	go func() {
		defer close(chanPokemonList)
		wg.Wait()
	}()

	for list := range chanPokemonList {
		pokemonList = append(pokemonList, list...)
	}

	// Get each pokemon details
	for _, pokemonFromList := range pokemonList {
		go func(name string) {
			wg.Add(1)
			defer wg.Done()

			pokemonDetail, err := src.GetPokemonDetailsByName(name)
			if err != nil {
				log.Fatal(err)
			}

			if *debugFlag {
				log.Printf("Get detail pokemon %s\n", pokemonDetail.Name)
			}

			err = SavePokemonDummy(pokemonDetail)
			if err != nil {
				log.Fatal(err)
			}
		}(pokemonFromList.Name)
	}

	go func() {
		wg.Wait()
	}()

	log.Printf("Fetched %d pokemons in %v!\n", len(pokemonList), time.Since(now))
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
