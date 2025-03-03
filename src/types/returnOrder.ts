export interface SearchOrderRequest {
  searchTerm: string; // สำหรับรับค่า OrderNo หรือ SoNo
}

export interface OrderHeadDetail {
  orderNo: string;
  soNo: string;
  channelId: number;
  customerId: string;
  orderDate: string;
  salesStatus: string;
  mkpStatus: string;
  srNo: string | null;
}

export interface OrderLineDetail {
  orderNo: string;
  sku: string;
  itemName: string;
  qty: number;
  price: number;
  warehouseId: number;
}

export interface SearchOrderResponse {
  success: boolean;
  message?: string;
  data?: {
    head: OrderHeadDetail;
    lines: OrderLineDetail[];
  }
}
