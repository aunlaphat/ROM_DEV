// src/redux/draftConfirm/types/response.ts

// Response Types (สอดคล้องกับ GoLang response types)

// OrderItem - รายการสินค้าใน Order
export interface OrderItem {
    orderNo: string;
    sku: string;
    itemName: string;
    qty: number;
    returnQty?: number;
    price: number;
    createBy?: string;
    createDate?: string;
    type?: 'system' | 'addon'; // Type สำหรับ UI เพื่อแยกประเภทรายการที่มาจากระบบหรือเพิ่มใหม่
  }
  
  // Order - รายการคำสั่ง
  export interface Order {
    orderNo: string;
    soNo: string;
    srNo: string | null;
    customerId: string;
    trackingNo: string;
    logistic: string;
    channelId: number;
    createDate: string;
    warehouseId: number;
    items?: OrderItem[];
  }
  
  // DraftConfirmResponse - รายละเอียดคำสั่งแบบ Draft/Confirm
  export interface DraftConfirmResponse {
    orderNo: string;
    soNo: string;
    srNo: string;
    items: OrderItem[];
  }
  
  // CodeR - รายการ CodeR
  export interface CodeR {
    sku: string;
    nameAlias: string;
  }
  
  // AddItemResponse - ผลลัพธ์การเพิ่มรายการสินค้า
  export interface AddItemResponse {
    orderNo: string;
    sku: string;
    itemName: string;
    qty: number;
    returnQty: number;
    price: number;
    createBy: string;
    createDate: string;
  }
  
  // UpdateOrderStatusResponse - ผลลัพธ์การอัพเดทสถานะคำสั่ง
  export interface UpdateOrderStatusResponse {
    orderNo: string;
    statusReturnID: number;
    statusConfID: number;
    confirmBy: string;
    confirmDate: string;
  }