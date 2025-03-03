import { RoleID } from '../constants/roles';
import { STATUS } from '../constants/returnOrder';

/**
 * Parameters interface for calculating return order status
 * @interface CalculateStatusParams
 * @property {number} roleID - Role ID of user (2=Accounting, 3=Warehouse)
 * @property {boolean} isCNCreated - Whether CN has been created
 * @property {boolean} isEdited - Whether order has been edited
 */
interface CalculateStatusParams {
  roleID: number;
  isCNCreated: boolean;
  isEdited: boolean;
}

/**
 * Calculate return order status based on role and conditions
 * 
 * Business Rules:
 * 1. For Accounting (roleID = 2):
 *    - If CN not created: Status = Pending/Draft
 *    - If CN created: Status = Booking/Confirmed
 * 
 * 2. For Warehouse (roleID = 3):  
 *    - If order edited: Status = Pending/Draft
 *    - If not edited: Status = Booking/Confirmed
 * 
 * @param {CalculateStatusParams} params - Parameters for calculation
 * @returns {Object} Status IDs for return and confirm statuses
 * @throws {Error} If invalid role ID provided
 */
export const calculateReturnStatus = ({
  roleID,
  isCNCreated,
  isEdited
}: CalculateStatusParams) => {
  // Handle Accounting Role (RoleID = 2)
  if (roleID === RoleID.ACCOUNTING) {
    // CN not created yet - set to pending status
    if (!isCNCreated) {
      return {
        statusReturnID: STATUS.RETURN.PENDING,  // 1 = pending
        statusConfID: STATUS.CONFIRM.DRAFT      // 1 = draft
      };
    }
    // CN created - can proceed to booking
    return {
      statusReturnID: STATUS.RETURN.BOOKING,    // 3 = booking
      statusConfID: STATUS.CONFIRM.CONFIRMED    // 2 = confirmed
    };
  }

  // Handle Warehouse Role (RoleID = 3)
  if (roleID === RoleID.WAREHOUSE) {
    // Order has been edited - needs review
    if (isEdited) {
      return {
        statusReturnID: STATUS.RETURN.PENDING,  // 1 = pending
        statusConfID: STATUS.CONFIRM.DRAFT      // 1 = draft 
      };
    }
    // No edits - can proceed to booking
    return {
      statusReturnID: STATUS.RETURN.BOOKING,    // 3 = booking
      statusConfID: STATUS.CONFIRM.CONFIRMED    // 2 = confirmed
    };
  }

  // Invalid role provided
  throw new Error('Invalid role ID');
};
