package models

import (
	"fiber-backend/helpers"
	"fmt"
	"time"
)

// Views. Containt view's information.
// It is the Items foreign key, related on ItemId column.
type Views struct {
	ID           uint        `gorm:"primaryKey" json:"id,omitempty"`
	PublishingId uint        `gorm:"index:,sort:asc" json:"publishing_id,omitempty"`
	Category     string      `gorm:"index" json:"category,omitempty"`
	Subcategory  string      `gorm:"index" json:"subcategory,omitempty"`
	Count        int         `gorm:"column:count" json:"count,omitempty"`
	StartDate    time.Time   `json:"start_date,omitempty"`
	EndDate      time.Time   `json:"end_date,omitempty"`
	CreatedAt    time.Time   `gorm:"column:created_at" json:"created_at,omitempty"`
	Publishings  Publishings `gorm:"foreignKey:PublishingId" json:"publishings,omitempty"`
}

//********************************************
//**************Client section***************
//********************************************

// Get view average using count column. Return publishigs associated
func (v Views) GetViews() ([]Views, error) {
	var vs []Views
	if err := DB.Preload("Publishings").Select("publishing_id, AVG(count) AS total").
		Group("publishing_id").Having("total > ?", 2).
		Limit(9).Find(&vs).Error; err != nil {
		go loggerError.Printf("Operation GetViews => %s", err)
		return nil, err
	}
	return vs, nil
}

// Save views for publishing
func (v *Views) SaveView() error {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
	v.StartDate = startDate
	v.EndDate = endDate
	resp := DB.Where("publishing_id = ? and start_date = ?", v.PublishingId, startDate).FirstOrCreate(v)
	if resp.Error != nil {
		go loggerError.Printf("Operation SaveView => %s", resp.Error)
		return resp.Error
	} else if resp.RowsAffected != 0 {
		return nil
	}
	v.Count += 1
	if err := DB.Updates(v).Error; err != nil {
		go loggerError.Printf("Operation SaveView => %s", err)
		return err
	}
	return nil
}

// Get top 9 publishings by category
func (v Views) GetTopByCategory(category, from, to string) ([]Views, error) {
	var vs []Views
	startDate, endDate, err := helpers.DateSelector(from, to)
	if err != nil {
		return nil, err
	}
	if err := DB.Preload("Publishings").
		Select("publishing_id, AVG(count) AS total").
		Where("category = ? and start_date >= ? and end_date <= ?", category, startDate, endDate).
		Group("publishing_id").Having("total > ?", 2).Limit(9).Find(&vs).Error; err != nil {
		go loggerError.Printf("Operation GetTopByCategory => %s", err)
		return nil, err
	}
	return vs, nil
}

// Get top 9 publishings by sub-category
func (v Views) GetTopBySubCategory(subcategory, from, to string) ([]Views, error) {
	var vs []Views
	startDate, endDate, err := helpers.DateSelector(from, to)
	if err != nil {
		return nil, err
	}
	if err := DB.Preload("Publishings").
		Select("publishing_id, AVG(count) AS total").
		Where("subcategory = ? and start_date >= ? and end_date <= ?", subcategory, startDate, endDate).
		Group("publishing_id").Having("total > ?", 2).Limit(9).Find(&vs).Error; err != nil {
		go loggerError.Printf("Operation GetTopBySubCategory => %s", err)
		return nil, err
	}
	return vs, nil
}

//********************************************
//**************End section***************
//********************************************

//********************************************
//**************Sellers section***************
//********************************************

// Feedback about pusblishings views. Only for sellers
func (v Views) GetViewsByPublishing(from, to string, pk int) ([]Views, error) {
	var arrayViews []Views
	startDate, endDate, err := helpers.DateSelector(from, to)
	if err != nil {
		return nil, err
	}
	if err := DB.Select("count", "start_date", "end_date").
		Where("publishing_id = ? and start_date >= ? and end_date <= ?", pk, startDate, endDate).
		Find(&arrayViews).Error; err != nil {
		go loggerError.Printf("Operation GetViewByPublishing => %s", err)
		return nil, err
	}
	return arrayViews, nil
}

//********************************************
//**************End section*******************
//********************************************

//********************************************
//**************Admin section*****************
//********************************************


// Delete views
func (v *Views) DeleteView() error {
	if err := DB.Delete(&v).Error; err != nil {
		go loggerError.Printf("Operation DeleteView => %s", err)
		return err
	}
	return nil
}

// Get views
func (v Views) GetViewsData(offset, size, pubId int, cat, sub string) ([]Views, error) {
	var vs []Views
	query := "id "
	if cat != "" {
		query += fmt.Sprintf("AND category LIKE %q ", "%"+cat+"%")
	}
	if sub != "" {
		query += fmt.Sprintf("AND subcategory LIKE %q ", "%"+sub+"%")
	}
	if pubId != 0 {
		query += fmt.Sprintf("AND publishing_id = %d ", pubId)
	}
	if err := DB.Where(query).Offset(offset).Limit(size).Find(&vs).Error; err != nil {
		go loggerError.Printf("Operation GetViewsData => %s", err)
		return nil, err
	}
	return vs, nil
}

//********************************************
//**************End section*******************
//********************************************
