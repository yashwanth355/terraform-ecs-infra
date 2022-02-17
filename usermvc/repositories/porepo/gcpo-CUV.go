package porepo
import (
	"context"
	"usermvc/model"
	// apputils "usermvc/utility/apputils"
	logger2 "usermvc/utility/logger"
	"strconv"
	"time"
	// "github.com/jinzhu/gorm"
	// "database/sql"
	// "encoding/json"
	// "fmt"
	// "log"
	// "bytes"
	// "net/smtp"
	// "text/template"

	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	// _ "github.com/lib/pq"
)

func (po poRepo) InsertGCPODetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error) {
	logger := logger2.GetLoggerWithContext(ctx)

	logger.Error("Entered Edit Module")
	return nil,nil
}



func (po poRepo) EditGCPODetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	
	var audit model.AuditLogGCPO

	if req.Update {
		logger.Error("Entered Edit Module")
		logger.Error("Created user id is: ", req.CreatedUserID)
		// Find created user username
		// sqlStatementCUser1 := `SELECT username
		// 					FROM dbo.users_master_newpg
		// 					where userid=$1`
		// rows, err = po.db.Raw(sqlStatementCUser1, po.CreatedUser)
		// for rows.Next() {
		// 	err = rows.Scan(&po.CreatedUserName)
		// }
		floatquan, _ := strconv.ParseFloat(req.TotalQuantity, 64)
		logger.Error("Formattted total quantity: ", floatquan)
		req.MT_Quantity = floatquan / 1000
		logger.Error("quantity in MT: ", req.MT_Quantity)
		req.TotalBalQuan = req.TotalQuantity

		if req.SupplierType == "1001" {
			//IMPORT MODULE
			logger.Error("Selected supplier type Import Code:", req.SupplierType)
			if req.CurrencyID == "" {
				req.CurrencyID = "HO-102"
			}

			sqlStatementImp1 := `update dbo.pur_gc_po_con_master_newpg
							set
							cid=$1,
							vendorid=$2,
							dispatchterms=$3,
							origin=$4,
							poloading=$5,
							insurance=$6,
							destination=$7,
							forwarding=$8,
							currencyid=$9,
							nocontainers=$10,
							payment_terms=$11,
							remarks=$12,
							billing_at_id=$13,
							delivery_at_id=$14,
							taxes_duties=$15,
							transport_mode=$16,
							transit_insurence=$17,
							packing_forward=$18,
							othercharges=$19,
							rate=$20,
							noofbags=$21,
							netweight=$22,
							container_type=$23,
							purchase_type=$24,
		   					terminal_month=$25,
		  					booked_term_rate=$26,
		  					booked_differential=$27, 
		   					fixed_term_rate=$28,
		  					fixed_differential=$29,
			  				purchase_price=$30,
			   				market_price=$31,
			   				po_margin=$32,
			   				total_price=$33,
							gross_price=$34,
							quantity_mt=$35,
							balance_quantity=$36							
							where pono=$37`
			_, err := po.db.Raw(sqlStatementImp1,
				req.Contract,
				req.SupplierID,
				req.IncoTerms,
				req.Origin,
				req.PortOfLoad,
				req.Insurance,
				req.PlaceOfDestination,
				req.Forwarding,
				req.CurrencyID,
				NewNullString(req.NoOfContainers),
				req.PaymentTerms,
				req.Comments,
				req.POBillTypeID,
				req.PODelTypeID,
				req.TaxDuties,
				req.ModeOfTransport,
				req.TransitInsurance,
				req.PackForward,
				req.OtherCharges,
				NewNullString(req.Rate),
				req.NoOfBags,
				req.NetWt,
				req.ContainerType,
				req.PurchaseType,
				req.TerminalMonth,
				NewNullString(req.BookedTerminalRate),
				NewNullString(req.BookedDifferential),
				NewNullString(req.FixedTerminalRate),
				NewNullString(req.FixedDifferential),
				NewNullString(req.PurchasePrice),
				NewNullString(req.MarketPrice),
				NewNullString(req.POMargin),
				NewNullString(req.TotalPrice),
				NewNullString(req.GrossPrice),
				req.MT_Quantity,
				req.TotalBalQuan,
				req.PoNO).Rows()
			logger.Error("Update into PO Table Executed")
			if err != nil {
				logger.Error(err.Error())
				  
			}

		} else if req.SupplierType == "1002" {
			logger.Error("Selected supplier type Domestic Code:", req.SupplierType)
			//-----------DOMESTIC INFO INSERT-----------------------
			//Green coffee id,name,quantity-Missing
			if req.CurrencyID == "" {
				req.CurrencyID = "HO-101"
			}

			sqlStatementImp1 := `Update dbo.pur_gc_po_con_master_newpg
									set
									vendorid=$1,
									currencyid=$2,
									advancetype=$3,
									advance=$4,
									payment_terms_days=$5,
									billing_at_id=$6,
									delivery_at_id=$7,
									taxes_duties=$8,
									transport_mode=$9,
									transit_insurence=$10,
									packing_forward=$11,
									othercharges=$12,
									rate=$13,
									purchase_type=$14,
									terminal_month=$15,
									terminal_price=$16,
									purchase_price=$17,
									market_price=$18,
									total_price=$19,
									remarks=$20,
									gross_price=$21,
									quantity_mt=$22,
									balance_quantity=$23
									where 
									pono=$24`
			_, err := po.db.Raw(sqlStatementImp1,
				req.SupplierID,
				req.CurrencyID,
				req.AdvanceType,
				req.Advance,
				req.PaymentTermsDays,
				req.POBillTypeID,
				req.PODelTypeID,
				req.TaxDuties,
				req.ModeOfTransport,
				req.TransitInsurance,
				req.PackForward,
				req.OtherCharges,
				NewNullString(req.Rate),
				req.PurchaseType,
				req.TerminalMonth,
				NewNullString(req.DTerminalPrice),
				NewNullString(req.PurchasePriceInr),
				NewNullString(req.MarketPriceInr),
				NewNullString(req.TotalPrice),
				req.Comments,
				NewNullString(req.GrossPrice),
				req.MT_Quantity,
				req.TotalBalQuan,
				req.PoNO).Rows()
			if err != nil {
				logger.Error("unable to Update details to PO", err)
			}

			//-----------Update Domestic Tax Information----------------
			//Delete other charges & taxes and insert freshly
			sqlStatementDelTax := `delete from dbo.pur_gc_po_details_taxes_newpg
									where pono=$1`
			_, errDel := po.db.Raw(sqlStatementDelTax, req.PoNO).Rows()
			if errDel != nil {
				logger.Error("unable to delete tax information")
			}
			//Find last Taxids
			//Find latest TaxID
			var lasttaxidsno, taxidsno int
			sqlStatementTGen1 := `SELECT taxidsno FROM dbo.pur_gc_po_details_taxes_newpg
									where taxidsno is not null
									 order by taxidsno desc limit 1`
			rowsTG1, errTG1 := po.db.Raw(sqlStatementTGen1).Rows()
			if errTG1 != nil {
				logger.Error("unable to find latest tax idsno", errTG1)
			}
			// var po InputAdditionalDetails
			for rowsTG1.Next() {
				errTG1 = rowsTG1.Scan(&lasttaxidsno)
			}
			taxidsno = lasttaxidsno + 1
			req.TaxId = "DTAX-" + strconv.Itoa(taxidsno)
			//Insert Tax info
			sqlStatementDTax1 := `INSERT INTO dbo.pur_gc_po_details_taxes_newpg(
				taxidsno,pono,taxid,itemid, sgst, cgst, igst,
				pack_forward, installation, freight, handling,
				misc, hamali, mandifee, full_tax, insurance)
				  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
			_, errDTax1 := po.db.Raw(sqlStatementDTax1,
				taxidsno,
				req.PoNO,
				req.TaxId,
				req.ItemID,
				NewNullString(req.SGST),
				NewNullString(req.CGST),
				NewNullString(req.IGST),
				NewNullString(req.DPackForward),
				NewNullString(req.DInstallation),
				NewNullString(req.DFreight),
				NewNullString(req.DHandling),
				NewNullString(req.DMisc),
				NewNullString(req.DHamali),
				NewNullString(req.DMandiFee),
				NewNullString(req.DFullTax),
				NewNullString(req.DInsurance)).Rows()

			logger.Error("Domestic Tax Insert Query Executed")
			if errDTax1 != nil {
				logger.Error("unable to insert Dometic tax info", errDTax1)
			}

			//Update new tax id to PO Table

			sqlStatementDTax2 := `update dbo.pur_gc_po_con_master_newpg
									set
									taxid=$1
									where pono=$2`
			_, errDTax2 := po.db.Raw(sqlStatementDTax2,
				req.TaxId,
				req.PoNO).Rows()

			logger.Error("Domestic Tax Update Query Executed")
			if errDTax2 != nil {
				logger.Error("unable to insert Dometic tax info", errDTax1)
			}

		}
		//---PO Item Update---------
		logger.Error("Updating Item Details")
		sqlStatementITU1 := `Update dbo.pur_gc_po_con_master_newpg
				set
				itemid=$1,
				total_quantity=$2
				where
				pono=$3`
		_, errITU1 := po.db.Raw(sqlStatementITU1,
			req.ItemID,
			req.TotalQuantity,
			req.PoNO).Rows()
		logger.Error("Item Details for PO updated")
		if errITU1 != nil {
			logger.Error("unable to update PO Item details", errITU1)
		}

		// 	//-------Delete dispatches for the PO--------------//
		sqlStatementDD1 := `delete from dbo.pur_gc_po_dispatch_master_newpg 
								where 
								pono=$1`
		_, errDD1 := po.db.Raw(sqlStatementDD1, req.PoNO).Rows()
		if errDD1 != nil {
			logger.Error("unable to insert GC Multi dispatch Details", errDD1)
		}
		logger.Error("Dispatches for the PO have been deleted", req.PoNO)
		//----------------Create Fresh Dispatches-------------//
		logger.Error("Entered PO Create Module")
		sqlStatementDT1 := `select detidsno from dbo.pur_gc_po_dispatch_master_newpg
							where detidsno is not null
							order by detidsno desc limit 1`

		rows, err := po.db.Raw(sqlStatementDT1).Rows()

		// var po InputAdditionalDetails
		for rows.Next() {
			err = rows.Scan(&req.LastDetIDSNo)
		}
		logger.Error("Last DETIDSNO from table:", req.LastDetIDSNo)
		req.DetIDSNo = req.LastDetIDSNo + 1
		logger.Error("New DETIDSNO from table:", req.DetIDSNo)
		req.DetID = "GCDIS-" + strconv.Itoa(req.DetIDSNo)
		logger.Error("New DETID from table:", req.DetID)
		logger.Error("Dispatch Type Selected:", req.DispatchType)

		sqlStatementDD2 := `delete from dbo.pur_gc_po_master_documents 
								where 
								poid=$1`
		_, errDD2 := po.db.Raw(sqlStatementDD2, req.PoId).Rows()
		if errDD2 != nil {
			logger.Error("unable to delete Po documents section", errDD1)
		}
		logger.Error("Po documents section is successfully deleted", req.PoId)

		for _, dis := range req.ItemDispatchDetails {
			logger.Error("Loop Entered")
			if req.DispatchType == "Single" {
				req.DispatchCount = "1"
				req.DispatchType = "Single"
			} else {
				req.DispatchType = "Multiple"
			}
			logger.Error("Values of dispatch details are ", dis)
			sqlStatementMD1 := `insert into dbo.pur_gc_po_dispatch_master_newpg(
						pono,
						detid,
						detidsno,
						itemid,
						quantity,
						dispatch_date,
						dispatch_count,
						dispatch_type,
						createdon,
						createdby) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
			_, errMD1 := po.db.Raw(sqlStatementMD1,
				req.PoNO,
				req.DetID,
				req.DetIDSNo,
				req.ItemID,
				dis.DispatchQuantity,
				dis.DispatchDate,
				req.DispatchCount,
				req.DispatchType,
				req.PoDate,
				NewNullString(req.CreatedUserID)).Rows()
			logger.Error("Row Inserted Successfully")
			if req.DocumentsSection != nil && len(req.DocumentsSection) > 0 {
				for i, document := range req.DocumentsSection {
					logger.Error("document in loop", i, document)

					sqlStatement := `select docidsno from dbo.pur_gc_po_master_documents order by docidsno DESC LIMIT 1`
					rows1, err1 := po.db.Raw(sqlStatement).Rows()

					if err1 != nil {
						logger.Error("Unable to get last updated id")
						  
					}

					var lastDoc int
					for rows1.Next() {
						err = rows1.Scan(&lastDoc)
					}

					docIdsno := lastDoc + 1
					docId := "FAC-" + strconv.Itoa(docIdsno)
					sqlStatement1 := `INSERT INTO dbo.pur_gc_po_master_documents (docid, docidsno, poid, dockind, required, dispatchid) VALUES ($1, $2, $3, $4, $5, $6)`
					rows, err = po.db.Raw(sqlStatement1, docId, docIdsno, req.PoId, document.DocKind, document.Required, req.DetID).Rows()
				}
			}
			if req.DispatchType == "Multiple" {
				req.DetIDSNo = req.DetIDSNo + 1
				req.DetID = "GCDIS-" + strconv.Itoa(req.DetIDSNo)
			}

			if errMD1 != nil {
				logger.Error("unable to insert GC Multi dispatch Details", errMD1)
			}
			logger.Error("Inserted details are :", req.ItemDispatchDetails)
		}

		//------------Quote if its a special Coffee type------------------------
		if req.QuotNo != "" {
			sqlStatementQuoteInfo1 := `update dbo.pur_gc_po_con_master_newpg
										set 
										quote_no=$1,
										quote_date=$2,
										quote_price=$3
											where pono=$4`
			defer rows.Close()
			rows, err = po.db.Raw(sqlStatementQuoteInfo1, req.QuotNo, req.QuotDate, req.QuotPrice, req.PoNO).Rows()
			if err != nil {
				logger.Error("unable to insert Quote details to PO", err)
			}
		}

		// Insert Audit Info.
		logger.Error("Entered Audit Module for PO Type")
		// Find created user username
		sqlStatementAUser1 := `SELECT u.userid 
								FROM dbo.users_master_newpg u
								inner join dbo.pur_gc_po_con_master_newpg po on po.createdby=u.userid
								where po.pono=$1`
		rows, err = po.db.Raw(sqlStatementAUser1, req.PoNO).Rows()
		for rows.Next() {
			err = rows.Scan(&req.GCCreatedUserID)
		}
		audit.CreatedUserid = req.GCCreatedUserID
		audit.CreatedDate = req.PoDate
		audit.Description = "PO Modified"
		// sd.InvoiceDate = time.Now().Format("2006-01-02")
		audit.ModifiedDate = time.Now().Format("2006-01-02")
		audit.ModifiedUserid = req.CreatedUserID

		sqlStatementADT := `INSERT INTO dbo.auditlog_pur_gc_master_newpg(
						pono,createdby, created_date, description,modifiedby, modified_date)
						VALUES($1,$2,$3,$4,$5,$6)`
		_, errADT := po.db.Raw(sqlStatementADT,
			req.PoNO,
			audit.CreatedUserid,
			audit.CreatedDate,
			audit.Description,
			audit.ModifiedUserid,
			audit.ModifiedDate).Rows()

		logger.Error("Audit Insert Query Executed")
		if errADT != nil {
			logger.Error("unable to insert Audit Details", errADT)
		}

	}


	
	return nil, nil
}

