/*
 * This file was last modified at 2024-07-31 14:33 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * logger.go
 * $Id$
 */
//!+

// Package alog TODO.
package alog

import (
	"log/slog"
	"sync"
)

var (
	onceLogger = new(sync.Once)
	logger     *slog.Logger
)

func GetLogger() *slog.Logger {
	return logger
}

func NewLogger(handler slog.Handler) {
	onceLogger.Do(func() {
		logger = slog.New(handler)
	})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
