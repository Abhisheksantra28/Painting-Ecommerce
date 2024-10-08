package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/niladri2003/PaintingEcommerce/app/models"
	"github.com/niladri2003/PaintingEcommerce/app/uploader"
	"github.com/niladri2003/PaintingEcommerce/pkg/middleware"
	"github.com/niladri2003/PaintingEcommerce/platform/database"
	"net/http"
	"strconv"
	"time"
)

func CreateProduct(c *fiber.Ctx) error {

	now := time.Now().Unix()

	//Get claims from jwt
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "token invalid",
		})
	}
	expires := claims.Expires

	//Checking if now time is greater than expiration from jwt
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "token expired",
		})
	}
	if claims.UserRole != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Only admin can create Product",
		})
	}
	// Parse product details
	title := c.FormValue("title")
	description := c.FormValue("description")
	originalPrice := c.FormValue("original_price")
	discountedPrice := c.FormValue("discounted_price")
	is_active := c.FormValue("is_active")
	categoryID := c.FormValue("category_id")

	// Validate inputs
	if title == "" || description == "" || originalPrice == "" || categoryID == "" || discountedPrice == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	is_activeBool, err := strconv.ParseBool(is_active)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid is_active flag"})
	}
	// Convert price to float
	originalPriceValue, err := strconv.ParseFloat(originalPrice, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid price format"})
	}
	discountedPriceValue, err := strconv.ParseFloat(discountedPrice, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid price format"})
	}
	if originalPrice < discountedPrice {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": true, "msg": "Original price is lower than discounted price"})
	}
	// Convert categoryID to uuid.UUID
	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID format"})
	}

	product := &models.Product{}

	// Generate UUID for the new product
	productId := uuid.New()

	product.ID = productId
	product.Title = title
	product.Description = description
	product.OriginalPrice = originalPriceValue
	product.DiscountedPrice = discountedPriceValue
	product.IsActive = is_activeBool
	product.CategoryID = categoryUUID
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	//Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Insert Product into database.
	if err := db.CreateProduct(product); err != nil {
		// Return status 500 and create category process error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	//Handle multipart files
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Could not parse multipart form",
		})
	}

	files := form.File["images"]
	if files == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No images uploaded"})
	}

	for _, file := range files {
		filePath := fmt.Sprintf("/tmp/%s", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save file"})
		}

		imageURL, err := uploader.UploadImage(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not upload image"})
		}
		//Create a struct of product image model
		productImage := &models.ProductImage{}
		//Initialize product image details
		productImage.ID = uuid.New()
		productImage.ProductID = productId
		productImage.ImageURL = imageURL
		productImage.CreatedAt = time.Now()

		// Insert image URL into the database
		if err := db.CreateProductImage(productImage); err != nil {
			// Return status 500 and create category process error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "Could not save image",
			})
		}
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"product": product,
		"meg":     "success",
	})
}
func DeleteProduct(c *fiber.Ctx) error {
	now := time.Now().Unix()

	//Get claims from jwt
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "token invalid",
		})
	}
	expires := claims.Expires

	//Checking if now time is greater than expiration from jwt
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "token expired",
		})
	}
	if claims.UserRole != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Only admin can create Product",
		})
	}

	// Validate the ID
	productID := c.Params("productId")

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	//Fetch product Images
	productImages, err := db.GetImagesByProduct(productUUID)
	if err != nil {
		fmt.Println("cloudinary delete error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "true", "msg": "Could not retrieve product images"})
	}
	for _, image := range productImages {
		if err := uploader.DeleteImage(image.ImageURL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "true", "msg": "Could not delete image from cloudinary"})
		}
	}

	// Delete product  images
	if err := db.DeleteProductImage(productUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Delete product
	if err := db.DeleteProduct(productUUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"error": "false", "msg": "Product and images deleted successfully"})
}
func GetProductDetails(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate the ID
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}
	fmt.Println(productUUID)
	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	productDetails, err := db.GetProduct(productUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	productImage, err := db.GetImagesByProduct(productUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	allProducts, err := db.GetProductsByCategory(productDetails.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Product details retrieved successfully", "data": productDetails, "images": productImage, "related_products": allProducts})

}
func GetAllProducts(c *fiber.Ctx) error {
	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	allProducts, err := db.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Product details retrieved successfully", "data": allProducts})

}
func GetTop6Products(c *fiber.Ctx) error {
	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	allProducts, err := db.GetTop6Products()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Product details retrieved successfully", "data": allProducts})
}
func GetTop5ProductsCategoryWise(c *fiber.Ctx) error {
	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	allProducts, err := db.GetTopProductsByCategory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Product details retrieved successfully", "data": allProducts})

}
func GetProductsByCategoryID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse the ID as uuid.UUID
	categoryId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid contact ID",
		})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	allProducts, err := db.GetProductsByCategory(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"msg": "Product details retrieved successfully", "data": allProducts})

}
func ChangeProductStatus(c *fiber.Ctx) error {

	type isActiveRequest struct {
		IsActive string `json:"is_active"`
	}
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": "token invalid"})
	}

	if time.Now().Unix() > claims.Expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "msg": "token expired"})
	}

	if claims.UserRole != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "msg": "only users can create order"})
	}
	id := c.Params("productId")
	// Parse the ID as uuid.UUID
	ProductId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "msg": "Invalid Product ID"})
	}
	var request isActiveRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Convert the is_active string to boolean
	isActiveBool, err := strconv.ParseBool(request.IsActive)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid is_active value"})
	}
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.ChnageProductStatus(ProductId, isActiveBool); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"msg": "Product status changed successfully"})

}
