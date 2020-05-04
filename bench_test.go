// The tests in this file are copied from the compress/flate package.

// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brotli

import (
	"bytes"
	"io"
	"io/ioutil"
	"runtime"
	"testing"
)

var suites = []struct{ name, file string }{
	// Digits is the digits of the irrational number e. Its decimal representation
	// does not repeat, but there are only 10 possible digits, so it should be
	// reasonably compressible.
	{"Digits", "testdata/e.txt"},
	// Newton is Isaac Newtons's educational text on Opticks.
	{"Newton", "testdata/Isaac.Newton-Opticks.txt"},
}

func BenchmarkDecode(b *testing.B) {
	doBench(b, func(b *testing.B, buf0 []byte, level, n int) {
		b.ReportAllocs()
		b.StopTimer()
		b.SetBytes(int64(n))

		compressed := new(bytes.Buffer)
		w := NewWriterLevel(compressed, level)
		for i := 0; i < n; i += len(buf0) {
			if len(buf0) > n-i {
				buf0 = buf0[:n-i]
			}
			io.Copy(w, bytes.NewReader(buf0))
		}
		w.Close()
		buf1 := compressed.Bytes()
		buf0, compressed, w = nil, nil, nil
		runtime.GC()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			io.Copy(ioutil.Discard, NewReader(bytes.NewReader(buf1)))
		}
	})
}

var levelTests = []struct {
	name  string
	level int
}{
	{"Speed", BestSpeed},
	{"Default", DefaultCompression},
	{"Compression", BestCompression},
}

var sizes = []struct {
	name string
	n    int
}{
	{"1e4", 1e4},
	{"1e5", 1e5},
	{"1e6", 1e6},
}

func doBench(b *testing.B, f func(b *testing.B, buf []byte, level, n int)) {
	for _, suite := range suites {
		buf, err := ioutil.ReadFile(suite.file)
		if err != nil {
			b.Fatal(err)
		}
		if len(buf) == 0 {
			b.Fatalf("test file %q has no data", suite.file)
		}
		for _, l := range levelTests {
			for _, s := range sizes {
				b.Run(suite.name+"/"+l.name+"/"+s.name, func(b *testing.B) {
					f(b, buf, l.level, s.n)
				})
			}
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	doBench(b, func(b *testing.B, buf0 []byte, level, n int) {
		b.ReportAllocs()
		b.StopTimer()
		b.SetBytes(int64(n))

		buf1 := make([]byte, n)
		for i := 0; i < n; i += len(buf0) {
			if len(buf0) > n-i {
				buf0 = buf0[:n-i]
			}
			copy(buf1[i:], buf0)
		}
		buf0 = nil
		w := NewWriterLevel(ioutil.Discard, level)
		runtime.GC()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			w.Reset(ioutil.Discard)
			w.Write(buf1)
			w.Close()
		}
	})
}