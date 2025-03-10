// src/screens/Orders/Marketplace/hooks/useReturnOrderConfirm.ts
import { useCallback } from 'react';
import { Modal, message } from 'antd';
import { OrderData } from '../../../../redux/orders/types/state';
import { getRoleName } from '../utils/validation';
import { useNavigate } from 'react-router-dom';

/**
 * Custom hook สำหรับจัดการการยืนยัน Return Order
 */
export const useReturnOrderConfirm = (
  orderData: OrderData | null,
  auth: any,
  updateStatus: (payload: any) => void,
  setHasConfirmedOrder: (value: boolean) => void
) => {
  // จัดการการยืนยันคำสั่งคืนสินค้า
  const handleConfirm = useCallback(() => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      if (!auth.userID) {
        message.error("ไม่พบข้อมูลผู้ใช้งาน กรุณาเข้าสู่ระบบใหม่");
        return;
      }

      // แสดง Modal ยืนยันคำสั่งคืนสินค้า
      Modal.confirm({
        title: "ยืนยันคำสั่งคืนสินค้า",
        content: (
          <div>
            <p>คุณต้องการยืนยันคำสั่งคืนสินค้าใช่หรือไม่?</p>
            <p>Order No: {orderData.head.orderNo}</p>
            <p>SR No: {orderData.head.srNo}</p>
            <p style={{ color: '#1890ff' }}>
              หมายเหตุ: สถานะจะถูกอัพเดตตามสิทธิ์การใช้งานของคุณ ({getRoleName(auth.roleID)})
            </p>
          </div>
        ),
        okText: "ยืนยัน",
        cancelText: "ยกเลิก",
        onOk: () => {
          const confirmPayload = {
            orderNo: orderData.head.orderNo,
            roleID: auth.roleID || 1,
            userID: auth.userID,
          };

          console.log(`[ReturnOrder] Confirming return order: ${JSON.stringify(confirmPayload)}`);
          setHasConfirmedOrder(true);
          
          // แสดง loading indicator ระหว่างอัพเดตสถานะ
          message.loading({
            content: 'กำลังอัพเดตสถานะ...',
            key: 'confirmStatus',
            duration: 0
          });
          
          // เรียก API สำหรับอัพเดตสถานะ
          updateStatus(confirmPayload);
        }
      });
    } catch (error: any) {
      console.error('[ReturnOrder] Confirm return order error:', error);
      setHasConfirmedOrder(false);
      message.error({
        content: error.message || "ไม่สามารถยืนยันคำสั่งคืนสินค้าได้",
        key: 'confirmStatus'
      });
    }
  }, [orderData, auth, updateStatus, setHasConfirmedOrder]);

  return {
    handleConfirm
  };
};