import { ReactNode } from 'react';
import { RoleID } from '../constants/roles';

// interface สำหรับ Menu Items
export interface SubItem {
  title: string;
  key: string;
  icon?: ReactNode;
  link: string;
  roles?: RoleID[]; // เปลี่ยนจาก role เป็น roles และใช้ RoleID enum
}

export interface MenuItemProps {
  title: string;
  key: string;
  icon?: ReactNode;
  link: string;
  roles?: RoleID[]; // เปลี่ยนจาก role เป็น roles และใช้ RoleID enum
  subItems?: SubItem[];
}

// interface สำหรับ User
export interface User {
  userId: string;
  password: string;
  fullName?: string;
  roleId?: RoleID;
}

// interface สำหรับ Icon function
export interface IconFunction {
  (props?: React.SVGProps<SVGSVGElement>): JSX.Element;
}

// interface สำหรับ Icon object
export interface IconCollection {
  [key: string]: IconFunction;
}