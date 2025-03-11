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
    Price: number;
  }
  
  export interface Product {
    Key: string;
    sku: string;
    nameAlias: string;
    size: string;
  }

  export {};