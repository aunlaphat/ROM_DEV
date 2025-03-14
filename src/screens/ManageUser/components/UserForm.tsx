// src/screens/ManageUser/components/UserForm.tsx
import React, { useEffect, useState } from "react";
import { Modal, Form, Input, Select, Typography, Tag, Tooltip } from "antd";
import { UserOutlined, TagsOutlined, BankOutlined } from "@ant-design/icons";
import {
  UserResponse,
  RoleResponse,
  WarehouseResponse,
  AddUserRequest,
  EditUserRequest,
} from "../../../redux/users/types";

const { Option } = Select;
const { Text } = Typography;

interface UserFormProps {
  visible: boolean;
  user: UserResponse | null;
  roles: RoleResponse[];
  warehouses: WarehouseResponse[];
  onSave: (values: AddUserRequest | EditUserRequest) => void;
  onCancel: () => void;
  confirmLoading?: boolean;
}

/**
 * Modal Form สำหรับเพิ่มหรือแก้ไขข้อมูลผู้ใช้
 */
const UserForm: React.FC<UserFormProps> = ({
  visible,
  user,
  roles,
  warehouses,
  onSave,
  onCancel,
  confirmLoading = false,
}) => {
  const [form] = Form.useForm();
  const isEditMode = !!user;

  // State สำหรับเก็บค่า ID ที่เลือก
  const [roleID, setRoleID] = useState<number | null>(null);
  const [warehouseID, setWarehouseID] = useState<string | number | null>(null);

  // Mapping functions สำหรับแปลง ID เป็นชื่อ และกลับกัน
  const getRoleName = (id: number): string => {
    const role = roles.find((r) => r.roleID === id);
    return role ? role.roleName : `บทบาท ${id}`;
  };

  const getWarehouseName = (id: string | number): string => {
    const warehouse = warehouses.find(
      (w) => String(w.warehouseID) === String(id)
    );
    return warehouse ? warehouse.warehouseName : `คลังสินค้า ${id}`;
  };

  const findRoleIDByName = (name: string): number | undefined => {
    const role = roles.find((r) => r.roleName === name);
    return role?.roleID;
  };

  const findWarehouseIDByName = (name: string): string | number | undefined => {
    const warehouse = warehouses.find((w) => w.warehouseName === name);
    return warehouse?.warehouseID;
  };

  // เลือกสีของ Tag สำหรับแต่ละบทบาท
  const getRoleColor = (roleID: number): string => {
    switch (roleID) {
      case 1:
        return "magenta"; // Admin
      case 2:
        return "blue"; // Accounting
      case 3:
        return "green"; // Warehouse
      case 4:
        return "orange"; // Trade
      case 5:
        return "purple"; // Viewer
      default:
        return "default";
    }
  };

  // Reset form และตั้งค่าเริ่มต้นเมื่อ modal เปิดหรือผู้ใช้เปลี่ยน
  useEffect(() => {
    if (visible) {
      form.resetFields();

      if (user) {
        // กำหนดค่า state
        setRoleID(user.roleID);
        setWarehouseID(user.warehouseID);

        // กำหนดค่าเริ่มต้นให้ form
        form.setFieldsValue({
          userID: user.userID,
          roleName: user.roleName, // ใช้ชื่อแทน ID
          warehouseName: user.warehouseName, // ใช้ชื่อแทน ID
        });
      } else {
        // รีเซ็ต state เมื่อเปิด modal ใหม่
        setRoleID(null);
        setWarehouseID(null);

        form.setFieldsValue({
          userID: "",
          roleName: undefined,
          warehouseName: undefined,
        });
      }
    }
  }, [visible, user, form]);

  // อัพเดต handleSubmit เพื่อตรวจสอบให้แน่ใจว่า roleID และ warehouseID เป็นปัจจุบัน
  const handleSubmit = async () => {
    try {
      // ตรวจสอบฟิลด์ที่จำเป็น
      await form.validateFields(["userID", "roleName", "warehouseName"]);

      const userIDValue = form.getFieldValue("userID");
      const roleNameValue = form.getFieldValue("roleName");
      const warehouseNameValue = form.getFieldValue("warehouseName");

      console.log("Form values before submit:", {
        userID: userIDValue,
        roleName: roleNameValue,
        warehouseName: warehouseNameValue,
      });

      // ค้นหา ID จากชื่อโดยตรงอีกครั้ง เพื่อให้แน่ใจว่าเป็นค่าล่าสุด
      const selectedRole = roles.find((r) => r.roleName === roleNameValue);
      const selectedWarehouse = warehouses.find(
        (w) => w.warehouseName === warehouseNameValue
      );

      if (!selectedRole || !selectedWarehouse) {
        console.error("Cannot find role or warehouse:", {
          roleNameValue,
          warehouseNameValue,
          selectedRole,
          selectedWarehouse,
        });
        return;
      }

      const finalRoleID = selectedRole.roleID;
      const finalWarehouseID = selectedWarehouse.warehouseID;

      console.log("Final IDs for submission:", {
        roleID: finalRoleID,
        warehouseID: finalWarehouseID,
      });

      if (isEditMode) {
        // สำหรับการแก้ไข
        const editRequest: EditUserRequest = {
          userID: userIDValue,
          roleID: Number(finalRoleID),
          warehouseID: Number(finalWarehouseID),
        };

        console.log("Edit User Request:", editRequest);
        onSave(editRequest);
      } else {
        // สำหรับการเพิ่ม
        const addRequest: AddUserRequest = {
          userID: userIDValue,
          roleID: Number(finalRoleID),
          warehouseID: Number(finalWarehouseID),
        };

        console.log("Add User Request:", addRequest);
        onSave(addRequest);
      }
    } catch (error) {
      console.error("Form validation failed:", error);
    }
  };

  // อัพเดท roleID เมื่อเลือกชื่อบทบาท
  const handleRoleChange = (roleName: string) => {
    console.log("Role changed to:", roleName);

    const role = roles.find((r) => r.roleName === roleName);
    if (role) {
      console.log("Found role:", role);
      setRoleID(role.roleID);

      // Log ค่า state หลังจากการเปลี่ยนแปลง
      setTimeout(() => {
        console.log("Updated roleID state:", roleID);
      }, 0);
    } else {
      console.warn("Role not found:", roleName);
    }
  };

  // อัพเดท warehouseID เมื่อเลือกชื่อคลังสินค้า
  const handleWarehouseChange = (warehouseName: string) => {
    console.log("Warehouse changed to:", warehouseName);

    const warehouse = warehouses.find((w) => w.warehouseName === warehouseName);
    if (warehouse) {
      console.log("Found warehouse:", warehouse);
      setWarehouseID(warehouse.warehouseID);

      // Log ค่า state หลังจากการเปลี่ยนแปลง
      setTimeout(() => {
        console.log("Updated warehouseID state:", warehouseID);
      }, 0);
    } else {
      console.warn("Warehouse not found:", warehouseName);
    }
  };

  // แสดงบทบาทปัจจุบัน (สำหรับโหมดแก้ไข)
  const getCurrentRoleDisplay = () => {
    if (!user) return null;

    return (
      <div style={{ marginTop: 8 }}>
        <Text type="secondary">บทบาทปัจจุบัน: </Text>
        <Tag color={getRoleColor(user.roleID)}>{user.roleName}</Tag>
      </div>
    );
  };

  // แสดงคลังสินค้าปัจจุบัน (สำหรับโหมดแก้ไข)
  const getCurrentWarehouseDisplay = () => {
    if (!user) return null;

    return (
      <div style={{ marginTop: 8 }}>
        <Text type="secondary">คลังสินค้าปัจจุบัน: </Text>
        <Tag color="cyan">{user.warehouseName}</Tag>
      </div>
    );
  };

  return (
    <Modal
      title={
        <span>
          <UserOutlined />{" "}
          {isEditMode ? "แก้ไขผู้ใช้งาน" : "เพิ่มผู้ใช้งานใหม่"}
        </span>
      }
      open={visible}
      onOk={handleSubmit}
      onCancel={onCancel}
      okText={isEditMode ? "บันทึกการแก้ไข" : "เพิ่มผู้ใช้"}
      cancelText="ยกเลิก"
      confirmLoading={confirmLoading}
      forceRender
      width={500}
    >
      <Form form={form} layout="vertical" requiredMark="optional">
        <Form.Item
          name="userID"
          label={
            <span>
              <UserOutlined /> รหัสผู้ใช้
            </span>
          }
          rules={[{ required: true, message: "กรุณาระบุรหัสผู้ใช้" }]}
          tooltip={
            isEditMode
              ? "ไม่สามารถแก้ไขรหัสผู้ใช้ได้"
              : "กรุณาระบุรหัสผู้ใช้เพื่อเข้าสู่ระบบ"
          }
        >
          <Input placeholder="กรุณาระบุรหัสผู้ใช้" disabled={isEditMode} />
        </Form.Item>

        <Form.Item
          name="roleName"
          label={
            <span>
              <TagsOutlined /> บทบาท
            </span>
          }
          rules={[{ required: true, message: "กรุณาเลือกบทบาท" }]}
          tooltip="เลือกบทบาทที่กำหนดสิทธิ์การเข้าถึง"
          extra={getCurrentRoleDisplay()}
        >
          <Select
            placeholder="เลือกบทบาท"
            showSearch
            optionFilterProp="label"
            onChange={handleRoleChange}
            filterOption={(input, option) =>
              (option?.label?.toString() || "")
                .toLowerCase()
                .includes(input.toLowerCase())
            }
          >
            {roles.map((role) => (
              <Option
                key={role.roleID}
                value={role.roleName}
                label={role.roleName}
              >
                {role.roleName}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          name="warehouseName"
          label={
            <span>
              <BankOutlined /> คลังสินค้า
            </span>
          }
          rules={[{ required: true, message: "กรุณาเลือกคลังสินค้า" }]}
          tooltip="เลือกคลังสินค้าที่ผู้ใช้สังกัด"
          extra={getCurrentWarehouseDisplay()}
        >
          <Select
            placeholder="เลือกคลังสินค้า"
            showSearch
            optionFilterProp="label"
            onChange={handleWarehouseChange}
            filterOption={(input, option) =>
              (option?.label?.toString() || "")
                .toLowerCase()
                .includes(input.toLowerCase())
            }
          >
            {warehouses.map((warehouse) => (
              <Option
                key={warehouse.warehouseID}
                value={warehouse.warehouseName}
                label={warehouse.warehouseName}
              >
                {warehouse.warehouseName}
              </Option>
            ))}
          </Select>
        </Form.Item>

        {isEditMode && (
          <div
            style={{
              marginTop: 16,
              padding: 12,
              backgroundColor: "#f9f9f9",
              borderRadius: 4,
            }}
          >
            <Text type="secondary">
              หมายเหตุ: เฉพาะฟิลด์ที่มีการเปลี่ยนแปลงเท่านั้นที่จะถูกอัปเดต
            </Text>
          </div>
        )}
      </Form>
    </Modal>
  );
};

export default UserForm;
