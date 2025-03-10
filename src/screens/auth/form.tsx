export const LoginForm = () => {
  return [
    {
      key: "1",
      span: 24,
      name: "userName",
      label: "ชื่อผู้ใช้ (username)",
      type: "INPUT",
      placeholder: "กรอกชื่อผู้ใช้",
      title: "ชื่อผู้ใช้",
      rules: { required: true },
    },
    {
      key: "2",
      span: 24,
      name: "password",
      label: "รหัสผ่าน (password)",
      type: "INPUT_PASSWORD",
      placeholder: "กรอกรหัสผ่าน",
      title: "รหัสผ่าน",
      rules: {
        required: true,
        maxLength: {
          value: 15,
          message: 15,
        },
      },
    },
  ];
};
