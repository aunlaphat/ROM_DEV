export interface SubItem {
  title: string;
  key: string;
  icon?: React.ReactNode;
  link: string;
  role?: number[];
}

export interface MenuItemProps {
  title: string;
  key: string;
  icon?: React.ReactNode;
  link: string;
  role?: number[];
  subItems?: SubItem[];
}
export interface User {
  userId: string;
  password: string;
}
