import React, { useState } from "react";
import { Modal, Form, Input, message } from "antd";
import { useDraftConfirm } from "../../../redux/draftConfirm/hook";

interface CancelOrderModalProps {
  visible: boolean;
  orderNo: string | null;
  onClose: () => void;
}

export const CancelOrderModal: React.FC<CancelOrderModalProps> = ({
  visible,
  orderNo,
  onClose,
}) => {
  const [cancelReason, setCancelReason] = useState<string>("");
  const [form] = Form.useForm();
  const { cancelOrder, loading } = useDraftConfirm();

  // Handle confirmation
  const handleConfirm = () => {
    if (orderNo && cancelReason) {
      // Send cancel request
      cancelOrder({ orderNo, cancelReason });
      // Reset and close
      setCancelReason("");
      onClose();
    } else {
      message.warning("Please provide a reason for cancellation");
    }
  };

  // Handle cancel
  const handleCancel = () => {
    setCancelReason("");
    onClose();
  };

  return (
    <Modal
      title="Cancel Order"
      visible={visible}
      onOk={handleConfirm}
      onCancel={handleCancel}
      okText="Confirm"
      confirmLoading={loading}
      okButtonProps={{
        style: { backgroundColor: "#ff4d4f", borderColor: "#ff4d4f" },
        disabled: !cancelReason,
      }}
    >
      <Form form={form} layout="vertical">
        <Form.Item label="Order No" required={false}>
          <Input value={orderNo || ""} disabled />
        </Form.Item>
        <Form.Item
          label="Reason for cancellation"
          required
          rules={[
            { required: true, message: "Please enter cancellation reason" },
          ]}
        >
          <Input.TextArea
            rows={4}
            value={cancelReason}
            onChange={(e) => setCancelReason(e.target.value)}
            placeholder="Please provide a reason for cancellation"
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};