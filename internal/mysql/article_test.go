package mysql

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMysql(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/article?charset=utf8&parseTime=True&loc=Local")
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
			}
			{
				as, err := m.SelectArticles(10, 10)
				So(err, ShouldBeNil)
				So(len(as), ShouldEqual, 10)
			}
		})

		Convey("select article by id", func() {
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

func TestMysql_UpdateArticle(t *testing.T) {
	m, err := NewMysql("hatlonely:keaiduo1@tcp(test-mysql:3306)/article?charset=utf8&parseTime=True&loc=Local")
	Convey("test article", t, func() {
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)

		So(m.db.Delete(&Article{ID: 12}).Error, ShouldBeNil)
		So(m.db.Create(&Article{
			ID:       12,
			AuthorID: 456,
			Author:   "hatlonely",
			Title:    "标题12",
			Content:  "hello world",
		}).Error, ShouldBeNil)

		Convey("update article", func() {
			err := m.UpdateArticle(&Article{
				ID:      12,
				Title:   "标题124",
				Content: "hello golang",
			})
			So(err, ShouldBeNil)

			article := &Article{}
			So(m.db.Where("id=?", 12).First(article).Error, ShouldBeNil)
			So(article.ID, ShouldEqual, 12)
			So(article.Title, ShouldEqual, "标题124")
			So(article.Content, ShouldEqual, "hello golang")
		})
	})
}
