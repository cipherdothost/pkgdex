// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: EUPL-1.2

// Package errors contains common errors and error codes for different
// categories of errors.
package errors

const (
	_ = iota
	ErrCodeUnknown
	ErrCodeOutput
	ErrCodeUsage
	ErrCodeInterrupted
	ErrCodeConfiguration
	ErrCodeDatabase
	ErrCodeServer
	ErrCodeFilesystem
	ErrCodeNetwork
	ErrCodeAuth
	ErrCodeValidation
)
