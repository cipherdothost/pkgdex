// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package main // import "go.cipher.host/pkgdex"

import (
	"os"

	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/app"
	"go.cipher.host/pkgdex/internal/meta"
)

func main() {
	cli := cmdkit.New(meta.Name, meta.Description, meta.Version)

	cli.AddCommand(app.NewStartCommand())
	cli.AddCommand(app.NewStopCommand())
	cli.AddCommand(app.NewGenerateKeyCommand())
	cli.AddCommand(app.NewGenerateUsageCommand())

	cli.Examples = []string{
		meta.Name + " start --config '/etc/pkgdex/config.json'",
		meta.Name + " stop --config '/etc/pkgdex/config.json'",
		meta.Name + " generate-key > '/etc/pkgdex/credential/pkgdex-key'",
		meta.Name + " generate-usage '/path/to/example.go'",
	}

	if err := cli.Run(os.Args[1:]); err != nil {
		cli.Errorf("error: %v\n", err)

		cmdkit.Exit(err)
	}
}
