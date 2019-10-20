package mysql

import (
	"github.com/jinzhu/gorm"
)

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

func (m *Mysql) SelectTagsByArticle(articleID int) ([]*Tag, error) {
	var tags []*Tag

	if err := m.db.Where("article_id=?", articleID).Find(&tags).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return tags, nil
}

func (m *Mysql) UpdateTagsByArticle(articleID int, newTagStr []string) error {
	oldTags, err := m.SelectTagsByArticle(articleID)
	if err != nil {
		return err
	}
	oldTagStrSet := map[string]struct{}{}
	for _, tag := range oldTags {
		oldTagStrSet[tag.Tag] = struct{}{}
	}
	newTagStrSet := map[string]struct{}{}
	for _, t := range newTagStr {
		newTagStrSet[t] = struct{}{}
	}

	var tagIDToDelete []int
	for _, oldTag := range oldTags {
		if _, ok := newTagStrSet[oldTag.Tag]; !ok {
			tagIDToDelete = append(tagIDToDelete, oldTag.ID)
		}
	}

	if err := m.db.Where("id IN (?)", tagIDToDelete).Delete(&Tag{}).Error; err != nil {
		return err
	}

	var tagToInsert []string
	for _, newT := range newTagStr {
		if _, ok := oldTagStrSet[newT]; !ok {
			tagToInsert = append(tagToInsert, newT)
		}
	}

	for _, t := range tagToInsert {
		if err := m.InsertTag(t, articleID); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mysql) SelectArticlesByTag(tag string, offset, limit int) ([]*Article, error) {
	var articles []*Article

	if err := m.db.Table("tags").
		Select(`articles.id AS id, articles.title AS title, articles.tags AS tags, 
			articles.author_id AS author_id, articles.ctime AS ctime, 
			articles.utime AS utime, articles.brief AS brief`).
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

type TagCountPair struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

func (m *Mysql) CountTag() ([]*TagCountPair, error) {
	var tagCloud []*TagCountPair
	if err := m.db.Raw(`
	SELECT tag, count(*) AS count FROM tags GROUP BY tag ORDER BY count DESC LIMIT 50;
`).Scan(&tagCloud).Error; err != nil {
		return nil, err
	}

	return tagCloud, nil
}
