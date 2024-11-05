import { useController } from "react-hook-form";
import { Row, Select as SelectAntd } from "antd";
import { TextXSMall, TextInputLabel } from "../../text";
import { Required } from "../../../style";
import { renderTypeError } from "..";
import { Fragment, useMemo } from "react";

const { Option } = SelectAntd;

const Dropdown = ({ control, item, handleChange, handleClear }: any) => {
  const {
    rules,
    name,
    defaultValue,
    placeholder,
    label,
    mode,
    disabled,
    properties,
    ...propsInput
  } = item;
  const { options, valueKey, labelKey } = properties;
  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

  const opt = useMemo(() => {
    return options.map((e: any, i: any) => {
      const val = valueKey ? e[valueKey || "value"] : e.value;
      const lab = labelKey ? e[labelKey || "label"] : e.label;
      return (
        <Option key={`${e.label}_${i + 1}`} value={val}>
          {lab}
        </Option>
      );
    });
  }, [options]);

  return (
    <Fragment>
      <Row>
        {label && <TextInputLabel text={label} />}
        {rules?.required && <Required>*</Required>}
      </Row>
      <SelectAntd
        showSearch
        allowClear
        id={name}
        value={value || []}
        placeholder={placeholder}
        optionFilterProp="children"
        onChange={(e) => handleChange(e, item)}
        onClear={() => handleClear({}, item)}
        style={{ width: "100%" }}
        mode={mode || undefined}
        disabled={disabled}
        {...propsInput}
      >
        {opt}
      </SelectAntd>
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedDropdown = Dropdown;
