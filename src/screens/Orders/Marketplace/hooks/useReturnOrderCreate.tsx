// src/screens/Orders/Marketplace/hooks/useReturnOrderCreate.ts
import { useCallback } from 'react';
import { message, Modal, notification } from 'antd';
import { FormInstance } from 'antd/lib/form';
import { OrderData } from '../../../../redux/orders/types/state';

// กำหนดค่า default
const DEFAULT_CHANNEL_ID = 1;
const DEFAULT_CUSTOMER_ID = "Customer-Test";

/**
 * Custom hook สำหรับจัดการการสร้าง Return Order
 */
export const useReturnOrderCreate = (
  orderData: OrderData | null,
  form: FormInstance,
  returnItems: { [key: string]: number },
  createReturnOrder: (payload: any) => void,
  setHasCreatedOrder: (value: boolean) => void,
  setStepLoading: (value: boolean) => void,
  calculateTotalAmount: (items: any[]) => number
) => {
  // จัดการการสร้างคำสั่งคืนสินค้า
  const handleCreateReturnOrder = useCallback(() => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      // ตรวจสอบค่าในฟอร์ม
      form.validateFields().then(formValues => {
        // กรองรายการสินค้าที่มีจำนวนคืนมากกว่า 0
        const returnItemsList = orderData.lines
          .filter((item) => returnItems[item.sku] > 0)
          .map((item) => ({
            orderNo: orderData.head.orderNo,
            sku: item.sku,
            itemName: item.itemName,
            qty: Math.abs(item.qty),
            returnQty: returnItems[item.sku],
            price: Math.abs(item.price),
            trackingNo: formValues.trackingNo,
          }));

        if (returnItemsList.length === 0) {
          message.error("กรุณาระบุจำนวนสินค้าที่ต้องการคืน");
          return;
        }

        const warehouseID = Number(formValues.warehouseFrom);
        if (isNaN(warehouseID)) {
          message.error("กรุณาเลือกคลังสินค้าที่ถูกต้อง");
          return;
        }

        // สร้าง payload
        const createReturnPayload = {
          orderNo: orderData.head.orderNo,
          soNo: orderData.head.soNo,
          channelID: DEFAULT_CHANNEL_ID,
          customerID: DEFAULT_CUSTOMER_ID,
          reason: formValues.reason || "Return",
          warehouseID: warehouseID,
          returnDate: formValues.returnDate.toISOString(),
          trackingNo: formValues.trackingNo,
          logistic: formValues.transportType,
          soStatus: orderData.head.salesStatus,
          mkpStatus: orderData.head.mkpStatus,
          items: returnItemsList,
        };

        setStepLoading(true);

        // แสดง Modal ยืนยันการสร้างคำสั่ง
        Modal.confirm({
          title: "ยืนยันการสร้างคำสั่งคืนสินค้า",
          content: (
            <div>
              <p>Order No: {orderData.head.orderNo}</p>
              <p>SO No: {orderData.head.soNo}</p>
              <p>จำนวนรายการที่จะคืน: {returnItemsList.length} รายการ</p>
              <p>Tracking No: {formValues.trackingNo}</p>
              <p>ขนส่ง: {formValues.transportType}</p>
              <p>วันที่คืน: {formValues.returnDate.format("DD/MM/YYYY HH:mm")}</p>
              <p>มูลค่ารวม: ฿{calculateTotalAmount(returnItemsList).toLocaleString()}</p>
            </div>
          ),
          okText: "สร้างคำสั่งคืนสินค้า",
          cancelText: "ยกเลิก",
          onOk: () => {
            try {
              console.log(`[ReturnOrder] Creating return order: ${JSON.stringify(createReturnPayload)}`);
              setHasCreatedOrder(true);
              createReturnOrder(createReturnPayload);
            } catch (error: any) {
              notification.error({
                message: "เกิดข้อผิดพลาด",
                description: error.message || "ไม่สามารถสร้างคำสั่งคืนสินค้าได้",
              });
              setHasCreatedOrder(false);
            } finally {
              setStepLoading(false);
            }
          },
          onCancel: () => {
            setStepLoading(false);
          }
        });
      });
    } catch (error: any) {
      console.error('[ReturnOrder] Create return order error:', error);
      setStepLoading(false);
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message || "ไม่สามารถสร้างคำสั่งคืนสินค้าได้",
      });
    }
  }, [orderData, form, returnItems, createReturnOrder, setHasCreatedOrder, setStepLoading, calculateTotalAmount]);

  return {
    handleCreateReturnOrder
  };
};