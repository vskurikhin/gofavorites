/*
 * This file was last modified at 2024-07-31 00:15 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * json_helpers.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"database/sql"
	"time"

	"github.com/goccy/go-json"
)

type JsonNullBool struct {
	sql.NullBool
}

type JsonNullTime struct {
	sql.NullTime
}

func FromNullBool(b sql.NullBool) JsonNullBool {
	return JsonNullBool{NullBool: struct {
		Bool  bool
		Valid bool
	}{Bool: b.Bool, Valid: b.Valid}}
}

func FromNullTime(b sql.NullTime) JsonNullTime {
	return JsonNullTime{NullTime: struct {
		Time  time.Time
		Valid bool
	}{Time: b.Time, Valid: b.Valid}}
}

func (v JsonNullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullBool) UnmarshalJSON(data []byte) error {

	var x *bool

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Bool = *x
	} else {
		v.Valid = false
	}
	return nil
}

func (v JsonNullBool) ToNullBool() sql.NullBool {
	return sql.NullBool{
		Bool:  v.Bool,
		Valid: v.Valid,
	}
}

func (v JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *JsonNullTime) UnmarshalJSON(data []byte) error {

	var x *time.Time

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}

func (v JsonNullTime) ToNullTime() sql.NullTime {
	return sql.NullTime{
		Time:  v.Time,
		Valid: v.Valid,
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
