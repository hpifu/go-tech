package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRunecut(t *testing.T) {
	Convey("test calculate brief", t, func() {
		So(runecut("你好世界", 2), ShouldEqual, "你好")
		So(runecut("你好 golang", 5), ShouldEqual, "你好gol")
		So(runecut("### 你好 golang", 5), ShouldEqual, "你好gol")
	})
}
