package utils

import (
	res "boilerplate-backend-go/dto/response"
	"fmt"
	"time"
)

func PrintOrderDetails(order *res.BeforeReturnOrderResponse) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", order.SrNo)
	fmt.Printf("📡 ChannelID: %d\n", order.ChannelID)
	fmt.Printf("🎐 Reason: %s\n", order.Reason)
	fmt.Printf("👤 CustomerID: %s\n", order.CustomerID)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", order.Logistic)
	fmt.Printf("🏢 WarehouseID: %d\n", order.WarehouseID)
	fmt.Printf("📄 SoStatus: %v\n", order.SoStatus)
	fmt.Printf("📊 MkpStatus: %v\n", order.MkpStatus)
	fmt.Printf("📅 ReturnDate: %v\n", order.ReturnDate)
	fmt.Printf("🪃  StatusReturnID: %d\n", order.StatusReturnID)
	fmt.Printf("✅ StatusConfID: %d\n", order.StatusConfID)
	fmt.Printf("👤 ConfirmBy: %v\n", order.ConfirmBy)
	fmt.Printf("👤 CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
	fmt.Printf("👤 UpdateBy: %v\n", order.UpdateBy)
	fmt.Printf("📅 UpdateDate: %v\n", order.UpdateDate)
	fmt.Printf("❌ CancelID: %v\n", order.CancelID)
}

func PrintOrderLineDetails(line *res.BeforeReturnOrderLineResponse) {
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("📲 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %v\n", line.CreateDate)
}

func PrintSaleOrderDetails(order *res.SaleOrderResponse) {
	fmt.Printf("🧾 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("📊 StatusMKP: %s\n", order.StatusMKP)
	fmt.Printf("📄 SalesStatus: %s\n", order.SalesStatus)
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate)
}

func PrintSaleOrderLineDetails(line *res.SaleOrderLineResponse) {
	fmt.Printf("🛒 SoNo: %s\n", line.SoNo)
	fmt.Printf("🧾 OrderNo: %s\n", line.OrderNo)
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
}

// ************************ Draft & Confirm ************************ //

// Draft & Confirm Head
func PrintDraftConfirmOrderDetails(draft *res.ListDraftConfirmOrdersResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("👤 Customer: %s\n", draft.CustomerID)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
	fmt.Printf("📦 TrackingNo: %s\n", draft.TrackingNo)
	fmt.Printf("🚚 Logistic: %s\n", draft.Logistic)
	fmt.Printf("📡 Channel: %d\n", draft.ChannelID)
	fmt.Printf("📅 CreateDate: %v\n", draft.CreateDate)
	fmt.Printf("🏢 Warehouse: %d\n", draft.WarehouseID)
}

// Modal Edit Draft & Modal Show Confirm
func PrintDraftOrderDetails(draft *res.DraftHeadResponse) {
	fmt.Printf("📦 OrderNo: %s\n", draft.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", draft.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", draft.SrNo)
}

