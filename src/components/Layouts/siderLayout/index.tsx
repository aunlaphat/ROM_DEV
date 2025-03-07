import { Layout, Menu, Divider } from "antd";
import { Link, useLocation } from "react-router-dom";
import { useMemo } from "react";
import { getAccessibleRoutes } from "../../../resources/routes";
import { logger } from "../../../utils/logger";
import { Icon, isIconKey } from "../../../resources/icon";
import logo from "../../../assets/images/logo_org.png";
import { ScrollMenu } from "./style";
import { useAuth } from "../../../hooks/auth";

const { Sider } = Layout;

interface SiderLayoutProps {
  collapsed: boolean;
  collapsedWidth?: number;
  toggle: () => void;
}

const SiderLayout: React.FC<SiderLayoutProps> = ({
  collapsed,
  collapsedWidth,
}) => {
  const location = useLocation();
  const { roleID } = useAuth();

  // ใช้ฟังก์ชัน getAccessibleRoutes จากไฟล์ routes.ts เพื่อกรองเมนูตามบทบาท
  const accessibleRoutes = useMemo(() => {
    const routes = getAccessibleRoutes(roleID);
    logger.log("info", `🔹 Sidebar Menu for Role ${roleID}:`, routes);
    return routes;
  }, [roleID]);

  // สร้างรายการเมนูจาก routes ที่เข้าถึงได้
  const menuItems = useMemo(() => {
    return accessibleRoutes.map((route) => {
      // ตรวจสอบว่ามี icon หรือไม่ และใช้ Icon component ที่เหมาะสม
      let iconComponent = null;

      if (route.icon && typeof route.icon === "string") {
        // ตรวจสอบว่า key มีอยู่ใน Icon object หรือไม่โดยใช้ type guard
        if (isIconKey(route.icon)) {
          const IconComponent = Icon[route.icon];
          iconComponent = <IconComponent />;
        } else {
          // แค่ log warning แต่ไม่ assign ค่าให้ iconComponent
          console.warn(
            `Icon "${route.icon}" not found for route "${route.PATH}"`
          );
        }
      }

      return {
        key: route.PATH,
        icon: iconComponent,
        label: <Link to={route.PATH}>{route.LABEL}</Link>,
      };
    });
  }, [accessibleRoutes]);

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={collapsed}
      collapsedWidth={collapsedWidth}
      theme="light"
      style={{
        position: "sticky",
        top: 0,
        height: "100vh",
        overflow: "auto",
      }}
    >
      <div className="logo">
        <img
          src={logo}
          alt="Company Logo"
          className="avatar"
          style={{
            width: collapsed ? "70%" : "90%",
            margin: "10px auto",
            display: "block",
            transition: "all 0.2s",
          }}
        />
      </div>
      <Divider style={{ margin: "0 0 8px 0" }} />
      <ScrollMenu>
        <Menu
          theme="light"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
        />
      </ScrollMenu>
    </Sider>
  );
};

export default SiderLayout;
