import { useController } from "react-hook-form";
import { ConfigProvider, TimePicker } from "antd";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import thTh from "antd/locale/th_TH";
import { renderTypeError } from "..";
import { Fragment } from "react";

const TimePickerComponent = ({
  control,
  item,
  handleChange,
  setValue,
}: any) => {
  const { rules, name, defaultValue, label, disabled, ...propsInput } = item;

  const { fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const format = "HH:mm";

  const onChange = (e: any) => {
    setValue(name, null);
    handleChange(e, item);
  };

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules?.required && <Required>*</Required>}
      </RowComponent>
      <ConfigProvider locale={thTh}>
        <TimePicker
          id={name}
          disabled={disabled}
          //defaultValue={dayjs('12:08', format)}
          allowClear
          format={format}
          onChange={onChange}
          {...propsInput}
        />
      </ConfigProvider>
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedTimePicker = TimePickerComponent;