func PrintDraftOrderLineDetails(draft *res.DraftLineResponse) {
	fmt.Printf("🔖 SKU: %s\n", draft.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", draft.ItemName)
	fmt.Printf("📱 QTY: %d\n", draft.QTY)
	fmt.Printf("💲 Price: %.2f\n", draft.Price)
}

// ************************ Return Order ************************ //

func PrintReturnOrderDetails(order *res.ReturnOrder) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", convertString(order.SrNo)) 
	fmt.Printf("📦 TrackingNo: %s\n", convertString(order.TrackingNo)) 
	fmt.Printf("🌍 PlatfID: %s\n", convertInt(order.PlatfID)) 
	fmt.Printf("📲 ChannelID: %s\n", convertInt(order.ChannelID)) 
	fmt.Printf("🔄 OptStatusID: %s\n", convertInt(order.OptStatusID)) 
	fmt.Printf("📊 AxStatusID: %s\n", convertInt(order.AxStatusID)) 
	fmt.Printf("🛍️  PlatfStatusID: %s\n", convertInt(order.PlatfStatusID)) 
	fmt.Printf("🗨️  Reason: %s\n", convertString(order.Reason)) 
	fmt.Printf("✔️  StatusCheckID: %s\n", convertInt(order.StatusCheckID)) 
	fmt.Printf("🕵️  CheckBy: %s\n", convertString(order.CheckBy)) 
	fmt.Printf("📅 CheckDate: %s\n", convertDate(order.CheckDate)) 
	fmt.Printf("🕵️  UpdateBy: %s\n", convertString(order.UpdateBy)) 
	fmt.Printf("📅 UpdateDate: %s\n", convertDate(order.UpdateDate)) 
	fmt.Printf("🕵️  CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %s\n", order.CreateDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("❌ CancelID: %s\n", convertInt(order.CancelID)) 
}

func PrintReturnOrderLineDetails(line *res.ReturnOrderLine) {
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	// fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("📲 ReturnQTY: %d\n", line.ReturnQTY)
	fmt.Printf("📲 ActualQTY: %d\n", line.ActualQTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
	fmt.Printf("📦 TrackingNo: %s\n", line.TrackingNo)
	fmt.Printf("📅 CreateDate: %s\n", line.CreateDate.Format("2006-01-02 15:04:05"))
}

func PrintCreateReturnOrder(order *res.CreateReturnOrder) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", order.SrNo)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo)
	fmt.Printf("🌍 PlatfID: %d\n", order.PlatfID) 
	fmt.Printf("📲 ChannelID: %d\n", order.ChannelID) 
	fmt.Printf("🔄 OptStatusID: %d\n", order.OptStatusID) 
	fmt.Printf("📊 AxStatusID: %d\n", order.AxStatusID) 
	fmt.Printf("🛍️  PlatfStatusID: %d\n", order.PlatfStatusID) 
	fmt.Printf("🗨️  Reason: %s\n", order.Reason)
	fmt.Printf("✔️  StatusCheckID: %d\n", order.StatusCheckID) 
	fmt.Printf("🕵️  CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %s\n", order.CreateDate.Format("2006-01-02 15:04:05"))
}

func PrintUpdateReturnOrder(order *res.UpdateReturnOrder) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", convertString(order.SrNo)) 
	fmt.Printf("📦 TrackingNo: %s\n", convertString(order.TrackingNo)) 
	fmt.Printf("🌍 PlatfID: %s\n", convertInt(order.PlatfID)) 
	fmt.Printf("📲 ChannelID: %s\n", convertInt(order.ChannelID)) 
	fmt.Printf("🔄 OptStatusID: %s\n", convertInt(order.OptStatusID)) 
	fmt.Printf("📊 AxStatusID: %s\n", convertInt(order.AxStatusID)) 
	fmt.Printf("🛍️  PlatfStatusID: %s\n", convertInt(order.PlatfStatusID)) 
	fmt.Printf("🗨️  Reason: %s\n", convertString(order.Reason)) 
	fmt.Printf("✔️  StatusCheckID: %s\n", convertInt(order.StatusCheckID)) 
	fmt.Printf("🕵️  CheckBy: %s\n", convertString(order.CheckBy)) 
	fmt.Printf("📅 CheckDate: %s\n", convertDate(order.CheckDate)) 
	fmt.Printf("🕵️  UpdateBy: %s\n", convertString(order.UpdateBy)) 
	fmt.Printf("📅 UpdateDate: %s\n", convertDate(order.UpdateDate)) 
	fmt.Printf("❌ CancelID: %s\n", convertInt(order.CancelID)) 
}

func PrintImportOrderDetails(order *res.ImportOrderResponse) {
	fmt.Printf("🧾 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("📦 TrackingNo: %s\n", order.TrackingNo) 
	fmt.Printf("📅 CreateDate: %v\n", order.CreateDate.Format("2006-01-02 15:04:05"))
}

func PrintImportOrderLineDetails(line *res.ImportOrderLineResponse) {
	fmt.Printf("🔖 SKU: %s\n", line.SKU)
	fmt.Printf("🏷️  ItemName: %s\n", line.ItemName)
	fmt.Printf("📱 QTY: %d\n", line.QTY)
	fmt.Printf("💲 Price: %.2f ฿\n", line.Price)
}

func PrintDraftTradeOrder (order *res.DraftTradeDetail) {
	fmt.Printf("📦 OrderNo: %s\n", order.OrderNo)
	fmt.Printf("🛒 SoNo: %s\n", order.SoNo)
	fmt.Printf("🔄 SrNo: %s\n", convertString(order.SrNo)) 
	fmt.Printf("📦 TrackingNo: %s\n", convertString(order.TrackingNo)) 
	fmt.Printf("📲 ChannelID: %s\n", convertInt(order.ChannelID)) 
	fmt.Printf("🗨️  Reason: %s\n", convertString(order.Reason)) 
	fmt.Printf("✔️  StatusCheckID: %d\n", order.StatusCheckID) 
	fmt.Printf("🕵️  CreateBy: %s\n", order.CreateBy)
	fmt.Printf("📅 CreateDate: %s\n", order.CreateDate.Format("2006-01-02 15:04:05"))
}

// ฟังก์ชันสำหรับแปลงค่า null
func convertString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func convertInt(i *int) string {
	if i != nil {
		return fmt.Sprintf("%d", *i)
	}
	return ""
}

func convertDate(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}
