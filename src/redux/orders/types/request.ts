// src/redux/orders/types/requests.ts

// Request Types (สอดคล้องกับ GoLang request types)

// SearchOrder - ค้นหา order
export interface SearchOrderRequest {
    soNo?: string;
    orderNo?: string;
  }
  
  // ReturnOrderItem - รายการสินค้าที่ต้องการคืน
  export interface ReturnOrderItemRequest {
    orderNo: string;
    sku: string;
    itemName: string;
    qty: number;
    returnQty: number;
    price: number;
    createBy?: string;
    trackingNo?: string;
    alterSKU?: string;
  }
  
  // CreateBeforeReturnOrder - สร้าง return order
  export interface CreateReturnOrderRequest {
    orderNo: string;
    soNo: string;
    channelID: number;
    customerID: string;
    reason: string;
    soStatus?: string;
    mkpStatus?: string;
    warehouseID: number;
    returnDate: string;
    trackingNo: string;
    logistic: string;
    items: ReturnOrderItemRequest[];
  }
  
  // UpdateSrNo - อัพเดท SR number
  export interface UpdateSrRequest {
    orderNo: string;
    srNo: string;
  }
  
  // UpdateStatus - อัพเดทสถานะ
  export interface UpdateStatusRequest {
    orderNo: string;
    roleID: number;
    userID: string;
  }
  
  // CancelOrder - ยกเลิก order
  export interface CancelOrderRequest {
    refID: string;
    sourceTable: string;
    cancelReason: string;
  }