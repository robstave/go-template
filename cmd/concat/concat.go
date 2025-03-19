package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Define a structure that matches the YAML file
type FileGroups map[string][]string

func main() {
	// Specify the YAML file containing the grouped file paths
	yamlFile := "file_list.yaml" // adjust this to your YAML file
	// Define the output directory (optional)
	outputDir := "output" // all combined files will be placed here

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Open the YAML file
	file, err := os.Open(yamlFile)
	if err != nil {
		log.Fatalf("Failed to open YAML file: %v", err)
	}
	defer file.Close()

	// Parse the YAML file
	var groups FileGroups
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&groups); err != nil {
		log.Fatalf("Failed to parse YAML file: %v", err)
	}

	// Iterate over each group
	for groupName, filePaths := range groups {
		// Define the output file path for this group
		outputFilePath := filepath.Join(outputDir, fmt.Sprintf("%s_combined.txt", groupName))

		// Create or truncate the output file
		outFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Printf("Failed to create output file for group '%s': %v", groupName, err)
			continue
		}

		log.Printf("Processing group: %s", groupName)

		// Iterate over each file in the current group
		for _, filePath := range filePaths {
			trimmedPath := strings.TrimSpace(filePath)

			// Skip empty lines or comments
			if trimmedPath == "" || strings.HasPrefix(trimmedPath, "#") {
				continue
			}

			log.Printf("Reading file: %s", trimmedPath)

			// Open the source file
			srcFile, err := os.Open(filepath.Clean(trimmedPath))
			if err != nil {
				log.Printf("Failed to open file '%s': %v", trimmedPath, err)
				continue
			}

			// Copy the contents to the output file
			if _, err := io.Copy(outFile, srcFile); err != nil {
				log.Printf("Failed to write to output file '%s': %v", outputFilePath, err)
			}

			// Close the source file
			srcFile.Close()

			// Optionally, add a separator between file contents
			separator := fmt.Sprintf("\n\n=== End of %s ===\n\n", trimmedPath)
			if _, err := outFile.WriteString(separator); err != nil {
				log.Printf("Failed to write separator to output file '%s': %v", outputFilePath, err)
			}
		}

		// Close the output file
		if err := outFile.Close(); err != nil {
			log.Printf("Failed to close output file '%s': %v", outputFilePath, err)
		} else {
			log.Printf("Successfully created '%s'", outputFilePath)
		}
	}

	fmt.Println("All files successfully concatenated into respective grouped output files.")
}
