// src/screens/Orders/Marketplace/hooks/useReturnOrderNavigation.ts
import { useCallback, useEffect } from 'react';
import { FormInstance } from 'antd/lib/form';
import { notification, Modal } from 'antd';
import { useNavigate } from 'react-router-dom';
import { OrderData } from '../../../../redux/orders/types/state';
import { validateStepTransition } from '../utils/validation';

/**
 * Custom hook สำหรับจัดการการนำทางในระหว่างขั้นตอน
 */
export const useReturnOrderNavigation = (
  orderData: OrderData | null,
  returnOrder: any,
  currentStep: string,
  form: FormInstance,
  loading: boolean,
  srCreated: boolean,
  setStep: (step: any) => void,
  setSelectedSalesOrder: (value: string) => void,
  setReturnItems: (value: { [key: string]: number }) => void,
  setStepLoading: (value: boolean) => void,
  returnItems: { [key: string]: number } // เพิ่มพารามิเตอร์นี้
) => {
  const navigate = useNavigate();

  // Log when step changes
  useEffect(() => {
    console.log(`[Navigation] Step changed to: ${currentStep}`);
    console.log(`[Navigation] Current returnItems:`, returnItems);
  }, [currentStep, returnItems]);

  // จัดการปุ่มย้อนกลับ
  const handleBack = useCallback(() => {
    const steps = ['search', 'create', 'sr', 'preview', 'confirm'];
    const currentIndex = steps.indexOf(currentStep);
    const prevStep = steps[currentIndex - 1];

    console.log(`[Navigation] Attempting to go back from ${currentStep} to ${prevStep}`);

    if (!prevStep) {
      navigate("/home");
      return;
    }

    if (validateStepTransition(currentStep, prevStep, orderData, returnOrder)) {
      console.log(`[Navigation] Going back to ${prevStep} with returnItems:`, returnItems);
      setStep(prevStep as any);
      if (prevStep === 'search') {
        form.resetFields();
        setSelectedSalesOrder("");
        setReturnItems({});
      }
    } else {
      notification.warning({
        message: 'ไม่สามารถย้อนกลับได้',
        description: 'กรุณาตรวจสอบข้อมูลให้ครบถ้วน'
      });
    }
  }, [currentStep, form, navigate, setStep, orderData, returnOrder, setSelectedSalesOrder, setReturnItems, returnItems]);

  // จัดการการยกเลิก
  const handleCancel = useCallback(() => {
    Modal.confirm({
      title: 'ยืนยันการยกเลิก',
      content: 'คุณต้องการยกเลิกการทำรายการนี้ใช่หรือไม่? ข้อมูลที่กรอกไว้จะหายไป',
      okText: 'ยืนยัน',
      cancelText: 'ยกเลิก',
      onOk: () => {
        form.resetFields();
        setSelectedSalesOrder("");
        setReturnItems({});
        if (currentStep !== 'search') {
          setStep('search');
        }
      }
    });
  }, [form, currentStep, setStep, setSelectedSalesOrder, setReturnItems]);

  // จัดการการดำเนินการต่อไปยังขั้นตอนถัดไป
  const handleNext = useCallback(() => {
    setStepLoading(true);
    try {
      const steps = ['search', 'create', 'sr', 'preview', 'confirm'];
      const currentIndex = steps.indexOf(currentStep);
      const nextStep = steps[currentIndex + 1];

      if (!nextStep) return;

      console.log(`[Navigation] Checking transition from ${currentStep} to ${nextStep}`);
      console.log(`[Navigation] ReturnItems:`, returnItems);
      
      if (validateStepTransition(currentStep, nextStep, orderData, returnOrder)) {
        console.log(`[Navigation] Moving to next step: ${nextStep}`);
        setStep(nextStep as any);
      } else {
        notification.warning({
          message: 'ไม่สามารถดำเนินการต่อได้',
          description: 'กรุณาตรวจสอบข้อมูลให้ครบถ้วน'
        });
      }
    } catch (error) {
      console.error('[Navigation] Next step error:', error);
    } finally {
      setStepLoading(false);
    }
  }, [currentStep, setStep, orderData, returnOrder, setStepLoading, returnItems]);

  // Render ปุ่มย้อนกลับ
  const renderBackButton = useCallback(() => {
    const buttonTexts = {
      'search': 'Back to Home',
      'create': 'Back to Search',
      'sr': 'Back to Create',
      'preview': 'Back to SR',
      'confirm': 'Back to Preview'
    };
    
    const buttonText = buttonTexts[currentStep as keyof typeof buttonTexts] || 'Back';

    return {
      onClick: handleBack,
      disabled: loading || (currentStep === "confirm" && !!orderData?.head.srNo),
      text: buttonText
    };
  }, [currentStep, handleBack, loading, orderData]);

  // สถานะของแต่ละขั้นตอน
  const getStepStatus = useCallback((stepKey: string): 'process' | 'finish' | 'wait' => {
    switch (stepKey) {
      case "search":
        return currentStep === "search" ? "process" : "finish";
      case "create":
        if (currentStep === "search") return "wait";
        return currentStep === "create" ? "process" : "finish";
      case "sr":
        if (currentStep === "search" || currentStep === "create") return "wait";
        return currentStep === "sr" ? "process" : "finish";
      case "preview":
        if (currentStep === "search" || currentStep === "create" || currentStep === "sr") {
          // ใช้เงื่อนไขเพิ่มเติม: ถ้า sr step เสร็จแล้ว แต่ยังไม่ได้อยู่ที่ preview
          if (currentStep === "sr" && orderData?.head.srNo && srCreated) {
            return "process"; // เตรียมไปขั้นตอน preview
          }
          return "wait";
        }
        return currentStep === "preview" ? "process" : "finish";
      case "confirm":
        if (currentStep === "confirm") return "process";
        // อนุญาตให้ "confirm" เป็น active ถ้ามี srNo และ currentStep = preview
        if (currentStep === "preview" && orderData?.head.srNo) return "wait";
        return "wait";
      default:
        return "wait";
    }
  }, [currentStep, orderData?.head.srNo, srCreated]);

  return {
    handleBack,
    handleCancel,
    handleNext,
    renderBackButton,
    getStepStatus
  };
};