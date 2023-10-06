package models

import "time"

type Trolley struct {
	ID        uint      `gorm:"primaryKey" json:"trolley_id,omitempty"`
	UserId    uint      `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Item      uint      `json:"item_id,omitempty"`
	Count     int       `json:"count,omitempty"`
	State     string    `gorm:"index"  json:"state,omitempty"`
	Payment   float64   `json:"payment,omitempty"`
}

//********************************************
//**************Client section***************
//********************************************

// Add item to trolley
func (t *Trolley) AddToTrolley() error {
	if err := DB.Create(t).Error; err != nil {
		go loggerError.Printf("Operation AddToTrolley => %s", err)
		return err
	}
	return nil
}

// Update item from trolley
func (t *Trolley) UpdateTrolley() error {
	if err := DB.Updates(t).Error; err != nil {
		go loggerError.Printf("Operation UpdateTrolley => %s", err)
		return err
	}
	return nil
}

// Delete item from trolley
func (t *Trolley) DeleteFromTrolley() error {
	if err := DB.Delete(t).Error; err != nil {
		go loggerError.Printf("Operation DeleteFromTrolley => %s", err)
		return err
	}
	return nil
}

// Get all items from trolley
func (t Trolley) GetAllTrolley(pk int) ([]Trolley, error) {
	var ts []Trolley
	if err := DB.Find(&ts, "user_id = ? and state <> 'bought'", pk).Error; err != nil {
		go loggerError.Printf("Operation GetAllTrolley => %s", err)
		return nil, err
	}
	return ts, nil
}

// Delete all trolley
func (t *Trolley) DeleteAll(pk int) error {
	if err := DB.Delete(t, "user_id = ?", pk).Error; err != nil {
		go loggerError.Printf("Operation DeleteAll => %s", err)
		return err
	}
	return nil
}

//********************************************
//**************End section***************
//********************************************
