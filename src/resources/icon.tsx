import { UnlockOutlined } from "@ant-design/icons";
import { SVGProps } from "react";
import {
  BsFillCloudArrowDownFill,
  BsSearch,
  BsFillPenFill,
  BsFillTrashFill,
  BsBookHalf,
  BsFillXCircleFill,
  BsSave2Fill,
  BsPrinterFill,
  BsFillCheckCircleFill,
  BsLockFill,
  BsArrowLeftCircleFill,
  BsFillPlusCircleFill,
  BsArrowClockwise,
  BsShieldFillCheck,
  BsPersonCheckFill,
  BsKanban,
  BsJournalBookmarkFill,
  BsLayoutSidebarInsetReverse,
  BsFileEarmarkText,
  BsHr,
  BsQrCode,
  BsSliders2Vertical,
  BsBoxSeam,
  BsGraphUp,
  BsBarChartLine,
  BsClipboardData,
  BsKanbanFill,
} from "react-icons/bs";
import { RiLogoutBoxRFill, RiLoginBoxFill } from "react-icons/ri";
import {
  IoSettingsOutline,
  IoBookmarks,
  IoAlertCircleSharp,
  IoCheckbox,
} from "react-icons/io5";
import { HiDotsCircleHorizontal } from "react-icons/hi";
import {
  FaChevronCircleDown,
  FaExternalLinkAlt,
  FaMap,
  FaMapMarkedAlt,
  FaCamera,
  FaRegUser,
  FaUserCog,
  FaUserPlus,
  FaRegEdit,
} from "react-icons/fa";
import { TbUserCog, TbTruckReturn, TbReportAnalytics } from "react-icons/tb";
import {
  MdAutorenew,
  MdWavingHand,
  MdOutlineClose,
  MdDashboard,
  MdAnalytics,
  MdInventory2,
  MdWarehouse,
  MdPendingActions,
} from "react-icons/md";
import { BiBarcode, BiCheckCircle } from "react-icons/bi";
import { FaBoxesPacking } from "react-icons/fa6";
import { SlCheck } from "react-icons/sl";
import { HiOutlineArchiveBox } from "react-icons/hi2";
import { HiOutlineArrowDownOnSquare } from "react-icons/hi2";
import { GrView, GrNext, GrHomeRounded } from "react-icons/gr";
import { GoHome } from "react-icons/go";
import {
  BsCheckCircle,
  BsBoxArrowInLeft,
  BsArrowDownSquare,
} from "react-icons/bs";
import { VscGoToFile } from "react-icons/vsc";
import {
  DashboardOutlined,
  AppstoreOutlined,
  ProfileOutlined,
  PieChartOutlined,
  RiseOutlined,
  FundOutlined,
  ShopOutlined,
  TeamOutlined,
  UserOutlined,
  DatabaseOutlined,
  CloudSyncOutlined,
  ApartmentOutlined,
  BarcodeOutlined,
  CrownOutlined,
  DeploymentUnitOutlined,
} from "@ant-design/icons";

// สร้าง Type สำหรับ Icon Function
export type IconFunction = (
  props?: SVGProps<SVGSVGElement> & Record<string, any>
) => JSX.Element;

// สร้าง Type สำหรับ Icon Object
export interface IconDictionary {
  [key: string]: IconFunction;
}

// สร้าง Type สำหรับ Icon Keys
export type IconKey =
  | "download"
  | "search"
  | "Report"
  | "Box1"
  | "Return"
  | "BoxArrow"
  | "edit"
  | "GoToFile"
  | "Check"
  | "Edit1"
  | "Home"
  | "remove"
  | "view"
  | "clear"
  | "cancel"
  | "logout"
  | "login"
  | "save"
  | "print"
  | "confirm"
  | "forgetPassword"
  | "back"
  | "create"
  | "reAct"
  | "recheck"
  | "register"
  | "settings"
  | "threeDots"
  | "downCircle"
  | "bookmark"
  | "alert"
  | "checkboxSquare"
  | "wavingHand"
  | "external"
  | "dashboard"
  | "overview"
  | "analytics"
  | "reports"
  | "statistics"
  | "trending"
  | "performance"
  | "warehouse"
  | "inventory"
  | "pending"
  | "tracking"
  | "management"
  | "database"
  | "sync"
  | "team"
  | "profile"
  | "admin"
  | "department"
  | "chart"
  | "analytics2"
  | "barcode"
  | "UnlockOutlined"
  | "WorkSheet"
  | "Box"
  | "BarCode"
  | "Scanner"
  | "QRcode"
  | "Slider"
  | "BoxInvent"
  | "Renew"
  | "Map"
  | "MapMark"
  | "viewdetail"
  | "close"
  | "camera"
  | "next"
  | "users"
  | "manageUser"
  | "addUser"
  | "Draft"
  | "Confirm";

// เพิ่มฟังก์ชันสำหรับตรวจสอบว่า key เป็น IconKey หรือไม่
export function isIconKey(key: string): key is IconKey {
  return key in Icon;
}

/**
 * React Icons --> https://react-icons.github.io/react-icons/
 */

