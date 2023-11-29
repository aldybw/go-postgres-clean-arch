package http

import (
	"go-postgres-clean-arch/domain"

	"github.com/labstack/echo"
)

type TagHandler struct {
	TUsecase domain.TagUseCase
}

// NewTagHandler will initialize the tags/ resources endpoint
func NewTagHandler(e *echo.Echo, tu domain.TagUseCase) *TagHandler {
	// return &TagHandler{TUsecase: tu}
	handler := &TagHandler{
		TUsecase: tu,
	}

	// e.GET("/tags", handler.FindByAll)
	// e.GET("/tags/:tagId", handler.FindById)
	// e.POST("/tags", handler.Create)
	// e.PATCH("/tags/:tagId", handler.Update)
	// e.DELETE("/tags/:tagId", handler.Delete)
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
