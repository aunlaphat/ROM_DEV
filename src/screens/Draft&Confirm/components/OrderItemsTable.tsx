import React from "react";
import { Table, Button, Popconfirm } from "antd";
import type { ColumnType } from 'antd/es/table';
import { DeleteOutlined } from "@ant-design/icons";
import { OrderItem } from "../../../redux/draftConfirm/types";
import { useDraftConfirm } from "../../../redux/draftConfirm/hook";

interface OrderItemsTableProps {
  items: OrderItem[];
  isDraftMode: boolean;
}

export const OrderItemsTable: React.FC<OrderItemsTableProps> = ({
  items,
  isDraftMode,
}) => {
  const { removeItem } = useDraftConfirm();

  // Handle delete item
  const handleDelete = (sku: string) => {
    if (items.length > 0) {
      removeItem({
        orderNo: items[0].orderNo,
        sku,
      });
    }
  };

  // Setup columns
  const getColumns = (): ColumnType<OrderItem>[] => {
    // Base columns used in both modes
    const baseColumns: ColumnType<OrderItem>[] = [
      {
        title: "SKU",
        dataIndex: "sku",
        key: "sku",
        render: (text: string) => (
          <span style={{ color: "#35465B" }}>{text}</span>
        ),
      },
      {
        title: "Name",
        dataIndex: "itemName",
        key: "itemName",
        render: (text: string) => (
          <span style={{ color: "#35465B" }}>{text}</span>
        ),
      },
      {
        title: "QTY",
        dataIndex: "qty",
        key: "qty",
        render: (text: number) => (
          <span style={{ color: "#35465B" }}>{text}</span>
        ),
      },
      {
        title: "Price",
        dataIndex: "price",
        key: "price",
        render: (text: number) => (
          <span style={{ color: "#35465B" }}>{text}</span>
        ),
      },
    ];

    // Add Action column only in draft mode
    if (isDraftMode) {
      baseColumns.push({
        title: "Action",
        key: "action",
        render: (_: any, record: OrderItem) =>
          record.type === "addon" ? (
            <Popconfirm
              title="Are you sure to delete this item?"
              onConfirm={() => handleDelete(record.sku)}
              okText="Yes"
              cancelText="No"
            >
              <Button
                type="link"
                icon={<DeleteOutlined style={{ color: "red" }} />}
              />
            </Popconfirm>
          ) : null,
      });
    }

    return baseColumns;
  };

  return (
    <Table
      id={`ItemTable${isDraftMode ? "3" : "4"}`}
      components={{
        header: {
          cell: (props: React.HTMLAttributes<HTMLElement>) => (
            <th
              {...props}
              style={{ backgroundColor: "#E9F3FE", color: "#35465B" }}
            />
          ),
        },
      }}
      dataSource={items}
      columns={getColumns()}
      rowKey="sku"
      pagination={false}
    />
  );
};