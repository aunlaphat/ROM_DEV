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
      // {
      //   title: ROUTES.ROUTE_RETURNORDER.LABEL,
      //   key: ROUTES.ROUTE_RETURNORDER.PATH,
      //   icon: Icon.Return(),
      //   link: ROUTES.ROUTE_RETURNORDER.PATH,
      //   roles: [RoleID.ADMIN, RoleID.TRADE_CONSIGN, RoleID.WAREHOUSE],
      // },
      // {
      //   title: ROUTES.ROUTE_CREATERETURN.LABEL,
      //   key: ROUTES.ROUTE_CREATERETURN.PATH,
      //   icon: Icon.Edit1(),
      //   link: ROUTES.ROUTE_CREATERETURN.PATH,
      //   roles: [RoleID.ADMIN, RoleID.TRADE_CONSIGN],
      // },
      {
        title: ROUTES.ROUTE_CREATEBLINDRETURN.LABEL,
        key: ROUTES.ROUTE_CREATEBLINDRETURN.PATH,
        icon: Icon.Edit1(),
        link: ROUTES.ROUTE_CREATEBLINDRETURN.PATH,
        role: [
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
        icon: Icon.Edit1(),
        link: ROUTES.ROUTE_IJ.PATH,
        role: [
          RoleID.ADMIN,
          RoleID.ACCOUNTING,
          RoleID.WAREHOUSE,
          RoleID.TRADE_CONSIGN,
          RoleID.VIEWER,
        ],
      },
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

  // const ROUTES_MENU: MenuItemProps[] = [
  //   // {
  //   //   title: ROUTES_PATH.ROUTE_MAIN.LABEL,
  //   //   key: ROUTES_PATH.ROUTE_MAIN.PATH,
  //   //   icon: Icon.dashboard(),
  //   //   link: ROUTES_PATH.ROUTE_MAIN.PATH,
  //   //   role: [1, 2],
  //   // },
  //   // {
  //   //   title: ROUTES_PATH.ROUTE_ADJUST.LABEL,
  //   //   key: ROUTES_PATH.ROUTE_ADJUST.PATH,
  //   //   icon: Icon.dashboard(),
  //   //   link: ROUTES_PATH.ROUTE_ADJUST.PATH,
  //   //   role: [1, 2],
  //   // },
  //   // {
  //   //   title: ROUTES_PATH.ROUTE_PLATFORM.LABEL,
  //   //   key: ROUTES_PATH.ROUTE_PLATFORM.PATH,
  //   //   icon: Icon.dashboard(),
  //   //   link: ROUTES_PATH.ROUTE_PLATFORM.PATH,
  //   //   role: [1, 2],
  //   // },
  //   // {
  //   //   title: ROUTES_PATH.ROUTE_MANAGEMENT.LABEL,
  //   //   key: ROUTES_PATH.ROUTE_MANAGEMENT.PATH,
  //   //   icon: Icon.dashboard(),
  //   //   link: ROUTES_PATH.ROUTE_MANAGEMENT.PATH,
  //   //   role: [1, 2],
  //   // },
  //   // {
  //   //   title: "Sync MKP",
  //   //   key: "",
  //   //   icon: Icon.dashboard(),
  //   //   link:"",
  //   //   role: [1],
  //   //   subItems: [
  //   //     {
  //   //       title: ROUTES_PATH.ROUTE_SYNC.LABEL,
  //   //       key: ROUTES_PATH.ROUTE_SYNC.PATH,
  //   //       icon: Icon.dashboard(),
  //   //       link: ROUTES_PATH.ROUTE_SYNC.PATH,
  //   //       role: [1, 2],
  //   //     },
  //   //     {
  //   //       title: "Uplift",
  //   //       key: "/sync/subitem2",
  //   //       icon: Icon.dashboard(), // Replace with your icon
  //   //       link: "/sync/subitem2",
  //   //     },
  //   //   ],
  //   // },
  //   {
  //     title: ROUTES_PATH.ROUTE_RETURNORDER.LABEL,
  //     key: ROUTES_PATH.ROUTE_RETURNORDER.PATH,
  //     icon: Icon.Home(),
  //     link: ROUTES_PATH.ROUTE_RETURNORDER.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_CREATERETURN.LABEL,
  //     key: ROUTES_PATH.ROUTE_CREATERETURN.PATH,
  //     icon: Icon.Edit1(),
  //     link: ROUTES_PATH.ROUTE_CREATERETURN.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_CREATETRADERETURN.LABEL,
  //     key: ROUTES_PATH.ROUTE_CREATETRADERETURN.PATH,
  //     icon: Icon.Edit1(),
  //     link: ROUTES_PATH.ROUTE_CREATETRADERETURN.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_DRAFTANDCONFIRM.LABEL,
  //     key: ROUTES_PATH.ROUTE_DRAFTANDCONFIRM.PATH,
  //     icon: Icon.Check(),
  //     link: ROUTES_PATH.ROUTE_DRAFTANDCONFIRM.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_CREATEBLIND.LABEL,
  //     key: ROUTES_PATH.ROUTE_CREATEBLIND.PATH,
  //     icon: Icon.Edit1(),
  //     link: ROUTES_PATH.ROUTE_CREATEBLIND.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_CONFIRMRETURNTRADE.LABEL,
  //     key: ROUTES_PATH.ROUTE_CONFIRMRETURNTRADE.PATH,
  //     icon: Icon.Check(),
  //     link: ROUTES_PATH.ROUTE_CONFIRMRETURNTRADE.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_IMPORTORDER.LABEL,
  //     key: ROUTES_PATH.ROUTE_IMPORTORDER.PATH,
  //     icon: Icon.BoxArrow(),
  //     link: ROUTES_PATH.ROUTE_IMPORTORDER.PATH,
  //     role: [1, 2],
  //   },
    
  //   {
  //     title: ROUTES_PATH.ROUTE_SALERETURN.LABEL,
  //     key: ROUTES_PATH.ROUTE_SALERETURN.PATH,
  //     icon: Icon.Return(),
  //     link: ROUTES_PATH.ROUTE_SALERETURN.PATH,
  //     role: [1, 2],
  //   },
  //   {
  //     title: ROUTES_PATH.ROUTE_OTHER.LABEL,
  //     key: ROUTES_PATH.ROUTE_OTHER.PATH,
  //     icon: Icon.Return(),
  //     link: ROUTES_PATH.ROUTE_OTHER.PATH,
  //     role: [1, 2],
  //   },
   
  //   {
  //     title: ROUTES_PATH.ROUTE_REPORT.LABEL,
  //     key: ROUTES_PATH.ROUTE_REPORT.PATH,
  //     icon: Icon.Report(),
  //     link: ROUTES_PATH.ROUTE_REPORT.PATH,
  //     role: [1, 2],
  //   },

  //   {
  //     title: ROUTES_PATH.ROUTE_MANAGEUSER.LABEL,
  //     key: ROUTES_PATH.ROUTE_MANAGEUSER.PATH,
  //     icon: Icon.manageUser(),
  //     link: ROUTES_PATH.ROUTE_MANAGEUSER.PATH,
  //     role: [1, 2],
  //   },
  //   // {
  //   //   title: ROUTES_PATH.   ROUTE_CREATEEXPENSE.LABEL,
  //   //   key: ROUTES_PATH.  ROUTE_CREATEEXPENSE.PATH,
  //   //   icon: Icon.dashboard(),
  //   //   link: ROUTES_PATH.   ROUTE_CREATEEXPENSE.PATH,
  //   //   role: [1, 2],
  //   // },
   
  //   // {
  //   //   title: "จัดการผู้ใช้งาน",
  //   //   key: "/management_user",
  //   //   icon: Icon.users(),
  //   //   link: "",
  //   //   role: [1],
  //   //   subItems: [],
  //   // },
  // ];
  
  /**
   * กรองเมนูตาม RoleID ของผู้ใช้
   */
  const filteredMenu = useMemo(() => {
    if (!roleID) return [];

    const menu = ROUTES_MENU.filter((item) => item.roles.includes(roleID));
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
            <Menu.Item key={item.key} icon={item.icon}>
              <Link to={item.link}>{item.title}</Link>
            </Menu.Item>
          ))}
        </Menu>
      </ScrollMenu>
    </Sider>
  );
};

export default SiderLayout;
