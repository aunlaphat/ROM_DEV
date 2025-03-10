// src/screens/Orders/Marketplace/components/ReturnOrderSteps.tsx

import React, { useMemo } from "react";
import { theme, Typography, Space, Tooltip, Badge } from "antd";
import {
  SearchOutlined,
  FormOutlined,
  NumberOutlined,
  CheckCircleOutlined,
  EyeOutlined,
  ArrowRightOutlined,
  CheckOutlined,
  ClockCircleOutlined,
  InfoCircleOutlined
} from "@ant-design/icons";
import { ReturnOrderState } from "../types";
import styles from "../styles/ReturnOrder.module.css";

const { Text, Title } = Typography;
const { useToken } = theme;

interface ReturnOrderStepsProps {
  currentStep: "search" | "create" | "sr" | "preview" | "confirm";
  orderData: ReturnOrderState["orderData"];
  getStepStatus: (stepKey: string) => "process" | "finish" | "wait";
}

const ReturnOrderSteps: React.FC<ReturnOrderStepsProps> = ({
  currentStep,
  orderData,
  getStepStatus,
}) => {
  const { token } = useToken();

  // คำนวณเปอร์เซ็นต์ความคืบหน้า
  const progress = useMemo(() => {
    const stepOrder = ['search', 'create', 'sr', 'preview', 'confirm'];
    const currentIndex = stepOrder.indexOf(currentStep);
    return ((currentIndex + 1) / stepOrder.length) * 100;
  }, [currentStep]);

  // ข้อมูล steps ที่ปรับปรุงแล้ว
  const steps = useMemo(() => [
    {
      key: "search",
      title: "ค้นหา",
      subtitle: "Search Order",
      icon: <SearchOutlined />,
      description: "ค้นหา Order ที่ต้องการคืนสินค้า",
      help: "กรอกเลข SO หรือ Order Number เพื่อค้นหาข้อมูลคำสั่งซื้อที่ต้องการคืน"
    },
    {
      key: "create",
      title: "สร้างคำสั่งคืน",
      subtitle: "Create Return",
      icon: <FormOutlined />,
      description: "ระบุข้อมูลและเลือกสินค้าที่ต้องการคืน",
      help: "กรอกข้อมูลการคืนและเลือกรายการสินค้าที่ต้องการคืน โดยระบุจำนวนในแต่ละรายการ"
    },
    {
      key: "sr",
      title: "สร้าง SR",
      subtitle: "Generate SR",
      icon: <NumberOutlined />,
      description: orderData?.head.srNo 
        ? `SR: ${orderData.head.srNo}` 
        : "สร้างเลข SR Number",
      help: "เลข SR (Sale Return) เป็นเลขอ้างอิงสำหรับการคืนสินค้า ซึ่งจะถูกสร้างโดยระบบ"
    },
    {
      key: "preview",
      title: "ตรวจสอบ",
      subtitle: "Review",
      icon: <EyeOutlined />,
      description: "ตรวจสอบข้อมูลก่อนยืนยัน",
      help: "ตรวจสอบความถูกต้องของข้อมูลการคืนสินค้าทั้งหมด ก่อนที่จะยืนยันการคืน"
    },
    {
      key: "confirm",
      title: "ยืนยัน",
      subtitle: "Confirm",
      icon: <CheckCircleOutlined />,
      description: "ยืนยันและเสร็จสิ้น",
      help: "ยืนยันการคืนสินค้าเพื่อเสร็จสิ้นกระบวนการ"
    },
  ], [orderData]);

  // สร้าง indicator สถานะของ step
  const renderStepIndicator = (status: "process" | "finish" | "wait", index: number) => {
    if (status === "finish") {
      return (
        <div className={styles.stepIndicatorFinish}>
          <CheckOutlined />
        </div>
      );
    }
    
    if (status === "process") {
      return (
        <div className={styles.stepIndicatorActive}>
          <span>{index + 1}</span>
        </div>
      );
    }
    
    return (
      <div className={styles.stepIndicatorWait}>
        <span>{index + 1}</span>
      </div>
    );
  };

  return (
    <div className={styles.stepsContainer}>
      {/* Progress Bar */}
      <div className={styles.progressBarContainer}>
        <div 
          className={styles.progressBarFill} 
          style={{ 
            width: `${progress}%`,
            backgroundColor: progress === 100 ? token.colorSuccess : token.colorPrimary 
          }}
        />
      </div>

      {/* Steps Content */}
      <div className={styles.stepsContentWrapper}>
        {steps.map((step, index) => {
          const status = getStepStatus(step.key);
          const isActive = currentStep === step.key;
          const isCompleted = status === "finish";
          const isWaiting = status === "wait";
          
          return (
            <div 
              key={step.key}
              className={`${styles.stepItem} ${isActive ? styles.stepItemActive : ''} ${isCompleted ? styles.stepItemCompleted : ''} ${isWaiting ? styles.stepItemWaiting : ''}`}
            >
              <Tooltip 
                title={
                  <div>
                    <div>{step.description}</div>
                    <div className={styles.tooltipHelp}>
                      <InfoCircleOutlined /> {step.help}
                    </div>
                  </div>
                }
                placement="bottom"
              >
                <div className={styles.stepContent}>
                  {/* Step Indicator */}
                  {renderStepIndicator(status, index)}
                  
                  {/* Step Info */}
                  <div className={styles.stepInfo}>
                    <div className={styles.stepTitle}>
                      <Space align="center" size={4}>
                        {step.icon}
                        <Text strong={isActive} className={isActive ? styles.activeStepTitle : ''}>
                          {step.title}
                        </Text>
                      </Space>
                    </div>
                    <div className={styles.stepSubtitle}>
                      {step.subtitle}
                    </div>
                  </div>
                  
                  {/* Step Status */}
                  <div className={styles.stepStatus}>
                    {status === "process" && (
                      <Badge status="processing" text="กำลังดำเนินการ" />
                    )}
                    {status === "finish" && (
                      <Badge status="success" text="เสร็จสิ้น" />
                    )}
                    {status === "wait" && (
                      <Badge status="default" text="รอดำเนินการ" />
                    )}
                  </div>
                  
                  {/* Connector */}
                  {index < steps.length - 1 && (
                    <div className={styles.stepConnector}>
                      <ArrowRightOutlined />
                    </div>
                  )}
                </div>
              </Tooltip>
              
              {/* Current Step Description (Mobile Only) */}
              {isActive && (
                <div className={styles.mobileStepDescription}>
                  <Text type="secondary">{step.description}</Text>
                </div>
              )}
            </div>
          );
        })}
      </div>
      
      {/* Mobile Progress Indicator
      <div className={styles.mobileProgressIndicator}>
        <div className={styles.mobileProgressText}>
          <Space align="center">
            <ClockCircleOutlined />
            <span>ขั้นตอนที่ {steps.findIndex(s => s.key === currentStep) + 1} จาก {steps.length}</span>
          </Space>
        </div>
        <div className={styles.mobileProgressBarContainer}>
          <div 
            className={styles.mobileProgressBarFill} 
            style={{ 
              width: `${progress}%`,
              backgroundColor: progress === 100 ? token.colorSuccess : token.colorPrimary 
            }}
          />
        </div>
      </div> */}
    </div>
  );
};

export default ReturnOrderSteps;