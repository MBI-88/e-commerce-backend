package models

import (
	"fmt"
	"math/rand"
	"time"
)

type ConfirmationCode struct {
	ID        uint      `gorm:"primaryKey" json:"code_id,omitempty"`
	Code      int       `gorm:"unique" json:"code,omitempty"`
	UserId    uint      `json:"user_id,omitempty"`
	Users     Users     `gorm:"foreignKey:UserId" json:"users,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

//********************************************
//**************Sellers / Client section***************
//********************************************

// Save  an user's confirmation code
func (c *ConfirmationCode) GenerateCode(user Users) error {
	c.Code = generateCode()
	c.UserId = user.ID
	if err := DB.Omit("Users").Create(c).Error; err != nil {
		go loggerError.Printf("Operation GenerateCode => %s", err)
		return err
	}
	return nil
}

// Return an user's confirmation code given an id
func (c *ConfirmationCode) CheckCode() error {
	defer c.DeleteCode()
	code := c.Code
	if err := DB.Preload("Users").First(c).Error; err != nil {
		go loggerError.Printf("Operation CheckCode => %s", err)
		return err
	}
	if !checkCode(*c, code) {
		return fmt.Errorf("Error in confirmation code. It's over 10 minutes")
	}
	return nil
}

// Update code
func (c *ConfirmationCode) UpdateCode() error {
	c.Code = generateCode()
	if err := DB.Updates(c).Error; err != nil {
		go loggerError.Printf("Operation UpdateCode => %s", err)
		return err
	}
	if err := DB.Preload("Users").First(c).Error; err != nil {
		go loggerError.Printf("Operation UpdateCode => %s", err)
		return err
	}
	return nil
}

// Delete code
func (c *ConfirmationCode) DeleteCode() error {
	if err := DB.Delete(c).Error; err != nil {
		go loggerError.Printf("Operation DeleteCode => %s", err)
		return err
	}
	return nil
}

// Generate confirmation code
func generateCode() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9000) + 1000
}

// Confirm confirmation code
func checkCode(c ConfirmationCode, code int) bool {
	return time.Now().Sub(c.CreatedAt).Minutes() <= 10 && code == c.Code
}
