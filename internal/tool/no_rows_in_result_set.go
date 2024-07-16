/*
 * This file was last modified at 2024-07-16 23:25 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * no_rows_in_result_set.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

func NoRowsInResultSet(err error) bool {
	return err.Error() == "no rows in result set"
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
