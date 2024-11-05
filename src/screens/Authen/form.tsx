export const LoginForm = () => {
  return [
    {
      key: "1",
      span: 24,
      name: "username",
      label: "ชื่อผู้ใช้ (username)",
      type: "INPUT",
      rules: { required: true },
    },
    {
      key: "2",
      span: 24,
      name: "password",
      label: "รหัสผ่าน",
      type: "INPUT_PASSWORD",
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
