import { Layout, Menu, Divider } from "antd";
import { Link, useLocation } from "react-router-dom";
import { useSelector } from "react-redux";
import { useMemo } from "react";
import { RoleID } from "../../../constants/roles";
import { ROUTES } from "../../../resources/routes";
import { logger } from "../../../utils/logger";
import { Icon } from "../../../resources/icon";
import logo from "../../../assets/images/logo_org.png";
import { ScrollMenu } from "../../Layouts/siderLayout/style";

const { Sider } = Layout;

const SiderLayout = ({ collapsed, collapsedWidth }: any) => {
  const location = useLocation();
  const auth = useSelector((state: any) => state.auth);
  const roleID: RoleID | undefined = auth?.user?.roleID;

  /**
   * เมนูทั้งหมดก่อนกรองตาม Role
   */
  const ROUTES_MENU = useMemo(
    () => [
      {
        title: ROUTES.ROUTE_MAIN.LABEL,
        key: ROUTES.ROUTE_MAIN.PATH,
        icon: Icon.Home(),
        link: ROUTES.ROUTE_MAIN.PATH,
        roles: [
          RoleID.ADMIN,
          RoleID.TRADE_CONSIGN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.VIEWER,
        ],
      },
      {
        title: 'สร้างรายการคืนสินค้า',
        key: 'create',
        icon: Icon.Edit1(),
        roles: [
          RoleID.ADMIN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.TRADE_CONSIGN,
          RoleID.VIEWER,
        ],
        children: [
          {
            title: ROUTES.ROUTE_CREATERETURNORDERMKP.LABEL,
            key: ROUTES.ROUTE_CREATERETURNORDERMKP.PATH,
            icon: Icon.Edit1(),
            link: ROUTES.ROUTE_CREATERETURNORDERMKP.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
          {
            title: ROUTES.ROUTE_CREATETRADERETURN.LABEL, // สร้างรายการคืนจากหน้าสาขา (Offline)
            key: ROUTES.ROUTE_CREATETRADERETURN.PATH,
            icon: Icon.Edit1(),
            link: ROUTES.ROUTE_CREATETRADERETURN.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
          {
            title: ROUTES.ROUTE_CREATEBLINDRETURN.LABEL,
            key: ROUTES.ROUTE_CREATEBLINDRETURN.PATH,
            icon: Icon.Edit1(),
            link: ROUTES.ROUTE_CREATEBLINDRETURN.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
          {
            title: ROUTES.ROUTE_IJ.LABEL,
            key: ROUTES.ROUTE_IJ.PATH,
            link: ROUTES.ROUTE_IJ.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
        ],
      },
      {
        title: 'ตรวจสอบ/ยืนยันการคืนสินค้า',
        key: 'check',
        icon: Icon.Confirm(),
        roles: [
          RoleID.ADMIN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.TRADE_CONSIGN,
          RoleID.VIEWER,
        ],
        children: [
          {
            title: ROUTES.ROUTE_DRAFTANDCONFIRM.LABEL, // รายการรออนุมัติ MKP
            key: ROUTES.ROUTE_DRAFTANDCONFIRM.PATH,
            icon: Icon.Draft(),
            link: ROUTES.ROUTE_DRAFTANDCONFIRM.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
          {
            title: ROUTES.ROUTE_CONFIRMRETURNTRADE.LABEL, // รายการรออนุมัติ Trade
            key: ROUTES.ROUTE_CONFIRMRETURNTRADE.PATH,
            icon: Icon.Confirm(),
            link: ROUTES.ROUTE_CONFIRMRETURNTRADE.PATH,
            roles: [
              RoleID.ADMIN,
              RoleID.ACCOUNTING,
              RoleID.WAREHOUSE,
              RoleID.TRADE_CONSIGN,
              RoleID.VIEWER,
            ],
          },
        ]
      },
      {
        title: ROUTES.ROUTE_IMPORTORDER.LABEL,
        key: ROUTES.ROUTE_IMPORTORDER.PATH,
        icon: Icon.BoxArrow(),
        link: ROUTES.ROUTE_IMPORTORDER.PATH,
        roles: [
          RoleID.ADMIN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.TRADE_CONSIGN,
          RoleID.VIEWER,
        ],
      },
      {
        title: 'รับเข้าสินค้าหน้าคลัง',
        key: 'return',
        icon: Icon.Return(),
        roles: [
          RoleID.ADMIN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.TRADE_CONSIGN,
          RoleID.VIEWER,
        ],
        children: [
            {
                title: ROUTES.ROUTE_SALERETURN.LABEL, // รับเข้าการคืนสินค้าของรายการสินค้าที่มีการกรอกข้อมูลเข้าระบบมา => Sale Return
                key: ROUTES.ROUTE_SALERETURN.PATH,
                icon: Icon.Return(),
                link: ROUTES.ROUTE_SALERETURN.PATH,
                roles: [
                  RoleID.ADMIN,
                  RoleID.ACCOUNTING,
                  RoleID.WAREHOUSE,
                  RoleID.TRADE_CONSIGN,
                  RoleID.VIEWER,
                ],
            },
            {
                title: ROUTES.ROUTE_OTHERRETURN.LABEL, // รับเข้าการคืนสินค้าของรายการสินค้าที่ไม่ทราบที่มา => Other Return
                key: ROUTES.ROUTE_OTHERRETURN.PATH,
                icon: Icon.Return(),
                link: ROUTES.ROUTE_OTHERRETURN.PATH,
                roles: [
                  RoleID.ADMIN,
                  RoleID.ACCOUNTING,
                  RoleID.WAREHOUSE,
                  RoleID.TRADE_CONSIGN,
                  RoleID.VIEWER,
                ],
            },
          ]
        },
        {
          title: ROUTES.ROUTE_REPORT.LABEL, // รายงาน
          key: ROUTES.ROUTE_REPORT.PATH,
          icon: Icon.Report(),
          link: ROUTES.ROUTE_REPORT.PATH,
          roles: [
            RoleID.ADMIN,
            RoleID.ACCOUNTING,
            RoleID.WAREHOUSE,
            RoleID.TRADE_CONSIGN,
            RoleID.VIEWER,
          ],
        },
        {
          title: ROUTES.ROUTE_MANAGEUSER.LABEL, // จัดการผู้ใช้งาน
          key: ROUTES.ROUTE_MANAGEUSER.PATH,
          icon: Icon.manageUser(),
          link: ROUTES.ROUTE_MANAGEUSER.PATH,
          roles: [RoleID.ADMIN],
        },
    ],
   []
  );
  
  /**
   * กรองเมนูตาม RoleID ของผู้ใช้
   */
  const filteredMenu = useMemo(() => {
    if (!roleID) return [];

    const menu = ROUTES_MENU.filter((item) => item.roles?.includes(roleID));
    logger.log("info", `🔹 Sidebar Menu for Role ${roleID}:`, menu);
    return menu;
  }, [roleID, ROUTES_MENU]);

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
          alt="Logo"
          className="avatar"
          style={{ width: "90%" }}
        />
      </div>
      <Divider />
      <ScrollMenu>
        <Menu theme="light" mode="inline" selectedKeys={[location.pathname]}>
          {filteredMenu.map((item) => (
            item.children ? (
              <Menu.SubMenu key={item.key} icon={item.icon} title={item.title}>
                {item.children.map((child) => (
                  <Menu.Item key={child.key}>
                    <Link to={child.link}>{child.title}</Link>
                  </Menu.Item>
                ))}
              </Menu.SubMenu>
            ) : (
              <Menu.Item key={item.key} icon={item.icon}>
                <Link to={item.link}>{item.title}</Link>
              </Menu.Item>
            )
          ))}
        </Menu>
      </ScrollMenu>
    </Sider>
  );
};

export default SiderLayout;
