// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package app

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/errors"
)

const (
	// ErrStopServer is the error message returned when the server cannot be
	// stopped.
	ErrStopServer cmdkit.Error = "failed to stop server"

	// ErrReadPIDFile is the error message returned when the PID file cannot be
	// read.
	ErrReadPIDFile cmdkit.Error = "cannot read PID file"

	// ErrParsePIDFile is the error message returned when the PID file cannot be
	// parsed.
	ErrParsePIDFile cmdkit.Error = "cannot parse PID file"

	// ErrRemovePIDFile is the error message returned when the PID file cannot be
	// removed.
	ErrRemovePIDFile cmdkit.Error = "cannot remove PID file"

	// ErrServerProcess is the error message returned when the server process
	// cannot be found.
	ErrServerProcess cmdkit.Error = "server process not found"

	// ErrStopSignal is the error message returned when the stop signal cannot be
	// sent to the server.
	ErrStopSignal cmdkit.Error = "cannot send stop signal to server"
)

// stopCommand represents the "stop" command.
type stopCommand struct {
	config string
}

// NewStopCommand returns the "stop" command, which is responsible for stopping
// the server.
func NewStopCommand() *cmdkit.Command {
	var (
		stop = &stopCommand{}
		cmd  = cmdkit.NewCommand("stop", "stop the server", "[flags]")
	)

	cmd.Flags.StringVar(&stop.config, "config", config.DefaultConfigLocation, "path to the configuration file")
	cmd.AddShorthand("config", "c")

	cmd.RunE = stop.Action

	return cmd
}

// Action is the action for the "stop" command.
func (c *stopCommand) Action(_ []string) error {
	cfg, err := config.Load(c.config)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStopServer, err),
			errors.ErrCodeConfiguration,
		)
	}

	pidBytes, err := os.ReadFile(cfg.Server.PID)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStopServer, ErrReadPIDFile, err),
			errors.ErrCodeFilesystem,
		)
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStopServer, ErrParsePIDFile, err),
			errors.ErrCodeFilesystem,
		)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStopServer, ErrServerProcess, err),
			errors.ErrCodeServer,
		)
	}

	if err = process.Signal(syscall.SIGTERM); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStopServer, ErrStopSignal, err),
			errors.ErrCodeServer,
		)
	}

	if err = os.Remove(cfg.Server.PID); err != nil && !os.IsNotExist(err) {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStopServer, ErrRemovePIDFile, err),
			errors.ErrCodeFilesystem,
		)
	}

	return nil
}
