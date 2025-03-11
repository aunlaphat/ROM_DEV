// src/screens/Orders/Marketplace/hooks/useReturnItemsManager.ts
import { useState, useCallback, useEffect, useRef } from 'react';
import { OrderData } from '../../../../redux/orders/types/state';

/**
 * Custom hook สำหรับจัดการข้อมูลสินค้าที่ต้องการคืน
 */
export const useReturnItemsManager = (orderData: OrderData | null) => {
  // State สำหรับเก็บจำนวนสินค้าที่จะคืน (format: { sku: quantity })
  const [returnItems, setReturnItems] = useState<{ [key: string]: number }>({});
  // ref สำหรับเก็บค่า orderData.orderNo ล่าสุด เพื่อตรวจสอบการเปลี่ยนแปลง
  const lastOrderNoRef = useRef<string | null>(null);

  // เริ่มต้นข้อมูลจำนวนสินค้าที่จะคืน
  const initializeReturnItems = useCallback(() => {
    if (!orderData?.lines) return;

    // เก็บค่า orderNo ปัจจุบัน
    const currentOrderNo = orderData.head.orderNo;
    
    // ถ้า orderNo เปลี่ยนหรือเป็น null ให้รีเซ็ต returnItems
    if (currentOrderNo !== lastOrderNoRef.current) {
      console.log(`[ReturnItems] Initializing for new order: ${currentOrderNo}`);
      
      const initialQty = orderData.lines.reduce(
        (acc, item) => ({
          ...acc,
          [item.sku]: 0,
        }),
        {}
      );
      
      setReturnItems(initialQty);
      lastOrderNoRef.current = currentOrderNo;
    } else {
      console.log('[ReturnItems] Same order, keeping existing return items');
    }
  }, [orderData]);

  // เริ่มต้นข้อมูลสินค้าที่จะคืนเมื่อได้รับข้อมูล order
  useEffect(() => {
    if (orderData?.lines) {
      // เรียกฟังก์ชัน initialize เมื่อมีข้อมูล lines
      initializeReturnItems();
      
      // Log สถานะปัจจุบันของ returnItems
      console.log('[ReturnItems] Current state:', {
        orderNo: orderData.head.orderNo,
        returnItems: returnItems,
        itemCount: Object.keys(returnItems).length,
        hasReturnItems: Object.values(returnItems).some(qty => qty > 0)
      });
    }
  }, [orderData, initializeReturnItems, returnItems]);

  // ดึงจำนวนสินค้าที่จะคืน
  const getReturnQty = useCallback((sku: string): number => {
    const qty = returnItems[sku];
    // ป้องกันการคืนค่า undefined หรือ NaN
    return (qty !== undefined && !isNaN(qty)) ? qty : 0;
  }, [returnItems]);

  // อัพเดตจำนวนสินค้าที่จะคืน
  const updateReturnQty = useCallback((sku: string, change: number) => {
    if (!orderData?.lines) return;

    setReturnItems((prev) => {
      const currentQty = prev[sku] || 0;
      const item = orderData.lines.find((item) => item.sku === sku);

      if (!item) return prev;

      const originalQty = Math.abs(item.qty);
      const newQty = Math.max(0, Math.min(originalQty, currentQty + change));
      
      console.log(`[ReturnItems] Updating ${sku}: ${currentQty} -> ${newQty} (change: ${change})`);
      
      return {
        ...prev,
        [sku]: newQty,
      };
    });
  }, [orderData]);

  // คำนวณมูลค่ารวมของสินค้าที่คืน
  const calculateTotalAmount = useCallback((items: any[]): number => {
    return items.reduce((sum, item) => {
      const returnQty = item.returnQty || 0;
      const price = Math.abs(item.price) || 0;
      return sum + (price * returnQty);
    }, 0);
  }, []);

  return {
    returnItems,
    setReturnItems,
    getReturnQty,
    updateReturnQty,
    initializeReturnItems,
    calculateTotalAmount
  };
};