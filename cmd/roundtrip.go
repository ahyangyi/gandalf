package main

import (
	"flag"
	"fmt"
	"github.com/ahyangyi/gandalf/magica"
	"log"
)

func main() {
	source := flag.String("source", "", "Path to the source file")
	destination := flag.String("destination", "", "Path to the destination file")
	var discards []string
	flag.Var((*stringsFlag)(&discards), "discard", "List of object names to delete")

	// Parse flags
	flag.Parse()

	obj, err := magica.FromFileWithLayersAndDisallowedLayerNames(*source, []int{}, discards)
	if err != nil {
		log.Fatalf("could not read file: %s", err)
	}

	err = obj.SaveToFile(*destination)
	if err != nil {
		log.Fatalf("could not write file: %s", err)
	}
}

type stringsFlag []string

func (sf *stringsFlag) String() string {
	return fmt.Sprintf("%v", *sf)
}

func (sf *stringsFlag) Set(value string) error {
	*sf = append(*sf, value)
	return nil
}
