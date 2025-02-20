import { Layout } from "antd";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import SiderLayout from "./siderLayout";
import HeaderBar from "./headerLayout";
import ContentLayout from "./contentLayout";
import Loading from "../loading";

const LayoutPage = ({ children }: any) => {
  const location = useLocation();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(window.innerWidth <= 767);
  const [collapsedWidth, setCollapsedWidth] = useState(window.innerWidth <= 767 ? 0 : undefined);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null); // ✅ ใช้ state ควบคุม Auth

  useEffect(() => {
    const handleResize = () => {
      setCollapsed(window.innerWidth <= 767);
      setCollapsedWidth(window.innerWidth <= 767 ? 0 : undefined);
    };

    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}/auth/`, {
          method: "GET",
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error("Unauthorized");
        }

        setIsAuthenticated(true); // ✅ ตั้งค่า Authenticated
      } catch (error) {
        console.warn("🚨 No JWT found → Redirecting to Login");
        setIsAuthenticated(false);
      }
    };

    checkAuth();
  }, []);

  // ✅ ถ้ายังไม่ได้ตรวจสอบ Auth → แสดง Loading ก่อน
  if (isAuthenticated === null) {
    return <Loading />;
  }

  // ✅ ถ้าไม่ได้ Authenticated → Redirect ไปหน้า Login (ป้องกัน Loop)
  if (!isAuthenticated && location.pathname !== "/") {
    navigate("/");
    return null;
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <SiderLayout collapsed={collapsed} collapsedWidth={collapsedWidth} toggle={() => setCollapsed(!collapsed)} />
      <Layout className="site-layout">
        <HeaderBar collapsed={collapsed} toggle={() => setCollapsed(!collapsed)} />
        <ContentLayout>{children}</ContentLayout>
      </Layout>
    </Layout>
  );
};

export default LayoutPage;
