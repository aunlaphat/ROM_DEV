package api

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) Excels(apiRouter *chi.Mux) {
	apiRouter.Route("/excel", func(r chi.Router) {
		// r.Post("/check-items-import", app.CheckItemsExcel)
	})
}

// // Check items import excel.
// // @Summary Validate imported Excel items.
// // @Description This endpoint checks and validates items from an Excel import and categorizes them as found or not found.
// // @ID excel-check-items-import
// // @Tags Excels
// // @Accept json
// // @Produce json
// // @Param CheckItems body service.CheckItemExcel true "Check items import excel"
// // @Success 200 {object} Response{result=[]service.CheckItemResponse} "Found and Not Found items"
// // @Failure 400 {object} Response "Bad Request - Invalid input"
// // @Failure 404 {object} Response "Not Found - Items not found"
// // @Failure 500 {object} Response "Internal Server Error"
// // @Router /excel/check-items-import [post]
// func (app *Application) CheckItemsExcel(w http.ResponseWriter, r *http.Request) {

// 	var req service.CheckItemExcel
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}
// 	processedItems, err := processItems(req.Items)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	// Wrap the processed items back into a CheckItemExcel struct
// 	processedInput := service.CheckItemExcel{
// 		Items: processedItems,
// 	}

// 	foundItems, notFoundItems, err := app.Service.Excel.CheckItemExcel(processedInput)
// 	if err != nil {
// 		HandleError(w, err)
// 		return
// 	}
// 	res := map[string]interface{}{
// 		"foundItems":    foundItems,
// 		"notFoundItems": notFoundItems,
// 	}

// 	handleResponse(w, true, "Items processed successfully", res, http.StatusOK)
// }

// func processItems(items []service.CheckItem) ([]service.CheckItem, error) {
// 	processedItems := make([]service.CheckItem, 0, len(items))

// 	for _, item := range items {
// 		// Normalize and clean the SKU
// 		item.Sku = handleSpecialSyntax(item.Sku)
// 		// Append the item to the processed list
// 		processedItems = append(processedItems, item)
// 	}

// 	return processedItems, nil
// }

// func handleSpecialSyntax(sku string) string {
// 	// Trim whitespace
// 	sku = strings.TrimSpace(sku)
// 	// Normalize case (optional, depending on requirements)
// 	sku = strings.ToUpper(sku)
// 	// Remove invalid characters (e.g., replacing non-alphanumeric characters)
// 	sku = removeInvalidCharacters(sku, "-_")
// 	return sku
// }

// func removeInvalidCharacters(sku string, allowedChars string) string {
// 	validSKU := strings.Builder{}
// 	for _, char := range sku {
// 		if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || strings.ContainsRune(allowedChars, char) {
// 			validSKU.WriteRune(char)
// 		}
// 	}
// 	return validSKU.String()
// }
