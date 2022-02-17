package controller

import (
	"github.com/gin-gonic/gin"
	"usermvc/model"
	"usermvc/repositories/porepo"
	logger2 "usermvc/utility/logger"
)

type PoController interface {
	GetPoCreationInfo(ctx *gin.Context)
	GetPOFormInfo(ctx *gin.Context)
	SenEmail(ctx *gin.Context)
	ViewGCPODetails(ctx *gin.Context)
	// ListPurchaseOrders(ctx *gin.Context)
	GetPortandOrigin(ctx *gin.Context)
	EditGCPODetails(ctx *gin.Context)
	GetBalQuoteQtyForPoOrder(ctx *gin.Context)
	InsertGCPODetails(ctx *gin.Context)
}

type poController struct {
	poRepo porepo.PoRepo
}

func NewPoController() PoController {
	return poController{
		poRepo: porepo.NewPoRepo(),
	}
}
func (po poController) GetPoCreationInfo(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var getPoCreationInfoRequestBody *model.Input
	logger.Info("started handling poController.GetPoCreationInfo")
	if err := ctx.ShouldBindJSON(&getPoCreationInfoRequestBody); err != nil {
		logger.Error("error while parsing the GetPoCreationInfo request body", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.GetPoCreationInfo(ctx, getPoCreationInfoRequestBody)
	if err != nil {
		ctx.JSON(503, "error while getting po creation info from repo")
		return
	}
	ctx.JSON(200, res)
	return
}
// func (po poController) 	ListPurchaseOrders(ctx *gin.Context) {
// 	logger := logger2.GetLoggerWithContext(ctx)
// 	var listPurchaseOrderRequest *model.ListPurchaseOrderRequest
// 	logger.Info("start hadling poController.listPurchaseOrderRequest")
// 	if err := ctx.ShouldBindJSON(&listPurchaseOrderRequest); err != nil {
// 		logger.Error("error while parsing the listPurchaseOrderRequest request body", err.Error())
// 		ctx.JSON(403, err.Error())
// 		return
// 	}
// 	res, err := po.poRepo.ListPurchaseOrders(ctx, listPurchaseOrderRequest)
// 	if err != nil {
// 		ctx.JSON(503, "error while getting po info from repo")
// 		return
// 	}
// 	ctx.JSON(200, res)
// 	return
// }
func (po poController) GetPOFormInfo(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var getPoFormInfoRequestBody *model.GetPoFormInfoRequestBody
	logger.Info("start hadling poController.GetPOFormInfo")
	if err := ctx.ShouldBindJSON(&getPoFormInfoRequestBody); err != nil {
		logger.Error("error while parsing the GetPOFormInfo request body", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.GetPOFormInfo(ctx, getPoFormInfoRequestBody)
	if err != nil {
		ctx.JSON(503, "error while getting po info from repo")
		return
	}
	ctx.JSON(200, res)
	return
}

func (po poController) SenEmail(ctx *gin.Context) {

}
func (po poController) GetPortandOrigin(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var port *model.Input
	logger.Info("start handling of get port and origin")
	if err := ctx.ShouldBindJSON(&port); err != nil {
		logger.Error("Error while parsing", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.GetPortandOrigin(ctx, port)
	if err != nil {
		ctx.JSON(503, "error while getting info")
		return
	}
	ctx.JSON(200, res)
	return

}
func (po poController) ViewGCPODetails(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var purchaseOrderDetailsRequest *model.PurchaseOrderDetails
	logger.Info("start hadling poController.ViewGCPODetails")
	if err := ctx.ShouldBindJSON(&purchaseOrderDetailsRequest); err != nil {
		logger.Error("error while parsing the ViewGCPODetails request body", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.ViewPoDetails(ctx, purchaseOrderDetailsRequest)
	if err != nil {
		ctx.JSON(503, "error while getting po info from repo")
		return
	}
	ctx.JSON(200, res)
	return
}

func (po poController) EditGCPODetails(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var gcpo *model.PurchaseOrderDetails
	logger.Info("start handling of edit gc po details")
	if err := ctx.ShouldBindJSON(&gcpo); err != nil {
		logger.Error("Error while parsing", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.EditGCPODetails(ctx, gcpo)
	if err != nil {
		ctx.JSON(503, "error while getting info")
		return
	}
	ctx.JSON(200, res)
	return

}
func (po poController) GetBalQuoteQtyForPoOrder(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var gbq *model.PurchaseOrderDetails
	logger.Info("start handling of edit gc po details")
	if err := ctx.ShouldBindJSON(&gbq); err != nil {
		logger.Error("Error while parsing", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.GetBalQuoteQtyForPoOrder(ctx, gbq)
	if err != nil {
		ctx.JSON(503, "error while getting info")
		return
	}
	ctx.JSON(200, res)
	return

}

func (po poController) InsertGCPODetails(ctx *gin.Context) {
	logger := logger2.GetLoggerWithContext(ctx)
	var ingcpo *model.PurchaseOrderDetails
	logger.Info("start handling of edit gc po details")
	if err := ctx.ShouldBindJSON(&ingcpo); err != nil {
		logger.Error("Error while parsing", err.Error())
		ctx.JSON(403, err.Error())
		return
	}
	res, err := po.poRepo.InsertGCPODetails(ctx, ingcpo)
	if err != nil {
		ctx.JSON(503, "error while getting info")
		return
	}
	ctx.JSON(200, res)
	return

}