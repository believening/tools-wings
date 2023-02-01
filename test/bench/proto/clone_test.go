package proto

import (
	"testing"

	gogoproto "github.com/gogo/protobuf/proto"
	legacyproto "github.com/golang/protobuf/proto"
	googleproto "google.golang.org/protobuf/proto"
)

var googlepb = ProtoBufGoogle{
	Name:     "foo",
	BirthDay: 56,
	Phone:    "bar",
	Siblings: 960,
	Spouse:   true,
	Money:    5000.0,
}

var gogopb = ProtoBufGoGo{
	Name:     "foo",
	BirthDay: 56,
	Phone:    "bar",
	Siblings: 960,
	Spouse:   true,
	Money:    5000.0,
}

func BenchmarkGoogleProtoClone(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		googleproto.Clone(&googlepb)
	}
}

func BenchmarkLegacyProtoClone(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		legacyproto.Clone(&googlepb)
	}
}

func BenchmarkGoGoProtoClone(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gogoproto.Clone(&gogopb)
	}
}

// type ProtoBufGo struct {
// 	Name     string
// 	BirthDay int64
// 	Phone    string
// 	Siblings int32
// 	Spouse   bool
// 	Money    float64
// }

// func clone(in *ProtoBufGo) *ProtoBufGo {
// 	if in == nil {
// 		return nil
// 	}
// 	out := new(ProtoBufGo)
// 	*out = *in
// 	return out
// }

// func BenchmarkClone(b *testing.B) {
// 	var gostruct = ProtoBufGo{
// 		Name:     "foo",
// 		BirthDay: 56,
// 		Phone:    "bar",
// 		Siblings: 960,
// 		Spouse:   true,
// 		Money:    5000.0,
// 	}
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		clone(&gostruct)
// 	}
// }


/*
goos: linux
goarch: amd64
pkg: github.com/believening/tools-wings/test/bench/proto
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkGoogleProtoClone        7687380               167.2 ns/op            96 B/op          1 allocs/op
BenchmarkGoogleProtoClone-2      6721846               167.5 ns/op            96 B/op          1 allocs/op
BenchmarkGoogleProtoClone-4      6596197               168.1 ns/op            96 B/op          1 allocs/op
BenchmarkGoogleProtoClone-8      6788445               168.0 ns/op            96 B/op          1 allocs/op
BenchmarkLegacyProtoClone        5497088               207.4 ns/op            96 B/op          1 allocs/op
BenchmarkLegacyProtoClone-2      5988266               190.6 ns/op            96 B/op          1 allocs/op
BenchmarkLegacyProtoClone-4      6122752               194.3 ns/op            96 B/op          1 allocs/op
BenchmarkLegacyProtoClone-8      6026874               198.9 ns/op            96 B/op          1 allocs/op
BenchmarkGoGoProtoClone          4672897               250.0 ns/op            96 B/op          1 allocs/op
BenchmarkGoGoProtoClone-2        4909786               237.9 ns/op            96 B/op          1 allocs/op
BenchmarkGoGoProtoClone-4        4926844               236.6 ns/op            96 B/op          1 allocs/op
BenchmarkGoGoProtoClone-8        4888125               232.3 ns/op            96 B/op          1 allocs/op
PASS
ok      github.com/believening/tools-wings/test/bench/proto     16.523s
*/
