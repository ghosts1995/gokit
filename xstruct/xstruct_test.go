/*
 * Copyright 2012-2019 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * A toolkit for Golang development
 * https://www.likexian.com/
 */

package xstruct

import (
	"reflect"
	"testing"

	"github.com/likexian/gokit/assert"
)

type Techer struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type Student struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Techer  Techer `json:"techer"`
	score   map[string]int
}

var techer = Techer{100, "techer.li", true}
var student = Student{1, "kexian.li", true, techer, map[string]int{}}

func TestVersion(t *testing.T) {
	assert.Contains(t, Version(), ".")
	assert.Contains(t, Author(), "likexian")
	assert.Contains(t, License(), "Apache License")
}

func TestIsStruct(t *testing.T) {
	var i interface{}
	tests := []struct {
		in  interface{}
		out bool
	}{
		{nil, false},
		{"", false},
		{1, false},
		{i, false},
		{student, true},
		{&student, true},
		{student.Techer, true},
		{student.Techer.Name, false},
	}

	for _, v := range tests {
		assert.Equal(t, IsStruct(v.in), v.out)
	}
}

func TestNew(t *testing.T) {
	assert.Panic(t, func() { New(nil) })
	assert.Panic(t, func() { New("") })
	assert.Panic(t, func() { New(map[string]interface{}{}) })

	s := New(student)
	assert.NotNil(t, s)
}

func TestName(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	name := s.Name()
	assert.Equal(t, name, "Student")

	name = Name(student)
	assert.Equal(t, name, "Student")
}

func TestStruct(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	assert.Panic(t, func() { s.Struct("not-exists") })
	assert.Panic(t, func() { s.Struct("Id") })

	ss := s.Struct("Techer")
	assert.NotNil(t, ss)
	assert.Equal(t, ss.Name(), "Techer")

	f, ok := ss.Field("Name")
	assert.True(t, ok)

	n := f.Name()
	assert.Equal(t, n, "Name")

	v := f.Value()
	assert.Equal(t, v, "techer.li")

	k := f.Kind()
	assert.Equal(t, k, reflect.String)

	b := f.IsAnonymous()
	assert.False(t, b)

	ss = Struct(student, "Techer")
	assert.NotNil(t, ss)
	assert.Equal(t, ss.Name(), "Techer")
}

func TestMap(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	v := s.Map()
	assert.Len(t, v, 4)
	assert.Equal(t, v["Name"], "kexian.li")

	v = Map(student)
	assert.Len(t, v, 4)
	assert.Equal(t, v["Name"], "kexian.li")
}

func TestNames(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	n := s.Names()
	assert.Len(t, n, 5)

	n = Names(student)
	assert.Len(t, n, 5)
}

func TestTags(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	m, err := s.Tags("json")
	assert.Nil(t, err)
	assert.Len(t, m, 4)
	assert.Equal(t, m["Name"], "name")

	m, err = Tags(student, "json")
	assert.Nil(t, err)
	assert.Len(t, m, 4)
	assert.Equal(t, m["Name"], "name")
}

func TestValues(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	v := s.Values()
	assert.Len(t, v, 4)

	v = Values(student)
	assert.Len(t, v, 4)
}

func TestFields(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	f := s.Fields()
	assert.Len(t, f, 5)

	f = Fields(student)
	assert.Len(t, f, 5)
}

func TestField(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	_, ok := s.Field("not-exists")
	assert.False(t, ok)

	f, ok := s.Field("Name")
	assert.True(t, ok)

	n := f.Name()
	assert.Equal(t, n, "Name")

	v := f.Value()
	assert.Equal(t, v, "kexian.li")

	k := f.Kind()
	assert.Equal(t, k, reflect.String)

	b := f.IsAnonymous()
	assert.False(t, b)

	f, ok = Field(student, "Name")
	assert.True(t, ok)
	n = f.Name()
	assert.Equal(t, n, "Name")
}

