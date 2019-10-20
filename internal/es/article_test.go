package es

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test(t *testing.T) {
	es, err := NewES("http://test-elasticsearch:9200", "article", 5000*time.Millisecond)
	Convey("test article", t, func() {
		So(err, ShouldBeNil)
		So(es, ShouldNotBeNil)

		So(es.InsertArticle(&Article{
			ID:      123,
			Title:   "标题",
			Author:  "hatlonely",
			Tags:    "c++,java",
			Brief:   "hello",
			Content: "hello world",
		}), ShouldBeNil)

		So(es.UpdateArticle(&Article{
			ID:      123,
			Title:   "标题",
			Author:  "hatlonely",
			Tags:    "c++,java",
			Brief:   "hello",
			Content: "你好世界",
		}), ShouldBeNil)

		So(es.DeleteArticle(123), ShouldBeNil)
	})
}
