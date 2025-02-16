// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.cipher.host/cmdkit"
	"go.cipher.host/pkgdex/internal/config"
	"go.cipher.host/pkgdex/internal/database"
	serrors "go.cipher.host/pkgdex/internal/errors"
	"go.cipher.host/pkgdex/internal/server"
)

const (
	// ErrStartServer is the error message returned when the server cannot be
	// started.
	ErrStartServer cmdkit.Error = "failed to start server"

	// ErrServerRunning is the error message returned when the server is already
	// running.
	ErrServerRunning cmdkit.Error = "server already running; PID file exists"

	// ErrCreatePIDFile is the error message returned when the PID file cannot be
	// created.
	ErrCreatePIDFile cmdkit.Error = "cannot create PID file"

	// ErrWritePIDFile is the error message returned when the PID file cannot be
	// written to.
	ErrWritePIDFile cmdkit.Error = "cannot write PID file"
)

// startCommand represents the "start" command.
type startCommand struct {
	config string
}

// NewStartCommand returns the "start" command, which is responsible for
// starting the server.
func NewStartCommand() *cmdkit.Command {
	var (
		start = &startCommand{}
		cmd   = cmdkit.NewCommand("start", "start the server", "[flags]")
	)

	cmd.Flags.StringVar(&start.config, "config", config.DefaultConfigLocation, "path to the configuration file")
	cmd.AddShorthand("config", "c")

	cmd.RunE = start.Action

	return cmd
}

// Action is the action for the "start" command.
func (c *startCommand) Action(_ []string) error {
	cfg, err := config.Load(c.config)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeConfiguration,
		)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	pidFile, err := os.OpenFile(cfg.Server.PID, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return cmdkit.NewExitError(
				fmt.Errorf("%w: %w: %w", ErrStartServer, ErrServerRunning, err),
				serrors.ErrCodeServer,
			)
		}

		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStartServer, ErrCreatePIDFile, err),
			serrors.ErrCodeFilesystem,
		)
	}

	defer func() {
		if err = pidFile.Close(); err != nil {
			logger.Error(
				"cannot close PID file",
				slog.Any("error", err),
			)
		}
	}()

	if _, err = pidFile.WriteString(strconv.Itoa(os.Getpid())); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w: %w", ErrStartServer, ErrWritePIDFile, err),
			serrors.ErrCodeFilesystem,
		)
	}

	defer func() {
		if err = os.Remove(cfg.Server.PID); err != nil && !os.IsNotExist(err) {
			logger.Error(
				"cannot remove PID file",
				slog.Any("error", err),
			)
		}
	}()

	if err = cfg.Database.CreateIfNotExist(); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeFilesystem,
		)
	}

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeDatabase,
		)
	}

	if err = db.HealthCheck(); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeDatabase,
		)
	}

	ctx := context.Background()

	srv, err := server.NewInstance(ctx, cfg, db, logger)
	if err != nil {
		if errors.Is(err, server.ErrParseHTMLTemplate) {
			return cmdkit.NewExitError(
				fmt.Errorf("%w: %w", ErrStartServer, err),
				serrors.ErrCodeFilesystem,
			)
		}

		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeServer,
		)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan

		if err = srv.Stop(ctx); err != nil {
			logger.Error(
				"failed to stop server",
				slog.Any("error", err),
			)
		}
	}()

	if err = srv.Start(); err != nil {
		return cmdkit.NewExitError(
			fmt.Errorf("%w: %w", ErrStartServer, err),
			serrors.ErrCodeServer,
		)
	}

	return nil
}