func TestMustField(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	assert.Panic(t, func() { s.MustField("not-exists") })

	f := s.MustField("Name")

	n := f.Name()
	assert.Equal(t, n, "Name")

	v := f.Value()
	assert.Equal(t, v, "kexian.li")

	k := f.Kind()
	assert.Equal(t, k, reflect.String)

	b := f.IsAnonymous()
	assert.False(t, b)

	n = s.MustField("Enabled").Name()
	assert.Equal(t, n, "Enabled")

	v = s.MustField("Enabled").Value()
	assert.Equal(t, v, true)

	f = MustField(student, "Name")
	n = f.Name()
	assert.Equal(t, n, "Name")
}

func TestHasField(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	b := s.HasField("not-exists")
	assert.False(t, b)

	b = s.HasField("Id")
	assert.True(t, b)

	b = s.HasField("Techer")
	assert.True(t, b)
}

func TestFieldTag(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	f, ok := s.Field("Name")
	assert.True(t, ok)

	n := f.Tag("not-exists")
	assert.Equal(t, n, "")

	n = f.Tag("json")
	assert.Equal(t, n, "name")
}

func TestFieldIsExport(t *testing.T) {
	s := New(student)
	assert.NotNil(t, s)

	f, ok := s.Field("Name")
	assert.True(t, ok)
	b := f.IsExport()
	assert.True(t, b)

	f, ok = s.Field("score")
	assert.True(t, ok)
	b = f.IsExport()
	assert.False(t, b)
}

func TestFieldIsZero(t *testing.T) {
	s := New(Student{})
	assert.NotNil(t, s)

	f, ok := s.Field("Name")
	assert.True(t, ok)
	b := f.IsZero()
	assert.True(t, b)

	s = New(student)
	assert.NotNil(t, s)

	f, ok = s.Field("Name")
	assert.True(t, ok)
	b = f.IsZero()
	assert.False(t, b)

	f, ok = s.Field("score")
	assert.True(t, ok)
	assert.Panic(t, func() { f.IsZero() })
}

func TestFieldSet(t *testing.T) {
	s := New(&student)
	assert.NotNil(t, s)

	f, ok := s.Field("score")
	assert.True(t, ok)

	err := f.Set(0)
	assert.Equal(t, err, ErrNotExported)

	f, ok = s.Field("Name")
	assert.True(t, ok)

	err = f.Set("lkx")
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "lkx")

	err = f.Set(0)
	assert.NotNil(t, err)

	err = s.Set("not-exists", 0)
	assert.Equal(t, err, ErrNotField)

	err = s.Set("score", 0)
	assert.Equal(t, err, ErrNotExported)

	err = s.Set("Name", "likexian")
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "likexian")

	s = New(student)
	assert.Nil(t, err)

	f, ok = s.Field("Name")
	assert.True(t, ok)

	err = f.Set("lkx")
	assert.Equal(t, err, errNotSettable)

	err = Set(&student, "Name", "likexian")
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "likexian")
}

func TestFieldZero(t *testing.T) {
	s := New(&student)
	assert.NotNil(t, s)

	f, ok := s.Field("score")
	assert.True(t, ok)

	err := f.Zero()
	assert.Equal(t, err, ErrNotExported)

	f, ok = s.Field("Id")
	assert.True(t, ok)

	err = f.Zero()
	assert.Nil(t, err)
	assert.Equal(t, student.Id, int64(0))
	assert.True(t, f.IsZero())

	f, ok = s.Field("Name")
	assert.True(t, ok)

	err = f.Zero()
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "")
	assert.True(t, f.IsZero())

	err = s.Zero("not-exists")
	assert.Equal(t, err, ErrNotField)

	err = s.Zero("score")
	assert.Equal(t, err, ErrNotExported)

	err = s.Zero("Name")
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "")

	err = Zero(&student, "Name")
	assert.Nil(t, err)
	assert.Equal(t, student.Name, "")
}

func TestFieldIsStruct(t *testing.T) {
	s := New(&student)
	assert.NotNil(t, s)

	b := s.IsStruct("not-exists")
	assert.False(t, b)

	b = s.IsStruct("Id")
	assert.False(t, b)

	b = s.IsStruct("Techer")
	assert.True(t, b)

	b = s.Struct("Techer").IsStruct("Id")
	assert.False(t, b)
}
