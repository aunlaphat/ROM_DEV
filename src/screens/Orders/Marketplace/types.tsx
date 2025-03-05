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
  head: {
    orderNo: string;
    soNo: string;
    srNo: string | null;
    salesStatus: string;
    mkpStatus: string;
    locationTo: string;
  };
  lines: OrderLineItem[];
}

export interface OrderLineItem {
  sku: string;
  itemName: string;
  qty: number;
  price: number;
  warehouse: string;
}
