import { Spin } from "antd";
import { Fragment, useEffect, useState } from "react";
import { CardTheme } from "../../components/card/cardTheme";
import { TextLogoLogin, TextXSMall } from "../../components/text";
import imgLogin from "../../assets/images/login.png";
import imgLogo from "../../assets/images/logo.jpeg";
import { RenderForm } from "../../components/forms";
import { useForm } from "react-hook-form";
import { LoginForm } from "./form";
import { ContainerButton } from "../../style";
import { ButtonTheme } from "../../components/buttons";
import { useAuthLogin } from "../../hooks/useAuth";

const LoginScene = () => {
  const [disabled, setDisabled] = useState(true);
  const [loading, setLoading] = useState(false);
  const { onLogin } = useAuthLogin();

  const { setValue, getValues, control } = useForm();

  const onChange = () => {
    const { username, password } = getValues();
    const d = !username || !password;
    setDisabled(d);
  };

  const onSubmit = async () => {
    try {
      const { username, password } = getValues();
      onLogin({ username, password });
      setLoading(true);
    } catch (error) {
      console.log(error);
    }
  };

  // useEffect(() => {
  //   checkLoginToken();
  // }, []);

  return (
    <Spin spinning={loading} delay={500}>
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
            <img alt={imgLogin} src={imgLogin} width={420} />
          </div>
          <div style={{ margin: "auto", width: "500px", padding: "10px" }}>
            <CardTheme
              className="card-login"
              content={
                <Fragment>
                  <div style={{ display: "flex", justifyContent: "center" }}>
                    <img alt={imgLogo} src={imgLogo} width={300} />
                  </div>
                  <TextLogoLogin
                    text={"Boiler Page System"}
                    bold
                    size={"28px"}
                    align={"center"}
                    color="#00477d"
                  />

                  <form>
                    <RenderForm
                      control={control}
                      setValue={setValue}
                      getValues={getValues}
                      forms={LoginForm()}
                      onChange={onChange}
                      renderButton={
                        <div style={{ width: "100%", marginTop: -20 }}>
                          <div style={{ display: "flex" }}>
                            <ContainerButton align={"center"}>
                              <ButtonTheme
                                useFor="FORGET_PASSWORD"
                                type="submit"
                                onClick={() =>
                                  window.location.replace(
                                    "https://dcapp.site/re-password.php"
                                  )
                                }
                              />
                              <ButtonTheme
                                useFor="LOGIN"
                                type="submit"
                                onClick={onSubmit}
                                disabled={disabled}
                              />
                            </ContainerButton>
                          </div>
                          <TextXSMall
                            text={`Boiler Page Digitalcommerce`}
                            align={"center"}
                            color="grey"
                          />
                        </div>
                      }
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

export default LoginScene;
