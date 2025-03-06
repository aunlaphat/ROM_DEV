import { Popconfirm, Button, Col, ConfigProvider, DatePicker, Form, FormInstance, Input, InputNumber, Layout, Row, Select, Table, notification, Modal, Upload, Divider, Tooltip, } from "antd";
import { SearchOutlined, DeleteOutlined, LeftOutlined, PlusCircleOutlined, UploadOutlined, CloseOutlined, QuestionCircleOutlined, } from "@ant-design/icons";
import { debounce } from "lodash";
import { useEffect, useState } from "react";
import * as XLSX from "xlsx";
import Popup from "reactjs-popup";
import icon from "../../assets/images/document-text.png";
import axios from "axios";
import api from "../../utils/axios/axiosInstance"; // นำเข้า axios instance

const { Option } = Select;

interface Address {
  provicesTH: string;
  provinceCode: number;
  districtTH: string;
  districtCode: number;
  subdistrictTH: string;
  subdistrictCode: number;
  zipCode: number;
}

interface Customer {
  Key: number;
  customerID: string;
  customerName: string;
  address: string;
  taxID: string;
}

interface DataItem {
  key: number;
  SKU: string; 
  Name: string;
  QTY: number;
}

interface Product {
  Key: string;
  sku: string;
  nameAlias: string;
  size: string;
}

