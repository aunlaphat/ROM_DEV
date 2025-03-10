import React, { useState, useEffect } from "react";
import {
  Modal,
  Button,
  Form,
  Input,
  Row,
  Col,
  Select,
  InputNumber,
  Spin,
  Empty,
} from "antd";
import { PlusCircleOutlined } from "@ant-design/icons";
import { Order } from "../../../redux/draftConfirm/types";
import { useDraftConfirm } from "../../../redux/draftConfirm/hook";
import { OrderItemsTable } from "./OrderItemsTable";
import { message } from "antd";

interface EditOrderModalProps {
  visible: boolean;
  order: Order | null;
  activeTabKey: string;
  loading: boolean;
  onCancel: () => void;
}

export const EditOrderModal: React.FC<EditOrderModalProps> = ({
  visible,
  order,
  activeTabKey,
  loading,
  onCancel,
}) => {
  // Redux hooks
  const {
    codeRList,
    selectedOrderItems,
    addItem,
    confirmDraftOrder,
    fetchCodeR,
  } = useDraftConfirm();

  // Local state for form fields
  const [codeR, setCodeR] = useState<string | undefined>(undefined);
  const [nameR, setNameR] = useState<string | undefined>(undefined);
  const [qty, setQty] = useState<number | null>(null);
  const [price, setPrice] = useState<number | null>(null);
  const [loadingCodeR, setLoadingCodeR] = useState<boolean>(false);

  // Fetch CodeR list when modal is opened
  useEffect(() => {
    if (visible && activeTabKey === "1") {
      console.log('Modal opened - Fetching CodeR list');
      setLoadingCodeR(true);
      
      // Wrap in setTimeout to ensure it runs after Redux is ready
      setTimeout(() => {
        fetchCodeR();
        setLoadingCodeR(false);
      }, 100);
    }
  }, [visible, activeTabKey, fetchCodeR]);

  // Debug CodeR list
  useEffect(() => {
    if (codeRList) {
      console.log('CodeR List State:', {
        type: typeof codeRList,
        isArray: Array.isArray(codeRList),
        length: codeRList.length,
        sample: codeRList[0]
      });
    } else {
      console.log('CodeR List is null or undefined');
    }
  }, [codeRList]);

  // Handle code selection change - also update name
  const handleCodeChange = (value: string) => {
    setCodeR(value);
    
    // Find matching name for selected code
    const selectedItem = (codeRList || []).find(item => item.sku === value);
    if (selectedItem) {
      setNameR(selectedItem.nameAlias);
    }
  };

  // Handle name selection change - also update code
  const handleNameChange = (value: string) => {
    setNameR(value);
    
    // Find matching code for selected name
    const selectedItem = (codeRList || []).find(item => item.nameAlias === value);
    if (selectedItem) {
      setCodeR(selectedItem.sku);
    }
  };

  // Handle form submission
  const handleOk = () => {
    if (order) {
      if (activeTabKey === "1") {
        confirmDraftOrder({ orderNo: order.orderNo });
        message.success("Order updated successfully");
      }
      onCancel();
    }
  };

  // Reset form fields
  const resetFormFields = () => {
    setCodeR(undefined);
    setNameR(undefined);
    setQty(null);
    setPrice(null);
  };

  // Handle adding new item
  const handleAdd = () => {
    if (order && codeR && nameR && qty && price) {
      addItem({
        orderNo: order.orderNo,
        sku: codeR,
        itemName: nameR,
        qty: qty,
        returnQty: qty,
        price: price,
      });
      resetFormFields();
    } else {
      message.warning("Please fill in all fields");
    }
  };

  // ป้องกัน null/undefined ทุกครั้งที่ใช้ codeRList
  const safeCodeRList = codeRList || [];

  // Prepare options for selects
  const codeROptions = safeCodeRList.map((item) => ({
    value: item.sku,
    label: item.sku,
    id: item.sku,
  }));

  const codeNameOptions = safeCodeRList.map((item) => ({
    value: item.nameAlias,
    label: item.nameAlias,
    id: item.sku,
  }));

  const { Option } = Select;

  // เพิ่ม Modal เต็มจอ
  return (
    <Modal
      closable={false}
      width="80%"
      style={{ maxWidth: "1200px" }}
      centered={true}
      title={activeTabKey === "1" ? "Edit Order" : "Confirm Order"}
      visible={visible}
      onOk={handleOk}
      footer={
        <div style={{ display: "flex", justifyContent: "center" }}>
          {activeTabKey === "1" ? (
            <>
              <Button
                id="Update"
                onClick={handleOk}
                style={{
                  marginLeft: 8,
                  backgroundColor: "#14C11B",
                  color: "#FFF",
                }}
                loading={loading}
              >
                Update
              </Button>
              <Button
                id="Cancel"
                onClick={onCancel}
                style={{
                  marginLeft: 8,
                  background: "#D9D9D9",
                  color: "#909090",
                }}
                disabled={loading}
              >
                Cancel
              </Button>
            </>
          ) : (
            <Button
              id="Close"
              onClick={onCancel}
              style={{
                marginLeft: 8,
                background: "#D9D9D9",
                color: "#909090",
              }}
            >
              Close
            </Button>
          )}
        </div>
      }
    >
      {order && (
        <Spin spinning={loading}>
          <Form layout="vertical" style={{ marginTop: 20 }}>
            {/* Order information form fields */}
            <Row gutter={16}>
              <Col span={6}>
                <Form.Item
                  id="Order"
                  label={<span style={{ color: "#657589" }}>Order</span>}
                >
                  <Input
                    style={{ height: 40 }}
                    value={order.orderNo}
                    readOnly
                    disabled
                  />
                </Form.Item>
              </Col>
              <Col span={6}>
                <Form.Item
                  id="So/Inv"
                  label={<span style={{ color: "#657589" }}>SO/INV</span>}
                >
                  <Input
                    style={{ height: 40 }}
                    value={order.soNo}
                    disabled
                  />
                </Form.Item>
              </Col>
              <Col span={6}>
                <Form.Item
                  id="SR"
                  label={<span style={{ color: "#657589" }}>SR</span>}
                >
                  <Input
                    style={{ height: 40 }}
                    value={order.srNo || "-"}
                    disabled
                  />
                </Form.Item>
              </Col>
              <Col span={6}>
                <Form.Item
                  id="Customer"
                  label={<span style={{ color: "#657589" }}>Customer</span>}
                >
                  <Input
                    style={{ height: 40 }}
                    value={order.customerId || "-"}
                    disabled
                  />
                </Form.Item>
              </Col>
            </Row>

            {/* CodeR form fields (only in Edit mode) */}
            {activeTabKey === "1" && (
              <Row gutter={16}>
                <Col span={5}>
                  <Form.Item
                    id="codeR"
                    label={<span style={{ color: "#657589" }}>กรอกโค้ด R</span>}
                  >
                    <Select
                      style={{ height: 40 }}
                      value={codeR}
                      onChange={handleCodeChange}
                      showSearch
                      placeholder={loadingCodeR ? "กำลังโหลดข้อมูล..." : "เลือกโค้ด R"}
                      filterOption={(input, option) =>
                        (option?.label?.toString() || "")
                          .toLowerCase()
                          .includes(input.toLowerCase())
                      }
                      notFoundContent={
                        loadingCodeR ? (
                          <Spin size="small" />
                        ) : safeCodeRList.length === 0 ? (
                          <Empty description="ไม่พบข้อมูล" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                        ) : null
                      }
                      onFocus={() => {
                        if (!codeRList || codeRList.length === 0) {
                          setLoadingCodeR(true);
                          fetchCodeR();
                          setTimeout(() => setLoadingCodeR(false), 500);
                        }
                      }}
                    >
                      {codeROptions.map((code) => (
                        <Option key={code.value} value={code.value}>
                          {code.label}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={5}>
                  <Form.Item
                    id="NameR"
                    label={<span style={{ color: "#657589" }}>ชื่อของโค้ด R</span>}
                  >
                    <Select
                      style={{ height: 40 }}
                      value={nameR}
                      onChange={handleNameChange}
                      showSearch
                      placeholder={loadingCodeR ? "กำลังโหลดข้อมูล..." : "เลือกชื่อโค้ด R"}
                      filterOption={(input, option) =>
                        (option?.label?.toString() || "")
                          .toLowerCase()
                          .includes(input.toLowerCase())
                      }
                      notFoundContent={
                        loadingCodeR ? (
                          <Spin size="small" />
                        ) : safeCodeRList.length === 0 ? (
                          <Empty description="ไม่พบข้อมูล" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                        ) : null
                      }
                      onFocus={() => {
                        if (!codeRList || codeRList.length === 0) {
                          setLoadingCodeR(true);
                          fetchCodeR();
                          setTimeout(() => setLoadingCodeR(false), 500);
                        }
                      }}
                    >
                      {codeNameOptions.map((name) => (
                        <Option key={name.value} value={name.value}>
                          {name.label}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={5}>
                  <Form.Item
                    id="qty"
                    label={<span style={{ color: "#657589" }}>QTY:</span>}
                  >
                    <InputNumber
                      min={1}
                      max={100}
                      value={qty}
                      onChange={(value) => setQty(value)}
                      style={{
                        width: "100%",
                        height: "40px",
                        lineHeight: "40px",
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={5}>
                  <Form.Item
                    id="price"
                    label={<span style={{ color: "#657589" }}>Price:</span>}
                  >
                    <InputNumber
                      min={1}
                      max={100000}
                      value={price}
                      onChange={(value) => setPrice(value)}
                      step={0.01}
                      style={{
                        width: "100%",
                        height: "40px",
                        lineHeight: "40px",
                      }}
                    />
                  </Form.Item>
                </Col>
                <Col span={4}>
                  <Button
                    id="add"
                    type="primary"
                    style={{ width: "100%", height: "40px", marginTop: 30 }}
                    onClick={handleAdd}
                    disabled={!codeR || !nameR || !qty || !price}
                  >
                    <PlusCircleOutlined />
                    Add
                  </Button>
                </Col>
              </Row>
            )}
          </Form>

          {/* Items table - เพิ่มความสูงให้กับตาราง */}
          <div style={{ marginTop: 16, height: 'calc(65vh - 250px)', overflow: 'auto' }}>
            <OrderItemsTable 
              items={selectedOrderItems || []}
              isDraftMode={activeTabKey === "1"}
            />
          </div>
        </Spin>
      )}
    </Modal>
  );
};