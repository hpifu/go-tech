package mysql

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMysql_SelectArticles(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/hads?charset=utf8&parseTime=True&loc=Local")
	Convey("test article", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		article := &Article{
			Title:   "标题1",
			Content: "hello world",
		}

		for i := 0; i < 20; i++ {
			So(m.db.Delete(&Article{ID: i + 1}).Error, ShouldBeNil)
			So(m.db.Create(&Article{
				ID:      i + 1,
				Title:   fmt.Sprintf("%s-%v", article.Title, i+1),
				Content: article.Content,
			}).Error, ShouldBeNil)
		}

	})
}
