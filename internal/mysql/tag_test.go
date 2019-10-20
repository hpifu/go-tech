package mysql

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMysql_Tag(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/article?charset=utf8&parseTime=True&loc=Local")
	Convey("test tag", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		// clean
		So(m.db.Delete(&Article{ID: 123}).Error, ShouldBeNil)
		So(m.db.Delete(&Article{ID: 124}).Error, ShouldBeNil)
		So(m.DeleteTag("tag1", 123), ShouldBeNil)
		So(m.DeleteTag("tag2", 123), ShouldBeNil)

		So(m.db.Create(&Article{
			ID:       123,
			Title:    "标题1",
			AuthorID: 666,
			Author:   "hatlonely",
			Content:  "hello world",
		}).Error, ShouldBeNil)
		So(m.db.Create(&Article{
			ID:       124,
			Title:    "标题124",
			AuthorID: 666,
			Author:   "hatlonely",
			Content:  "hello world",
		}).Error, ShouldBeNil)
		So(m.InsertTag("tag1", 123), ShouldBeNil)
		So(m.InsertTag("tag1", 124), ShouldBeNil)
		So(m.InsertTag("tag2", 124), ShouldBeNil)
		{
			tags, err := m.SelectTagsByArticle(124)
			So(err, ShouldBeNil)
			So(len(tags), ShouldEqual, 2)
		}
		{
			articles, err := m.SelectArticlesByTag("tag1", 0, 20)
			So(err, ShouldBeNil)
			So(len(articles), ShouldEqual, 2)
		}
		{
			articles, err := m.SelectArticlesByTag("tag2", 0, 20)
			So(err, ShouldBeNil)
			So(len(articles), ShouldEqual, 1)
			So(articles[0].ID, ShouldEqual, 124)
			So(articles[0].Title, ShouldEqual, "标题124")
			So(articles[0].AuthorID, ShouldEqual, 666)
		}

		// clean
		So(m.db.Delete(&Article{ID: 123}).Error, ShouldBeNil)
		So(m.db.Delete(&Article{ID: 124}).Error, ShouldBeNil)
		So(m.DeleteTag("tag1", 123), ShouldBeNil)
		So(m.DeleteTag("tag2", 123), ShouldBeNil)
	})
}
