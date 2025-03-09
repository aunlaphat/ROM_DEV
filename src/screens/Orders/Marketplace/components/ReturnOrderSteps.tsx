// src/screens/Orders/Marketplace/components/ReturnOrderSteps.tsx

import React from "react";
import { Steps, Badge, Tooltip, Space, Typography, theme } from "antd";
import {
  SearchOutlined,
  FormOutlined,
  NumberOutlined,
  CheckCircleOutlined,
  EyeOutlined,
  RightOutlined,
} from "@ant-design/icons";
import { ReturnOrderState } from "../types";

const { Text } = Typography;
const { useToken } = theme;

interface ReturnOrderStepsProps {
  currentStep: "search" | "create" | "sr" | "confirm" | "preview";
  orderData: ReturnOrderState["orderData"];
  getStepStatus: (stepKey: string) => "process" | "finish" | "wait";
}

const ReturnOrderSteps: React.FC<ReturnOrderStepsProps> = ({
  currentStep,
  orderData,
  getStepStatus,
}) => {
  const { token } = useToken();

  // ฟังก์ชันสร้างชื่อขั้นตอนและคำอธิบาย
  const renderStepTitle = (title: string, subTitle?: string) => (
    <Space direction="vertical" size={0}>
      <Text strong>{title}</Text>
      {subTitle && (
        <Text type="secondary" style={{ fontSize: "12px" }}>
          {subTitle}
        </Text>
      )}
    </Space>
  );

  // สร้างข้อมูลขั้นตอน
  const steps = [
    {
      key: "search",
      title: renderStepTitle("ค้นหา", "Search Order"),
      icon: <SearchOutlined />,
      description: (
        <Badge
          status={currentStep === "search" ? "processing" : "success"}
          text="ค้นหา Order ที่ต้องการคืนสินค้า"
          style={{
            whiteSpace: "nowrap",
          }}
        />
      ),
      disabled: currentStep === "confirm" && !!orderData?.head.srNo,
    },
    {
      key: "create",
      title: renderStepTitle("สร้างคำสั่ง", "Create Return"),
      icon: <FormOutlined />,
      description: (
        <Badge
          status={
            getStepStatus("create") === "process"
              ? "processing"
              : getStepStatus("create") === "finish"
              ? "success"
              : "default"
          }
          text="ระบุข้อมูลและเลือกสินค้าที่ต้องการคืน"
          style={{
            whiteSpace: "nowrap",
          }}
        />
      ),
      disabled: currentStep === "confirm" && !!orderData?.head.srNo,
    },
    {
      key: "sr",
      title: renderStepTitle("สร้าง SR", "Generate SR"),
      icon: <NumberOutlined />,
      description: (
        <Tooltip
          title={
            orderData?.head.srNo
              ? `SR Number: ${orderData.head.srNo}`
              : "รอการสร้าง SR"
          }
        >
          <Badge
            status={
              orderData?.head.srNo
                ? "success"
                : getStepStatus("sr") === "process"
                ? "processing"
                : "default"
            }
            text={
              orderData?.head.srNo ? (
                <Space size={4}>
                  <Text type="success" style={{ fontWeight: "bold" }}>
                    {orderData.head.srNo}
                  </Text>
                  <CheckCircleOutlined style={{ color: token.colorSuccess }} />
                </Space>
              ) : (
                "สร้างเลข SR Number"
              )
            }
            style={{
              whiteSpace: "nowrap",
            }}
          />
        </Tooltip>
      ),
    },
    {
      key: "preview",
      title: renderStepTitle("ตรวจสอบ", "Preview"),
      icon: <EyeOutlined />,
      description: (
        <Badge
          status={
            getStepStatus("preview") === "process"
              ? "processing"
              : getStepStatus("preview") === "finish"
              ? "success"
              : "default"
          }
          text="ตรวจสอบข้อมูลก่อนยืนยัน"
          style={{
            whiteSpace: "nowrap",
          }}
        />
      ),
    },
    {
      key: "confirm",
      title: renderStepTitle("ยืนยัน", "Confirm"),
      icon: <CheckCircleOutlined />,
      description: (
        <Badge
          status={
            getStepStatus("confirm") === "process" ? "processing" : "default"
          }
          text="ยืนยันและเสร็จสิ้น"
          style={{
            whiteSpace: "nowrap",
          }}
        />
      ),
    },
  ];

  return (
    <div
      style={{
        background: "#fff",
        padding: "24px",
        borderRadius: "8px",
        boxShadow: "0 1px 3px rgba(0,0,0,0.05)",
        position: "relative",
        overflow: "auto",
      }}
    >
      <Steps
        type="navigation"
        current={steps.findIndex((s) => s.key === currentStep)}
        items={steps}
        responsive={true}
        style={{
          padding: "8px 0",
        }}
        className="return-order-steps"
      />

      {/* แถบความคืบหน้า */}
      <div
        style={{
          position: "absolute",
          bottom: 0,
          left: 0,
          height: "4px",
          backgroundColor: token.colorPrimary,
          width: `${(steps.findIndex((s) => s.key === currentStep) + 1) * 20}%`,
          transition: "width 0.3s ease",
        }}
      />
    </div>
  );
};

export default ReturnOrderSteps;
