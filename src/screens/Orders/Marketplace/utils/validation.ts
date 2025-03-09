// src/screens/Orders/Marketplace/utils/validation.ts
import { OrderData } from '../../../../redux/orders/types/state';
import { FormInstance } from 'antd';

/**
 * ตรวจสอบความถูกต้องในการเปลี่ยนขั้นตอน
 */
export const validateStepTransition = (
  fromStep: string, 
  toStep: string, 
  orderData: OrderData | null, 
  returnOrder: any
): boolean => {
  switch (toStep) {
    case 'create':
      return !!orderData;
    case 'sr':
      return !!returnOrder;
    case 'preview':
      return fromStep === 'sr' && !!orderData?.head.srNo;
    case 'confirm':
      return fromStep === 'preview' && !!orderData?.head.srNo;
    default:
      return true;
  }
};

/**
 * ตรวจสอบการ disable ปุ่ม Create Return Order
 */
export const isCreateReturnOrderDisabled = (
  orderData: OrderData | null, 
  returnItems: { [key: string]: number },
  form: FormInstance,
  loading: boolean,
  stepLoading: boolean
): boolean => {
  // 1. ตรวจสอบว่ามีข้อมูล Order หรือไม่
  if (!orderData?.head?.orderNo) return true;

  // 2. ตรวจสอบว่ามีการเลือกสินค้าที่จะคืนหรือไม่
  const hasSelectedItems = Object.values(returnItems).some(qty => qty > 0);
  if (!hasSelectedItems) return true;

  // 3. ตรวจสอบว่ากรอกข้อมูลจำเป็นครบถ้วนหรือไม่
  const formValues = form.getFieldsValue();
  const requiredFields = [
    'warehouseFrom',
    'returnDate',
    'trackingNo',
    'transportType'
  ];
  
  const hasAllRequiredFields = requiredFields.every(field => {
    const value = formValues[field];
    return value !== undefined && value !== null && value !== '';
  });

  // 4. ตรวจสอบว่ามี SR Number แล้วหรือไม่
  if (orderData.head.srNo) return true;

  // 5. ตรวจสอบสถานะ loading
  if (loading || stepLoading) return true;

  // คืนค่า false ถ้าผ่านทุกเงื่อนไข (สามารถกดปุ่มได้)
  return !(hasSelectedItems && hasAllRequiredFields);
};

/**
 * ตรวจสอบการ disable ปุ่ม Create SR
 */
export const isCreateSRDisabled = (
  returnOrder: any, 
  orderData: OrderData | null, 
  loading: boolean, 
  stepLoading: boolean
): boolean => {
  // 1. ตรวจสอบว่ามี returnOrder หรือไม่
  if (!returnOrder) return true;

  // 2. ตรวจสอบว่ามี SR Number แล้วหรือยัง
  if (orderData?.head.srNo) return true;

  // 3. ตรวจสอบสถานะ loading
  if (loading || stepLoading) return true;

  return false;
};

/**
 * รับชื่อบทบาทตาม roleID
 */
export const getRoleName = (roleID: number): string => {
  switch (roleID) {
    case 2: return 'Accounting';
    case 3: return 'Warehouse';
    default: return 'Staff';
  }
};