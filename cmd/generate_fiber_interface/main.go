//go:generate go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
//go:generate go run main.go
// github.com/getkin/kin-openapi v0.123.0

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/deepmap/oapi-codegen/v2/pkg/codegen"
	"github.com/deepmap/oapi-codegen/v2/pkg/util"
)

const yamlFiles = "./src/interface/fiber_server/spec/*.yaml"

func generateHTTPFiber() error {
	files, err := filepath.Glob(yamlFiles)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := filepath.Base(file[:len(file)-len(filepath.Ext(file))])
		packageName := sanitizePackageName(filename) + "_spec"
		outputPath := fmt.Sprintf("./src/interface/fiber_server/spec/%s/spec.gen.go", filename)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}

		swagger, err := util.LoadSwaggerWithCircularReferenceCount(file, 3)
		if err != nil {
			return err
		}

		conf := codegen.Configuration{
			Generate: codegen.GenerateOptions{
				EchoServer:   false,
				FiberServer:  true,
				Client:       false,
				Models:       true,
				EmbeddedSpec: true,
				Strict:       true,
			},
			PackageName: packageName,
		}

		gen, err := codegen.Generate(swagger, conf)
		if err != nil {
			return err
		}

		//	 Write to file `{{outputPath}}/spec.gen.go`
		if err := os.WriteFile(outputPath, []byte(gen), 0644); err != nil {
			return err
		}
	}

	return nil
}

func sanitizePackageName(input string) string {
	var result strings.Builder

	for i, char := range input {
		// Check if the character is an English letter (a-z, A-Z) or underscore
		if unicode.IsLetter(char) || char == '_' {
			// If the character is uppercase and not the first character, add an underscore before it
			if unicode.IsUpper(char) && i > 0 {
				result.WriteRune('_')
			}

			// Convert the character to lowercase and append to the result
			result.WriteRune(unicode.ToLower(char))
		}

		if char == '-' {
			result.WriteString("_")
		}
	}

	return result.String()
}

func main() {
	if err := generateHTTPFiber(); err != nil {
		fmt.Printf("Error generating HTTP Fiber code: %v\n", err)
		os.Exit(1)
	}
}
