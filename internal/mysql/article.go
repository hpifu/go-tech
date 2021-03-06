package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	ID       int       `gorm:"type:bigint(20) auto_increment;primary_key" json:"id"`
	AuthorID int       `gorm:"type:bigint(20);unique_index:author_title_idx;default:0" json:"authorID,omitempty"`
	Title    string    `gorm:"type:varchar(128);unique_index:author_title_idx;not null" json:"title,omitempty"`
	Tags     string    `gorm:"type:varchar(256);default:''" json:"tags,omitempty"`
	Brief    string    `gorm:"type:varchar(60);default:''" json:"brief,omitempty"`
	Content  string    `gorm:"type:longtext COLLATE utf8mb4_unicode_520_ci;not null" json:"content,omitempty"`
	CTime    time.Time `gorm:"type:timestamp;column:ctime;not null;default:CURRENT_TIMESTAMP;index:ctime_idx" json:"ctime,omitempty"`
	UTime    time.Time `gorm:"type:timestamp;column:utime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:utime_idx" json:"utime,omitempty"`
}

func (m *Mysql) SelectArticles(offset int, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Select("id, title, tags, author_id, ctime, utime, brief").
		Order("utime DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return articles, nil
}

func (m *Mysql) SelectArticlesByAuthor(authorID int, offset, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Select("id, title, tags, author_id, ctime, utime, brief").
		Where("author_id=?", authorID).
		Order("utime DESC").
		Offset(offset).Limit(limit).
		Find(&articles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return articles, nil
}

func (m *Mysql) SelectArticlesByIDs(ids []int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Select("id, title, tags, author_id, ctime, utime, brief").
		Where("id IN (?)", ids).
		Order("utime DESC").
		Find(&articles).Error; err != nil {
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

func (m *Mysql) DeleteArticleByID(id int) error {
	return m.db.Delete(&Article{ID: id}).Error
}
