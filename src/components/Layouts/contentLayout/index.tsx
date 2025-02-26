import { Content } from "antd/es/layout/layout";
import { useSelector } from "react-redux";
import { RoleID } from "../../../constants/roles";
import { ROUTES } from "../../../resources/routes";
import { logger } from "../../../utils/logger";

interface ContentLayoutProps {
  children: React.ReactNode;
}

// สร้าง `ROLE_ROUTES` ให้อ่านง่าย และลดโค้ดซ้ำซ้อน
const ROLE_ROUTES: Record<RoleID, (typeof ROUTES)[keyof typeof ROUTES][]> = {
  [RoleID.ADMIN]: Object.values(ROUTES),
  [RoleID.TRADE_CONSIGN]: [ROUTES.ROUTE_MAIN],
  [RoleID.ACCOUNTING]: [ROUTES.ROUTE_MAIN],
  [RoleID.WAREHOUSE]: [ROUTES.ROUTE_MAIN],
  [RoleID.VIEWER]: [ROUTES.ROUTE_MAIN],
};

const ContentLayout: React.FC<ContentLayoutProps> = ({ children }) => {
  const auth = useSelector((state: any) => state.auth);
  const roleID: RoleID | undefined = auth?.user?.roleID;

  const userRoutes = roleID ? ROLE_ROUTES[roleID] ?? [] : [];

  if (!roleID) {
    logger.auth(
      "warn",
      "⚠️ No valid roleID found in user data, rendering empty routes."
    );
  } else {
    logger.auth("info", `🔹 Routes Loaded for Role ${roleID}:`, userRoutes);
  }

  return (
    <Content
      className="site-layout-background"
      style={{ padding: 24, minHeight: 280 }}
    >
      {children}
    </Content>
  );
};

export default ContentLayout;
