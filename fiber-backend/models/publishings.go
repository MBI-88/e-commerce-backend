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

// Items. Containt item's iformation.
// It is the Users foreign key, related on UserId column.
// It has foreign keys with Comments and Views.
type Publishings struct {
	ID          uint       `gorm:"primaryKey" json:"publishing_id,omitempty"`
	UserId      uint       `json:"user_id,omitempty"`
	Category    string     `gorm:"size:20;index" json:"category,omitempty"`
	Subcategory string     `gorm:"size:20;index" json:"subcategory,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	ProductName string     `gorm:"size:80;index:,sort:desc,priority:1" json:"product_name,omitempty" validate:"omitempty,min=5,max=80"`
	Description string     `gorm:"size:1000" json:"description,omitempty" validate:"omitempty,max=1000"`
	Price       float64    `gorm:"index" json:"price,omitempty"`
	Stock       int        `json:"stock,omitempty"`
	Os          string     `gorm:"size:20;index" json:"os,omitempty"`
	Ram         int        `gorm:"index" json:"ram,omitempty"`
	Store       int        `gorm:"index" json:"store,omitempty"`
	Screen      float64    `gorm:"index" json:"screen,omitempty"`
	Brand       string     `gorm:"size:20;index" json:"brand,omitempty"`
	InOffer     float64    `gorm:"index" json:"in_offer,omitempty"`
	Resolution  string     `gorm:"size:20;index" json:"resolution,omitempty"`
	Gpu         string     `gorm:"size:20;index" json:"gpu,omitempty"`
	Micro       string     `gorm:"size:20;index" json:"micro,omitempty"`
	Long        float64    `gorm:"index" json:"long,omitempty"`
	Width       float64    `gorm:"index" json:"width,omitempty"`
	Height      float64    `gorm:"index" json:"height,omitempty"`
	Weight      float64    `gorm:"index" json:"weight,omitempty"`
	Comments    []Comments `gorm:"foreignKey:PublishingId" json:"comments,omitempty"`
	Views       []Views    `gorm:"foreignKey:PublishingId" json:"views,omitempty"`
	Image1      string     `gorm:"size:100" json:"image_1,omitempty"`
	Image2      string     `gorm:"size:100" json:"image_2,omitempty"`
	Image3      string     `gorm:"size:100" json:"image_3,omitempty"`
	Image4      string     `gorm:"size:100" json:"image_4,omitempty"`
	Image5      string     `gorm:"size:100" json:"image_5,omitempty"`
	Owner       Users      `gorm:"foreignKey:UserId" json:"owner,omitempty"`
}

//********************************************
//**************Client section***************
//********************************************

// Get a publishing
func (p *Publishings) GetPublishing(pk int) error {
	if err := DB.Preload("Owner").First(p, pk).Error; err != nil {
		go loggerError.Printf("Operation  GetPublishing => %s", err)
		return err
	}
	p.Owner.Password = ""
	p.Owner.Role = ""
	p.Owner.ID = 0
	return nil
}

// Get items by categories
func (p Publishings) GetPublishingsByCategories(category string, offset, size int) ([]Publishings, error) {
	var ps []Publishings
	if err := DB.Offset(offset).Limit(size).Where("category = ?", category).Find(&ps).Error; err != nil {
		go loggerError.Printf("Operation  GetPublishingsByCategories => %s", err)
		return nil, err
	}
	return ps, nil
}

// Get publishings by subcategory having category
func (p Publishings) GetPublishingsBySubCategory(offset, size int, category, subcategory string) ([]Publishings, error) {
	var ps []Publishings
	if err := DB.Where("category = ? AND subcategory = ?", category, subcategory).Offset(offset).
		Limit(size).Find(&ps).Error; err != nil {
		go loggerError.Printf("Operation  GetPublishingsBySubCategory => %s", err)
		return nil, err
	}
	return ps, nil
}

// Get count by categories
func (p Publishings) GetCountByCategories() (interface{}, error) {
	var result []struct {
		Category      string
		Subcategories int
	}
	if err := DB.Model(p).Select("category", "COUNT(*) AS subcategories").Group("category").
		Scan(&result).Error; err != nil {
		go loggerError.Printf("Operation GetCountByCategories => %s", err)
		return nil, err
	}
	return result, nil
}

// Get count by subcategories
func (p Publishings) GetCountBySubCategories(param string) (interface{}, error) {
	var result []struct {
		Subcategory string
		Count       int
	}
	if err := DB.Model(p).Select("subcategory", "COUNT(*) as count").Group("subcategory").
		Where("category", param).Scan(&result).Error; err != nil {
		go loggerError.Printf("Operation  GetCountBySubCategories => %s", err)
		return nil, err
	}
	return result, nil
}

// Get items by filters
func (p Publishings) GetPublishingsByFilters(q helpers.ParseQuery, offset, size int) ([]Publishings, error) {
	var ps []Publishings
	query := fmt.Sprintf("category = %q AND subcategory = %q ", q.Category, q.Subcategory)
	if q.Brand != "" {
		query += fmt.Sprintf("AND brand = %q ", q.Brand)
	}
	if q.Gpu != "" {
		query += fmt.Sprintf("AND gpu = %q ", q.Gpu)
	}
	if q.Micro != "" {
		query += fmt.Sprintf("AND micro = %q ", q.Micro)
	}
	if q.Os != "" {
		query += fmt.Sprintf("AND os = %q ", q.Os)
	}
	if q.ProductName != "" {
		query += fmt.Sprintf("AND product_name = %q ", q.ProductName)
	}
	if q.MinPrice != 0.0 && q.MaxPrice != 0.0 {
		query += fmt.Sprintf("AND price BETWEEN %f AND %f ", q.MinPrice, q.MaxPrice)
	}
	if q.MinRam != 0 && q.MaxRam != 0 {
		query += fmt.Sprintf("AND ram BETWEEN %d AND %d ", q.MinRam, q.MaxRam)
	}
	if q.Resolution != "" {
		query += fmt.Sprintf("AND resolution = %q ", q.Resolution)
	}
	if q.MinScreen != 0.0 && q.MaxScreen != 0.0 {
		query += fmt.Sprintf("AND screen BETWEEN %f AND %f ", q.MinScreen, q.MaxScreen)
	}
	if q.MinStore != 0 && q.MaxStore != 0 {
		query += fmt.Sprintf("AND store BETWEEN %d AND %d ", q.MinStore, q.MaxStore)
	}
	if q.MinLong != 0.0 && q.MaxLong != 0.0 {
		query += fmt.Sprintf("AND long BETWEEN %f AND %f ", q.MinLong, q.MaxLong)
	}
	if q.MinWidth != 0.0 && q.MaxWidth != 0.0 {
		query += fmt.Sprintf("AND width BETWEEN %f AND %f ", q.MinWidth, q.MaxWidth)
	}
	if q.MinWeight != 0.0 && q.MaxWeight != 0.0 {
		query += fmt.Sprintf("AND weight BETWEEN %f AND %f ", q.MinWeight, q.MaxWeight)
	}
	if q.MinHeight != 0.0 && q.MaxHeight != 0.0 {
		query += fmt.Sprintf("AND height BETWEEN %f AND %f", q.MinHeight, q.MaxHeight)
	}
	if err := DB.Where(query).Offset(offset).Limit(size).Find(&ps).Error; err != nil {
		go loggerError.Printf("Operation GetPublishingsByFilter => %s", err)
		return nil, err
	}
	return ps, nil
}

// Get all publishins
func (p Publishings) GetAll(q helpers.ParseQuery, offset, size int) ([]Publishings, error) {
	var ps []Publishings
	query := "id "
	if q.Category != "" {
		query += fmt.Sprintf("AND category LIKE %q ", "%"+q.Category+"%")
	}
	if q.Subcategory != "" {
		query += fmt.Sprintf("AND subcategory LIKE %q ", "%"+q.Subcategory+"%")
	}
	if q.Brand != "" {
		query += fmt.Sprintf("AND brand LIKE %q ", "%"+q.Brand+"%")
	}
	if q.ProductName != "" {
		query += fmt.Sprintf("AND product_name LIKE %q ", "%"+q.ProductName+"%")
	}
	if err := DB.Where(query).Offset(offset).Limit(size).Find(&ps).Error; err != nil  {
		go loggerError.Printf("Operation  GetAll => %s", err)
		return nil, err
	}
	return ps, nil
}

//********************************************
//**************End section***************
//********************************************

//********************************************
//**************Sellers section***************
//********************************************

// Get publishing only for sellers
func (p Publishings) GetPublishingBySeller(offset, size, userId int) ([]Publishings, error) {
	var ps []Publishings
	if err := DB.Offset(offset).Limit(size).Find(&ps, "user_id = ?", userId).Error; err != nil {
		go loggerError.Printf("Operation GetPublishingBySeller => %s", err)
		return nil, err
	}
	return ps, nil
}

// Create publishing only for sellers
func (p *Publishings) CreatePublishing(ctx *fiber.Ctx) error {
	singStruct(p, ctx)
	if errors := helpers.ValidateStruct(p); errors != nil {
		return fmt.Errorf("%v", errors)
	}
	if err := DB.Create(p).Error; err != nil {
		go loggerError.Printf("Operation CreatePublishing => %s", err)
		return err
	}
	return nil
}

// Update publishing only for sellers
func (p *Publishings) UpdatePublishing(ctx *fiber.Ctx) error {
	singStruct(p, ctx)
	if errors := helpers.ValidateStruct(p); errors != nil {
		return fmt.Errorf("%v", errors)
	}
	if err := DB.Updates(p).Error; err != nil {
		go loggerError.Printf("Operation UpdatePublishing => %s", err)
		return err
	}
	return nil
}

// Delete publishing only for sellers
func (p *Publishings) DeletePublishing() error {
	if err := DB.Delete(p).Error; err != nil {
		return err
	}
	if err := DB.Delete(&Comments{}, "publishing_id = ?", p.ID).Error; err != nil {
		go loggerError.Printf("Operation DeletePublishing => %s", err)
		return err
	}
	if err := DB.Delete(&Views{}, "publishing_id = ?", p.ID).Error; err != nil {
		go loggerError.Printf("Operation DeletePublishing => %s", err)
		return err
	}
	return nil
}

//********************************************
//**************End section*******************
//********************************************

// Helper for serializing
func singStruct(p *Publishings, ctx *fiber.Ctx) {
	var (
		imageName1 string
		imageName2 string
		imageName3 string
		imageName4 string
		imageName5 string
	)
	image1, err1 := ctx.FormFile("image_1")
	image2, err2 := ctx.FormFile("image_2")
	image3, err3 := ctx.FormFile("image_3")
	image4, err4 := ctx.FormFile("image_4")
	image5, err5 := ctx.FormFile("image_5")

	if id, err := strconv.Atoi(ctx.FormValue("publishing_id", "")); err == nil {
		p.ID = uint(id)
	}
	if user_id, err := strconv.Atoi(ctx.FormValue("user_id", "")); err == nil {
		p.UserId = uint(user_id)
	}
	if price, err := strconv.ParseFloat(ctx.FormValue("price", ""), 64); err == nil {
		p.Price = price
	}
	if ram, err := strconv.Atoi(ctx.FormValue("ram", "")); err == nil {
		p.Ram = ram
	}
	if stock, err := strconv.Atoi(ctx.FormValue("stock", "")); err == nil {
		p.Stock = stock
	}
	if store, err := strconv.Atoi(ctx.FormValue("store", "")); err == nil {
		p.Store = store
	}
	if screen, err := strconv.ParseFloat(ctx.FormValue("screen", ""), 64); err == nil {
		p.Screen = float64(screen)
	}
	if InOffer, err := strconv.ParseFloat(ctx.FormValue("in_offer", ""), 64); err == nil {
		p.InOffer = InOffer
	}
	if long, err := strconv.ParseFloat(ctx.FormValue("long", ""), 64); err == nil {
		p.Long = long
	}
	if width, err := strconv.ParseFloat(ctx.FormValue("width", ""), 64); err == nil {
		p.Width = width
	}
	if height, err := strconv.ParseFloat(ctx.FormValue("height", ""), 64); err == nil {
		p.Height = height
	}
	if weight, err := strconv.ParseFloat(ctx.FormValue("weight", ""), 64); err == nil {
		p.Weight = weight
	}
	p.ProductName = ctx.FormValue("product_name", "")
	p.Category = ctx.FormValue("category", "")
	p.Subcategory = ctx.FormValue("subcategory", "")
	p.Description = ctx.FormValue("description", "")
	p.Brand = ctx.FormValue("brand", "")
	p.Os = ctx.FormValue("os", "")
	p.Micro = ctx.FormValue("micro", "")
	p.Gpu = ctx.FormValue("gpu", "")
	p.Resolution = ctx.FormValue("resolution", "")

	timeNow := time.Now()
	dir := fmt.Sprintf("%d-%d-%d", timeNow.Year(), timeNow.Month(), timeNow.Day())
	localpath := fmt.Sprintf("./static/images/publishings/%s", dir)

	if err1 == nil {
		imgID := uuid.New()
		ext := strings.Split(image1.Filename, ".")[1]
		imageName1 = fmt.Sprintf("%s.%s", imgID, ext)
		p.Image1 = fmt.Sprintf("/static/images/publishings/%s/%s", dir, imageName1)
		go helpers.Serializer(ctx, localpath, imageName1, image1)
	}
	if err2 == nil {
		imgID := uuid.New()
		ext := strings.Split(image2.Filename, ".")[1]
		imageName2 = fmt.Sprintf("%s.%s", imgID, ext)
		p.Image2 = fmt.Sprintf("/static/images/publishings/%s/%s", dir, imageName2)
		go helpers.Serializer(ctx, localpath, imageName2, image2)
	}
	if err3 == nil {
		imgID := uuid.New()
		ext := strings.Split(image3.Filename, ".")[1]
		imageName3 = fmt.Sprintf("%s.%s", imgID, ext)
		p.Image3 = fmt.Sprintf("/static/images/publishings/%s/%s", dir, imageName3)
		go helpers.Serializer(ctx, localpath, imageName3, image3)
	}
	if err4 == nil {
		imgID := uuid.New()
		ext := strings.Split(image4.Filename, ".")[1]
		imageName4 = fmt.Sprintf("%s.%s", imgID, ext)
		p.Image4 = fmt.Sprintf("/static/images/publishings/%s/%s", dir, imageName4)
		go helpers.Serializer(ctx, localpath, imageName4, image4)
	}
	if err5 == nil {
		imgID := uuid.New()
		ext := strings.Split(image5.Filename, ".")[1]
		imageName5 = fmt.Sprintf("%s.%s", imgID, ext)
		p.Image5 = fmt.Sprintf("/static/images/publishings/%s/%s", dir, imageName5)
		go helpers.Serializer(ctx, localpath, imageName5, image5)
	}

}
