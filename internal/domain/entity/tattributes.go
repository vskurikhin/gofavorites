/*
 * This file was last modified at 2024-07-18 19:05 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * tattributes.go
 * $Id$
 */

package entity

import (
	"database/sql"
	"time"
)

type TAttributes struct {
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

func DefaultTAttributes() TAttributes {
	return TAttributes{deleted: sql.NullBool{}, createdAt: time.Time{}, updatedAt: sql.NullTime{}}
}

func MakeTAttributes(deleted sql.NullBool, createdAt time.Time, updatedAt sql.NullTime) TAttributes {
	return TAttributes{deleted: deleted, createdAt: createdAt, updatedAt: updatedAt}
}
