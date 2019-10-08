package mysql

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	ID      int       `gorm:"type:bigint(20) auto_increment;primary_key" json:"id"`
	Author  int       `gorm:"type:bigint(20);index:author_idx;default:0" json:"author,omitempty"`
	Title   string    `gorm:"type:varchar(128);index:title_idx" json:"title,omitempty"`
	Content string    `gorm:"type:longtext COLLATE utf8mb4_unicode_520_ci;not null" json:"content,omitempty"`
	CTime   time.Time `gorm:"type:timestamp;column:ctime;default:'1970-01-01 00:00:01'" json:"ctime,omitempty"`
	UTime   time.Time `gorm:"type:timestamp;column:utime;default:CURRENT_TIMESTAMP" json:"utime,omitempty"`
}

func (m *Mysql) SelectArticles(offset int, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Select("id, title, author").Order("id").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return articles, nil
}

func (m *Mysql) SelectArticleByID(id int) (*Article, error) {
	article := &Article{}
	if err := m.db.Where("id=?", id).First(article).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return article, nil
}

func (m *Mysql) InsertArticle(article *Article) error {
	return m.db.Create(article).Error
}
