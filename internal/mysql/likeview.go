package mysql

type Likeview struct {
	ID   int `gorm:"type:bigint(20);primary_key" json:"id"`
	View int `gorm:"type:bigint(20);no null;default 0" json:"view"`
	Like int `gorm:"type:bigint(20);no null;default 0" json:"like"`
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