const CreateTradeReturn = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [invoiceAddress, setInvoiceAddress] = useState("");
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm(); // หลัก
  const [formValid, setFormValid] = useState(false);
  const [formaddress] = Form.useForm(); // ใช้ใน modal
  const [dataSource, setDataSource] = useState<DataItem[]>([]);
  const [isInvoiceEnabled, setIsInvoiceEnabled] = useState(false);
  const [loading, setLoading] = useState(false);

  const [provinces, setProvinces] = useState<Address[]>([]);
  const [districts, setDistricts] = useState<Address[]>([]);
  const [subDistricts, setSubDistricts] = useState<Address[]>([]);
  const [postalCode, setPostalCode] = useState<any[]>([]); 
  const [province, setProvince] = useState<string | undefined>(undefined);
  const [district, setDistrict] = useState<string | undefined>(undefined);
  const [subDistrict, setSubDistrict] = useState<string | undefined>(undefined);
  const [selecteProvince, setSelectedProvince] = useState<string>(""); 
  const [selectedDistrict, setSelectedDistrict] = useState<string>(""); 
  const [selectedSubDistrict, setSelectedSubDistrict] = useState<string>(""); 

  const [customerAccounts, setCustomerAccounts] = useState<Customer[]>([]); // เก็บข้อมูล customer accounts
  const [selectedAccount, setSelectedAccount] = useState<Customer | null>(null); // เก็บข้อมูล customer ที่เลือก
  const [invoiceNames, setInvoiceNames] = useState<any[]>([]); // เก็บข้อมูล invoice names
  const [selectedInvoice, setSelectedInvoice] = useState<any | null>(null); // เก็บข้อมูล invoice ที่เลือก

  const [skuOptions, setSkuOptions] = useState<Product[]>([]); // To store SKU options
  const [nameOptions, setNameOptions] = useState<Product[]>([]); // To store Name Alias options
  const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
  const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
  const [price, setPrice] = useState<number | null>(null); 
  const [qty, setQty] = useState<number | null>(null); 

  /*** Customer&Invoice ***/
  useEffect(() => {
    const fetchCustomerAccounts = async () => {
      setLoading(true);
      try {
        const response = await api.get("/api/constants/get-customer-id");
        setCustomerAccounts(response.data.data); // เก็บข้อมูล customer accounts
      } catch (error) {
        notification.error({
          message: "Error",
          description: "Failed to display customer account",
        });
        setCustomerAccounts([]); // ตั้งค่ากลับเป็น array ว่างเมื่อเกิดข้อผิดพลาด
      } finally {
        setLoading(false);
      }
    };

    fetchCustomerAccounts();
  }, []); 

  // หลังเลือก Customer Account
  const handleAccountChange = async (value: string) => {
    try { // Reset invoice when changing the customer
      setSelectedAccount(null);
      setSelectedInvoice(null); 
      form.resetFields(["Invoice_name"]);

      const customerResponse = await api.get(
        `/api/constants/get-customer-info?customerID=${value}`,
      );

      const customerData = customerResponse.data.data;

      if (customerData && customerData.length > 0) {
        const firstCustomer = customerData[0];
        setSelectedAccount(firstCustomer);
  
        // Set available invoice names
        setInvoiceNames(customerData);
  
        form.setFieldsValue({
          Customer_name: firstCustomer.customerName,
          Address: firstCustomer.address,
          Tax: firstCustomer.taxID,
          Invoice_name: firstCustomer.customerName, // Set the first available invoice name (or leave empty if needed)
        });
      } else {
        notification.warning({
          message: "Data Not Found",
          description: "No invoice data found for this Customer Account",
        });
        form.setFieldsValue({
          Customer_name: "",
          Address: "",
          Tax: "",
          Invoice_name: "", 
        });
      }
    } catch (error) {
      notification.error({
        message: "Error",
        description: "Failed to display customer account",
      });
    }
  };

  const debouncedSearch = debounce(async (value: string) => {
    setLoading(true); 
    try {
      const response = await api.get(
        "/api/constants/search-invoice-names", 
        {
          params: {
            customerID: selectedAccount?.customerID, // ใช้ customerID ที่เลือก
            keyword: value, // ใช้ keyword ที่ค้นหา
            offset: 0,
            limit: 50, 
          },
        },
      );
      setInvoiceNames(response.data.data); 
    } catch (error) {
      console.error("Error fetching invoice names:", error);
      notification.error({
        message: "Error",
        description: "Failed to display invoice name",
      });
    } finally {
      setLoading(false); // หยุดการโหลดหลังจากเสร็จสิ้นการค้นหา
    }
  }, 1000); // ตั้งเวลา debounce การ search เป็น 1000ms (=1 วินาที)

  const handleInvoiceSearch = (value: string) => {
    debouncedSearch(value);
  };

  const handleInvoiceChange = async (value: string) => {
    const invoiceData = value.split("+"); // ใช้ + แยกข้อมูล
    const customerName = invoiceData[0].trim(); 
    const address = invoiceData
      .slice(1, invoiceData.length - 1)
      .join("+")
      .trim(); // แบ่งข้อมูลที่รวมกันอยู่หลายส่วนออกมา ด้วย + ที่เชื่อมกัน
    const taxID = invoiceData[invoiceData.length - 1].trim(); // taxID จะเป็นค่าที่แยกออกมาเป็นส่วนสุดท้าย
    const selectedInvoice = invoiceNames.find(  // ค้นหา selectedInvoice ที่ตรงกับ customerName, address และ taxID
      (invoice) =>
        invoice.customerName === customerName &&
        invoice.address === address &&
        invoice.taxID === taxID,
    );

    if (selectedInvoice) {
      setSelectedInvoice(selectedInvoice); 
      // อัปเดตฟอร์มด้วยข้อมูลที่เลือก
      form.setFieldsValue({
        Customer_name: selectedInvoice.customerName,
        Address: selectedInvoice.address,
        Tax: selectedInvoice.taxID,
      });
    } else {  // reset value หากไม่พบข้อมูลใบแจ้งหนี้ที่ตรงกับการเลือก
      notification.warning({
        message: "Data Not Found",
        description: "No information invoice data found for this Customer Account",
      });
      form.setFieldsValue({
        Customer_name: "", 
        Address: "", 
        Tax: "",
      });
    }
  };

  /*** Address ***/
  useEffect(() => {
    const fetchProvinces = async () => {
      setLoading(true);
      try {
        const response = await api.get("/api/constants/get-provinces");
        setProvinces(response.data.data);
      } catch (error) {
        console.error("Failed to fetch provinces", error);
      } finally {
        setLoading(false);
      }
    };
    fetchProvinces();
  }, []);

  useEffect(() => {
    if (province) {
      const fetchDistricts = async () => {
        setLoading(true);
        try {
          const response = await api.get(`/api/constants/get-district?provinceCode=${province}`);
          setDistricts(response.data.data);
        } catch (error) {
          console.error("Failed to fetch districts", error);
        } finally {
          setLoading(false);
        }
      };
      fetchDistricts();
    } else {
      setDistricts([]);
    } 
  }, [province]);

  useEffect(() => {
    if (district) {
      const fetchSubDistricts = async () => {
        setLoading(true);
        try {
          const response = await api.get(`/api/constants/get-sub-district?districtCode=${district}`);
          setSubDistricts(response.data.data);
        } catch (error) {
          console.error("Failed to fetch subdistricts", error);
        } finally {
          setLoading(false);
        }
      };
      fetchSubDistricts();
    } else {
      setSubDistricts([]);
    }
  }, [district]);

  useEffect(() => {
    if (subDistrict) {
      const fetchPostalCode = async () => {
        setLoading(true);
        try {
          const response = await api.get(`/api/constants/get-postal-code?subdistrictCode=${subDistrict}`);
          setPostalCode(response.data.data);
          formaddress.setFieldsValue({
            PostalCode: response.data.data.length > 0 ? response.data.data[0].zipCode : "",
          });
          console.log("Fetched Postal Codes: ", response.data.data); // ตรวจสอบข้อมูลที่ได้รับ
        } catch (error) {
          console.error("Failed to fetch postal code", error);
        } finally {
          setLoading(false);
        }
      };
      fetchPostalCode();
    } else {
      setPostalCode([]);
      formaddress.setFieldsValue({
        PostalCode: "",
      });
    }
  }, [subDistrict]); 

  const handleProvinceChange = (value: string) => {
    if (!value) return; // ป้องกันค่า undefined
    setProvince(value);
    setSelectedProvince(value); 

    formaddress.resetFields(["District", "SubDistrict", "PostalCode"]); // รีเซ็ตค่าในฟอร์ม
    setDistrict(undefined); 
    setSubDistrict(undefined);
    setPostalCode([]); 
    setSelectedDistrict(""); 
    setSelectedSubDistrict(""); 
  };
  
  const handleDistrictChange = (value: string) => {
    if (!value) return; // ป้องกันค่า undefined
    setDistrict(value);
    setSelectedDistrict(value); 

    formaddress.resetFields(["SubDistrict", "PostalCode"]); // รีเซ็ตค่าในฟอร์ม
    setSubDistrict(undefined); 
    setPostalCode([]); 
    setSelectedSubDistrict(""); 
  };

  const handleSubDistrictChange = (value: string) => {
    if (!value) return; // ป้องกันค่า undefined
    setSubDistrict(value);
    setSelectedSubDistrict(value); 
  };

  const handleOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);

    // ให้ฟอร์มรีเซ็ตเฉพาะในกรณีที่ไม่ได้กดบันทึก
    formaddress.resetFields(["HouseNo","Province","District", "SubDistrict", "PostalCode"]);
    setProvince(undefined); 
    setDistrict(undefined); 
    setSubDistrict(undefined); 
    setPostalCode([]); 
    setSelectedProvince(""); 
    setSelectedDistrict(""); 
    setSelectedSubDistrict(""); 
  };

  const handleSelectChange = (value: any) => {
    // เมื่อเลือกจังหวัดแล้วปิด Popup
    // setOpen(false);
  };

  /*** SKU&NameAlias ***/

  // ค้นหา Product (SKU หรือ NAMEALIAS)
  const debouncedSearchSKU = debounce(async (value: string, searchType: string) => {
    setLoading(true);
    try {
      const response = await api.get("/api/constants/search-product", {
        params: {
          keyword: value,
          searchType,
          offset: 0,
          limit: 5,
        },
      });

      const products = response.data.data;

      if (searchType === "SKU") {
        setSkuOptions(products.map((product: Product) => ({
          sku: product.sku,
          nameAlias: product.nameAlias,
          size: product.size,
        })));
      } else if (searchType === "NAMEALIAS") {
        setNameOptions(products.map((product: Product) => ({
          sku: product.sku,
          nameAlias: product.nameAlias,
          size: product.size,
        })));
      }
    } catch (error) {
      console.error("Error fetching products:", error);
      notification.error({
        message: "Error",
        description: "There was an error fetching product data.",
      });
    } finally {
      setLoading(false);
    }
  }, 1000);

  const handleSearchSKU = (value: string) => {
    debouncedSearchSKU(value, "SKU");
  };

  const handleSearchNameAlias = (value: string) => {
    debouncedSearchSKU(value, "NAMEALIAS");
  };

  // เมื่อเลือก Name Alias แล้วใช้ `/api/constants/get-sku` เพื่อหา SKU
  const handleNameChange = async (value: string) => {
    const [nameAlias, size] = value.split("+"); // แยกค่า nameAlias และ size โดยใช้ `+`

    try {
      setLoading(true);
      const response = await api.get("/api/constants/get-sku", {
        params: { nameAlias, size },
      });

      // เก็บผลลัพธ์จาก API เพื่อแสดงหลาย SKU
      const products = response.data.data;

      if (products.length > 0) {
        setSkuOptions(products.map((product: Product) => ({
          sku: product.sku,
          nameAlias: product.nameAlias,
          size: product.size,
        })));
        form.setFieldsValue({
          SKU: products[0].sku, // ตั้งค่า SKU ตัวแรกที่พบ
        });
      } else {
        console.warn("No SKU found for:", nameAlias, size);
        setSkuOptions([]); 
        setNameOptions([]); 
        form.setFieldsValue({ SKU: "", SKU_Name: "" }); // เคลียร์ค่าในฟอร์ม
      }
    } catch (error) {
      console.error("Error fetching SKU:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSKUChange = (value: string) => {
    const selected = skuOptions.find((option) => option.sku === value);
    
    if (selected) {
      form.setFieldsValue({
        SKU: selected.sku,
        SKU_Name: selected.nameAlias,
      });
      setSelectedSKU(selected.sku);
      setSelectedName(selected.nameAlias);

     // อัปเดต nameOptions ตาม SKU ที่เลือก
     const filteredNameOptions = skuOptions
     .filter((option) => option.sku === selected.sku) // กรองเฉพาะ SKU ที่ตรงกับที่เลือก
     .map((option) => ({
       ...option,  // คัดลอกค่าเดิม
       Key: option.sku,  // เพิ่มคีย์ Key ที่ต้องการ
     }));
     setNameOptions(filteredNameOptions);  // อัปเดต nameOptions
    } else { // เคลียร์ค่าเมื่อไม่มี SKU ที่ตรงกัน
      setSkuOptions([]); 
      setNameOptions([]); 
      setSelectedSKU("");
      setSelectedName("");
    }
  };

  const handleSubmit = () => {
    if (dataSource.length === 0) {
      notification.warning({
        message: "ไม่สามารถส่งข้อมูลได้",
        description: "กรุณาเพิ่มข้อมูลในตารางก่อนส่ง!",
      });
      return; // หยุดการทำงานของฟังก์ชัน
    }

    // ส่งข้อมูลและรีเซ็ตฟอร์มและตาราง
    console.log("Table Data:", dataSource);

    form.resetFields();  // รีเซ็ตฟอร์มและตาราง
    formaddress.resetFields(); // รีเซ็ตฟอร์ม address
    setSelectedAccount(null);
    setSelectedInvoice(null);
    setInvoiceNames([]);
    setDataSource([]); // หรือปรับเป็นค่าเริ่มต้นที่คุณต้องการได้

    notification.success({
      message: "ส่งข้อมูลสำเร็จ",
      description: "ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!",
    });
  };

  // handle value Invoice_Name ของ New Invoice Address
  useEffect(() => {
    if (formaddress && (selectedInvoice || selectedAccount)) {
      formaddress.setFieldsValue({
        Invoice_Name: selectedInvoice?.customerName || selectedAccount?.customerName || "",
      });
    }
  }, [selectedInvoice, selectedAccount]);

  // save new invoice address
  const handleSave = async () => {
    try {
      const values = await formaddress.validateFields();
      console.log("Form Values:", values);

      const ProvicesTH = provinces.find(
        (item) => item.provinceCode.toString() === values.Province
      )?.provicesTH;

      const DistrictTH = districts.find(
        (item) => item.districtCode.toString() === values.District
      )?.districtTH;

      const SubdistrictTH = subDistricts.find(
        (item) => item.subdistrictCode.toString() === values.SubDistrict
      )?.subdistrictTH;

      const PostalCode = postalCode.find(
        (item) => item.zipCode === values.PostalCode
      )?.zipCode;

      // Update form values in the main form
      form.setFieldsValue({
        Invoice_Name: values.Invoice_Name, 
        Address:
        values.HouseNo +
        " " +
        ProvicesTH +
        " " +
        DistrictTH +
        " " +
        SubdistrictTH +
        " " +
        PostalCode,
    
      });

      setIsSaving(true);
      setIsSaving(false);

      notification.success({
        message: "Update Success",
        description: "Update Invoice Address Success!",
      });

      handleClose(); // Close modal after save
    } catch (error) {
      console.error("Failed to save:", error);
    }
  };

  const onSearch = (value: string) => {
    console.log("search:", value);
  };

  const columns = [
    { title: "SKU", dataIndex: "SKU", key: "SKU", id: "SKU" },
    { title: "Name", dataIndex: "Name", key: "Name", id: "Name" },
    { title: "QTY", dataIndex: "QTY", key: "QTY", id: "QTY" },
    {
      title: "Action",
      id: "Action",
      dataIndex: "Action",
      key: "Action",
      render: (_: any, record: { key: number }) => (
        <Popconfirm
          title="คุณแน่ใจหรือไม่ว่าต้องการลบข้อมูลนี้?"
          onConfirm={() => handleDelete(record.key)} // เรียกใช้ฟังก์ชัน handleDelete เมื่อกดยืนยัน
          okText="ใช่"
          cancelText="ไม่"
        >
          <DeleteOutlined
            style={{ cursor: "pointer", color: "red", fontSize: "20px" }}
          />
        </Popconfirm>
      ),
    },
  ];

  const handleDownloadTemplate = () => {
    const templateColumns = columns.filter((col) => col.key !== "Action"); // กรองออก action column
    const ws = XLSX.utils.json_to_sheet([]);
    XLSX.utils.sheet_add_aoa(ws, [templateColumns.map((col) => col.title)]);

    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, "Template");

    XLSX.writeFile(wb, "Template.xlsx");
  };

  const handleUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const data = new Uint8Array(e.target?.result as ArrayBuffer);
      const workbook = XLSX.read(data, { type: "array" });
      const worksheet = workbook.Sheets[workbook.SheetNames[0]];
      const json = XLSX.utils.sheet_to_json<DataItem>(worksheet);

      // กรองข้อมูลเฉพาะที่มี SKU และ QTY
      const filteredData = json.filter((item) => item.SKU && item.QTY);

      // อัปเดต dataSource ด้วยข้อมูลที่กรอง
      setDataSource(filteredData);

      notification.success({
        message: "อัปโหลดสำเร็จ",
        description: "ข้อมูลจากไฟล์ Excel ถูกนำเข้าเรียบร้อยแล้ว!",
      });
    };
    reader.readAsArrayBuffer(file);
  };

  const uploadProps = {
    beforeUpload: (file: File) => {
      handleUpload(file);
      return false; // ป้องกันไม่ให้ Ant Design ทำการอัปโหลด
    },
  };

  const handleAdd = () => {
    form
      .validateFields()
      .then((values) => {
        const [nameAlias, size] = values.SKU_Name.split('+');  // แยกค่า nameAlias กับ size
        // ตรวจสอบว่า SKU ที่กรอกมีอยู่ใน dataSource หรือไม่
        const isSKUExist = dataSource.some((item) => item.SKU === values.SKU);

        if (isSKUExist) {
          // แสดงข้อความเตือนว่า SKU ซ้ำ
          notification.warning({
            message: "มีข้อผิดพลาด",
            description: "SKU นี้ถูกเพิ่มไปแล้วในรายการ!",
          });
          return; // ไม่ทำการเพิ่มข้อมูล
        }

        // ถ้า SKU ยังไม่ซ้ำ เพิ่มข้อมูลใหม่
        const newData = {
          key: dataSource.length + 1,
          SKU: values.SKU,
          Name: nameAlias,
          QTY: values.QTY,
          Price: values.Price,
          Tax: values.Tax || "", // ถ้าไม่มี Tax ให้ใส่ค่าว่าง
        };

        setDataSource([...dataSource, newData]); // เพิ่มข้อมูลใหม่ไปยัง dataSource

        // แสดงข้อความเมื่อเพิ่มข้อมูลสำเร็จ
        notification.success({
          message: "เพิ่มสำเร็จ",
          description: "ข้อมูลของคุณถูกเพิ่มเรียบร้อยแล้ว!",
        });

        // ล้างฟิลด์ในฟอร์มหลังจากเพิ่มข้อมูลเสร็จ
        form.resetFields(["SKU", "SKU_Name", "QTY", "Price"]);
        setSkuOptions([]);
        setNameOptions([]);
        setSelectedSKU("");
        setSelectedName("");
      })
      .catch((info) => {
        console.log("Validate Failed:", info);
        notification.warning({
          message: "มีข้อสงสัย",
          description: "กรุณากรอกข้อมูลให้ครบก่อนเพิ่ม!",
        });
      });
  };

  const handleDelete = (key: number) => {
    setDataSource(dataSource.filter((item) => item.key !== key));
    notification.success({
      message: "ลบข้อมูลสำเร็จ",
      description: "ข้อมูลของคุณถูกลบออกเรียบร้อยแล้ว.",
    });
  };

  const onChange = () => {
    const values = form.getFieldsValue();
    const { Date, SKU, QTY } = values;

    // Set form validity based on required fields
    setFormValid(Date && SKU && QTY);
  };

  return (
    <ConfigProvider>
      <div
        style={{
          marginLeft: "28px",
          fontSize: "25px",
          fontWeight: "bold",
          color: "DodgerBlue",
        }}
      >
        Create Trade Return
      </div>
      <Layout>
        <Layout.Content
          style={{
            margin: "24px",
            padding: 36,
            minHeight: 360,
            background: "#fff",
            borderRadius: "8px",
          }}
        >
          <Form
            form={form}
            layout="vertical"
            style={{ width: "100%", padding: "30px" }}
          >
            <div>
              <Divider
                style={{ color: "#657589", fontSize: "22px", margin: 30 }}
                orientation="left"
              >
                Sale Order Information
              </Divider>
              <Row gutter={16}>
                <Col span={8}>
                  <Form.Item
                    id="Tracking"
                    label={
                      <span style={{ color: "#657589" }}>
                        กรอกเลข Tracking:&nbsp;
                        <Tooltip title="เลขTracking จากขนส่ง">
                          <QuestionCircleOutlined
                            style={{ color: "#657589" }}
                          />
                        </Tooltip>
                      </span>
                    }
                    name="Tracking"
                    rules={[{ required: true, message: "กรอกเลข Tracking" }]}
                  >
                    <Input style={{ height: 40 }} />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    id="Logistic"
                    label={
                      <span style={{ color: "#657589" }}>
                        กรอก Logistic:&nbsp;
                        <Tooltip title="ผู้ให้บริการขนส่ง">
                          <QuestionCircleOutlined
                            style={{ color: "#657589" }}
                          />
                        </Tooltip>
                      </span>
                    }
                    name="Logistic"
                    rules={[
                      { required: true, message: "กรอก Logistic" },
                    ]}
                  >
                    <Input style={{ height: 40 }} />
                  </Form.Item>
                </Col>
                <Col span={8}></Col>
              </Row>
              <Row
                gutter={16}
                align="middle"
                justify="center"
                style={{ marginTop: "20px", width: "100%" }}
              >
                <Col span={8}>
                  <Form.Item
                    id="Doc"
                    label={
                      <span style={{ color: "#657589" }}>
                        กรอกเอกสารอ้างอิง:&nbsp;
                        <Tooltip title="ตัวอย่างเอกสาร SOA2410-00234">
                          <QuestionCircleOutlined
                            style={{ color: "#657589" }}
                          />
                        </Tooltip>
                      </span>
                    }
                    name="Doc"
                  >
                    <Input
                      style={{ width: "100%", height: "40px" }}
                      placeholder="ตัวอย่างเอกสาร SOA2410-00234"
                    />
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="Customer account"
                    name="Customer_account"
                    rules={[{ required: true }]}
                  >
                    <Select
                      showSearch
                      placeholder="Select Customer Account"
                      onChange={handleAccountChange}
                      loading={loading}
                      listHeight={160} // ปรับให้พอดีกับ 4 รายการ
                      virtual // ทำให้ค้นหาไวขึ้น
                    >
                      {customerAccounts.length > 0 ? (
                        customerAccounts.map((account) => (
                          <Option
                            key={account.customerID}
                            value={account.customerID}
                          >
                            {account.customerID}
                          </Option>
                        ))
                      ) : (
                        <Option disabled>No customer accounts available</Option>
                      )}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item
                    label="Customer Name"
                    name="Customer_name"
                    rules={[{ required: true }]}
                  >
                    <Input
                      value={selectedAccount?.customerName || "-"}
                      disabled
                    />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16} style={{ marginTop: "10px" }}>
                <Col span={8}>
                  <Form.Item
                    label="Invoice Name"
                    name="Invoice_name"
                    rules={[{ required: true }]}
                  >
                    <Select
                      showSearch
                      value={selectedAccount?.customerName || "-"}
                      placeholder="Select Invoice Name"
                      onSearch={handleInvoiceSearch} // เรียกฟังก์ชันค้นหาตอนพิมพ์
                      onChange={handleInvoiceChange}
                      loading={loading}
                      listHeight={160} // ปรับให้พอดีกับ 4 รายการ
                      virtual
                    >
                      {invoiceNames.map((invoice) => (
                        <Option
                          key={`${invoice.customerName}-${invoice.address}-${invoice.taxID}`}
                          value={`${invoice.customerName}+${invoice.address}+${invoice.taxID}`}
                        >
                          {invoice.customerName}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={8}>
                  <Form.Item label="Tax ID" name="Tax">
                    <Input value={selectedAccount?.taxID || "-"} disabled />
                  </Form.Item>
                </Col>
              </Row>
              <Divider
                style={{ color: "#657589", fontSize: "22px", margin: 30 }}
                orientation="left"
              >
                Address Information
              </Divider>
              <Row gutter={16} style={{ marginTop: "10px" }}>
                <Col span={18}>
                  <Form.Item
                    label="Invoice Address"
                    name="Address"
                    // rules={[{ required: true }]}
                  >
                    <Input value={selectedAccount?.address || "-"} disabled />
                  </Form.Item>
                </Col>
                <Col span={6}>
                  <Button
                    id="NewInvoiceAddress"
                    type="primary"
                    onClick={handleOpen}
                    style={{ width: "100%", height: "40px", marginTop: 30 }}
                  >
                    New invoice address
                  </Button>
                </Col>
              </Row>
              <Divider
                style={{ color: "#657589", fontSize: "22px", margin: 30 }}
                orientation="left"
              > 
                {" "}
                SKU information
              </Divider>
              <Row gutter={16} style={{ marginTop: "10px", width: "100%" }}>
                <Col span={6}>
                  <Form.Item
                    id="SKU"
                    label={<span style={{ color: "#657589" }}>กรอก SKU :</span>}
                    name="SKU"
                    rules={[{ required: true, message: "กรุณากรอก SKU" }]}
                  >
                    <Select
                      showSearch
                      style={{ width: "100%", height: "40px" }}
                      placeholder="Search by SKU"
                      value={selectedSKU} // ใช้ค่าที่เลือก
                      onSearch={handleSearchSKU} // ใช้สำหรับค้นหา SKU
                      onChange={handleSKUChange} // เมื่อเลือก SKU
                      loading={loading}
                      listHeight={160} // ปรับให้พอดีกับ 4 รายการ
                      virtual // ทำให้ค้นหาไวขึ้น
                      dropdownStyle={{ minWidth: 200 }}
                    >
                      {skuOptions.map((option) => (
                        <Option 
                          key={`${option.sku}-${option.size}`} 
                          value={option.sku}
                        >
                          {option.sku}
                      </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>

                <Col span={7}>
                  <Form.Item
                    id="SKU_Name"
                    label={
                      <span style={{ color: "#657589" }}>กรอก SKU Name:</span>
                    }
                    name="SKU_Name"
                    rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
                  >
                    <Select
                      showSearch
                      style={{ width: "100%", height: "40px" }}
                      placeholder="Search by Product Name"
                      value={selectedName} // ใช้ค่าที่เลือก
                      onSearch={handleSearchNameAlias} // ใช้สำหรับค้นหา Name Alias
                      onChange={handleNameChange} // เมื่อเลือก Name Alias
                      loading={loading}
                      listHeight={160} // ปรับให้พอดีกับ 4 รายการ
                      virtual // ทำให้ค้นหาไวขึ้น
                      dropdownStyle={{ minWidth: 300 }}
                    >
                      {nameOptions.map((option) => (
                        <Option 
                          key={`${option.nameAlias}-${option.size}`} 
                          value={`${option.nameAlias}+${option.size}`}
                        >
                          {option.nameAlias}
                        </Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={4}>
                  <Form.Item
                    id="qty"
                    label={<span style={{ color: "#657589" }}>QTY:</span>}
                    name="QTY"
                    rules={[{ required: true, message: "กรุณากรอก QTY" }]}
                  >
                    <InputNumber
                      min={1}
                      max={100}
                      value={qty}
                      onChange={(value) => setQty(value)}
                      style={{
                        width: "100%",
                        height: "40px",
                        lineHeight: "40px",
                      }}
                    />
                  </Form.Item>
                </Col>

                <Col span={4}>
                  <Form.Item
                    id="price"
                    label={<span style={{ color: "#657589" }}>Price:</span>}
                    name="Price"
                    rules={[{ required: true, message: "กรุณากรอก Price" }]}
                  >
                    <InputNumber
                      min={1}
                      max={100000}
                      value={price}
                      onChange={(value) => setPrice(value)}
                      step={0.01}
                      style={{
                        width: "100%",
                        height: "40px",
                        lineHeight: "40px",
                      }}
                    />
                  </Form.Item>
                </Col>

                <Col span={3}>
                  <Button
                    id="add"
                    type="primary"
                    style={{ width: "100%", height: "40px", marginTop: 30 }}
                    onClick={handleAdd} // เรียกใช้ฟังก์ชัน handleAdd
                  >
                    <PlusCircleOutlined />
                    Add
                  </Button>
                </Col>
              </Row>
            </div>
          </Form>
          <Modal
            open={open}
            onClose={handleClose}
            closeIcon={false}
            footer={null}
          >
            <div
              style={{
                display: "flex",
                justifyContent: "flex-end",
                marginBottom: "10px",
                overflow: "auto",
              }}
            >
              <Button
                id="Closeicon"
                type="text"
                onClick={handleClose}
                icon={<CloseOutlined style={{ fontSize: "24px" }} />}
                danger
              />
            </div> 
            <div style={{ fontSize: "20px", color: "#35465B" }}>
              New Invoice Address
            </div>
            <Form
              form={formaddress}
              layout="vertical"
              style={{ width: "100%", display: "flex", padding: 20 }}
              onFinish={handleSave}
            >
              <Row
                gutter={16}
                style={{ marginTop: "10px", justifyContent: "center" }}
              >
                <Col>
                  <Form.Item
                    id="Invoice Name"
                    label={
                      <span style={{ color: "#657589" }}>Invoice name:</span>
                    }
                    name="Invoice_Name"
                    rules={[{ required: true, message: "Please Select Invoice name" }]}
                  >
                  <Input
                    style={{ width: "400px", height: "40px" }}
                    value={form.getFieldValue("Invoice_Name")} // ใช้ค่าจาก form
                    disabled
                  />
                  </Form.Item>
                </Col>

                <Col>
                  <Form.Item
                    id="้houseno"
                    label={
                      <span style={{ color: "#657589" }}>บ้านเลขที่:</span>
                    }
                    name="HouseNo"
                    rules={[{ required: true, message: "Please Input House no." }]}
                  >
                    <Input
                      style={{ width: "400px", height: "40px" }}
                      placeholder="กรอกบ้านเลขที่"
                    />
                  </Form.Item>
                </Col>

                {/* Province */}
                <Col>
                  <Form.Item
                    id="SelectProvince"
                    label={<span style={{ color: "#657589" }}>จังหวัด:</span>}
                    name="Province"
                    rules={[{ required: true, message: "Please Select Province" }]}
                  >
                  <Select
                      showSearch
                      placeholder="Select Province"
                      value={selecteProvince}
                      onChange={handleProvinceChange}
                      loading={loading}
                      listHeight={160}
                      virtual
                      style={{ width: "400px", height: "40px" }}
                      options={provinces.map(p => ({
                          label: p.provicesTH,
                          value: p.provinceCode.toString()
                      }))}
                      filterOption={(input, option) => {
                          if (!option) return false;
                          return option.label.toLowerCase().includes(input.toLowerCase());
                      }}
                  />
                  </Form.Item>
                </Col>

                {/* District */}
                <Col>
                  <Form.Item
                    id="SelectDistrict"
                    label={<span style={{ color: "#657589" }}>เขต:</span>}
                    name="District"
                    rules={[{ required: true, message: "Please Select District" }]}
                  >
                    <Select
                        showSearch
                        placeholder="Select District"
                        value={selectedDistrict}
                        onChange={handleDistrictChange}
                        loading={loading}
                        listHeight={160}
                        virtual
                        style={{ width: "400px", height: "40px" }}
                        options={districts.map(d => ({
                            label: d.districtTH,
                            value: d.districtCode.toString()
                        }))}
                        filterOption={(input, option) => {
                            if (!option) return false;
                            return option.label.toLowerCase().includes(input.toLowerCase());
                        }}
                    />

                  </Form.Item>
                </Col>

                {/* SubDistrict */}
                <Col>
                  <Form.Item
                    id="SelectSubDistrict"
                    label={<span style={{ color: "#657589" }}>แขวง:</span>}
                    name="SubDistrict"
                    rules={[{ required: true, message: "Please Select Sub-district" }]}
                  >
                      <Select
                          showSearch
                          placeholder="Select Sub-District"
                          value={selectedSubDistrict}
                          onChange={handleSubDistrictChange}
                          loading={loading}
                          listHeight={160}
                          virtual
                          style={{ width: "400px", height: "40px" }}
                          options={subDistricts.map(s => ({
                              label: s.subdistrictTH,
                              value: s.subdistrictCode.toString()
                          }))}
                          filterOption={(input, option) => {
                              if (!option) return false;
                              return option.label.toLowerCase().includes(input.toLowerCase());
                          }}
                      />
                  </Form.Item>
                </Col>

                {/* Postal Code */}
                <Col>
                  <Form.Item
                    id="PostalCode"
                    label={
                      <span style={{ color: "#657589" }}>รหัสไปรษณีย์:</span>
                    }
                    name="PostalCode"
                    rules={[{ required: true, message: "Please Select Postcode" }]}
                  >
                  <Input
                    style={{ width: "400px", height: "40px" }}
                    value={postalCode.length > 0 ? postalCode[0].zipCode : ""}
                    disabled
                  />
                  </Form.Item>
                </Col>

                {/* Save Button */}
                <Col>
                  <Button
                    id="save"
                    type="primary"
                    htmlType="submit"
                    disabled={isSaving}
                  >
                    Save
                  </Button>
                </Col>
              </Row>
            </Form>
          </Modal>

          <Row gutter={20} style={{ marginBottom: 20, marginLeft: 20 }}>
            <Col>
              <Button id=" Download Template" onClick={handleDownloadTemplate}>
                <img
                  src={icon}
                  alt="Download Icon"
                  style={{ width: 16, height: 16, marginRight: 8 }}
                />
                Download Template
              </Button>
            </Col>

            <Col>
              <Upload {...uploadProps} showUploadList={false}>
                <Button
                  id=" Import Excel"
                  icon={<UploadOutlined />}
                  style={{
                    background: "#7161EF",
                    color: "#FFF",
                    marginBottom: 10,
                  }}
                >
                  Import Excel
                </Button>
              </Upload>
            </Col>
          </Row>
          <div>
            <Table
              dataSource={dataSource}
              columns={columns}
              rowKey="key"
              pagination={false} // Disable pagination if necessary
              style={{ width: "100%", tableLayout: "fixed" }} // Ensure the table takes full width and is fixed layout
              scroll={{ x: "max-content" }}
            />
          </div>
          <Row justify="center" gutter={16}>
            <Popconfirm
              id="popconfirmSubmit"
              title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
              onConfirm={handleSubmit} // เรียกใช้ฟังก์ชัน handleSubmit เมื่อกดยืนยัน
              okText="ใช่"
              cancelText="ไม่"
            >
              <Button
                id="Submit"
                style={{
                  color: "#fff",
                  backgroundColor: "#14C11B",
                  width: 100,
                  height: 40,
                  margin: 20,
                }}
              >
                Submit
              </Button>
            </Popconfirm>
          </Row>
        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
};

export default CreateTradeReturn;
