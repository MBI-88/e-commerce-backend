package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// *********************************************
// *********************************************
// *********************************************
// Error Response
type errorRespnse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

// Validate structs
func ValidateStruct(data interface{}) []*errorRespnse {
	var errors []*errorRespnse
	if err := validate.Struct(data); err != nil {
		for _, er := range err.(validator.ValidationErrors) {
			var element errorRespnse
			element.FailedField = er.StructNamespace()
			element.Tag = er.Tag()
			element.Value = er.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

//**********************************************
// *********************************************
// *********************************************

// **********************************************
// *********************************************
// *********************************************
// Parse publishings query
type ParseQuery struct {
	Category    string  `query:"category"`
	Subcategory string  `query:"subcategory"`
	ProductName string  `query:"product_name"`
	MinPrice    float64 `query:"min_price"`
	MaxPrice    float64 `query:"max_price"`
	Os          string  `query:"os"`
	MinRam      int     `query:"min_ram"`
	MaxRam      int     `query:"max_ram"`
	MinStore    int     `query:"min_store"`
	MaxStore    int     `query:"max_store"`
	MinScreen   float64 `query:"min_screen"`
	MaxScreen   float64 `query:"max_screen"`
	Brand       string  `query:"brand"`
	Resolution  string  `query:"resolution"`
	Gpu         string  `query:"gpu"`
	Micro       string  `query:"micro"`
	MinLong     float64 `query:"min_long"`
	MaxLong     float64 `query:"max_long"`
	MinWidth    float64 `query:"min_width"`
	MaxWidth    float64 `query:"max_width"`
	MinHeight   float64 `query:"min_height"`
	MaxHeight   float64 `query:"max_height"`
	MinWeight   float64 `query:"min_weight"`
	MaxWeight   float64 `query:"max_weight"`
}

//**********************************************
// *********************************************
// *********************************************

// **********************************************
// *********************************************
// *********************************************
// Pagination
func Page(page, size int) (int, int) {
	if size <= 0 {
		size = 50
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	return offset, size
}

//**********************************************
// *********************************************
// *********************************************

// **********************************************
// *********************************************
// *********************************************
// Serialise file into static dir
func Serializer(ctx *fiber.Ctx, localpath, imageName string, image *multipart.FileHeader) {
	if _, err := os.Stat(localpath); os.IsNotExist(err) {
		os.Mkdir(localpath, 0775)
		ctx.SaveFile(image, fmt.Sprintf("%s/%s", localpath, imageName))
	} else {
		ctx.SaveFile(image, fmt.Sprintf("%s/%s", localpath, imageName))
	}
}

// **********************************************
// *********************************************
// *********************************************
// Date selector
func DateSelector(from, to string) (time.Time, time.Time,error) {
	var (
		startDate time.Time
		endDate   time.Time
	)
	if from == "" || to == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		next := startDate.AddDate(0, 0, 29)
		endDate = time.Date(next.Year(), next.Month(), next.Day(), 23, 59, 59, 59, next.Location())
		return startDate, endDate,nil
	}
	arrayFrom := strings.Split(from, "-")
	arrayTo := strings.Split(to, "-")
	year, err := strconv.Atoi(arrayFrom[0])
	if err != nil {
		return startDate, endDate, err
	}
	month, err := strconv.Atoi(arrayFrom[1])
	if err != nil {
		return startDate, endDate, err
	}
	day, err := strconv.Atoi(arrayFrom[2])
	if err != nil {
		return startDate, endDate, err
	}
	startDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())
	year, err = strconv.Atoi(arrayTo[0])
	if err != nil {
		return startDate, endDate, err
	}
	month, err = strconv.Atoi(arrayTo[1])
	if err != nil {
		return startDate, endDate, err
	}
	day, err = strconv.Atoi(arrayTo[2])
	if err != nil {
		return startDate, endDate, err
	}
	endDate = time.Date(year, time.Month(month), day, 23, 59, 59, 59, time.Now().Location())
	return startDate, endDate, nil
}
