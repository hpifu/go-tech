package mysql

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMysql(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/hads?charset=utf8&parseTime=True&loc=Local")
	Convey("test article", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		article := &Article{
			Title:    "标题1",
			AuthorID: 666,
			Author:   "hatlonely",
			Content:  "hello world",
		}

		for i := 0; i < 20; i++ {
			So(m.db.Delete(&Article{ID: i + 1}).Error, ShouldBeNil)
			So(m.db.Create(&Article{
				ID:       i + 1,
				AuthorID: article.AuthorID,
				Author:   article.Author,
				Title:    fmt.Sprintf("%s-%v", article.Title, i+1),
				Content:  article.Content,
			}).Error, ShouldBeNil)
		}

		Convey("select articles", func() {
			{
				as, err := m.SelectArticles(0, 10)
				So(err, ShouldBeNil)
				So(len(as), ShouldEqual, 10)
				for i := 0; i < 10; i++ {
					So(as[i].ID, ShouldEqual, i+1)
					So(as[i].Title, ShouldEqual, fmt.Sprintf("%s-%v", article.Title, i+1))
					So(as[i].AuthorID, ShouldEqual, article.AuthorID)
					So(as[i].Author, ShouldEqual, article.Author)
				}
			}
			{
				as, err := m.SelectArticles(10, 20)
				So(err, ShouldBeNil)
				So(len(as), ShouldEqual, 10)
				for i := 0; i < 10; i++ {
					So(as[i].ID, ShouldEqual, i+11)
					So(as[i].Title, ShouldEqual, fmt.Sprintf("%s-%v", article.Title, i+11))
					So(as[i].AuthorID, ShouldEqual, article.AuthorID)
					So(as[i].Author, ShouldEqual, article.Author)
				}
			}
		})

		Convey("select ancient by id", func() {
			for i := 0; i < 20; i++ {
				a, err := m.SelectArticleByID(i + 1)
				So(err, ShouldBeNil)
				So(a.ID, ShouldEqual, i+1)
				So(a.Title, ShouldEqual, fmt.Sprintf("%s-%v", article.Title, i+1))
				So(a.AuthorID, ShouldEqual, article.AuthorID)
				So(a.Content, ShouldEqual, article.Content)
				So(a.Author, ShouldEqual, article.Author)
			}

			a, err := m.SelectArticleByID(21)
			So(err, ShouldBeNil)
			So(a, ShouldBeNil)
		})

	})
}
