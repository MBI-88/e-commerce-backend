package models

import (
	"fiber-backend/helpers"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Users. Contain user's information, it's the main table
// It has foreign key with Items table
type Users struct {
	ID          uint          `gorm:"primaryKey" json:"user_id,omitempty"`
	Username    string        `gorm:"unique;size:20" json:"username,omitempty" validate:"omitempty,min=5,max=20"`
	Password    string        `gorm:"size:100" json:"password,omitempty" validate:"omitempty,min=8,max=12"`
	Email       string        `gorm:"unique;size:50;index:,sort:desc" json:"email,omitempty" validate:"omitempty,email,min=8,max=50"`
	Tel         string        `gorm:"unique;size:15" json:"tel,omitempty" validate:"omitempty,e164"`
	Address     string        `gorm:"size:200" json:"address,omitempty" validate:"omitempty,max=200"`
	Image       string        `gorm:"size:100" json:"image,omitempty" validate:"omitempty,filepath"`
	Role        string        `gorm:"size:10;default:client" json:"role,omitempty" validate:"omitempty,max=10"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	Active      bool          `json:"active,omitempty"`
	Publishings []Publishings `gorm:"foreignKey:UserId" json:"publishings,omitempty"`
	Trolley     []Trolley     `gorm:"foreignKey:UserId" json:"trolley,omitempty"`
}

//********************************************
//**************Sellers / Client section***************
//********************************************

// Insert an user into the database
func (u *Users) CreateUser() error {
	password, err := GenerateHasKey(u.Password)
	if err != nil {
		return err
	}
	u.Password = password
	u.Active = false
	if err := DB.Omit("Publishings").Create(u).Error; err != nil {
		go loggerError.Printf("Operation CreateUser => %s", err)
		return err
	}
	return nil
}

// Update an user in the database
func (u *Users) UpdateUser(ctx *fiber.Ctx) error {
	var imgName string
	image, err := ctx.FormFile("image")
	u.Username = ctx.FormValue("username", "")
	u.Tel = ctx.FormValue("tel", "")
	u.Address = ctx.FormValue("address", "")
	if id, err := strconv.Atoi(ctx.FormValue("user_id", "")); err == nil {
		u.ID = uint(id)
	}
	if err == nil {
		imgID := uuid.New()
		ext := strings.Split(image.Filename, ".")[1]
		imgName = fmt.Sprintf("%s.%s", imgID, ext)
		timeNow := time.Now()
		dir := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
		localpath := fmt.Sprintf("./static/images/users/%s", dir)
		go helpers.Serializer(ctx, localpath, imgName, image)
		u.Image = fmt.Sprintf("/static/images/users/%s/%s", dir, imgName)
	}
	if errors := helpers.ValidateStruct(u); errors != nil {
		return fmt.Errorf("%v", errors)
	}
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation UpdateUser => %s", err)
		return err
	}
	return nil
}

// Delete an user in the database (first delete publishings next delete user, done by front-end order)
func (u *Users) DeleteUser() error {
	if err := DB.Delete(&Publishings{}, "user_id = ?", u.ID).Error; err != nil {
		go loggerError.Printf("Operation  DeleteUser => %s", err)
		return err
	}
	if err := DB.Delete(u).Error; err != nil {
		go loggerError.Printf("Operation  DeleteUser => %s", err)
		return err
	}
	return nil
}

// Delete user related to an invalid email. Must be used only in register's section.
func (u *Users) DeleteInvalidUser(pk uint) error {
	if err := DB.Delete(u, "id = ?", pk).Error; err != nil {
		go loggerError.Printf("Operation  DeleteInvalidUser => %s", err)
		return err
	}
	return nil
}

// Log an user in the system
func (u *Users) LogUser() error {
	password := u.Password
	if err := DB.Select("id", "password", "email", "role", "image", "username").
		First(u, "email = ? AND active = ?", u.Email, true).Error; err != nil {
		go loggerError.Printf("Operation  loggerErrorUser => %s", err)
		return err
	}
	err := CheckHashPassword(u.Password, password)
	return err
}

// Get user's profile
func (u *Users) GetUser(pk int) error {
	if err := DB.Omit("password", "active").Where("active = ? AND id = ?", true, pk).First(u).Error; err != nil {
		go loggerError.Printf("Operation  GetUser => %s", err)
		return err
	}
	return nil
}

// Change user's password
func (u *Users) ChangePassword() error {
	hashpassword, err := GenerateHasKey(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashpassword
	u.Active = false
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation  ChangePassword => %s", err)
		return err
	}
	return nil
}

// Change user's email
func (u *Users) ChangeEmail() error {
	u.Active = false
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation ChangeEmail => %s", err)
		return err
	}
	return nil
}

// Activate user
func (u *Users) ActiveUser() error {
	u.Active = true
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation ActiveUser => %s", err)
		return err
	}
	return nil
}

// Get user by email
func (u *Users) GetUserbyEmail() error {
	if err := DB.Select("id", "email", "role").Where("active = ? AND email = ?", true, u.Email).
		First(u).Error; err != nil {
		go loggerError.Printf("Operation  GetUserbyEmail => %s", err)
		return err
	}
	return nil
}

// Change role
func (u *Users) ChangeRole() error {
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation ChangeRole => %s", err)
		return err
	}
	return nil
}

//********************************************
//**************End section***************
//********************************************

//********************************************
//**************Admin section***************
//********************************************

// Return all user from database
func (u Users) GetUsers(offset, size int) ([]Users, error) {
	var us []Users
	if err := DB.Offset(offset).Limit(size).Find(&us).Error; err != nil {
		go loggerError.Printf("Operation GetUsers => %s", err)
		return nil, err
	}
	return us, nil
}

// Create all fields of the user table
func (u *Users) MakeUser(ctx *fiber.Ctx) error {
	var imgName string
	u.Username = ctx.FormValue("username", "")
	u.Tel = ctx.FormValue("tel", "")
	u.Address = ctx.FormValue("address", "")
	u.Role = ctx.FormValue("role", "")
	u.Email = ctx.FormValue("email", "")
	u.Password = ctx.FormValue("password", "")
	if active := ctx.FormValue("active", "false"); active != "false" {
		u.Active = true
	} else {
		u.Active = false
	}
	image, err := ctx.FormFile("image")
	if err == nil {
		imgID := uuid.New()
		ext := strings.Split(image.Filename, ".")[1]
		imgName = fmt.Sprintf("%s.%s", imgID, ext)
		timeNow := time.Now()
		dir := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
		localpath := fmt.Sprintf("./static/images/users/%s", dir)
		go helpers.Serializer(ctx, localpath, imgName, image)
		u.Image = fmt.Sprintf("/static/images/users/%s/%s", dir, imgName)
	}
	if errors := helpers.ValidateStruct(u); errors != nil {
		return fmt.Errorf("%v", errors)
	}
	if u.Password != "" {
		u.Password, err = GenerateHasKey(u.Password)
		if err != nil {
			return err
		}
	}
	if err := DB.Create(u).Error; err != nil {
		go loggerError.Printf("Operation MakeUser => %s", err)
		return err
	}
	return nil
}

// Edit all fields of users
func (u *Users) EditUsers(ctx *fiber.Ctx) error {
	var imgName string
	u.Username = ctx.FormValue("username", "")
	u.Tel = ctx.FormValue("tel", "")
	u.Address = ctx.FormValue("address", "")
	u.Role = ctx.FormValue("role", "")
	u.Email = ctx.FormValue("email", "")
	password := ctx.FormValue("password", "")
	if id, err := strconv.Atoi(ctx.FormValue("user_id", "")); err == nil {
		u.ID = uint(id)
	}
	if active := ctx.FormValue("active", "false"); active != "false" {
		u.Active = true
	} else {
		u.Active = false
	}
	image, err := ctx.FormFile("image")
	if err == nil {
		imgID := uuid.New()
		ext := strings.Split(image.Filename, ".")[1]
		imgName = fmt.Sprintf("%s.%s", imgID, ext)
		timeNow := time.Now()
		dir := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
		localpath := fmt.Sprintf("./static/images/users/%s", dir)
		go helpers.Serializer(ctx, localpath, imgName, image)
		u.Image = fmt.Sprintf("/static/images/users/%s/%s", dir, imgName)
	}
	if errors := helpers.ValidateStruct(u); errors != nil {
		return fmt.Errorf("%v", errors)
	}
	if len(password) <= 12 {
		u.Password, err = GenerateHasKey(password)
		if err != nil {
			return err
		}
	}
	if err := DB.Updates(u).Error; err != nil {
		go loggerError.Printf("Operation EditUsers => %s", err)
		return err
	}
	return nil
}

// Search fields by filters
func (u Users) SearchUsers(username, email, role string, offset, size int) ([]Users, error) {
	var users []Users
	query := "id "
	if username != "" {
		query += fmt.Sprintf("And username LIKE %q ", "%"+username+"%")
	}
	if email != "" {
		query += fmt.Sprintf("And email LIKE %q ", "%"+email+"%")
	}
	if role != "" {
		query += fmt.Sprintf("And role LIKE %q ", "%"+role+"%")
	}

	if err := DB.Where(query).Offset(offset).Limit(size).Find(&users).Error; err != nil {
		go loggerError.Printf("Operation SearchUsers => %s", err)
		return nil, err
	}
	return users, nil
}

// Search users by role
func (u Users) GetUsersByRole() (interface{}, error) {
	var result []struct{
		Role string 
		Count int
	}
	if err := DB.Model(u).Select("role, COUNT(*) as count").Group("role").Scan(&result).Error; err != nil {
		go loggerError.Printf("Operation GetUsersByRole => %s", err)
		return nil, err
	}
	return result, nil

}

//********************************************
//**************End section***************
//********************************************
