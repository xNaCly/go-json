package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var depths = map[string]int{
	"1K":   1000,
	"10K":  10000,
	"100K": 100000,
	"1M":   1000000,
	"10M":  10000000,
}

func main() {
	for depthName, depth := range depths {
		fmt.Printf("Generating %d depth object\n", depth)

		jsonString := "{"

		for i := 1; i < depth; i++ {
			jsonString += `"next":{`
		}

		jsonString += `"next":null`

		for i := 0; i < depth; i++ {
			jsonString += "}"
		}

		fileName := fmt.Sprintf("%s_recursion.json", depthName)
		filePath := filepath.Join(filepath.Dir(os.Args[0]), fileName)

		err := os.WriteFile(filePath, []byte(jsonString), 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			continue
		}

		fmt.Printf("File for depth %d saved as %s\n", depth, fileName)
	}
}
