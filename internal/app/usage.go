// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	jsoniter "github.com/json-iterator/go"
	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/errors"
	"go.cipher.host/x/xunsafe"
)

const (
	// ErrGenerateUsage is the error message returned when the usage example cannot
	// be generated.
	ErrGenerateUsage cmdkit.Error = "failed to generate usage example from file; ensure the file exists and contains valid Go code"

	// ErrMissingGoFile is returned when the Go file to convert is not provided.
	ErrMissingGoFile cmdkit.Error = "missing Go file to convert"

	// ErrInvalidGoFile is returned when the file to convert is not a valid Go file.
	ErrInvalidGoFile cmdkit.Error = "invalid example Go file"
)

// NewGenerateUsageCommand returns the "generate-usage" command, which is
// responsible for converting Go code into usage examples.
func NewGenerateUsageCommand() *cmdkit.Command {
	cmd := cmdkit.NewCommand("generate-usage", "convert a Go file to a usage example", "<file>")

	cmd.RunE = GenerateUsageAction

	return cmd
}

// GenerateUsageAction is the action for the "generate-usage" command.
func GenerateUsageAction(args []string) error {
	if len(args) < 1 {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, ErrMissingGoFile),
			errors.ErrCodeValidation,
		)
	}

	file := args[0]

	if !strings.HasSuffix(file, ".go") {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, ErrInvalidGoFile),
			errors.ErrCodeValidation,
		)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, err),
			errors.ErrCodeFilesystem,
		)
	}

	contentString := xunsafe.BytesToString(content)

	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Analyse(contentString) //nolint:misspell // that's how the method is named
	}

	lexer = chroma.Coalesce(lexer)

	style := styles.Get("github")

	formatter := html.New(
		html.WithLineNumbers(true),
		html.WithLinkableLineNumbers(true, "usage-line-"),
		html.TabWidth(4),
		html.WithClasses(true),
		html.Standalone(false),
	)

	var buf strings.Builder

	buf.Grow(len(content) * 2)

	iterator, err := lexer.Tokenise(nil, contentString)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, err),
			errors.ErrCodeUnknown,
		)
	}

	if err = formatter.Format(&buf, style, iterator); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, err),
			errors.ErrCodeUnknown,
		)
	}

	codeString := buf.String()
	codeString = strings.ReplaceAll(codeString, "user-select:none;", "user-select:none;-webkit-user-select: none;")

	var (
		json    = jsoniter.ConfigCompatibleWithStandardLibrary
		jsonBuf bytes.Buffer
	)

	encoder := json.NewEncoder(&jsonBuf)
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(codeString); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrGenerateUsage, err),
			errors.ErrCodeUnknown,
		)
	}

	jsonBytes := jsonBuf.Bytes()

	if len(jsonBytes) >= 3 {
		if _, err = os.Stdout.Write(jsonBytes); err != nil {
			return cmdkit.NewExitError(
				fmt.Errorf("%w: %w", ErrGenerateUsage, err),
				errors.ErrCodeOutput,
			)
		}
	}

	return nil
}
