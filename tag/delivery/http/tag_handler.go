package http

import (
	"fmt"
	"go-postgres-clean-arch/domain"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type TagHandler struct {
	TUsecase domain.TagUseCase
}

// NewTagHandler will initialize the tags/ resources endpoint
func NewTagHandler(e *echo.Echo, tu domain.TagUseCase) {
	handler := &TagHandler{
		TUsecase: tu,
	}

	baseRouter := e.Group("/api")
	tagsRouter := baseRouter.Group("/tags")
	tagsRouter.GET("/welcome", func(ctx echo.Context) (err error) { return ctx.JSON(http.StatusOK, "welcome aldy") })
	tagsRouter.GET("", handler.FetchTag)
	tagsRouter.GET("/:tagId", handler.GetByID)
	tagsRouter.POST("", handler.Store)
	tagsRouter.PATCH("/:tagId", handler.Update)
	tagsRouter.DELETE("/:tagId", handler.Delete)
}

// FetchTag will fetch the tag based on given params
func (t *TagHandler) FetchTag(c echo.Context) error {
	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()

	listTag, nextCursor, err := t.TUsecase.Fetch(ctx, cursor, int64(num))
	// listTag, nextCursor, err := t.TUsecase.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listTag)
}

// GetByID will get tag by given id
func (t *TagHandler) GetByID(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		fmt.Println("Error pada ID: " + err.Error())
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	tag, err := t.TUsecase.FetchByID(ctx, id)
	if err != nil {
		fmt.Println("Error gagal fetch by ID: " + err.Error())
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tag)
}

// // GetByName will get tag by given name
// func (t *TagHandler) GetByName(c echo.Context) error {
// 	name := c.QueryParam("name")

// 	ctx := c.Request().Context()

// 	tags, err := t.TUsecase.FetchByName(ctx, name)
// 	if err != nil {
// 		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, tags)
// }

// Store will store the tag by given request body
func (t *TagHandler) Store(c echo.Context) (err error) {
	var tag domain.Tag
	err = c.Bind(&tag)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&tag); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = t.TUsecase.Store(ctx, &tag)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	var tagResponse TagResponse
	tagResponse.Name = tag.Name

	return c.JSON(http.StatusCreated, tagResponse)
}

// Update will update the tag by request body based on param id
func (t *TagHandler) Update(c echo.Context) (err error) {
	idP, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()
	_, err = t.TUsecase.FetchByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	var tag domain.Tag
	err = c.Bind(&tag)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&tag); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var tagResponse TagResponse
	tagResponse.Name = tag.Name

	err = t.TUsecase.Update(ctx, id, &tag)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tagResponse)
}

// Delete will delete tag by given param
func (t *TagHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	id := int64(idP)
	ctx := c.Request().Context()

	err = t.TUsecase.Delete(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

type ResponseError struct {
	Message string `json:"message"`
}

type TagResponse struct {
	Name string `json:"name"`
}

func isRequestValid(m *domain.Tag) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// // Create Tag Handler
// func (t *TagHandler) Create(ctx *gin.Context) {
// 	createTagRequest := domain.CreateTagRequest{}
// 	err := ctx.ShouldBindJSON(&createTagRequest)
// 	helper.ErrorPanic(err)

// 	errResponse := t.TUsecase.Store(&createTagRequest)
// 	if errResponse != nil {
// 		webResponse := helper.Response{
// 			Code:   http.StatusBadRequest,
// 			Status: "Error",
// 			Data:   nil,
// 		}
// 		ctx.Header("Content-Type", "application/json")
// 		ctx.JSON(http.StatusOK, webResponse)
// 	}
// 	webResponse := helper.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// // Update Tag Hanlder
// func (t *TagHandler) Update(ctx *gin.Context) {
// 	updateTagRequest := domain.UpdateTagRequest{}
// 	err := ctx.ShouldBindJSON(&updateTagRequest)
// 	helper.ErrorPanic(err)

// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)
// 	updateTagRequest.ID = id

// 	errResponse := t.TUsecase.Update(&updateTagRequest)
// 	if errResponse != nil {
// 		webResponse := helper.Response{
// 			Code:   http.StatusBadRequest,
// 			Status: "Error",
// 			Data:   nil,
// 		}
// 		ctx.Header("Content-Type", "application/json")
// 		ctx.JSON(http.StatusOK, webResponse)
// 	}
// 	webResponse := helper.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// // Delete Tag Handler
// func (t *TagHandler) Delete(ctx *gin.Context) {
// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)

// 	errResponse := t.TUsecase.Delete(id)
// 	if errResponse != nil {
// 		webResponse := helper.Response{
// 			Code:   http.StatusBadRequest,
// 			Status: "Error",
// 			Data:   nil,
// 		}
// 		ctx.Header("Content-Type", "application/json")
// 		ctx.JSON(http.StatusOK, webResponse)
// 	}
// 	webResponse := helper.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// // FindById Tag Handler
// func (t *TagHandler) FindById(ctx *gin.Context) {
// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)

// 	tagResponse, errResponse := t.TUsecase.FetchByID(id)

// 	if errResponse != nil {
// 		webResponse := helper.Response{
// 			Code:   http.StatusBadRequest,
// 			Status: "Error",
// 			Data:   nil,
// 		}
// 		ctx.Header("Content-Type", "application/json")
// 		ctx.JSON(http.StatusOK, webResponse)
// 	}

// 	webResponse := helper.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   tagResponse,
// 	}

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// // FindByAll Tag Handler
// func (t *TagHandler) FindByAll(ctx *gin.Context) {
// 	tagResponse, err := t.TUsecase.Fetch()

// 	if err != nil {
// 		webResponse := helper.Response{
// 			Code:   http.StatusBadRequest,
// 			Status: "Error",
// 			Data:   nil,
// 		}
// 		ctx.Header("Content-Type", "application/json")
// 		ctx.JSON(http.StatusOK, webResponse)
// 	}

// 	webResponse := helper.Response{
// 		Code:   http.StatusOK,
// 		Status: "Ok",
// 		Data:   tagResponse,
// 	}

// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }
