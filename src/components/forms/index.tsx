import { Fragment } from "react";
import { MemoizedCheckbox } from "./checkbox";
import {
  InputArea,
  InputNumber,
  MemoizedInputLabel,
  MemoizedInputPassword,
} from "./input";
import { Col, Row } from "antd";
import dayjs from "dayjs";
import { MemoizedDropdown } from "./dropdown";
import { TextLabel } from "./text";
import { MemoizedUpload } from "./upload";
import { MemoizedSwitch } from "./switch";
import { MemoizedRadio } from "./radio";
import { MemoizedTimePicker } from "./timepicker";
import { MemoizedDatePicker, MemoizedDateRange } from "./datepicker";

export const RenderForm = ({
  forms = [],
  setValue,
  getValues,
  control,
  styleRow,
  onChange,
  onClear,
  renderButton,
  spanSearch,
}: any) => {
  function handleChange(e: any, item: any, option: any) {
    if (!e && item.type !== "SWITCH") {
      return;
    }
    switch (item.type) {
      case "DATE_PICKER":
        setValue(
          item.name,
          dayjs(new Date(e)).format(item.format || "DD/MM/YYYY")
        );
        break;
      case "DATE_RANGE":
        setValue(
          item.name,
          e.map((date: any) =>
            dayjs(new Date(date)).format(item.format || "DD/MM/YYYY")
          )
        );
        break;
      case "TIME_PICKER":
        setValue(item.name, dayjs(new Date(e)).format("HH:mm"));
        break;
      case "DROPDOWN":
      case "SELECT":
      case "SELECT_MULTI":
      case "SWITCH":
      case "CHECKBOX":
        setValue(item.name, e);
        break;
      case "FILE_UPLOAD":
      case "IMAGE_UPLOAD":
        setValue(item.name, e.fileList);
        break;
      case "SELECT_MODAL":
        setValue(item.name, e);
        break;
      default:
        if (e.target) {
          setValue(item.name, e.target.value);
        } else {
          setValue(item.name, e);
        }
        break;
    }

    onChange && onChange(e, item, option);
  }

  function handleClear(e: any, item: any) {
    setValue(item.name, "");
    onClear && onClear(e, item);
  }

  const rest = {
    control,
    setValue,
    getValues,
    handleChange,
    handleClear,
  };

  return (
    <Fragment>
      <Row gutter={[8, 8]} style={styleRow}>
        {forms.map((f: any, i: any) => {
          return (
            <Col
              key={`colForm${i}`}
              xs={{ span: 24 }}
              md={{ span: 24 }}
              xl={{ span: f.span }}
              lg={{ span: f.span }}
              style={{ ...f.style }}
            >
              {renderInputType(f, rest)}
            </Col>
          );
        })}
      </Row>
      <Col
        xs={{ span: 24 }}
        sm={{ span: 24 }}
        xl={{ span: spanSearch }}
        lg={{ span: spanSearch }}
        style={{ padding: 0, display: "flex", marginTop: "5%" }}
      >
        {renderButton}
      </Col>
      {/* <div style={{ marginBottom: "16px" }} /> */}
    </Fragment>
  );
};

export const renderTypeError = (item: any, error: any) => {
  if (error?.message) {
    return error.message;
  }
  switch (error.type) {
    case "required":
      return `โปรดระบุ`;
    case "pattern":
      return `รูปแบบไม่ถูกต้อง`;
    case "maxLength":
      return `ระบุไม่เกิน ${error.message} ตัวอักษร`;
    case "max":
      return `ระบุจำนวนไม่เกิน ${error.message} ตัวอักษร`;
    case "minLength":
      return `ระบุไม่น้อยกว่า ${error.message} ตัวอักษร`;
    case "min":
      return `ระบุไม่น้อยกว่า ${error.message} ตัวอักษร`;
  }
};

export function renderInputType(item: any, globalProps: any) {
  const { type } = item;
  switch (type) {
    case "TEXT":
      return <TextLabel item={{ ...item }} {...globalProps} />;
    case "INPUT":
      return <MemoizedInputLabel item={{ ...item }} {...globalProps} />;
    case "INPUT_AREA":
      return <InputArea item={{ ...item }} {...globalProps} />;
    case "INPUT_PASSWORD":
      return <MemoizedInputPassword item={{ ...item }} {...globalProps} />;
    case "INPUT_NUMBER":
      return <InputNumber item={{ ...item }} {...globalProps} />;
    case "DROPDOWN":
      return <MemoizedDropdown item={item} {...globalProps} />;
    case "UPLOAD":
      return <MemoizedUpload item={item} {...globalProps} />;
    case "RADIO":
      return <MemoizedRadio item={item} {...globalProps} />;
    case "CHECKBOX":
      return <MemoizedCheckbox item={item} {...globalProps} />;
    case "SWITCH":
      return <MemoizedSwitch item={item} {...globalProps} />;
    case "DATE_PICKER":
      return <MemoizedDatePicker item={item} {...globalProps} />;
    case "DATE_RANGE":
      return <MemoizedDateRange item={item} {...globalProps} />;
    case "TIME_PICKER":
      return <MemoizedTimePicker item={item} {...globalProps} />;
    default:
      return <div />;
  }
}
