import { useController } from "react-hook-form";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import { Switch } from "antd";
import { renderTypeError } from "..";
import { Fragment, useEffect } from "react";

const SwitchComponent = ({
  control,
  item,
  handleChange,
  setValue,
  getValues,
}: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    disabled,
    propperties,
    ...propsInput
  } = item;

  const { fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;

  useEffect(() => {
    setValue(name, item.defaultValue || false);
  }, []);

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules?.required && <Required>*</Required>}
      </RowComponent>
      <Switch
        id={name}
        disabled={disabled}
        {...propperties}
        defaultChecked={item.defaultValue || false}
        onChange={(e) => handleChange(e, item)}
        {...propsInput}
      />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedSwitch = SwitchComponent; //memo(, compareRender);
