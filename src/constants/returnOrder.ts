export const ROLE_IDS = {
  ACCOUNTING: 2,
  WAREHOUSE: 3
} as const;

export const ROLE_NAMES = {
  ACCOUNTING: 'Accounting',
  WAREHOUSE: 'Warehouse'
} as const;

export const STATUS = {
  RETURN: {
    PENDING: 1,
    BOOKING: 3
  },
  CONFIRM: {
    DRAFT: 1,
    CONFIRMED: 2
  }
} as const;

export const TRANSPORT_TYPES = [
  { label: 'SPX Express', value: 'SPX' },
  { label: 'J&T Express', value: 'JNT' },
  { label: 'Flash', value: 'FLASH' },
  { label: 'NocNoc', value: 'NOCNOC' }
];

export const WAREHOUSES = [
  { label: 'RBN (สินค้าชิ้นใหญ่)', value: 'RBN' },
  { label: 'MMT (สินค้าชิ้นเล็ก)', value: 'MMT' }
];

export const RETURN_STATUS = {
  PENDING: 1,
  DRAFT: 1,
  BOOKING: 3,
  CONFIRM: 2
};
