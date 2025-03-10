import { Modal, notification } from "antd";
import { useSelector } from "react-redux";
import { RootState } from "../../redux/types"; // ✅ ใช้ RootState จาก Redux
import { delay } from "../../function";
import type { NotificationArgsProps } from "antd";
import { useEffect } from "react";
import { closeAlert } from "./useAlert";

type NotificationPlacement = NotificationArgsProps["placement"];
type NotificationType = "success" | "error" | "info" | "warning";
type ModalType = "info" | "success" | "error" | "warning" | "confirm";

interface ModalProps {
  title: string;
  content: string;
  onCancel: () => void;
  onOk: () => void;
  type: ModalType;
}

const Alert: React.FC = () => {
  const [api, contextHolder] = notification.useNotification();
  const alert = useSelector((state: RootState) => state.alert); // ✅ ใช้ RootState จาก Redux

  useEffect(() => {
    if (alert?.open) {
      const delayRemoveAlert = async () => {
        await delay(3000); // ✅ ใช้ await เพื่อรอให้ delay ทำงาน
      };

      const openNotification = (
        placement: NotificationPlacement,
        type: NotificationType
      ) => {
        api[type]({
          message: alert.title || "แจ้งเตือน",
          description: alert.message || "",
          placement,
        });
      };

      const showModal = ({
        title,
        content,
        onCancel,
        onOk,
        type,
      }: ModalProps) => {
        Modal[type]({
          title: title || "แจ้งเตือน", // ✅ กำหนดค่าเริ่มต้น
          content: content,
          onCancel: onCancel,
          onOk: onOk,
        });
      };

      const renderAlert = () => {
        switch (alert?.model ?? "notification") { // ✅ ใช้ `??` ป้องกัน model เป็น `undefined`
          case "notification":
            delayRemoveAlert();
            openNotification("topLeft", alert.type as NotificationType);
            break;
          case "modal":
            const content: ModalProps = {
              type: (alert.type as ModalType) || "info",
              content: alert.message || "",
              onCancel: () => closeAlert(),
              onOk: () => closeAlert(),
              title: alert.title || "แจ้งเตือน",
            };
            return showModal(content);
          default:
            delayRemoveAlert();
            notification.info({ message: "แจ้งเตือน", description: "" });
            break;
        }
      };

      renderAlert();
    }
  }, [alert]);

  return contextHolder;
};

export default Alert;
