package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	ID       int       `gorm:"type:bigint(20) auto_increment;primary_key" json:"id"`
	Author   string    `gorm:"type:varchar(64);index:author_idx;default:''" json:"author,omitempty"`
	AuthorID int       `gorm:"type:bigint(20);unique_index:author_title_idx;default:0" json:"authorID,omitempty"`
	Title    string    `gorm:"type:varchar(128);unique_index:author_title_idx;not null" json:"title,omitempty"`
	Tags     string    `gorm:"type:varchar(256);default:''" json:"tags,omitempty"`
	Content  string    `gorm:"type:longtext COLLATE utf8mb4_unicode_520_ci;not null" json:"content,omitempty"`
	CTime    time.Time `gorm:"type:timestamp;column:ctime;not null;default:CURRENT_TIMESTAMP" json:"ctime,omitempty"`
	UTime    time.Time `gorm:"type:timestamp;column:utime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"utime,omitempty"`
}

func (m *Mysql) SelectArticles(offset int, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Select("id, title, tags, author, author_id, ctime, utime, content").Order("id").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
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

func (m *Mysql) SelectArticleByAuthorAndTitle(authorID int, title string) (*Article, error) {
	article := &Article{}
	if err := m.db.Where(&Article{
		AuthorID: authorID,
		Title:    title,
	}).Find(article).Error; err != nil {
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

func (m *Mysql) UpdateArticle(article *Article) error {
	return m.db.Model(article).Where("id=?", article.ID).Updates(article).Error
}
