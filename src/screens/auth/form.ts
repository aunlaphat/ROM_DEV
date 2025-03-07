// Type สำหรับข้อมูล login
export interface LoginFormValues {
  username: string;
  password: string;
}

// กำหนดกฎสำหรับ validation ให้ใช้ใน React Hook Form
export const loginValidationRules = {
  username: {
    required: { value: true, message: 'กรุณากรอกชื่อผู้ใช้' },
  },
  password: {
    required: { value: true, message: 'กรุณากรอกรหัสผ่าน' },
    maxLength: { value: 15, message: 'รหัสผ่านต้องไม่เกิน 15 ตัวอักษร' }
  }
};

// ข้อมูลสำหรับสร้าง form fields
export const LoginForm = () => {
  return [
    {
      key: "1",
      span: 24,
      name: "username",
      label: "ชื่อผู้ใช้ (username)",
      type: "INPUT",
      placeholder: "กรุณากรอกชื่อผู้ใช้",
      title: "ชื่อผู้ใช้",
      rules: { required: true },
      autoComplete: "username",
    },
    {
      key: "2",
      span: 24,
      name: "password",
      label: "รหัสผ่าน (password)",
      type: "INPUT_PASSWORD",
      placeholder: "กรุณากรอกรหัสผ่าน",
      title: "รหัสผ่าน",
      rules: {
        required: true,
        maxLength: {
          value: 15,
          message: 'รหัสผ่านต้องไม่เกิน 15 ตัวอักษร',
        },
      },
      autoComplete: "current-password",
    },
  ];
};