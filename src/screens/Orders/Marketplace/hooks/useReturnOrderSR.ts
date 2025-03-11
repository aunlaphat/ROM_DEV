// src/screens/Orders/Marketplace/hooks/useReturnOrderSR.ts
import { useCallback } from "react";
import { FormInstance } from "antd/lib/form";
import { message, notification } from "antd";
import { OrderData } from "../../../../redux/orders/types/state";

/**
 * Custom hook สำหรับจัดการการสร้าง SR
 */
export const useReturnOrderSR = (
  orderData: OrderData | null,
  form: FormInstance,
  generateSr: (orderNo: string) => void,
  setHasGeneratedSr: (value: boolean) => void,
  setStepLoading: (value: boolean) => void
) => {
  // จัดการการสร้าง SR Number
  const handleCreateSr = useCallback(() => {
    try {
      if (!orderData?.head.orderNo) {
        message.error("ไม่พบเลขที่ Order");
        return;
      }

      form
        .validateFields([
          "warehouseFrom",
          "returnDate",
          "trackingNo",
          "transportType",
        ])
        .then(() => {
          // สร้าง SR Number
          console.log(
            `[ReturnOrder] Generating SR for order: ${orderData.head.orderNo}`
          );
          setHasGeneratedSr(true);
          // ไม่ต้องกำหนด stepLoading ที่นี่ เพราะจะถูกจัดการโดย useEffect ที่ติดตาม redux loading
          generateSr(orderData.head.orderNo);
        })
        .catch((error) => {
          setHasGeneratedSr(false);
          notification.error({
            message: "ไม่สามารถสร้าง SR ได้",
            description: "กรุณาตรวจสอบข้อมูลและลองใหม่อีกครั้ง",
          });
        });
    } catch (error: any) {
      console.error("[ReturnOrder] Create SR error:", error);
      setHasGeneratedSr(false);
      notification.error({
        message: "เกิดข้อผิดพลาด",
        description: error.message || "ไม่สามารถสร้างเลข SR ได้",
      });
    }
  }, [orderData, form, generateSr, setHasGeneratedSr]);

  return {
    handleCreateSr,
  };
};
