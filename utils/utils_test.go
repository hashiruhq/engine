package utils_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.com/around25/products/matching-engine/utils"
)

func BenchmarkDivideUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.Divide(1278543132023424178, 1244999636900000001, 8, 8, 8) // * 10^8 = 127.900.753.966,108
	}
}

func BenchmarkMultiplyUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.Multiply(1278543132023424178, 99900000, 8, 8, 8)
	}
}

func TestMultiply(t *testing.T) {
	Convey("Two uint64 numbers representing two decimal numbers with 8 digit precision", t, func() {
		Convey("I should be able to multiply them and receive a uint64 with the same precision", func() {
			So(utils.Multiply(1000000000, 300000000, 8, 8, 8), ShouldEqual, 3000000000)
			So(utils.Multiply(0, 3, 8, 8, 8), ShouldEqual, 0)
			So(utils.Multiply(10000000, 652934000000, 8, 9, 8), ShouldEqual, 6529340000)
			So(utils.Multiply(100000000, 652934000000, 8, 1, 1), ShouldEqual, 652934000000)
			So(utils.Multiply(1000000000, 652934000000, 8, 8, 8), ShouldEqual, 6529340000000)
			So(utils.Multiply(10000000000, 652934000000, 8, 8, 8), ShouldEqual, 65293400000000)
			So(utils.Multiply(100000000000, 652934000000, 8, 8, 8), ShouldEqual, 652934000000000)
			So(utils.Multiply(1000000000000, 652934000000, 8, 8, 8), ShouldEqual, 6529340000000000)
			So(utils.Multiply(10000000000000, 652934000000, 8, 8, 8), ShouldEqual, 65293400000000000)
			So(utils.Multiply(100000000000000, 652934000000, 8, 8, 8), ShouldEqual, 652934000000000000)
			So(utils.Multiply(100000000000000, 6529340000000, 8, 8, 8), ShouldEqual, 6529340000000000000)
			So(utils.Multiply(12785400000000, 3036900000000, 8, 8, 8), ShouldEqual, 388279812600000000)
			// check if result exceeds uint boundary
			So(utils.Multiply(12785400000000, 377036900000000, 8, 8, 8), ShouldEqual, 0)
		})
	})
}

func TestDivide(t *testing.T) {
	Convey("Two uint64 numbers representing two decimal numbers with 8 digit precision", t, func() {
		Convey("I should be able to divide them and receive a uint64 with the same precision", func() {
			So(utils.Divide(10, 3, 8, 8, 8), ShouldEqual, 333333333)
			So(utils.Divide(0, 3, 8, 8, 8), ShouldEqual, 0)
			So(utils.Divide(1, 3, 8, 8, 8), ShouldEqual, 33333333)
			So(utils.Divide(15, 5, 8, 8, 8), ShouldEqual, 300000000)
			So(utils.Divide(1278543132023424178, 999636900000000, 8, 8, 8), ShouldEqual, 127900753966)
		})
	})
}

func TestMaxMin(t *testing.T) {
	Convey("Two uint64 numbers representing two decimal numbers with 8 digit precision", t, func() {
		Convey("I should be able to check which one is bigger/smaller", func() {
			So(utils.Max(10, 3), ShouldEqual, 10)
			So(utils.Min(10, 3), ShouldEqual, 3)
			So(utils.Max(3, 10), ShouldEqual, 10)
			So(utils.Min(3, 10), ShouldEqual, 3)
		})
	})
}
