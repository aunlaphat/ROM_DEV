import React, { Fragment, useState } from "react";
import { PlusOutlined } from "@ant-design/icons";
import { useController } from "react-hook-form";
import { Button, Upload, Modal } from "antd";
import { Required, RowComponent } from "../../../style";
import { TextXSMall, TextInputLabel } from "../../text";
import { UploadOutlined } from "@ant-design/icons";
import { renderTypeError } from "..";

const UploadFiles = ({ control, item, setValue, getValues }: any) => {
  const { rules, name, defaultValue, label, disabled, ...propsInput } = item;
  const [previewOpen, setPreviewOpen] = useState(false);
  const [previewImage, setPreviewImage] = useState("");
  const [previewTitle, setPreviewTitle] = useState("");

  const { fieldState } = useController({
    control,
    name,
    rules,
    defaultValue,
  });
  const { error } = fieldState;

  const onChange = (event: any) => {
    setValue(name, event.fileList);
  };

  const getBase64 = (file: any) =>
    new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });

  const handleCancel = () => setPreviewOpen(false);
  const handlePreview = async (file: any) => {
    if (!file.url && !file.preview) {
      file.preview = await getBase64(file.originFileObj);
    }
    setPreviewImage(file.url || file.preview);
    setPreviewOpen(true);
    setPreviewTitle(
      file.name || file.url.substring(file.url.lastIndexOf("/") + 1)
    );
  };

  const uploadButton = (
    <div>
      <PlusOutlined />
      <div
        style={{
          marginTop: 8,
        }}
      >
        Upload
      </div>
    </div>
  );

  return (
    <Fragment>
      <RowComponent>
        {label && <TextInputLabel text={label} />}
        {rules?.required && <Required>*</Required>}
      </RowComponent>
      <Upload
        id={name}
        beforeUpload={() => false}
        listType={item.listType || "picture"} /** picture-card, picture */
        disabled={disabled || false}
        maxCount={item.maxCount || 1}
        multiple={item.maxCount > 1 || false}
        accept={item.accept || "*"}
        onPreview={handlePreview}
        fileList={getValues(name)}
        onChange={onChange}
        {...propsInput}
      >
        {item.listType !== "picture-card" ? (
          <Button icon={<UploadOutlined />}>Upload</Button>
        ) : (
          uploadButton
        )}
      </Upload>
      <Modal
        destroyOnClose={true}
        open={previewOpen}
        title={previewTitle}
        footer={null}
        onCancel={handleCancel}
      >
        <img
          alt="example"
          style={{
            width: "100%",
          }}
          src={previewImage}
        />
      </Modal>
      {error && <TextXSMall text={renderTypeError(item, error)} color="red" />}
    </Fragment>
  );
};

export const MemoizedUpload = UploadFiles; //memo(, compareRender);
