import { useController } from "react-hook-form";
import { RowComponent } from "../../../style";
import { TextSmall, TextInputLabel } from "../../text";
import { Fragment } from "react";

export const TextLabel = ({ control, item }: any) => {
  const { name, defaultValue, label, ...propsInput } = item;

  const { field } = useController({
    control,
    name,
    defaultValue,
  });
  const { value } = field;

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} key={name} />}
      </RowComponent>
      <TextSmall
        text={value}
        style={{ marginLeft: `20px`, color: `#585858` }}
        {...propsInput}
      />
    </Fragment>
  );
};
