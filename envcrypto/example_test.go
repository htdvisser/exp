package envcrypto_test

import (
	"fmt"
	"os"
	"path/filepath"

	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"golang.org/x/exp/slices"
	"htdvisser.dev/exp/envcrypto"
)

func Example() {
	var allFiles []string
	matches, err := filepath.Glob("testdata/example/*.env")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		if !slices.Contains(allFiles, match) {
			allFiles = append(allFiles, match)
		}
	}
	envFilesSource, err := envcrypto.NewEnvFilesSource(nil, allFiles...)
	if err != nil {
		panic(err)
	}

	os.Setenv("SOPS_AGE_KEY_FILE", "testdata/example/age-key.txt")

	fileSource := envFilesSource.GetFile("testdata/example/sops.env")
	if fileSource == nil {
		fileSource, err = envcrypto.NewEnvFileSource(nil, "testdata/example/sops.env")
		if err != nil {
			panic(err)
		}
		envFilesSource.AppendSource(fileSource)
	}
	err = fileSource.Replace(func(data []byte) ([]byte, error) {
		return decrypt.DataWithFormat(data, formats.Dotenv)
	})
	if err != nil {
		panic(err)
	}

	box, err := envcrypto.Open(envFilesSource)
	if err != nil {
		panic(err)
	}

	exampleValue, err := box.Get("DOTENV_EXAMPLE")
	if err != nil {
		panic(err)
	}

	fmt.Println(exampleValue)

	// Output:
	// hello dotenv
}
