// src/screens/Orders/Marketplace/hooks/useReturnOrderSearch.ts
import { useState, useCallback } from 'react';
import { message } from 'antd';
import { logger } from '../../../../utils/logger';

/**
 * Custom hook สำหรับจัดการการค้นหา Order
 */
export const useReturnOrderSearch = (
  searchOrder: (payload: any) => void, 
  setStep: (step: any) => void,
  setHasSearched: (value: boolean) => void
) => {
  const [selectedSalesOrder, setSelectedSalesOrder] = useState('');
  const [stepLoading, setStepLoading] = useState(false);

  // จัดการการเปลี่ยนแปลงข้อมูลช่องค้นหา
  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSelectedSalesOrder(e.target.value.trim());
  }, []);

  // จัดการการค้นหา
  const handleSearch = useCallback(() => {
    if (!selectedSalesOrder) {
      message.error("กรุณากรอกเลข SO/Order");
      return;
    }

    setStepLoading(true);
    try {
      const isSoNo = selectedSalesOrder.startsWith("SO");
      const searchPayload = {
        [isSoNo ? "soNo" : "orderNo"]: selectedSalesOrder.trim(),
      };

      logger.log('info', `[ReturnOrder] Searching order: ${JSON.stringify(searchPayload)}`);
      setHasSearched(true);
      searchOrder(searchPayload);
    } catch (error) {
      logger.error('[ReturnOrder] Search error:', error);
      setHasSearched(false);
    } finally {
      setStepLoading(false);
    }
  }, [selectedSalesOrder, searchOrder, setHasSearched]);

  return {
    selectedSalesOrder,
    setSelectedSalesOrder,
    stepLoading,
    setStepLoading,
    handleInputChange,
    handleSearch
  };
};