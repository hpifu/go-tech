package mysql

import "github.com/jinzhu/gorm"

type Tag struct {
	ID        int    `gorm:"type:bigint(20) auto_increment;primary_key" json:"id,omitempty"`
	ArticleID int    `gorm:"type:bigint(20);index:article_idx;unique_index:article_tag_idx;not null" json:"articleID,omitempty"`
	Tag       string `gorm:"type:varchar(128);index:tag_idx;unique_index:article_tag_idx;not null" json:"tag,omitempty"`
}

func (m *Mysql) InsertTag(tag string, articleID int) error {
	return m.db.Create(&Tag{
		Tag:       tag,
		ArticleID: articleID,
	}).Error
}

func (m *Mysql) DeleteTag(tag string, articleID int) error {
	return m.db.Delete(&Tag{
		Tag:       tag,
		ArticleID: articleID,
	}).Error
}

func (m *Mysql) SelectArticlesByTag(tag string, offset, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Table("tags").
		Select(`articles.id AS id, articles.title AS title, articles.tags AS tags, 
			articles.author AS author, articles.author_id AS author_id, articles.ctime AS ctime, 
			articles.utime AS utime, articles.content AS content`).
		Joins("LEFT JOIN articles ON articles.id=tags.article_id").
		Where("tags.tag=?", tag).
		Order("articles.utime DESC").
		Offset(offset).Limit(limit).
		Scan(&articles).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return articles, nil
}
