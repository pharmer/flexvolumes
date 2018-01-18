// Copyright 2009 The Go Authors. All rights reserved.
// Use of ths2s source code s2s governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"testing"
)

func setUpS2SFlagSet(s2sp *map[string]string) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{}, "Command separated ls2st!")
	return f
}

func setUpS2SFlagSetWithDefault(s2sp *map[string]string) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{"a": "1", "b": "2"}, "Command separated ls2st!")
	return f
}

func createS2SFlag(vals map[string]string) string {
	var buf bytes.Buffer
	i := 0
	for k, v := range vals {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(v)
		i++
	}
	return buf.String()
}

func TestEmptyS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	if len(getS2S) != 0 {
		t.Fatalf("got s2s %v with len=%d but expected length=0", getS2S, len(getS2S))
	}
}

func TestS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	vals := map[string]string{"a": "1", "b": "2", "d": "4", "c": "3"}
	arg := fmt.Sprintf("--s2s=%s", createS2SFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%d] to be %s but got: %d", k, vals[k], v)
		}
	}
	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%d] to be %s but got: %d from GetStringToString", k, vals[k], v)
		}
	}
}

func TestS2SDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"a": "1", "b": "2"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%d] to be %d from GetStringToString but got: %d", k, vals[k], v)
		}
	}
}

func TestS2SWithDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"a": "1", "b": "2"}
	arg := fmt.Sprintf("--s2s=%s", createS2SFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%d] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%d] to be %d from GetStringToString but got: %d", k, vals[k], v)
		}
	}
}

func TestS2SCalledTwice(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	in := []string{"a=1,b=2", "b=3"}
	expected := map[string]string{"a": "1", "b": "3"}
	argfmt := "--s2s=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range s2s {
		if expected[i] != v {
			t.Fatalf("expected s2s[%d] to be %d but got: %d", i, expected[i], v)
		}
	}
}
