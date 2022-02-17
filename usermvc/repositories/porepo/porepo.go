package porepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	// "strconv"
	"database/sql"
	// "reflect"
	"usermvc/entity"
	"usermvc/model"
	"usermvc/repositories"
	logger2 "usermvc/utility/logger"
	
)

type PoRepo interface {
	GetPoCreationInfo(ctx context.Context, req *model.Input) (interface{}, error)
	GetPOFormInfo(ctx context.Context, req *model.GetPoFormInfoRequestBody) (interface{}, error)
	ViewPoDetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error)
	// ListPurchaseOrders(ctx context.Context, req *model.ListPurchaseOrderRequest) (interface{}, error)
	GetPortandOrigin(ctx context.Context, req *model.Input) (interface{}, error)
	GetBalQuoteQtyForPoOrder(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error)
	EditGCPODetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error)
	InsertGCPODetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error)

}
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
// func NewNullTime(t time.Time) sql.NullTime {
// 	if t.IsZero() {
// 		return sql.NullTime{}
// 	}
// 	return sql.NullTime{
// 		Time:  t,
// 		Valid: true,
// 	}
// }

type poRepo struct {
	db *gorm.DB
}

func NewPoRepo() PoRepo {
	newDb, err := repositories.NewDb()
	if err != nil {
		panic(err)
	}
	newDb.AutoMigrate(&entity.User{})
	return &poRepo{
		db: newDb,
	}
}
func (po poRepo) GetPoCreationInfo(ctx context.Context, req *model.Input) (interface{}, error) {
	const (
		CONTAINERTYPES = "containerTypes"
		)
	logger := logger2.GetLoggerWithContext(ctx)
	if req.Type == CONTAINERTYPES {
		rows, err := po.db.Raw("select conttypeid, conttypename from dbo.sales_container_types").Rows()
		if err != nil {
			logger.Error("error while fetching records from dbo.sales_container_types ", err.Error())
		}
		var allContainerTypes []model.ContainerTypesList
		defer rows.Close()
		for rows.Next() {
			var ct model.ContainerTypesList
			err = rows.Scan(&ct.ConttypeId, &ct.ConttypeName)
			allContainerTypes = append(allContainerTypes, ct)
			
		}
		return allContainerTypes, err
		
	}
	return nil,nil
}

