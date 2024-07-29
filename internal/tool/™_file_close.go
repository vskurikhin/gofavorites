/*
 * This file was last modified at 2024-07-29 14:04 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * ™_file_close.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"fmt"
	"log/slog"
	"os"
)

func FileClose(file *os.File) {
	if err := file.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error file close %v", err)
	}
}

func FileCloseAndLog(file *os.File) {
	if err := file.Close(); err != nil {
		slog.Error(MSG+" FileClose", "err", err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
