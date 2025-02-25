import { useEffect, useState, ReactNode } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { RootState, AppDispatch } from "../../redux/store";
import { checkAuthen } from "../../redux/auth/action";
import SiderLayout from "./siderLayout";
import HeaderBar from "./headerLayout";
import ContentLayout from "./contentLayout";
import Loading from "../loading";
import { Layout } from "antd";
import { logger } from '../../utils/logger';

const LayoutPage = ({ children }: any ) => {
  const location = useLocation();
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const authState = useSelector((state: RootState) => state.auth);
  const [collapsed, setCollapsed] = useState(window.innerWidth <= 767);
  const [collapsedWidth, setCollapsedWidth] = useState(window.innerWidth <= 767 ? 0 : undefined);
  const userRole = authState?.user?.roleID;

  // 1. Resize effect
  useEffect(() => {
    const handleResize = () => {
      setCollapsed(window.innerWidth <= 767);
      setCollapsedWidth(window.innerWidth <= 767 ? 0 : undefined);
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  // Single auth check effect
  useEffect(() => {
    const validateAuth = async () => {
      logger.auth('debug', 'Validating auth state', {
        isAuth: authState?.isAuthenticated,
        hasUser: !!authState?.user,
        path: location.pathname,
        token: !!localStorage.getItem('access_token')
      });

      // ถ้ามี token แต่ยังไม่มี user data
      if (localStorage.getItem('access_token') && !authState?.user) {
        logger.auth('info', 'Has token but no user data, fetching user info');
        await dispatch(checkAuthen());
        return;
      }

      // ถ้าไม่มี auth และไม่ได้อยู่หน้า login
      if (!authState?.isAuthenticated && location.pathname !== '/') {
        logger.route('info', 'Unauthorized, redirecting to login');
        navigate('/');
        return;
      }
    };

    validateAuth();
  }, [authState?.isAuthenticated, authState?.user, location.pathname]);

  // ลบ effects ที่ไม่จำเป็นออก และปรับปรุง loading check
  if (authState?.loading) {
    logger.auth('debug', 'Loading auth state');
    return <Loading />;
  }

  // ถ้าไม่มี auth และไม่ได้อยู่หน้า login ให้ redirect
  if (!authState?.isAuthenticated && location.pathname !== '/') {
    return null;
  }

  // Render layout
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
