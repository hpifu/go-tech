package mysql

import "github.com/jinzhu/gorm"

type Likeview struct {
	ID   int `gorm:"type:bigint(20);primary_key" json:"id"`
	View int `gorm:"type:bigint(20);no null;default 0" json:"view"`
	Like int `gorm:"type:bigint(20);no null;default 0" json:"like"`
}

func (m *Mysql) SelectLikeviewByID(id int) (*Likeview, error) {
	likeview := &Likeview{}
	if err := m.db.Where("id=?", id).First(likeview).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &Likeview{
				ID:   id,
				Like: 0,
				View: 0,
			}, nil
		}

		return nil, err
	}

	return likeview, nil
}

func (m *Mysql) SelectLikeviewsByArticles(articles []*Article) (map[int]*Likeview, error) {
	var ids []int
	for _, article := range articles {
		ids = append(ids, article.ID)
	}

	var likeviews []*Likeview

	if err := m.db.Where("id IN (?)", ids).Find(&likeviews).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	res := map[int]*Likeview{}
	for _, lv := range likeviews {
		res[lv.ID] = lv
	}

	return res, nil
}

func (m *Mysql) View(id int) error {
	var count int
	if err := m.db.Model(&Likeview{}).Where(&Likeview{ID: id}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := m.db.Create(&Likeview{ID: id}).Error; err != nil {
			return err
		}
	}

	if err := m.db.Exec("UPDATE likeviews SET `view`=`view`+1 WHERE id=?", id).Error; err != nil {
		return err
	}

	return nil
}

func (m *Mysql) Like(id int) error {
	var count int
	if err := m.db.Model(&Likeview{}).Where(&Likeview{ID: id}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := m.db.Create(&Likeview{ID: id}).Error; err != nil {
			return err
		}
	}

	if err := m.db.Exec("UPDATE likeviews SET `like`=`like`+1 WHERE id=?", id).Error; err != nil {
		return err
	}

	return nil
}
