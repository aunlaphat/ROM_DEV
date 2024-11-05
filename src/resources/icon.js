import { UnlockOutlined } from "@ant-design/icons";
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
} from "react-icons/fa";
import { TbUserCog,TbTruckReturn,TbReportAnalytics } from "react-icons/tb";
import { MdAutorenew, MdWavingHand, MdOutlineClose } from "react-icons/md";
import { BiBarcode ,BiCheckCircl,} from "react-icons/bi";
import { FaBoxesPacking } from "react-icons/fa6";
import { SlCheck } from "react-icons/sl";
import { HiOutlineArchiveBox } from "react-icons/hi2";
import { HiOutlineArrowDownOnSquare } from "react-icons/hi2";
import { GrView, GrNext,GrHomeRounded } from "react-icons/gr";
import { FaRegEdit,HiOutlineHome } from "react-icons/fa";
import { GoHome } from "react-icons/go";
import { BsCheckCircle,BsBoxArrowInLeft,BsArrowDownSquare ,} from "react-icons/bs";
import { VscGoToFile } from "react-icons/vsc";
/**
 * React Icons --> https://react-icons.github.io/react-icons/
 */
export const Icon = {
 
  download: (props) => <BsFillCloudArrowDownFill {...props} />,
  search: (props) => <BsSearch {...props} />,
  Report: (props) => <TbReportAnalytics {...props} />,
  Box1: (props) => <HiOutlineArchiveBox {...props} />,
  Return: (props) => <TbTruckReturn {...props} />,
  BoxArrow: (props) => <BsArrowDownSquare {...props} />,
  edit: (props) => <BsFillPenFill {...props} />,
  GoToFile: (props) => <VscGoToFile {...props} />,
  Check: (props) => <BsCheckCircle   {...props} />,
  Edit1:(props) => <FaRegEdit {...props} />,
  Home:(props) => <GrHomeRounded {...props} />,
  remove: (props) => <BsFillTrashFill {...props} />,
  view: (props) => <BsBookHalf {...props} />,
  clear: (props) => <BsFillXCircleFill {...props} />,
  cancel: (props) => <BsFillXCircleFill {...props} />,
  logout: (props) => <RiLogoutBoxRFill {...props} />,
  login: (props) => <RiLoginBoxFill {...props} />,
  save: (props) => <BsSave2Fill {...props} />,
  print: (props) => <BsPrinterFill {...props} />,
  confirm: (props) => <BsFillCheckCircleFill {...props} />,
  forgetPassword: (props) => <BsLockFill {...props} />,
  back: (props) => <BsArrowLeftCircleFill {...props} />,
  create: (props) => <BsFillPlusCircleFill {...props} />,
  reAct: (props) => <BsArrowClockwise {...props} />,
  recheck: (props) => <BsShieldFillCheck {...props} />,
  register: (props) => <BsPersonCheckFill {...props} />,
  settings: (props) => <IoSettingsOutline {...props} />,
  threeDots: (props) => <HiDotsCircleHorizontal {...props} />,
  downCircle: (props) => <FaChevronCircleDown {...props} />,
  bookmark: (props) => <IoBookmarks {...props} />,
  alert: (props) => <IoAlertCircleSharp {...props} />,
  checkboxSquare: (props) => <IoCheckbox {...props} />,
  wavingHand: (props) => <MdWavingHand {...props} />,
  external: (props) => <FaExternalLinkAlt {...props} />,
  /** MENU ICON */
  dashboard: (props) => <BsKanban {...props} />,
  todo: (props) => <BsJournalBookmarkFill {...props} />,
  formExample: (props) => <BsLayoutSidebarInsetReverse {...props} />,
  UnlockOutlined: (props) => <UnlockOutlined {...props} />,
  WorkSheet: (props) => <BsFileEarmarkText {...props} />,
  Box: (props) => <FaBoxesPacking {...props} />,
  BarCode: (props) => <BiBarcode {...props} />,
  Scanner: (props) => <BsHr {...props} />,
  QRcode: (props) => <BsQrCode {...props} />,
  Slider: (props) => <BsSliders2Vertical {...props} />,
  BoxInvent: (props) => <BsBoxSeam {...props} />,
  Renew: (props) => <MdAutorenew {...props} />,
  Map: (props) => <FaMap {...props} />,
  MapMark: (props) => <FaMapMarkedAlt {...props} />,
  viewdetail: (props) => <GrView {...props} />,
  close: (props) => <MdOutlineClose {...props} />,
  camera: (props) => <FaCamera {...props} />,
  next: (props) => <GrNext {...props} />,
  users: (props) => <FaRegUser {...props} />,
  manageUser: (props) => <FaUserCog {...props} />,
  addUser: (props) => <FaUserPlus {...props} />,
};
