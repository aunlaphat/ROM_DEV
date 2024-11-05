import { Fragment } from "react";
import { useController } from "react-hook-form";
import { ConfigProvider, DatePicker } from "antd";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import thTh from "antd/locale/th_TH";
import th from "dayjs/locale/th"; /**NOTE: ทำให้ datepicker แสดงเดือนภาษาไทย มันไม่ต้องใช้เป็น props แค่ประกาศก็ใช้ได้เลย */
import { renderTypeError } from "..";
import dayjs from "dayjs";

const DatePickerComponent = ({
  control,
  item,
  handleChange,
  setValue,
}: any) => {
  const { rules, name, defaultValue, label, disabled, ...propsInput } = item;

  console.log("th", th);

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

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
        <DatePicker
          id={name}
          style={{ width: "100%" }}
          format={item.format || "DD/MM/YYYY"}
          allowClear={true}
          disabled={disabled || false}
          disabledDate={item.disabledDate}
          defaultValue={
            value ? dayjs(value, item.format || "DD/MM/YYYY") : null
          }
          onChange={(e) => onChange(e)}
          {...propsInput}
        />
      </ConfigProvider>
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedDatePicker = DatePickerComponent; //memo(, compareRender);

const DateRangeComponent = ({ control, item, handleChange, setValue }: any) => {
  const { rules, name, defaultValue, label, disabled, ...propsInput } = item;

  const { field, fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;
  const { value } = field;

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
        <DatePicker.RangePicker
          id={name}
          style={{ width: "100%" }}
          format={item.format || "DD/MM/YYYY"}
          allowClear={true}
          disabled={disabled || false}
          onChange={(e) => onChange(e)}
          value={value || ""}
          {...propsInput}
        />
      </ConfigProvider>
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedDateRange = DateRangeComponent;