func (po poRepo) ViewPoDetails(ctx context.Context, req *model.PurchaseOrderDetails) (interface{}, error) {
	logger := logger2.GetLoggerWithContext(ctx)
	var viewPOResp []model.PurchaseOrderDetails
	if req.PoNO != "" {
		logger.Error("Entered PO View Module")
		logger.Error("selected PO NO:", req.PoNO)
		//check if po is import or domestic
		sqlStatementIDC1 := `SELECT posubcat FROM dbo.pur_gc_po_con_master_newpg
						where pono=$1`
		rows, err := po.db.Raw(sqlStatementIDC1, req.PoNO).Rows()
		
		if err != nil {
			logger.Error("Fetching PO Details from DB failed")
			logger.Error(err.Error())
		}
		// defer rows.Close()
		for rows.Next() {
			var cat model.PurchaseOrderDetails
			err = rows.Scan(&cat.POSubCategory)
			logger.Error("Scanned pocat is",cat)
			viewPOResp=append(viewPOResp,cat)
			logger.Error(viewPOResp)
		}
		if req.POSubCategory == "Import" {
			req.SupplierType = "Import"
			sqlStatementPOV1 := `SELECT total_quantity,cid,poid, podate, pocat,vendorid,itemid,
						billing_at_id, delivery_at_id,currencyid,status,dispatchterms, origin,
						poloading, insurance, destination, forwarding,nocontainers,container_type,
						payment_terms,remarks,taxes_duties, transport_mode, transit_insurence, packing_forward,
						othercharges,rate,noofbags,netweight,
						purchase_type, terminal_month, booked_term_rate,booked_differential, fixed_term_rate, fixed_differential,
						purchase_price, market_price, po_margin, total_price,gross_price,fixationdate,quantity_mt
						FROM dbo.pur_gc_po_con_master_newpg
						where pono=$1`
			rows, err := po.db.Raw(sqlStatementPOV1, req.PoNO).Rows()
			logger.Error("PO Master Query Executed")
			if err!= nil {
				logger.Error("Fetching PO Details from DB failed")
				logger.Error(err.Error())
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&req.TotalQuantity, &req.Contract, &req.PoId, &req.PoDate, &req.POCategory,
					&req.SupplierID, &req.ItemID, &req.POBillTypeID, &req.PODelTypeID, &req.CurrencyID, &req.Status,
					&req.IncoTermsID, &req.Origin, &req.PortOfLoad, &req.Insurance, &req.PlaceOfDestination, &req.Forwarding, &req.NoOfContainers,
					&req.ContainerType, &req.PaymentTerms, &req.Comments, &req.TaxDuties, &req.ModeOfTransport, &req.TransitInsurance,
					&req.PackForward, &req.OtherCharges, &req.Rate, &req.NoOfBags, &req.NetWt,
					&req.PurchaseType, &req.TerminalMonth, &req.BookedTerminalRate, &req.BookedDifferential, &req.FixedTerminalRate, &req.FixedDifferential, &req.PurchasePrice, &req.MarketPrice,
					&req.POMargin, &req.TotalPrice, &req.GrossPrice, &req.FixationDate, &req.MTQuantity)

			}	
				
			//Fetch Incoterms details:
			if req.IncoTermsID != "" {
				logger.Error("get incoterms for id :", req.IncoTermsID)
				sqlStatementIT1 := `SELECT incoterms FROM dbo.cms_incoterms_master where incotermsid=$1`
				rows, err := po.db.Raw(sqlStatementIT1, req.IncoTermsID).Rows()
				if err != nil {
					logger.Error("Fetching Incoterms Details from DB failed")

				}
				defer rows.Close()
				for rows.Next() {
					var inc model.PurchaseOrderDetails
					err = rows.Scan(&inc.IncoTerms)
					viewPOResp=append(viewPOResp,inc)
				}
			}

		} else {
			//DOMESTIC PO VIEW
			req.SupplierType = "Domestic"

			sqlStatementDPOV1 := `SELECT poid, podate, pocat, posubcat, 
							vendorid,itemid, billing_at_id, delivery_at_id,
							currencyid,status,
							advancetype, advance, payment_terms_days, 
							taxes_duties, transport_mode, transit_insurence, 
							packing_forward,othercharges,rate,remarks,
							purchase_type,terminal_month,terminal_price,
							purchase_price,market_price,total_price,gross_price,total_quantity,fixationdate
							FROM dbo.pur_gc_po_con_master_newpg
							where pono=$1`
			rows, err := po.db.Raw(sqlStatementDPOV1, req.PoNO).Rows()
			logger.Error("PO Master Query Executed")
			if err != nil {
				logger.Error("Fetching PO Details from DB failed")
				logger.Error(err.Error())
				// return events.APIGatewayProxyResponse{500, headers, nil, errd1.Error(), false}, nil
			}
			// defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&req.PoId, &req.PoDate, &req.POCategory, &req.POSubCategory,
					&req.SupplierID, &req.ItemID, &req.POBillTypeID, &req.PODelTypeID, &req.CurrencyID,
					&req.Status, &req.AdvanceType, &req.Advance,
					&req.PaymentTermsDays, &req.TaxDuties, &req.ModeOfTransport, &req.TransitInsurance,
					&req.DPackForward, &req.OtherCharges, &req.Rate, &req.Comments, &req.PurchaseType,
					&req.TerminalMonth, &req.DTerminalPrice, &req.PurchasePriceInr,
					&req.MarketPriceInr, &req.TotalPrice, &req.GrossPrice, &req.TotalQuantity, &req.FixationDate)
			}
		}
		// ------COMMON to IMPORT && DOMESTIC-------//
		
		// ---------------_Fetch Billing Address Info------------------------
		logger.Error("Entered Billing Module")
		sqlStatementPOVB2 := `SELECT 
						 potypeid,
						 initcap(bdi.potypename),
						 initcap(bdi.potypefullname)||','||initcap(bdi.address) as billingaddress
						 from dbo.pur_po_types bdi
						 where 
						 bdi.potypeid=(select pom.billing_at_id from dbo.pur_gc_po_con_master_newpg pom where pom.pono=$1)`
		rows, err = po.db.Raw(sqlStatementPOVB2, req.PoNO).Rows()
		logger.Error("PO Types Query Executed")
		if err != nil {
			logger.Error("Issue in fetching billing address from DB failed")
			logger.Error(err.Error())
			// return events.APIGatewayProxyResponse{500, headers, nil, errb2.Error(), false}, nil
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&req.POBillTypeID, &req.POBillTypeName, &req.POBillAddress)
			logger.Error(req.POBillAddress)			
		}
		//---------------_Fetch Delivery Address Info------------------------
		logger.Error("Entered PO Delivery Module")
		sqlStatementPOVD2 := `SELECT 
						  initcap(bdi.potypename),
						 initcap(bdi.potypefullname)||','||initcap(bdi.address) as billingaddress
						 from dbo.pur_po_types bdi
						 where 
						 bdi.potypeid=(select pom.delivery_at_id from dbo.pur_gc_po_con_master_newpg pom where pom.pono=$1)`
		rows, err = po.db.Raw(sqlStatementPOVD2, req.PoNO).Rows()
		logger.Error("PO Delivery Address Query Executed")
		if err != nil {
			logger.Error("Fetching PO Delivery Details from DB failed")
			logger.Error(err.Error())			
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&req.PODelTypeName, &req.PODelAddress)
		}
		//-------__Fetch Vendor Information---------------------
		logger.Error("Entered PO Vendor Module")
		sqlStatementPOV3 := `SELECT				
						vm.vendortypeid,
						vm.country,
						initcap(vm.vendorname),
						initcap(vm.address1)||','||initcap(vm.address2)||','||initcap(vm.city)||','||pincode||','||initcap(vm.state)||' -'||SUBSTRING (vm.gstin, 1 , 2)||','||'Phone:'||vm.phone||','||'Mobile:'||vm.mobile||','||'GST NO:'||vm.gstin||','||'PAN NO:'||vm.panno,
						vm.email
						from 
						dbo.pur_vendor_master_newpg vm
						where vm.vendorid=(select pom.vendorid from dbo.pur_gc_po_con_master_newpg pom where pom.pono=$1)`
		
		rows, err  = po.db.Raw(sqlStatementPOV3, req.PoNO).Rows()
		logger.Error("Vendor Details fetch Query Executed")
		if err != nil {
			logger.Error("Fetching Vendor Details from DB failed")
			logger.Error(err.Error())
			//   
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&req.SupplierTypeID, &req.SupplierCountry, &req.SupplierName, &req.SupplierAddress, &req.SupplierEmail)
		}
	
		//-------------Fetch Currencuy Info----------------------------
		logger.Error("Entered Currency Fetch Module")
		sqlStatementPOV4 := `SELECT currencyname,currencycode
							from dbo.project_currency_master 
							where currencyid=$1`
		rows, err  = po.db.Raw(sqlStatementPOV4, req.CurrencyID).Rows()
		logger.Error("Currency Details fetch Query Executed")
		if err != nil {
			logger.Error("Fetching Currency Details from DB failed")
			logger.Error(err.Error())
			// return events.APIGatewayProxyResponse{500, headers, nil, err4.Error(), false}, nil
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&req.CurrencyName, &req.CurrencyCode)
		}
		if req.AdvanceType == "101" {
			req.AdvanceType = "Percentage"
		} else {
			req.AdvanceType = "Amount"
		}
		logger.Error("Currency Name & Code are: ", req.CurrencyName, req.CurrencyCode)

		//----------_Fetch Green Coffee Item Information--------------------
		if req.ItemID != "" {
			logger.Error("Entered GC Item Fetch Module")
			sqlStatementPOV5 := `SELECT im.itemid,initcap(im.itemname),im.cat_type
							from dbo.inv_gc_item_master_newpg im
							where
							im.itemid=$1`
			rows, err  = po.db.Raw(sqlStatementPOV5, req.ItemID).Rows()
			logger.Error("GC Details fetch Query Executed")
			if err != nil {
				logger.Error("Fetching GC Details from DB failed")
				logger.Error(err.Error())
				//   
			}
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&req.ItemID, &req.ItemName, &req.GCCoffeeType)
			}
		}

		// ---------------------Fetch GC Composition Details--------------------------------------//
		logger.Error("The GC Composition for the Item #", req.ItemID)
		sqlStatementPOGC1 := `SELECT density, moisture, browns, blacks, brokenbits, insectedbeans, bleached, husk, sticks, stones, beansretained
						FROM dbo.pur_gc_po_composition_master_newpg where itemid=$1`
		rows, err  = po.db.Raw(sqlStatementPOGC1, req.ItemID).Rows()
		logger.Error("GC Fetch Query Executed")
		if err != nil {
			logger.Error("Fetching GC Composition Details from DB failed")
			logger.Error(err.Error())
			//   
		}

		for rows.Next() {
			err = rows.Scan(&req.Density, &req.Moisture, &req.Browns, &req.Blacks, &req.BrokenBits, &req.InsectedBeans, &req.Bleached, &req.Husk, &req.Sticks,
				&req.Stones, &req.BeansRetained)

		}

		// ---------------------Fetch Multiple Dispatch Info-------------------------------------//
		logger.Error("Fetching Single/Multiple Dispatch Information the Contract #")
		sqlStatementMDInfo1 := `select d.detid,d.dispatch_date,d.quantity, d.dispatch_type,d.dispatch_count,
							m.delivered_quantity, (m.expected_quantity-m.delivered_quantity) as balance_quantity
							from dbo.pur_gc_po_dispatch_master_newpg d
							left join dbo.inv_gc_po_mrin_master_newpg as m on m.detid=d.detid
							where d.pono=$1`
		rows, err  = po.db.Raw(sqlStatementMDInfo1, req.PoNO).Rows()
		
		logger.Error("Multi Dispatch Info Fetch Query Executed")
		if err != nil {
			logger.Error("Multi Dispatch Info Fetch Query failed")
			logger.Error(err.Error())
			//   
		}
		var mid model.ItemDispatch
		
		
		for rows.Next() {
			
			err = rows.Scan(&mid.DispatchID, &mid.DispatchDate, &mid.DispatchQuantity, &req.DispatchType, &req.DispatchCount, &mid.DeliveredQuantity, &mid.BalanceQuantity)
			// itemDisp = append(itemDisp, mid)
			gcMultiDispatch := append(req.ItemDispatchDetails, mid)
			req.ItemDispatchDetails = gcMultiDispatch
		}
		if err != nil {
			
			logger.Error(err.Error())
			//   
		}
		// logger.Error("Multi Dispatch Details:", req.ItemDispatchDetails)

		//---------------Fetch Domestic Tax info for Domestic PO-------------------

		if req.POSubCategory == "Domestic" {
			logger.Error("Selected supplier type Domestic Code:", req.POSubCategory)
			sqlStatementDTax1 := `SELECT sgst, cgst, igst,pack_forward, installation,
							 freight, handling, misc, hamali, mandifee, full_tax,
							  insurance FROM dbo.pur_gc_po_details_taxes_newpg 
							  where pono=$1`
			rows, err  = po.db.Raw(sqlStatementDTax1, req.PoNO).Rows()
			logger.Error("Domestic Tax Info Fetch Query Executed")
			if err != nil {
				logger.Error("Domestic Tax Info Fetch Query failed")
				logger.Error(err.Error())
				// return events.APIGatewayProxyResponse{500, headers, nil, errDTax1.Error(), false}, nil
			}

			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&req.SGST, &req.CGST, &req.IGST, &req.DPackForward, &req.DInstallation, &req.DFreight,
					&req.DHandling, &req.DMisc, &req.DHamali, &req.DMandiFee, &req.DFullTax, &req.DInsurance)
			}	
			if err != nil {
			
				logger.Error(err.Error())
			}
		}
		//----------Quote Info for Speciality Green Coffee Item Information--------------------
		if req.GCCoffeeType != "regular" {
			logger.Error("Entered Quote date & Quote Info Fetch Module for speciaity Coffee")
			sqlStatementSPQ := `SELECT 
							 pom.quote_no,
							 pom.quote_date,
							 pom.quote_price
							 from dbo.pur_gc_po_con_master_newpg pom
							 where pom.pono=$1`
			rows, err  = po.db.Raw(sqlStatementSPQ, req.PoNO).Rows()
			logger.Error("Quote Info Fetch Module for speciaity Coffee Query Executed")
			if err != nil {
				logger.Error("Quote Info Fetch Module for speciaity Coffee from DB failed")
				logger.Error(err.Error())
				// return events.APIGatewayProxyResponse{500, headers, nil, errSPQ.Error(), false}, nil
			}
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&req.QuotNo, &req.QuotDate, &req.QuotPrice)
			}
			if err != nil {
			
				logger.Error(err.Error())
			}
			logger.Error(req.QuotNo, req.QuotDate)
		}
		//------Consolidated Finance Status------------------//
		sqlStatementCFS := `SELECT accpay_status,qc_status,payable_amount
						FROM dbo.pur_gc_po_con_master_newpg
						where pono=$1`
		rows, err  = po.db.Raw(sqlStatementCFS, req.PoNO).Rows()
		logger.Error("Consolidated Finance Status Query Executed")
		if err != nil {
			logger.Error("Fetching Consolidated Finance Status from DB failed")
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&req.QCStatus, &req.APStatus, &req.PayableAmount)
		}
		if err != nil {
			
			logger.Error(err.Error())
		}
		//---------------------Fetch Audit Log Info-------------------------------------//
		logger.Error("Fetching Audit Log Info #")
		sqlStatementAI := `select u.username as createduser, gc.created_date,
			gc.description, v.username as modifieduser, gc.modified_date
   			from dbo.auditlog_pur_gc_master_newpg gc
   			inner join dbo.users_master_newpg u on gc.createdby=u.userid
  			left join dbo.users_master_newpg v on gc.modifiedby=v.userid
   			where gc.pono=$1 order by logid desc limit 1`
		rows, err  = po.db.Raw(sqlStatementAI, req.PoNO).Rows()
		logger.Error("Audit Info Fetch Query Executed")
		if err != nil {
			logger.Error("Audit Info Fetch Query failed")
			logger.Error(err.Error())
			//   
		}

		for rows.Next() {
			var al model.AuditLogGCPO
			err = rows.Scan(&al.CreatedUserid, &al.CreatedDate, &al.Description, &al.ModifiedUserid, &al.ModifiedDate)
			auditDetails := append(req.AuditLogDetails, al)
			req.AuditLogDetails = auditDetails
			logger.Error("added one")
		}
		if err != nil {
			
			logger.Error(err.Error())
		}
		logger.Error("Audit Details:", req.AuditLogDetails)
		return req,err
		
	} else {
		logger.Error("Couldnt find po")
		// return events.APIGatewayProxyResponse{200, headers, nil, string("Couldn't find PO Details"), false}, nil
	}
	// return events.APIGatewayProxyResponse{200, headers, nil, string("success"), false}, nil
	return nil,nil
}