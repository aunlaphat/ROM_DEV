import { Modal, notification } from "antd";
import { useSelector } from "react-redux";
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

const Alert = () => {
  const [api, contextHolder] = notification.useNotification();
  const alert = useSelector((state: any) => state.alert);

  useEffect(() => {
    if (alert.open) {
      const delayRemoveAlert = () => {
        delay(3000);
      };
      const openNotification = (
        placement: NotificationPlacement,
        type: NotificationType
      ) => {
        api[type]({
          message: `${alert.title}`,
          description: `${alert.message}`,
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
        // Use the provided props to show the modal
        Modal[type]({
          title: title, // Here you can specify the title property
          content: content,
          onCancel: onCancel,
          onOk: onOk,
        });
      };

      const renderAlert = () => {
        switch (alert.model) {
          case "notification":
            delayRemoveAlert();
            openNotification("topLeft", alert.type as NotificationType);
            break;
          case "modal":
            let content: any = {
              type: alert.type,
              content: alert.message,
              onCancel: () => {
                closeAlert();
              },
              onOk: () => {
                closeAlert();
              },
            };
            if (alert.title) {
              content = { ...content, title: alert.title };
            }
            return showModal(content);
          default:
            delayRemoveAlert();
            notification.info({ message: "" });
            break;
        }
      };

      //   alert.open && renderAlert();
      renderAlert();
    }
  }, [alert]);

  return contextHolder;
};

export default Alert;
