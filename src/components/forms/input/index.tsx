import { useController } from "react-hook-form";
import { Input as InputAntd, InputNumber as InputNumberAntD } from "antd";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import { renderTypeError } from "..";
import { EyeInvisibleOutlined, EyeTwoTone } from "@ant-design/icons";
import { Fragment } from "react";

const InputLabel = ({ control, item, handleChange }: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    placeholder,
    disabled,
    onEnter,
    inputType = "text",
    inputStep = "0.01",
    ...propsInput
  } = item;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter") {
      const { value: v } = event.target;
      onEnter && onEnter(v, true);
    }
  };

  return (
    <Fragment key={name}>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules?.required && <Required>*</Required>}
      </RowComponent>
      <InputAntd
        id={name}
        name={name}
        value={value}
        type={inputType}
        step={inputStep}
        disabled={disabled}
        onChange={(e) => handleChange(e, item)}
        onKeyDown={(e) => handleKeyDown(e)}
        placeholder={placeholder}
        autoComplete={"off"}
        {...propsInput}
      />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const InputArea = ({ control, item, handleChange }: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    placeholder,
    disabled,
    onEnter,
    inputType = "text",
    inputStep = "0.01",
    ...propsInput
  } = item;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter") {
      const { value: v } = event.target;
      onEnter && onEnter(v, true);
    }
  };

  return (
    <Fragment key={name}>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules && rules.required && <Required>*</Required>}
      </RowComponent>
      <InputAntd.TextArea
        id={name}
        name={name}
        value={value}
        type={inputType}
        step={inputStep}
        disabled={disabled}
        onChange={(e) => handleChange(e, item)}
        onKeyDown={(e) => handleKeyDown(e)}
        placeholder={placeholder}
        autoComplete={"off"}
        {...propsInput}
      />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const InputNumber = ({ control, item, handleChange }: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    placeholder,
    disabled,
    onEnter,
    inputPrecision = 2,
    ...propsInput
  } = item;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter") {
      const { value: v } = event.target;
      onEnter && onEnter(v, true);
    }
  };

  return (
    <Fragment key={name}>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules && rules.required && <Required>*</Required>}
      </RowComponent>
      <InputNumberAntD
        style={{ width: "100%", textAlign: "right", ...item.style }}
        id={name}
        name={name}
        value={value}
        precision={inputPrecision}
        disabled={disabled}
        onChange={(e) => handleChange(e, item)}
        onKeyDown={(e) => handleKeyDown(e)}
        placeholder={placeholder}
        {...propsInput}
      />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedInputLabel = InputLabel; //memo(, compareRender);

const InputPassword = ({ control, item, handleChange }: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    placeholder,
    disabled,
    onEnter,
    ...propsInput
  } = item;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const handleKeyDown = (event: any) => {
    if (event.key === "Enter") {
      const { value: v } = event.target;
      onEnter && onEnter(v, true);
    }
  };

  return (
    <Fragment>
      <RowComponent>
        {label && <TextXSMall text={label} />}
        {rules && rules.required && <Required>*</Required>}
      </RowComponent>
      <InputAntd.Password
        iconRender={(visible) =>
          visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
        }
        id={name}
        name={name}
        value={value}
        disabled={disabled}
        onChange={(e) => handleChange(e, item)}
        onKeyDown={(e) => handleKeyDown(e)}
        placeholder={placeholder}
        autoComplete={"off"}
        {...propsInput}
      />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedInputPassword = InputPassword;