export const Icon: Record<IconKey, IconFunction> = {
  download: (props) => <BsFillCloudArrowDownFill {...(props as any)} />,
  search: (props) => <BsSearch {...(props as any)} />,
  Report: (props) => <TbReportAnalytics {...(props as any)} />,
  Box1: (props) => <HiOutlineArchiveBox {...(props as any)} />,
  Return: (props) => <TbTruckReturn {...(props as any)} />,
  BoxArrow: (props) => <BsArrowDownSquare {...(props as any)} />,
  edit: (props) => <BsFillPenFill {...(props as any)} />,
  GoToFile: (props) => <VscGoToFile {...(props as any)} />,
  Check: (props) => <BsCheckCircle {...(props as any)} />,
  Edit1: (props) => <FaRegEdit {...(props as any)} />,
  Home: (props) => <GrHomeRounded {...(props as any)} />,
  remove: (props) => <BsFillTrashFill {...(props as any)} />,
  view: (props) => <BsBookHalf {...(props as any)} />,
  clear: (props) => <BsFillXCircleFill {...(props as any)} />,
  cancel: (props) => <BsFillXCircleFill {...(props as any)} />,
  logout: (props) => <RiLogoutBoxRFill {...(props as any)} />,
  login: (props) => <RiLoginBoxFill {...(props as any)} />,
  save: (props) => <BsSave2Fill {...(props as any)} />,
  print: (props) => <BsPrinterFill {...(props as any)} />,
  confirm: (props) => <BsFillCheckCircleFill {...(props as any)} />,
  forgetPassword: (props) => <BsLockFill {...(props as any)} />,
  back: (props) => <BsArrowLeftCircleFill {...(props as any)} />,
  create: (props) => <BsFillPlusCircleFill {...(props as any)} />,
  reAct: (props) => <BsArrowClockwise {...(props as any)} />,
  recheck: (props) => <BsShieldFillCheck {...(props as any)} />,
  register: (props) => <BsPersonCheckFill {...(props as any)} />,
  settings: (props) => <IoSettingsOutline {...(props as any)} />,
  threeDots: (props) => <HiDotsCircleHorizontal {...(props as any)} />,
  downCircle: (props) => <FaChevronCircleDown {...(props as any)} />,
  bookmark: (props) => <IoBookmarks {...(props as any)} />,
  alert: (props) => <IoAlertCircleSharp {...(props as any)} />,
  checkboxSquare: (props) => <IoCheckbox {...(props as any)} />,
  wavingHand: (props) => <MdWavingHand {...(props as any)} />,
  external: (props) => <FaExternalLinkAlt {...(props as any)} />,
  /** MENU ICON */
  dashboard: (props) => <DashboardOutlined {...(props as any)} />,
  overview: (props) => <AppstoreOutlined {...(props as any)} />,
  analytics: (props) => <MdAnalytics {...(props as any)} />,
  reports: (props) => <BsBarChartLine {...(props as any)} />,
  statistics: (props) => <PieChartOutlined {...(props as any)} />,
  trending: (props) => <RiseOutlined {...(props as any)} />,
  performance: (props) => <FundOutlined {...(props as any)} />,
  warehouse: (props) => <MdWarehouse {...(props as any)} />,
  inventory: (props) => <MdInventory2 {...(props as any)} />,
  pending: (props) => <MdPendingActions {...(props as any)} />,
  tracking: (props) => <BsKanbanFill {...(props as any)} />,
  management: (props) => <ApartmentOutlined {...(props as any)} />,
  database: (props) => <DatabaseOutlined {...(props as any)} />,
  sync: (props) => <CloudSyncOutlined {...(props as any)} />,
  team: (props) => <TeamOutlined {...(props as any)} />,
  profile: (props) => <UserOutlined {...(props as any)} />,
  admin: (props) => <CrownOutlined {...(props as any)} />,
  department: (props) => <DeploymentUnitOutlined {...(props as any)} />,
  chart: (props) => <BsGraphUp {...(props as any)} />,
  analytics2: (props) => <BsClipboardData {...(props as any)} />,
  barcode: (props) => <BarcodeOutlined {...(props as any)} />,
  UnlockOutlined: (props) => <UnlockOutlined {...(props as any)} />,
  WorkSheet: (props) => <BsFileEarmarkText {...(props as any)} />,
  Box: (props) => <FaBoxesPacking {...(props as any)} />,
  BarCode: (props) => <BiBarcode {...(props as any)} />,
  Scanner: (props) => <BsHr {...(props as any)} />,
  QRcode: (props) => <BsQrCode {...(props as any)} />,
  Slider: (props) => <BsSliders2Vertical {...(props as any)} />,
  BoxInvent: (props) => <BsBoxSeam {...(props as any)} />,
  Renew: (props) => <MdAutorenew {...(props as any)} />,
  Map: (props) => <FaMap {...(props as any)} />,
  MapMark: (props) => <FaMapMarkedAlt {...(props as any)} />,
  viewdetail: (props) => <GrView {...(props as any)} />,
  close: (props) => <MdOutlineClose {...(props as any)} />,
  camera: (props) => <FaCamera {...(props as any)} />,
  next: (props) => <GrNext {...(props as any)} />,
  users: (props) => <FaRegUser {...(props as any)} />,
  manageUser: (props) => <FaUserCog {...(props as any)} />,
  addUser: (props) => <FaUserPlus {...(props as any)} />,
  Draft: (props) => <BsFileEarmarkText {...(props as any)} />,
  Confirm: (props) => <BsCheckCircle {...(props as any)} />,
};