import React, { useState, useEffect, useCallback } from "react";
import {
  Button,
  ConfigProvider,
  Form,
  Layout,
  Row,
  Tabs,
  Col,
  message,
  Spin,
} from "antd";
import { DatePicker } from "antd";
import dayjs, { Dayjs } from "dayjs";
import isSameOrAfter from "dayjs/plugin/isSameOrAfter";
import isSameOrBefore from "dayjs/plugin/isSameOrBefore";
import isBetween from "dayjs/plugin/isBetween";
import "../Return.css";
import { useDraftConfirm } from "../../redux/draftConfirm/hook";
import { OrdersTable } from "./components/OrdersTable";
import { EditOrderModal } from "./components/EditOrderModal";
import { CancelOrderModal } from "./components/CancelOrderModal";
import { useLocation, useNavigate } from "react-router-dom";

// Initialize dayjs plugins
dayjs.extend(isSameOrAfter);
dayjs.extend(isSameOrBefore);
dayjs.extend(isBetween);

export const DraftAndConfirm = () => {
  const location = useLocation();
  const navigate = useNavigate();

  // ดึง tab จาก URL parameter
  const getTabFromUrl = useCallback(() => {
    const params = new URLSearchParams(location.search);
    return params.get("tab") || "1";
  }, [location.search]); // ระบุ dependency ให้กับ useCallback

  // Redux setup
  const {
    orders,
    selectedOrder,
    loading,
    error,
    fetchOrders,
    fetchOrderDetails,
    clearOrder,
  } = useDraftConfirm();

  // Local state
  const [activeTabKey, setActiveTabKey] = useState<string>(getTabFromUrl());
  const [dates, setDates] = useState<[Dayjs, Dayjs] | null>(null);
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [cancelModalVisible, setCancelModalVisible] = useState(false);
  const [orderToCancel, setOrderToCancel] = useState<string | null>(null);

  // Handle error messages
  useEffect(() => {
    if (error) {
      message.error(error);
    }
  }, [error]);

  // Event handlers
  const handleEdit = (orderNo: string) => {
    const statusConfID = activeTabKey === "1" ? 1 : 2; // 1 for Draft, 2 for Confirm
    fetchOrderDetails({ orderNo, statusConfID });
    setIsEditModalVisible(true);
  };

  const handleCancelEdit = () => {
    setIsEditModalVisible(false);
    clearOrder();
  };

  const handleCancelOrder = (orderNo: string) => {
    setOrderToCancel(orderNo);
    setCancelModalVisible(true);
  };

  const handleCancelModalClose = () => {
    setCancelModalVisible(false);
    setOrderToCancel(null);
  };

  const handleSearch = () => {
    if (dates && dates[0] && dates[1]) {
      const statusConfID = activeTabKey === "1" ? 1 : 2; // 1 for Draft, 2 for Confirm
      fetchOrders({
        statusConfID,
        startDate: dates[0].format("YYYY-MM-DD"),
        endDate: dates[1].format("YYYY-MM-DD"),
      });
    } else {
      message.warning("Please select date range");
    }
  };

  // อัพเดต URL เมื่อ tab เปลี่ยน
  const onTabChange = (key: string) => {
    setActiveTabKey(key);
    navigate(`/draft-and-confirm?tab=${key}`);

    // โหลดข้อมูลตาม tab ถ้ามีการเลือกวันที่แล้ว
    if (dates && dates[0] && dates[1]) {
      const statusConfID = key === "1" ? 1 : 2;
      fetchOrders({
        statusConfID,
        startDate: dates[0].format("YYYY-MM-DD"),
        endDate: dates[1].format("YYYY-MM-DD"),
      });
    }
  };

  // เมื่อโหลดหน้า ให้ตรวจสอบ tab จาก URL
  useEffect(() => {
    const tabKey = getTabFromUrl();
    if (tabKey !== activeTabKey) {
      setActiveTabKey(tabKey);
    }
  }, [location, activeTabKey, getTabFromUrl]);

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setDates([dates[0], dates[1]]);
    } else {
      setDates(null);
    }
  };

  const { RangePicker } = DatePicker;

  return (
    <ConfigProvider>
      <div
        style={{
          marginLeft: "28px",
          fontSize: "25px",
          fontWeight: "bold",
          color: "DodgerBlue",
        }}
      >
        Draft & Confirm Order MKP
      </div>

      <Layout>
        <Layout.Content
          style={{
            margin: "24px",
            padding: 36,
            minHeight: 360,
            background: "#fff",
            borderRadius: "8px",
            overflow: "auto",
          }}
        >
          {/* Tab Selection */}
          <Tabs
            id="card"
            activeKey={activeTabKey}
            onChange={onTabChange}
            type="card"
            items={[
              { label: "Draft", key: "1" },
              { label: "Confirm Draft", key: "2" },
            ]}
          />

          {/* Date Range Selector */}
          <Row
            gutter={8}
            align="middle"
            justify="center"
            style={{ marginTop: "20px" }}
          >
            <Col>
              <Form.Item
                id={`Selectdate${activeTabKey}`}
                layout="vertical"
                label="Select date"
                name="Select date"
                rules={[
                  {
                    required: true,
                    message: "Please select the Select date!",
                  },
                ]}
              >
                <RangePicker
                  value={dates}
                  style={{ height: "40px" }}
                  onChange={handleDateChange}
                />
              </Form.Item>
            </Col>
            <Col style={{ marginTop: "4px" }}>
              <Button
                id={`Search${activeTabKey}`}
                type="primary"
                style={{
                  height: "40px",
                  width: "100px",
                  background: "#32ADE6",
                }}
                onClick={handleSearch}
                loading={loading}
              >
                Search
              </Button>
            </Col>
          </Row>

          {/* Orders Table */}
          <Spin spinning={loading}>
            <OrdersTable
              orders={orders}
              activeTabKey={activeTabKey}
              onEdit={handleEdit}
              onCancel={handleCancelOrder}
            />
          </Spin>
        </Layout.Content>
      </Layout>

      {/* Modals */}
      <EditOrderModal
        visible={isEditModalVisible}
        order={selectedOrder}
        activeTabKey={activeTabKey}
        loading={loading}
        onCancel={handleCancelEdit}
      />

      <CancelOrderModal
        visible={cancelModalVisible}
        orderNo={orderToCancel}
        onClose={handleCancelModalClose}
      />
    </ConfigProvider>
  );
};

export default DraftAndConfirm;
