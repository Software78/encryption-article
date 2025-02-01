package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	middleware "github.com/Software78/encryption-test/src/middleware"
	models "github.com/Software78/encryption-test/src/models"
	services "github.com/Software78/encryption-test/src/services"

	// utils "github.com/Software78/encryption-test/src/utils"
	"github.com/gin-gonic/gin"
)



type UserController struct {
	userService services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{userService: service}
}

// Login godoc
//
//	@Summary		Login a user
//	@Description	Login a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.Login	true	"User object that needs to be created"
//	@Success		200		{object}	utils.SuccessResponse
//	@Router			/auth/login [post]
func (h *UserController) Login(c *gin.Context) {
	login := &models.Login{}
	if err := c.ShouldBindJSON(login); err != nil {
		c.Error(err)
		return
	}
	decryptedJSON := c.MustGet("decryptedJSON").(map[string]interface{})
	decryptedLogin := models.Login{
		Email:    decryptedJSON["email"].(string),
		Password: decryptedJSON["password"].(string),
	}
	user, err := h.userService.Login(&decryptedLogin)
	if err != nil {
		c.Error(err)
		return
	}
	crypto, _ := middleware.NewCryptoMiddlewareFromEnv( `/docs/`)
	encryptedUser , err :=	crypto.EncryptValues(user)
	fmt.Println(encryptedUser)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, json.RawMessage(encryptedUser))
}


// Register godoc
//
//	@Summary		Register a user
//	@Description	Register a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.Register	true	"User object that needs to be created"
//	@Success		200		{object}	utils.SuccessResponse
//	@Router			/auth/register [post]
func (h *UserController) Register(c *gin.Context) {
	register := &models.Register{}
	if err := c.ShouldBindJSON(register); err != nil {
		c.Error(err)
		return
	}
	decryptedJSON := c.MustGet("decryptedJSON").(map[string]interface{})
	decryptedRegister := models.Register{
		FirstName: decryptedJSON["first_name"].(string),
		LastName:  decryptedJSON["last_name"].(string),
		Email:     decryptedJSON["email"].(string),
		Password:  decryptedJSON["password"].(string),
	}
	_, err := h.userService.Register(&decryptedRegister)
	if err != nil {
		c.Error(err)
		return
	}
	registeredUser,err := h.userService.GetUserByEmail(decryptedRegister.Email)

	if err != nil {
		c.Error(err)
		return
	}
	crypto, _ := middleware.NewCryptoMiddlewareFromEnv( `/docs/`)
	encryptedUser , err :=	crypto.EncryptValues(registeredUser)

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK,  json.RawMessage(encryptedUser))
}