func (po poRepo) GetPOFormInfo(ctx context.Context, req *model.GetPoFormInfoRequestBody) (interface{}, error) {
	const (
		POSUBCATEGORY = "posubcategory"
		SUPPLIERINFO  = "supplierinfo"
		BILLINGINFO   = "billinginfo"
		DELIVERYINFO  = "deliveryinfo"
		ALLSUPPLIERS  = "allsuppliers"
		GREENCOFFEE   = "greencoffee"
		GCCOMPOSITION = "gccomposition"
	)
	logger := logger2.GetLoggerWithContext(ctx)
	if req.Type == POSUBCATEGORY {
		rows, err := po.db.Raw("select vendortypeid,initcap(vendortypename) from dbo.pur_vendor_types").Rows()
		if err != nil {
			logger.Error("error while getting records from dbo.pur_vendor_types ", err.Error())
		}
		var allPOSubs []model.GetPoFormInfoRequestBody
		defer rows.Close()
		for rows.Next() {
			var pos model.GetPoFormInfoRequestBody
			err = rows.Scan(&pos.SupplierTypeID, &pos.SupplierTypeName)
			allPOSubs = append(allPOSubs, pos)
			return allPOSubs, nil
		}
	}

	if req.Type == SUPPLIERINFO {
		sqlStatement2 := `select country,vendorid,initcap(vendorname),initcap(address1)||','||initcap(address2)||','||initcap(city)||','||pincode||','||initcap(state)||','||'Phone:'||phone||','||'Mobile:'||mobile||','||'GST NO:'||gstin address 
							from dbo.pur_vendor_master where vendorid=$1`

		rows, err := po.db.Raw(sqlStatement2, req.SupplierID).Rows()
		if err != nil {
			logger.Error("error while getting records from dbo.pur_vendor_master  ", err.Error())
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&req.SupplierCoun, &req.SupplierID, &req.SupplierName, &req.SupplierAddress)
		}

		if req.SupplierCoun != "INDIA" {
			req.SupplierType = "International"
		} else {
			req.SupplierType = "Domestic"
			return req, nil
		}
	}
	if req.Type == BILLINGINFO {
		logger.Info("getting billing info details")
		sqlStatement := `select potypeid,initcap(potypename),initcap(potypefullname)||','||initcap(address) as fulladdress from dbo.pur_po_types`

		rows, err := po.db.Raw(sqlStatement).Rows()
		if err != nil {
			logger.Error("error while getting records from dbo.pur_vendor_master  ", err.Error())
			return nil, err
		}
		defer rows.Close()
		var allPCA []model.POCreatedAt
		defer rows.Close()
		for rows.Next() {
			var pca model.POCreatedAt
			err = rows.Scan(&pca.POTypeID, &pca.POTypeName, &pca.POAddress)
			allPCA = append(allPCA, pca)
		}
		return allPCA, nil
	}
	if req.Type == DELIVERYINFO {
		logger.Info("get deliveryinfo details:")
		sqlStatement := `select potypeid,initcap(potypename),initcap(potypefullname)||','||initcap(address) as fulladdress from dbo.pur_po_types`
		rows, err := po.db.Raw(sqlStatement).Rows()
		if err != nil {
			logger.Error("error while getting  deliveryinfo records from dbo.dbo.pur_po_types  ", err.Error())
			return nil, err
		}
		logger.Info("Query executed")
		var allPCF []model.POCreatedFor
		defer rows.Close()
		for rows.Next() {
			var pcf model.POCreatedFor
			err = rows.Scan(&pcf.POTypeID, &pcf.POTypeName, &pcf.POAddress)
			allPCF = append(allPCF, pcf)
		}
		return allPCF, nil
	}
	if req.Type == ALLSUPPLIERS {
		logger.Info("get vendors", req.Type)
		logger.Info("get vendors based on supplierID", req.SupplierTypeID)
		sqlStatementSup := `select vendorid,initcap(vendorname) from dbo.pur_vendor_master 
							where ((groupid='FAC-2') or (groupid='FAC-3') or (groupid='FAC-4') or (groupid='FAC-9')) and vendortypeid=$1`
		rows, err := po.db.Raw(sqlStatementSup, &req.SupplierTypeID).Rows()
		if err != nil {
			logger.Error("error while getting  deliveryinfo records from dbo.dbo.pur_po_types  ", err.Error())
			return nil, err
		}
		logger.Info("Query executed")
		var allSuppliers []model.GetPoFormInfoRequestBody
		defer rows.Close()
		for rows.Next() {
			var a model.GetPoFormInfoRequestBody
			err = rows.Scan(&a.SupplierID, &a.SupplierName)
			allSuppliers = append(allSuppliers, a)
		}
		return allSuppliers, nil
	}
	if req.Type == GREENCOFFEE {
		logger.Info("get GC", req.Type)
		sqlStatement4 := `select itemid,initcap(itemname),cat_type from dbo.inv_gc_item_master`
		rows, err := po.db.Raw(sqlStatement4).Rows()
		if err != nil {
			if err != nil {
				logger.Error("error while getting  GC records from dbo.dbo.pur_po_types  ", err.Error())
				return nil, err
			}
		}
		logger.Info("GC Query executed")
		var allGc []model.GreenCoffee
		defer rows.Close()
		for rows.Next() {
			var gc model.GreenCoffee
			err = rows.Scan(&gc.ItemID, &gc.ItemName, &gc.GCCoffeeType)
			allGc = append(allGc, gc)
		}

		return allGc, err
	}

	if req.Type == GCCOMPOSITION {
		logger.Info("get GC new composition based on the GD ID", req.Type)
		logger.Info("Entered Item id is:", &req.ItemID)
		sqlStatement5 := `select itemid,density,moisture,browns,blacks,brokenbits,insectedbeans,bleached,husk,sticks,stones,beansretained from dbo.pur_gc_po_composition_master_newpg
							where itemid=$1`
		rows, err := po.db.Raw(sqlStatement5, &req.ItemID).Rows()
		if err != nil {
			if err != nil {
				logger.Error("error while getting  getting   GC new composition based on the GD I ", err.Error())
				return nil, err
			}
		}
		logger.Info("GC Query executed")
		var allGcComp []model.GreenCoffee
		defer rows.Close()
		for rows.Next() {
			var gc model.GreenCoffee
			err = rows.Scan(&gc.ItemID, &gc.Density, &gc.Moisture, &gc.Browns, &gc.Blacks, &gc.BrokenBits, &gc.InsectedBeans,
				&gc.Bleached, &gc.Husk, &gc.Sticks, &gc.Stones, &gc.BeansRetained)
			allGcComp = append(allGcComp, gc)
		}
		logger.Info(allGcComp)
		return allGcComp, nil
	}
	errInfo := fmt.Sprintf("req type should be either of  %s %s %s %s %s %s  %s", POSUBCATEGORY, SUPPLIERINFO, BILLINGINFO, DELIVERYINFO, ALLSUPPLIERS, GREENCOFFEE, GCCOMPOSITION)
	return nil, errors.New(errInfo)
}

