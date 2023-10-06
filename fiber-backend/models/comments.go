package models

import (
	"time"

	"gorm.io/gorm/clause"
)

// Comments. Containt comment's informtation
// It is the Items foreign key, related on ItemId column.
type Comments struct {
	ID           uint        `gorm:"primaryKey" json:"comment_id,omitempty"`
	PublishingId uint        `gorm:"index" json:"publishing_id,omitempty"`
	OwnerId      uint        `gorm:"index" json:"owner_id,omitempty"`
	Rating       int         `gorm:"index" json:"rating,omitempty" validate:"omitempty,min=0,max=5"`
	Description  string      `gorm:"size:1000" json:"description,omitempty" validate:"omitempty,max=1000"`
	CreatedAt    time.Time   `json:"created_at,omitempty"`
	Active       bool        `json:"active,omitempty"`
	Name         string      `gorm:"size:20" json:"name,omitempty" validate:"required,max=20"`
	Email        string      `gorm:"size:50" json:"email,omitempty" validate:"required,max=50,email"`
	Pub          Publishings `gorm:"foreignKey:PublishingId" json:"pub,omitempty"`
	Owner        Users       `gorm:"foreignKey:OwnerId" json:"owner,omitempty"`
}

//********************************************
//**************Client section***************
//********************************************

// Create comments by publishing
func (c *Comments) CreateComent() error {
	c.Active = true
	if err := DB.Create(c).Error; err != nil {
		go loggerError.Printf("Operation CreateComment => %s", err)
		return err
	}
	return nil
}

//********************************************
//**************End section***************
//********************************************

//********************************************
//**************Sellers section***************
//********************************************

// Get comments for a publishing
func (c Comments) GetComments(pub, offset, size int) ([]Comments, error) {
	var cs []Comments
	if err := DB.Offset(offset).Limit(size).Find(&cs, "publishing_id = ?", pub).Error; err != nil {
		go loggerError.Printf("Operation GetComments => %s", err)
		return nil, err
	}
	return cs, nil
}

// Get top 4 seller by top average rating
func (c Comments) GetTopSellerByRating() ([]*Comments, error) {
	var cs []*Comments
	if err := DB.Preload("Owner").Preload(clause.Associations).
		Select("owner_id, AVG(rating) as total").Group("owner_id").
		Having("total > ?", 2).Limit(4).Find(&cs).Error; err != nil {
		go loggerError.Printf("Operation GetTopSellerByRating => %s", err)
		return nil, err
	}
	for _, obj := range cs {
		obj.Owner.Password = ""
		obj.Owner.Email = ""
		obj.Owner.Address = ""
		obj.Owner.Tel = ""
		obj.Owner.Role = ""
		obj.Owner.ID = 0
	}
	return cs, nil
}

//********************************************
//**************End section***************
//********************************************
