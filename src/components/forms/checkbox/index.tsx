import { Checkbox, Col, Row } from "antd";
import { Fragment, useMemo } from "react";
import { useController } from "react-hook-form";
import { Required, RowComponent } from "../../../style";
import { TextInputLabel, TextXSMall } from "../../text";
import { renderTypeError } from "..";

const CheckboxComponent = ({ control, item, handleChange }: any) => {
  const { rules, name, defaultValue, label, properties, ...propsInput } = item;

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
    return options.map((el: any) => (
      <Checkbox key={`${name}${el.value}`} value={el.value}>
        {el.label}
      </Checkbox>
    ));
  }, [options]);

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules?.reqired && <Required>*</Required>}
      </RowComponent>
      <Checkbox.Group
        key={name}
        disabled={item.disabled || false}
        onChange={(e) => handleChange(e, item)}
        value={value}
        {...propsInput}
      >
        <Row>
          <Col span={24}>{opt}</Col>
        </Row>
      </Checkbox.Group>
      <br />
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedCheckbox = CheckboxComponent;
