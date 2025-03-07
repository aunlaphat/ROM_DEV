import { Spin, notification } from "antd";
import { Fragment, useState, useEffect } from "react";
import { CardTheme } from "../../components/card/cardTheme";
import { TextLogoLogin, TextXSMall } from "../../components/text";
import imgLogin from "../../assets/images/login.png";
import imgLogo from "../../assets/images/logo.jpeg";
import { RenderForm } from "../../components/forms";
import { useForm } from "react-hook-form";
import { LoginForm, LoginFormValues, loginValidationRules } from "./form";
import { ContainerButton } from "../../style";
import { ButtonTheme } from "../../components/buttons";
import { useAuth } from "../../hooks/auth"; // ปรับใช้ path ตามที่ refactor แล้ว
import { logger } from "../../utils/logger";

export const Login = () => {
  // State
  const [disabled, setDisabled] = useState(true);
  const [loginError, setLoginError] = useState<string | null>(null);
  
  // Custom hooks
  const { login, loading } = useAuth(); // ไม่ใช้ error จาก useAuth แล้ว
  
  // React Hook Form without Zod
  const { 
    setValue, 
    getValues, 
    control, 
    formState, 
    trigger,
    register
  } = useForm<LoginFormValues>({
    mode: 'onChange',
    defaultValues: {
      username: '',
      password: ''
    }
  });

  // Effect to handle button disable state based on form validation
  useEffect(() => {
    const { isValid, isDirty } = formState;
    setDisabled(!(isValid && isDirty));
  }, [formState]);

  // Effect to display error notifications
  useEffect(() => {
    if (loginError) {
      notification.error({
        message: 'เข้าสู่ระบบไม่สำเร็จ',
        description: loginError,
      });
      
      // ล้าง error หลังจากแสดงแล้ว
      setTimeout(() => setLoginError(null), 3000);
    }
  }, [loginError]);

  // Handler when form fields change
  const onChange = async () => {
    await trigger(); // Trigger validation
    const { username, password } = getValues();
    const d = !username || !password;
    setDisabled(d);
  };

  // Submit handler
  const onSubmit = async () => {
    try {
      const { username, password } = getValues();
      logger.log('info', `Login Initiated`, {
        username,
        timestamp: new Date().toISOString()
      });
      
      await login({ username, password });
    } catch (error: any) {
      // จัดการ error ด้วย state ภายใน component แทน
      setLoginError(error?.message || 'เกิดข้อผิดพลาดในการเข้าสู่ระบบ');
      logger.error('Login Error:', error);
    }
  };

  // Component for copyright text
  const CopyrightText = () => (
    <TextXSMall
      text={`© ${new Date().getFullYear()} All Rights Reserved.`}
      align={"center"}
      color="grey"
    />
  );

  // Component for login action buttons
  const LoginActions = () => (
    <div style={{ width: "100%", marginTop: -20 }}>
      <div style={{ display: "flex" }}>
        <ContainerButton align={"center"}>
          <ButtonTheme
            useFor="FORGET_PASSWORD"
            type="submit"
            onClick={() => window.location.replace("https://dcapp.site/re-password.php")}
          />
          <ButtonTheme
            useFor="LOGIN"
            type="submit"
            onClick={onSubmit}
            disabled={disabled}
            loading={loading}
          />
        </ContainerButton>
      </div>
      <CopyrightText />
    </div>
  );

  // Component for the login logo
  const LoginLogo = () => (
    <>
      <div style={{ display: "flex", justifyContent: "center" }}>
        <img alt="Company Logo" src={imgLogo} width={300} />
      </div>
      <TextLogoLogin
        text={"Return Order Management System"}
        bold
        size={"28px"}
        align={"center"}
        color="#00477d"
      />
    </>
  );

  return (
    <Spin spinning={loading} delay={300}>
      <div className="bg-login">
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            minHeight: "100vh",
          }}
        >
          <div className="img-login">
            <img alt="Login Illustration" src={imgLogin} width={420} />
          </div>
          <div style={{ margin: "auto", width: "500px", padding: "10px" }}>
            <CardTheme
              className="card-login"
              content={
                <Fragment>
                  <LoginLogo />
                  <form onSubmit={(e) => { e.preventDefault(); onSubmit(); }}>
                    <RenderForm
                      control={control}
                      setValue={setValue}
                      getValues={getValues}
                      forms={LoginForm()}
                      onChange={onChange}
                      renderButton={<LoginActions />}
                    />
                  </form>
                </Fragment>
              }
            />
          </div>
        </div>
      </div>
    </Spin>
  );
};

export default Login;