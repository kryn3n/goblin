package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/kryn3n/goblin/internal/arcgis"
	"github.com/kryn3n/goblin/internal/misc"
)

func main() {
	now := time.Now()

	defer func() {
		fmt.Println(time.Since(now))
	}()

	var server arcgis.Server

	flag.IntVar(&server.RecordLimit, "m", 1000, "record limit")
	flag.IntVar(&server.ConcurrencyLimit, "c", 100, "concurrency limit")
	flag.IntVar(&server.LayerId, "l", -1, "layer ID")
	flag.Parse()

	server.URL = flag.Arg(0)
	fileName := flag.Arg(1)

	if server.URL == "" {
		error := errors.New("please provide server URL as first argument")
		log.Fatalln(error)
	} else if fileName == "" {
		error := errors.New("please provide output filename as second argument")
		log.Fatalln(error)
	}

	misc.WelcomeMessage()

	server.GetServerInfo()

	if server.LayerId == -1 {
		fmt.Println("Which layer would you like to query?")
		for i := range server.Layers {
			server.LayerIds = append(server.LayerIds, server.Layers[i].Id)
			fmt.Printf("* %d %s\n", server.Layers[i].Id, server.Layers[i].Name)
		}
		fmt.Println("Layer ID: ")
		fmt.Scanln(&server.LayerId)
		if !slices.Contains(server.LayerIds, server.LayerId) {
			log.Fatal("Layer does not exist")
		}
	}
	server.GetLayerURL()
	server.GetRecordCount()
	server.GetObjectIds()

	fmt.Printf("Total Records: %d\n", server.RecordCount)
	fmt.Printf("Total Object IDs: %d\n", len(server.ObjectIds))
	fmt.Printf("Object ID Field: %s\n", server.ObjectIdField)
	fmt.Printf("Max Records Per Request: %d\n", server.MaxRecordCount)

	server.GetBatches()
	server.GetData(fileName)
}
