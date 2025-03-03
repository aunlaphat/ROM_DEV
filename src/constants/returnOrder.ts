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
