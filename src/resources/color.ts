// src/resources/color.ts (เปลี่ยนจาก .js เป็น .ts)
export const color = {
  theme: "#2c4e98",
  secondTheme: "#f7941d",
  edit: "#ff7b54",
  submit: "#7eca9c",
  clear: "#707070",
  remove: "#d35d6e",
  search: "#1890ff",
  red: "#FF0000",
  logout: "#737373",
  yellow: "#F4D35E",
  reset: "#ffc107",
};

// ส่งออก type สำหรับ color keys
export type ColorKeys = keyof typeof color;