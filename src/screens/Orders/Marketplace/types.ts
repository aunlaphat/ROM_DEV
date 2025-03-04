export interface ReturnOrderState {
  orderData: OrderData | null;
  orderLines: OrderLineItem[];
  loading: boolean;
  srCreated: boolean;
  error: string | null;
  currentStep: 'search' | 'create' | 'confirm';
}

export interface OrderData {
  orderNo: string;
  soNo: string;
  srNo: string | null;
  isCNCreated: boolean;
  isEdited: boolean;
  // ...other fields
}

export interface OrderLineItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
  warehouse: string;
}
