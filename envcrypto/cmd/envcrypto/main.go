package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/decrypt"
	"golang.org/x/exp/slices"
	"htdvisser.dev/exp/envcrypto"
)

func main() {
	app := &App{}
	cmd := app.rootCmd()
	if err := cmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

type App struct {
	printNewline bool

	envFilesSource *envcrypto.EnvFilesSource

	box *envcrypto.Box
}

func (app *App) writeString(s string) {
	if app.printNewline {
		fmt.Println(s)
	} else {
		fmt.Print(s)
	}
}

func (app *App) loadConfig(cmd *cobra.Command, _ []string) error {
	doNotPrintNewline, err := cmd.Flags().GetBool("n")
	if err != nil {
		return err
	}
	app.printNewline = !doNotPrintNewline
	return nil
}

func (app *App) loadFiles(cmd *cobra.Command, _ []string) error {
	files, err := cmd.Flags().GetStringSlice("files")
	if err != nil {
		return err
	}
	var allFiles []string
	for _, file := range files {
		matches, err := filepath.Glob(file)
		if err != nil {
			return err
		}
		for _, match := range matches {
			if !slices.Contains(allFiles, match) {
				allFiles = append(allFiles, match)
			}
		}
	}
	app.envFilesSource, err = envcrypto.NewEnvFilesSource(nil, allFiles...)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) loadSOPS(cmd *cobra.Command, _ []string) error {
	sops, err := cmd.Flags().GetString("sops")
	if err != nil {
		return err
	}
	if sops != "" {
		fileSource := app.envFilesSource.GetFile(sops)
		if fileSource == nil {
			fileSource, err = envcrypto.NewEnvFileSource(nil, sops)
			if err != nil {
				return err
			}
			app.envFilesSource.AppendSource(fileSource)
		}
		err = fileSource.Replace(func(data []byte) ([]byte, error) {
			return decrypt.DataWithFormat(data, formats.Dotenv)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) openBox(_ *cobra.Command, _ []string) error {
	box, err := envcrypto.Open(envcrypto.MultiSource{
		envcrypto.EnvSource{},
		app.envFilesSource,
	})
	if err != nil {
		return err
	}
	app.box = box
	return nil
}

func multiRunE(funcs ...func(cmd *cobra.Command, _ []string) error) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		for _, f := range funcs {
			if err := f(cmd, nil); err != nil {
				return err
			}
		}
		return nil
	}
}

func (app *App) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "envcrypto",
		Short:         "envcrypto is a tool for encrypting and decrypting environment variables",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolP("n", "n", false, "do not print the trailing newline character")

	cmd.AddCommand(app.generateCmd())
	cmd.AddCommand(app.encryptCmd())
	cmd.AddCommand(app.decryptCmd())
	cmd.AddCommand(app.getCmd())

	return cmd
}

func (app *App) generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a new env crypto",
		PreRunE: multiRunE(
			app.loadConfig,
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			m, err := envcrypto.New()
			if err != nil {
				return err
			}
			env, err := godotenv.Marshal(m)
			if err != nil {
				return err
			}
			app.writeString(env)
			return nil
		},
	}
	return cmd
}

func (app *App) encryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt a value",
		Args:  cobra.ExactArgs(1),
		PreRunE: multiRunE(
			app.loadConfig,
			app.loadFiles,
			app.loadSOPS,
			app.openBox,
		),
		RunE: func(_ *cobra.Command, args []string) error {
			encryptedValue, err := app.box.Encrypt(args[0])
			if err != nil {
				return err
			}
			app.writeString(encryptedValue)
			return nil
		},
	}
	cmd.Flags().StringSliceP("files", "f", nil, "files to read environment variables from")
	cmd.Flags().String("sops", "", "decrypt environment file with sops")
	return cmd
}

func (app *App) decryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt a value",
		PreRunE: multiRunE(
			app.loadConfig,
			app.loadFiles,
			app.loadSOPS,
			app.openBox,
		),
		RunE: func(_ *cobra.Command, args []string) error {
			decryptedValue, err := app.box.Decrypt(args[0])
			if err != nil {
				return err
			}
			app.writeString(decryptedValue)
			return nil
		},
	}
	cmd.Flags().StringSliceP("files", "f", nil, "files to read environment variables from")
	cmd.Flags().String("sops", "", "decrypt environment file with sops")
	return cmd
}

func (app *App) getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a value and decrypt if needed",
		PreRunE: multiRunE(
			app.loadConfig,
			app.loadFiles,
			app.loadSOPS,
			app.openBox,
		),
		RunE: func(_ *cobra.Command, args []string) error {
			value, err := app.box.Get(args[0])
			if err != nil {
				return err
			}
			app.writeString(value)
			return nil
		},
	}
	cmd.Flags().StringSliceP("files", "f", nil, "files to read environment variables from")
	cmd.Flags().String("sops", "", "decrypt environment file with sops")
	return cmd
}
