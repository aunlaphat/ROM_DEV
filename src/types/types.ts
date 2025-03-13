import { QrReader, QrReaderProps } from 'react-qr-reader';

export const TRANSPORT_TYPES = [
  { value: 'SPX', label: 'SPX Express' },
  { value: 'JNT', label: 'J&T Express' },
  { value: 'DHL', label: 'DHL Express' },
  { value: 'FLASH', label: 'Flash' },
  { value: 'THAIPOST', label: 'Thai Post' },
  { value: 'NOCNOC', label: 'NocNoc' },
  { value: 'OTHER', label: 'อื่นๆ' },
];

export interface Address {
    provicesTH: string;
    provinceCode: number;
    districtTH: string;
    districtCode: number;
    subdistrictTH: string;
    subdistrictCode: number;
    zipCode: number;
  }
  
  export interface Customer {
    Key: number;
    customerID: string;
    customerName: string;
    address: string;
    taxID: string;
  }
  
  export interface DataItem {
    key: number;
    SKU: string;
    Name: string;
    QTY: number;
    ReturnQTY: number;
    PricePerUnit: number;
    Price: number;
  }

  export interface DataItemBlind {
    key: number; 
    SKU: string;
    Name: string;
    QTY: number;
}
  
  export interface Product {
    Key: string;
    sku: string;
    nameAlias: string;
    size: string;
  }

  export interface Order {
    Order: string;
    SO_INV: string;
    Customer: string;
    SR: string;
    Transport: string;
    ReturnTracking: string;
    Channel: string;
    Date_Create: string;
    Warehouse: string;
    data: SKUData[];  // ใช้ SKUData ที่มีการกำหนด Type ที่ถูกต้อง
    codeR?: string;
    nameR?: string;
  }

  export interface SKUData {
    OrderNo: string;
    SKU: string;
    Name: string;
    QTY: number;
    Price: string;
    Action: string;
    Type: 'system' | 'addon';  // เพิ่ม Type เพื่อระบุว่ามาจากในระบบหรือเป็น addon
  }

  export interface SelectedRecord {
    data: SKUData[];
  }
  
  export interface OrderDetail {
      orderNo: string;
      soNo: string;
      customerId: string;
      srNo: string;
      trackingNo: string;
      // logistic: string;
      channelName: string;
      createDate: string;
      warehouseName: string;
      data: OrderLine[]; 
  }
  
  export interface OrderLine {
      sku: string;
      itemName: string;
      qty: number;
      price: string;
      Type: 'system' | 'addon';
  }

  export interface CustomQrReaderProps extends QrReaderProps {
      onScan: (result: string | null) => void;
      onError: (error: any) => void;
  }
  
  export interface ReceiptOrder {
      orderNo: string;
      trackingNo: string;
      data: ReceiptOrderLine[];
  }
  
  export interface ReceiptOrderLine {
      key: string;
      sku: string;
      itemName: string;
      qty: number;
      receivedQty: number;
      price: string;
      image: string | null;
      filePath: string;
  }

  export {};