func (po poRepo) GetPortandOrigin(ctx context.Context, req *model.Input) (interface{}, error) {
	const (
		ORIGIN   = "originDetails"
		PORTLOAD = "portLoadingDetails"
	)
	logger := logger2.GetLoggerWithContext(ctx)
	if req.Type == ORIGIN {
		var origins []model.Origin
		rows, err := po.db.Raw("select distinct initcap(origin) from dbo.pur_gc_contract_master where origin != '';").Rows()
		if err != nil {
			logger.Error("error while getting record", err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var origin model.Origin
			err = rows.Scan(&origin.Origin)
			origins = append(origins, origin)
		}
		return origins, nil
	}
	if req.Type == PORTLOAD {
		var items []model.PortLoading
		rows, err := po.db.Raw("select distinct initcap(poloading) from dbo.pur_gc_contract_master where poloading != '';").Rows()
		if err != nil {
			logger.Error("error while getting record", err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var portsload model.PortLoading
			err = rows.Scan(&portsload.Port)
			items = append(items, portsload)
		}
		return items, nil
	}
	return nil, nil
}



func (po poRepo) findLatestkey(ctx context.Context,s... string)(i int){
	logger := logger2.GetLoggerWithContext(ctx)
	rows, err := po.db.Raw("select vendortypeid,initcap(vendortypename) from dbo.pur_vendor_types").Rows()
	if err != nil {
		logger.Error("error while getting records from dbo.pur_vendor_types ", err.Error())
	}
	
	defer rows.Close()
	for rows.Next() {
		
	}
//
return i
}



func (po poRepo) GetBalQuoteQtyForPoOrder(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	var input model.Input
	// var res []byte
	
	if input.Type == "getBalqtyforPo" {
		 logger.Error("get ordered qty on quotation id", input.Type)
		sqlStatement := `select sum(total_quantity) from dbo.pur_gc_po_con_master_newpg where quote_no=$1`
		rows, err := po.db.Raw(sqlStatement, input.QuotationId).Rows()
		if err != nil {
			logger.Error("unable to insert Audit Details", err)
		}
		var g model.QtyforPo
		var t sql.NullString
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&t)
		}
		g.OrderQty = t.String

	}

	return nil,nil
}