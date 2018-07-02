package xtime

import (
	"testing"
	"time"
)

func Test_Fmt(t *testing.T) {
	tt := time.Unix(1234571490, 0)
	cet, _ := time.LoadLocation("CET")
	tt = tt.In(cet)
	if Fmt("%Y-%m-%d %H:%M:%S", tt) != "2009-02-14 01:31:30" {
		t.Fatalf("%s != %s", Fmt("%Y-%m-%d %H:%M:%S", tt), "2009-02-14 01:31:30")
	}
}

func Test_FmtNow(t *testing.T) {
	if FmtNow("%Y%m%d") != time.Now().Format("20060102") {
		t.Fatalf("FmtNow(\"%%Y%%m%%d\") returns %s", FmtNow("%Y%m%d"))
	}
}

func Test_Fmt2(t *testing.T) {
	if FmtNow("%&") != "%&" {
		t.Fatalf("FmtNow(\"%%&\") returns %s", FmtNow("%&"))
	}
}

func Benchmark_Fmt(b *testing.B) {
	t := time.Unix(1333333333, 0)
	cet, _ := time.LoadLocation("CET")
	t = t.In(cet)
	for i := 0; i < b.N; i++ {
		if Fmt("DATE: %Z %X %x", t) != "DATE: CEST 04:22:13 2012-04-02" {
			b.Logf("%s != %s", Fmt("DATE: %Z %X %x", t), "DATE: CEST 04:22:13 2012-04-02")
			b.Fail()
		}
	}
}
