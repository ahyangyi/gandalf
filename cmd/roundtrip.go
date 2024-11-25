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
	var keeps []string
	flag.Var((*stringsFlag)(&keeps), "keep", "List of object names to keep -- only applicable if non-empty")

	// Parse flags
	flag.Parse()

	obj, err := magica.FromFileWithLayersAndFilter(*source, []int{}, func(layerName string) bool {
		if len(keeps) > 0 {
			for _, keptName := range keeps {
				if keptName == layerName {
					return true
				}
			}
			return false
		}
		for _, discardedName := range discards {
			if discardedName == layerName {
				return false
			}
		}
		return true
	})
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
