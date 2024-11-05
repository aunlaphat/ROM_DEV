import { useController } from "react-hook-form";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import { Radio, Space } from "antd";
import { renderTypeError } from "..";
import { Fragment, useMemo } from "react";

const RadioComponent = ({ control, item, handleChange }: any) => {
  const {
    rules,
    name,
    defaultValue,
    label,
    disabled,
    properties,
    ...propsInput
  } = item;

  const { options } = properties;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const opt = useMemo(() => {
    return options.map((e: any) => {
      return (
        <Radio key={`radio${name}_${e.value}`} value={e.value}>
          {e.label}
        </Radio>
      );
    });
  }, [options]);

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules && rules.required && <Required>*</Required>}
      </RowComponent>
      <Radio.Group
        id={name}
        disabled={disabled}
        value={value}
        optionType={properties.optionType || "default"}
        buttonStyle={properties.buttonStyle || "outline"}
        onChange={(e) => handleChange(e, item)}
        {...propsInput}
      >
        <Space direction={properties.direction || "horizontal"}>{opt}</Space>
      </Radio.Group>
      <br />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedRadio = RadioComponent;
