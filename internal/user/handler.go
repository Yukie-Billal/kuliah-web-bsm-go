package user

import (
	"kuliah-web-bsm-go/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Index(c *gin.Context) {

	users, err := h.service.FindAll()

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, users)
}

func (h *Handler) Show(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.service.FindById(uint(id))

	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

func (h *Handler) Store(c *gin.Context) {

	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.service.Create(&user)

	if err != nil {
		response.InternalServerError(c, err.Error())

		return
	}

	response.Success(c, user)
}

func (h *Handler) Update(c *gin.Context) {

	id, _ := strconv.Atoi(
		c.Param("id"),
	)

	existingUser, err :=
		h.service.FindById(uint(id))

	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	var payload User

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	existingUser.Username = payload.Username
	existingUser.Email = payload.Email

	err = h.service.Update(existingUser)

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, existingUser)
}

func (h *Handler) Delete(c *gin.Context) {

	id, _ := strconv.Atoi(
		c.Param("id"),
	)

	err := h.service.Delete(
		uint(id),
	)

	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, "")
}
