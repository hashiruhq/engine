package conv_test

import (
	"testing"

	"github.com/hashiruhq/engine/conv"

	. "github.com/smartystreets/goconvey/convey"
)

func BenchmarkConvertToUnits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conv.ToUnits("101000101.33232313", 8)
	}
}

func BenchmarkConvertFromUnits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conv.FromUnits(10100010133232313, 8)
	}
}

func TestConvertToUnits(t *testing.T) {
	Convey("Given a string representation of a float number", t, func() {
		Convey("I should be able to convert it into units with a fixed precision", func() {
			So(conv.ToUnits("", 8), ShouldEqual, 0)
			So(conv.ToUnits("0.0", 8), ShouldEqual, 0)
			So(conv.ToUnits("1", 8), ShouldEqual, 100000000)
			So(conv.ToUnits("9340", 8), ShouldEqual, 934000000000)
			So(conv.ToUnits("9996369", 8), ShouldEqual, 999636900000000)
			So(conv.ToUnits("9996369.12", 8), ShouldEqual, 999636912000000)
			So(conv.ToUnits("0.00000001", 8), ShouldEqual, 1)
			So(conv.ToUnits("12785431320.23424178", 8), ShouldEqual, 1278543132023424178)
			So(conv.ToUnits("12785431320.234241781222", 8), ShouldEqual, 1278543132023424178)
		})
	})
}

func TestConvertFromUnits(t *testing.T) {
	Convey("Given a unit representation of a float number with a given precision", t, func() {
		Convey("I should be able to convert it into a string representation of a float", func() {
			So(conv.FromUnits(0, 8), ShouldEqual, "0.00000000")
			So(conv.FromUnits(100000000, 8), ShouldEqual, "1.00000000")
			So(conv.FromUnits(1, 8), ShouldEqual, "0.00000001")
			So(conv.FromUnits(1278543132023424178, 8), ShouldEqual, "12785431320.23424178")
			So(conv.FromUnits(934000000000, 8), ShouldEqual, "9340.00000000")
			So(conv.FromUnits(999636912000000, 8), ShouldEqual, "9996369.12000000")
			// max uint value
			So(conv.FromUnits(18446744073709551615, 8), ShouldEqual, "184467440737.09551615")
		})
	})
}
