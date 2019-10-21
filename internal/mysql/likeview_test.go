package mysql

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMysql_LikeView(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/article?charset=utf8&parseTime=True&loc=Local")
	Convey("test likeview", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		So(m.db.Delete(&Likeview{ID: 1}).Error, ShouldBeNil)
		{
			likeview := &Likeview{ID: 1}
			So(m.Like(1), ShouldBeNil)
			So(m.db.Where(likeview).First(likeview).Error, ShouldBeNil)
			So(likeview.View, ShouldEqual, 0)
			So(likeview.Like, ShouldEqual, 1)
		}
		{
			likeview := &Likeview{ID: 1}
			So(m.View(1), ShouldBeNil)
			So(m.db.Where(likeview).First(likeview).Error, ShouldBeNil)
			So(likeview.View, ShouldEqual, 1)
			So(likeview.Like, ShouldEqual, 1)
		}

		So(m.db.Delete(&Likeview{ID: 1}).Error, ShouldBeNil)
	})
}
