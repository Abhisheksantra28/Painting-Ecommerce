package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/niladri2003/PaintingEcommerce/app/models"
	"github.com/niladri2003/PaintingEcommerce/pkg/middleware"
	"github.com/niladri2003/PaintingEcommerce/platform/database"
	"time"
)

// CreateAddress handles creating a new address for a specific user.
func CreateAddress(c *fiber.Ctx) error {

	var addressForm models.AddressForm

	now := time.Now().Unix()

	// Get claims from JWT
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "token invalid",
		})
	}
	expires := claims.Expires

	// Check if token is expired
	if now > expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "token expired",
		})
	}
	if claims.UserRole != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Only user can create address",
		})
	}

	// Parse request body
	if err := c.BodyParser(&addressForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid request body",
		})
	}

	if addressForm.FirstName == "" || addressForm.LastName == "" || addressForm.Country == "" || addressForm.StreetAddress == "" || addressForm.TownCity == "" || addressForm.State == "" || addressForm.PinCode == "" || addressForm.MobileNumber == "" || addressForm.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "All fields are required",
		})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	addresses, err := db.GetAddressesByUserID(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error while fetching address in CreateAddress",
		})
	}

	// Handle SetAsDefault logic
	var setAsDefault *bool
	if addressForm.SetAsDefault != nil {
		// If SetAsDefault is provided in the request, use it
		defaultVal := *addressForm.SetAsDefault
		setAsDefault = &defaultVal
	} else {
		// If SetAsDefault is not provided, default to false
		defaultVal := false
		setAsDefault = &defaultVal

		// However, if this is the user's first address, set it to true
		if len(addresses) == 0 {
			defaultVal = true
			setAsDefault = &defaultVal
		}
	}

	// Create a new address model
	address := &models.Address{
		ID:            uuid.New(),
		UserID:        claims.UserID,
		FirstName:     addressForm.FirstName,
		LastName:      addressForm.LastName,
		Country:       addressForm.Country,
		StreetAddress: addressForm.StreetAddress,
		TownCity:      addressForm.TownCity,
		State:         addressForm.State,
		PinCode:       addressForm.PinCode,
		MobileNumber:  addressForm.MobileNumber,
		Email:         addressForm.Email,
		OrderNotes:    addressForm.OrderNotes, // Handle as optional
		SetAsDefault:  setAsDefault,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	// Create a new address in the database
	id, err := db.CreateAddress(address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Address created successfully",
		"id":    id.String(),
	})
}

func UpdateAddressByAddressId(c *fiber.Ctx) error {
	// Get address ID from the URL parameters
	addressId := c.Params("addressId")
	parsedAddressId, err := uuid.Parse(addressId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid address ID",
		})
	}

	// Get user claims from JWT
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Check if token is expired
	now := time.Now().Unix()
	if now > claims.Expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Token expired",
		})
	}
	if claims.UserRole != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Only users can update addresses",
		})
	}

	// Parse the request body
	var addressForm models.AddressForm
	if err := c.BodyParser(&addressForm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid request body",
		})
	}

	// Validate required fields
	if addressForm.FirstName == "" || addressForm.LastName == "" || addressForm.Country == "" || addressForm.StreetAddress == "" ||
		addressForm.TownCity == "" || addressForm.State == "" || addressForm.PinCode == "" ||
		addressForm.MobileNumber == "" || addressForm.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "All fields are required",
		})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get the existing address to check the current SetAsDefault value
	existingAddress, err := db.GetAddressByAddresId(parsedAddressId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error while fetching address",
		})
	}

	// Determine the SetAsDefault value
	setAsDefault := existingAddress[0].SetAsDefault
	if addressForm.SetAsDefault != nil {
		// If SetAsDefault is provided, use the new value
		setAsDefault = addressForm.SetAsDefault
	}

	// Update the address with the new data
	address := models.Address{
		ID:            parsedAddressId,
		UserID:        claims.UserID,
		FirstName:     addressForm.FirstName,
		LastName:      addressForm.LastName,
		Country:       addressForm.Country,
		StreetAddress: addressForm.StreetAddress,
		TownCity:      addressForm.TownCity,
		State:         addressForm.State,
		PinCode:       addressForm.PinCode,
		MobileNumber:  addressForm.MobileNumber,
		Email:         addressForm.Email,
		OrderNotes:    addressForm.OrderNotes,
		SetAsDefault:  setAsDefault,
		UpdatedAt:     time.Now().Format(time.RFC3339),
	}

	// Update the address in the database
	if err := db.UpdateAddressByAddressId(parsedAddressId, &address); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Fetch the updated address to return in the response
	updatedAddress, err := db.GetAddressByAddresId(parsedAddressId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error while fetching updated address",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     "Address updated successfully",
		"address": updatedAddress,
	})
}

func GetAddressByUserID(c *fiber.Ctx) error {
	// Get user ID from JWT claims
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userId := claims.UserID

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get address by user ID
	address, err := db.GetAddressesByUserID(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     "Address retrieved successfully",
		"address": address,
	})
}

// GetAllAddresses retrieves all addresses from the database
func GetAllAddresses(c *fiber.Ctx) error {
	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all addresses from the database
	addresses, err := db.GetAddresses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":     false,
		"msg":       "Addresses retrieved successfully",
		"addresses": addresses,
	})
}

// DeleteAddressByUserID deletes an address for a specific user by address ID.
func DeleteAddressByUserID(c *fiber.Ctx) error {
	// Get user ID from JWT claims
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "token invalid",
		})
	}

	userID := claims.UserID

	// Check if token is expired
	now := time.Now().Unix()
	if now > claims.Expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Token expired",
		})
	}
	if claims.UserRole != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Only users can delete addresses",
		})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Delete the address by userID
	err = db.DeleteAddressByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Address deleted successfully",
	})
}

// DeleteAddress deletes an address by address ID.
func DeleteAddressByAddressId(c *fiber.Ctx) error {

	addressId := c.Params("addressId")
	parsedAddressId, err := uuid.Parse(addressId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid address ID",
		})
	}

	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Internal server error",
		})
	}

	// Check if token is expired
	now := time.Now().Unix()
	if now > claims.Expires {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Token expired",
		})
	}

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Delete the address by address ID
	err = db.DeleteAddress(parsedAddressId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "Address deleted successfully",
	})
}

func SetDefaultAddress(c *fiber.Ctx) error {
	addressID := c.Params("addressId")
	// Convert addressID to uuid.UUID
	parsedAddressId, err := uuid.Parse(addressID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid address ID",
		})
	}

	// Get user ID from JWT claims
	claims, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userId := claims.UserID

	// Create database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set the default address
	err = db.SetDefaultAddressByID(parsedAddressId, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to set default address",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Default address updated successfully",
	})
}
