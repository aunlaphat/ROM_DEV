// src/screens/Orders/Marketplace/hooks/useReturnItemsManager.ts
import { useState, useCallback, useEffect } from 'react';
import { OrderData } from '../../../../redux/orders/types/state';

/**
 * Custom hook สำหรับจัดการข้อมูลสินค้าที่ต้องการคืน
 */
export const useReturnItemsManager = (orderData: OrderData | null) => {
  // State สำหรับเก็บจำนวนสินค้าที่จะคืน (format: { sku: quantity })
  const [returnItems, setReturnItems] = useState<{ [key: string]: number }>({});

  // เริ่มต้นข้อมูลจำนวนสินค้าที่จะคืน
  const initializeReturnItems = useCallback(() => {
    if (!orderData?.lines) return;
    
    const initialQty = orderData.lines.reduce(
      (acc, item) => ({
        ...acc,
        [item.sku]: 0,
      }),
      {}
    );
    setReturnItems(initialQty);
  }, [orderData]);

  // เริ่มต้นข้อมูลสินค้าที่จะคืนเมื่อได้รับข้อมูล order
  useEffect(() => {
    if (orderData?.lines) {
      initializeReturnItems();
    }
  }, [orderData, initializeReturnItems]); // เพิ่ม initializeReturnItems เป็น dependency

  // ดึงจำนวนสินค้าที่จะคืน
  const getReturnQty = useCallback((sku: string): number => {
    return returnItems[sku] || 0;
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
      
      return {
        ...prev,
        [sku]: newQty,
      };
    });
  }, [orderData]);

  // คำนวณมูลค่ารวมของสินค้าที่คืน
  const calculateTotalAmount = useCallback((items: any[]): number => {
    return items.reduce((sum, item) => sum + (Math.abs(item.price) * item.returnQty), 0);
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