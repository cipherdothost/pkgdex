// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package app

import (
	"fmt"
	"os"

	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/errors"
	"go.cipher.host/pkgdex/internal/key"
)

const (
	// ErrGenerateKey is the error message returned when the key cannot be
	// written to the output.
	ErrGenerateKey cmdkit.Error = "failed to generate authentication key"
)

// NewGenerateKeyCommand returns the "generate-key" command, which is
// responsible for generating a valid API key.
func NewGenerateKeyCommand() *cmdkit.Command {
	cmd := cmdkit.NewCommand("generate-key", "generate a valid API key", "")

	cmd.RunE = GenerateKeyAction

	return cmd
}

// GenerateKeyAction is the action for the "generate-key" command.
func GenerateKeyAction(_ []string) error {
	secret := key.Generate()

	if _, err := fmt.Fprintf(os.Stdout, "%s\n", secret); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateKey, err),
			errors.ErrCodeOutput,
		)
	}

	return nil
}
