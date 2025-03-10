import React from "react";
import { Table, Button, Tooltip, Space } from "antd";
import { FormOutlined, StopOutlined } from "@ant-design/icons";
import dayjs from "dayjs";
import { Order } from "../../../redux/draftConfirm/types";

interface OrdersTableProps {
  orders: Order[];
  activeTabKey: string;
  onEdit: (orderNo: string) => void;
  onCancel: (orderNo: string) => void;
}

export const OrdersTable: React.FC<OrdersTableProps> = ({
  orders,
  activeTabKey,
  onEdit,
  onCancel,
}) => {
  // Columns configuration
  const columns = [
    {
      title: "Order",
      dataIndex: "orderNo",
      key: "orderNo",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "SO/INV",
      dataIndex: "soNo",
      key: "soNo",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "Customer",
      dataIndex: "customerId",
      key: "customerId",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "SR",
      dataIndex: "srNo",
      key: "srNo",
      render: (text: string | null) => (
        <span style={{ color: "#35465B" }}>{text || "-"}</span>
      ),
    },
    {
      title: "Return Tracking",
      dataIndex: "trackingNo",
      key: "trackingNo",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "Transport",
      dataIndex: "logistic",
      key: "logistic",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "Channel",
      dataIndex: "channelId",
      key: "channelId",
      render: (text: number) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "Date Create",
      dataIndex: "createDate",
      key: "createDate",
      render: (text: string) => (
        <span style={{ color: "#35465B" }}>
          {dayjs(text).format("YYYY-MM-DD")}
        </span>
      ),
    },
    {
      title: "Warehouse",
      dataIndex: "warehouseId",
      key: "warehouseId",
      render: (text: number) => (
        <span style={{ color: "#35465B" }}>{text}</span>
      ),
    },
    {
      title: "Action",
      key: "action",
      render: (_: any, record: Order) => (
        <Space>
          <Tooltip title={activeTabKey === "1" ? "Edit" : "View Details"}>
            <Button
              type="link"
              icon={<FormOutlined />}
              onClick={() => onEdit(record.orderNo)}
              style={{ color: "gray" }}
            />
          </Tooltip>
          {activeTabKey === "1" && (
            <Tooltip title="Cancel">
              <Button
                type="link"
                danger
                icon={<StopOutlined />}
                onClick={() => onCancel(record.orderNo)}
                style={{ color: "red" }}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Table
      id={`Table${activeTabKey === "2" ? "2" : ""}`}
      components={{
        header: {
          cell: (props: React.HTMLAttributes<HTMLElement>) => (
            <th
              {...props}
              style={{
                backgroundColor: "#E9F3FE",
                color: "#35465B",
              }}
            />
          ),
        },
      }}
      pagination={false}
      style={{ width: "100%", tableLayout: "fixed" }}
      scroll={{ x: "max-content" }}
      dataSource={orders}
      columns={columns}
      rowKey="orderNo"
    />
  );
};