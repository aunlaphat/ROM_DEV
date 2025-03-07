import { GET, POST } from '../index';
import { SEARCHORDER, CREATEBEFORERETURNORDER, GENERATESR, UPDATESTATUS } from '../path';

/**
 * API Service สำหรับ Return Order
 */
export const returnOrderAPI = {
  /**
   * ค้นหา Order จาก SO No หรือ Order No
   */
  searchOrder: async (soNo: string, orderNo: string) => {
    const response = await POST(SEARCHORDER, { soNo, orderNo });
    return response.data;
  },

  /**
   * สร้าง Return Order
   */
  createReturnOrder: async (data: any) => {
    const response = await POST(CREATEBEFORERETURNORDER, data);
    return response.data;
  },

  /**
   * Generate SR Number
   */
  generateSR: async (orderNo: string) => {
    const response = await POST(GENERATESR, { orderNo });
    return response.data;
  },

  /**
   * อัปเดตสถานะของ Return Order
   */
  confirmReturnOrder: async (data: any) => {
    const response = await POST(UPDATESTATUS, data);
    return response.data;
  }
};
