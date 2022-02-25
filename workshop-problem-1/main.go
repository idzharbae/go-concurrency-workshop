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
	MaxAsyncProcessGetAvailablePromoLogistic := 5
	var pokemonList []src.PokemonResult

	var (
		wg        sync.WaitGroup
		asyncChan = make(chan bool, MaxAsyncProcessGetAvailablePromoLogistic)
	)

	limit := 10
	offset := 0

	// Get all pokemons
	for {
		pokemonListResponse, err := src.ListPokemon(limit, offset)
		if err != nil {
			log.Fatal(err)
		}

		pokemonList = append(pokemonList, pokemonListResponse.Results...)

		if *debugFlag {
			log.Printf("Fetched %d pokemons out of %d\n", len(pokemonList), pokemonListResponse.Count)
		}

		offset += 10
		if offset > pokemonListResponse.Count {
			break
		}
	}

	// Get each pokemon details
	for _, pokemonFromList := range pokemonList {

		wg.Add(1)

		go func(pokemon *src.PokemonResult, wg *sync.WaitGroup) {
			defer func() {
				if errPanic := recover(); errPanic != nil {
					log.Fatal(errPanic)
				}
				wg.Done()
			}()

			pokemonDetail, err := src.GetPokemonDetailsByName(pokemon.Name)
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
			asyncChan <- true

			if *debugFlag {
				log.Printf("Get detail pokemon %s\n", pokemonDetail.Name)
			}

			<-asyncChan

		}(&pokemonFromList, &wg)

	}
	wg.Wait()

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
