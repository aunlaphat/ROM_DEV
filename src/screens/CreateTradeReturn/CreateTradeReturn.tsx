import { Popconfirm, Button, Col, ConfigProvider, DatePicker, Form, FormInstance, Input, InputNumber, Layout, Row, Select, Table, notification, Modal, Upload, Divider, Tooltip, Pagination, } from "antd";
import { SearchOutlined, DeleteOutlined, LeftOutlined, PlusCircleOutlined, UploadOutlined, CloseOutlined, QuestionCircleOutlined, } from "@ant-design/icons";
import { debounce } from "lodash";
import { useEffect, useState } from "react";
import * as XLSX from "xlsx";
import Popup from "reactjs-popup";
import icon from "../../assets/images/document-text.png";
import api from "../../utils/axios/axiosInstance"; 
import { useSelector } from 'react-redux';
import { RootState } from "../../redux/types";
import { TRANSPORT_TYPES, Address, Customer, DataItem, Product } from '../../types/types';
import '../../style/styles.css';
const { Option } = Select;

const CreateTradeReturn = () => {
  const [isSaving, setIsSaving] = useState(false);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm(); 
  const [formValid, setFormValid] = useState(false);
  const [formaddress] = Form.useForm(); 
  const [dataSource, setDataSource] = useState<DataItem[]>([]);
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

  const [customerAccounts, setCustomerAccounts] = useState<Customer[]>([]); 
  const [selectedAccount, setSelectedAccount] = useState<Customer | null>(null); 
  const [invoiceNames, setInvoiceNames] = useState<any[]>([]); 
  const [selectedInvoice, setSelectedInvoice] = useState<any | null>(null); 

  const [skuOptions, setSkuOptions] = useState<Product[]>([]); 
  const [nameOptions, setNameOptions] = useState<Product[]>([]); 
  const [selectedSKU, setSelectedSKU] = useState<string | undefined>(undefined);
  const [selectedName, setSelectedName] = useState<string | undefined>(undefined);
  const [price, setPrice] = useState<number | null>(null); 
  const [qty, setQty] = useState<number | null>(null); 

  const [returnQty, setReturnQty] = useState<number | null>(null);
  const [pricePerUnit, setPricePerUnit] = useState<number | null>(null);

  // ดึงข้อมูลผู้ใช้ที่เข้าสู่ระบบ
  const auth = useSelector((state: RootState) => state.auth);
  const userID = auth?.user?.userID;
  const token = localStorage.getItem("access_token");

  const [isModalVisible, setIsModalVisible] = useState(false);

  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(5);

  // ฟังก์ชันสำหรับเปลี่ยนหน้า
  const handlePageChange = (page: number, pageSize: number) => {
    setCurrentPage(page);
    setPageSize(pageSize); // ถ้าผู้ใช้เลือกจำนวนรายการต่อหน้าใหม่
  };

  // คำนวณจำนวนหน้าทั้งหมดจากจำนวนรายการทั้งหมด
  const totalPages = Math.ceil(dataSource.length / pageSize);

  // ตรวจสอบว่า pagination ควรแสดงหรือไม่ (ให้แสดงเสมอแม้ว่า dataSource จะมีน้อยกว่า pageSize)
  const showPagination = dataSource.length > 0;

  const showModal = () => {
      setIsModalVisible(true);
  };

  const handleOk = () => {
      setIsModalVisible(false);
      handleSubmit(); // เรียกใช้ฟังก์ชัน handleSubmit
  };

  const handleCancel = () => {
      setIsModalVisible(false);
  };
  
  useEffect(() => {
    if (returnQty !== null && pricePerUnit !== null) {
      const calculatedPrice = returnQty * pricePerUnit;
      setPrice(calculatedPrice);
      form.setFieldsValue({ Price: calculatedPrice });
    }
  }, [returnQty, pricePerUnit, form]);

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
    if (!value) return; 
    setDistrict(value);
    setSelectedDistrict(value); 

    formaddress.resetFields(["SubDistrict", "PostalCode"]);
    setSubDistrict(undefined); 
    setPostalCode([]); 
    setSelectedSubDistrict(""); 
  };

  const handleSubDistrictChange = (value: string) => {
    if (!value) return; 
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

  /*** Logistic Type ***/
  const [isOtherTransport, setIsOtherTransport] = useState(false);
  const [transportValue, setTransportValue] = useState<string | undefined>(undefined);
  const handleTransportChange = (value: string) => {
    if (value === 'OTHER') {
      setIsOtherTransport(true);
      form.resetFields(['Logistic']);
      setTransportValue(''); 
    } else {
      setIsOtherTransport(false);
      setTransportValue(value);
    }
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
    { title: "รหัสสินค้า", dataIndex: "SKU", key: "SKU", id: "SKU" },
    { title: "ชื่อสินค้า", dataIndex: "Name", key: "Name", id: "Name" },
    { title: "จำนวนเริ่มต้น", dataIndex: "QTY", key: "QTY", id: "QTY" },
    { title: "จำนวนที่คืน", dataIndex: "ReturnQTY", key: "ReturnQTY", id: "ReturnQTY" },
    { title: "ราคาต่อหน่วย", dataIndex: "PricePerUnit", key: "PricePerUnit", id: "PricePerUnit" },
    { title: "ราคารวม", dataIndex: "Price", key: "Price", id: "Price" },
    {
      title: "ลบรายการคืน",
      id: "Action",
      dataIndex: "Action",
      key: "Action",
      render: (_: any, record: { key: number }) => (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Popconfirm
            title="คุณแน่ใจหรือไม่ว่าต้องการลบข้อมูลนี้?"
            onConfirm={() => handleDelete(record.key)} 
            okText="ใช่"
            cancelText="ไม่"
          >
            <DeleteOutlined
              style={{ cursor: "pointer", color: "red", fontSize: "20px" }}
            />
          </Popconfirm>
        </div>
      ),
    },
  ];

  const handleDownloadTemplate = () => {
    const templateColumns = [
        { title: "ลำดับ", dataIndex: "key", key: "key" }, // เพิ่ม Column "ลำดับ"
        ...columns.filter((col) => col.key !== "Action"), // เพิ่ม Column อื่นๆ
    ];
    const ws = XLSX.utils.json_to_sheet([]);
    XLSX.utils.sheet_add_aoa(ws, [templateColumns.map((col) => col.title)]);

    const mappedDataSource = dataSource.map((item) => {
        return {
            key: item.key, // เพิ่ม key
            SKU: item.SKU,
            Name: item.Name,
            QTY: item.QTY,
            ReturnQTY: item.ReturnQTY,
            PricePerUnit: item.PricePerUnit,
            Price: item.Price,
        };
    });

    XLSX.utils.sheet_add_json(ws, mappedDataSource, { origin: "A2", skipHeader: true });

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

      const json = XLSX.utils.sheet_to_json<any>(worksheet);

      console.log("JSON Data:", json);

      if (Array.isArray(json) && json.length > 0) {
          const mappedData: DataItem[] = json.map((row, index) => {
              const dataItem: DataItem = {
                  key: index + 1,
                  SKU: row["รหัสสินค้า"] as string,
                  Name: row["ชื่อสินค้า"] as string,
                  QTY: row["จำนวนเริ่มต้น"] as number,
                  ReturnQTY: row["จำนวนที่คืน"] as number,
                  PricePerUnit: row["ราคาต่อหน่วย"] as number,
                  Price: row["ราคารวม"] as number,
              };
              return dataItem;
          }).filter((item) => item.SKU && item.QTY);

          setDataSource(mappedData);

          notification.success({
              message: "อัปโหลดสำเร็จ",
              description: "ข้อมูลจากไฟล์ Excel ถูกนำเข้าเรียบร้อยแล้ว!",
          });
      } else {
          notification.error({
              message: "อัปโหลดล้มเหลว",
              description: "ไฟล์ Excel ไม่มีข้อมูลที่ถูกต้อง!",
          });
      }
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
          ReturnQTY: values.ReturnQTY,
          PricePerUnit: values.PricePerUnit,
          Price: values.ReturnQTY * values.PricePerUnit,
        };

        setDataSource([...dataSource, newData]); // เพิ่มข้อมูลใหม่ไปยัง dataSource

        notification.success({
          message: "เพิ่มสำเร็จ",
          description: "ข้อมูลของคุณถูกเพิ่มเรียบร้อยแล้ว!",
        });

        // ล้างฟิลด์ในฟอร์มหลังจากเพิ่มข้อมูลเสร็จ
        form.resetFields(["SKU", "SKU_Name", "QTY", "ReturnQTY", "PricePerUnit", "Price"]);
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

     const handleSubmit = async () => {
      try {
        // ตรวจสอบว่ามีข้อมูลในตารางอย่างน้อยหนึ่งรายการ
        if (dataSource.length === 0) {
          notification.warning({
            message: "ไม่สามารถส่งข้อมูลได้",
            description: "กรุณาเพิ่มข้อมูลในตารางก่อนส่ง!",
          });
          return; 
        }
    
        // ดึงค่าจากฟอร์ม
        const values = await form.validateFields();
    
        // เตรียมข้อมูลสำหรับส่งไปยัง API
        const requestData = {
          OrderNo: values.Order,
          SoNo: values.SO, 
          ChannelID: 2, 
          CustomerID: selectedAccount?.customerID,
          TrackingNo: values.Tracking,
          Logistic: values.Logistic,
          ReturnDate: values.Date,
          StatusReturnID: 3, 
          CreateBy: userID, 
          BeforeReturnOrderLines: dataSource.map(item => ({
            SKU: item.SKU,
            ItemName: item.Name,
            QTY: item.QTY,
            ReturnQTY: item.ReturnQTY, 
            Price: item.Price,
            TrackingNo: values.Tracking,
          })),
        };
    
        // ดึงโทเค็นจาก Local Storage
        const token = localStorage.getItem('access_token')
        const response = await api.post('/api/trade-return/create-trade', requestData, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
    
        if (response.status === 200) {
          notification.success({
            message: "ส่งข้อมูลสำเร็จ",
            description: "ข้อมูลของคุณถูกส่งเรียบร้อยแล้ว!",
          });
          form.resetFields();
          formaddress.resetFields();
          setSelectedAccount(null);
          setSelectedInvoice(null);
          setInvoiceNames([]);
          setDataSource([]);
        } else {
          notification.error({
            message: "เกิดข้อผิดพลาด",
            description: "ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง",
          });
        }
      } catch (error) {
        console.error("Error submitting data:", error);
        notification.error({
          message: "เกิดข้อผิดพลาด",
          description: "ไม่สามารถส่งข้อมูลได้ กรุณาลองใหม่อีกครั้ง",
        });
      }
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
              Order Information
            </Divider>
            <Row gutter={16} style={{ marginTop: "10px" }}>
              <Col span={8}>
                <Form.Item
                    id="Order"
                    label={
                      <span style={{ color: '#657589' }}>
                        Order Number &nbsp;
                        <Tooltip title="เลขที่ออเดอร์ ตัวอย่าง: 241215USC2YBX5">
                          <QuestionCircleOutlined
                            style={{ color: "#657589" }}
                          />
                        </Tooltip>
                      </span>}
                    name="Order"
                    rules={[{ required: true, message: "กรุณากรอกเลข Order" }]}
                >
                  <Input 
                    style={{ height: 40, width: "100%"}} 
                    placeholder="ตัวอย่าง 241215USC2YBX5" 
                  />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  id="SO"
                  label={
                    <span style={{ color: "#657589" }}>
                      Sale Order &nbsp;
                      <Tooltip title="เลขที่คำสั่งซื้อในระบบ ตัวอย่าง: SOA2410-00234">
                        <QuestionCircleOutlined
                          style={{ color: "#657589" }}
                        />
                      </Tooltip>
                    </span>
                  }
                  name="SO"
                  style={{ color: "#657589" }}
                  rules={[{ required: true }]}
                >
                  <Input
                    style={{ width: "100%", height: "40px" }}
                    placeholder="ตัวอย่าง SOA2410-00234"
                  />
                </Form.Item>
              </Col>
              <Col span={8}>
                  <Form.Item
                  id="Date"
                      label={
                        <span style={{ color: '#657589' }}>
                          วันที่ส่งคืน &nbsp;
                          <Tooltip title="วันที่คาดว่าขนส่งน่าจะถึงภายในวันนั้นหรือก่อนวันนั้น">
                            <QuestionCircleOutlined
                              style={{ color: "#657589" }}
                            />
                          </Tooltip>
                        </span>
                      }
                      name="Date"
                      rules={[{ required: true, message: 'กรุณาเลือกวันที่คืน' }]}
                  >
                      <DatePicker 
                        style={{ width: '100%', height: '40px', }} 
                        placeholder="เลือกวันที่คืน"
                      />
                  </Form.Item>
              </Col>
              <Col span={8}></Col>
            </Row>

            <Row gutter={16} style={{ marginTop: "10px" }}>
              <Col span={8}>
                <Form.Item
                  id="Tracking"
                  label={
                    <span style={{ color: "#657589" }}>
                      เลขพัสดุ &nbsp;
                      {/* <Tooltip title="เลขTracking จากขนส่ง">
                        <QuestionCircleOutlined
                          style={{ color: "#657589" }}
                        />
                      </Tooltip> */}
                    </span>
                  }
                  name="Tracking"
                  rules={[{ required: true, message: "กรุณากรอกเลขพัสดุ" }]}
                >
                  <Input 
                    style={{ height: 40, width: "100%" }}
                    placeholder="กรอกเลขพัสดุ"
                  />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  id="Logistic"
                  label={<span style={{ color: "#657589" }}>ประเภทขนส่ง</span>}
                  name="Logistic"
                  rules={[{ required: true, message: "กรุณาเลือกขนส่ง" }]}
                >
                  {isOtherTransport ? (
                    <Input
                      placeholder="กรอกประเภทขนส่ง"
                      value={transportValue}
                      onChange={(e) => setTransportValue(e.target.value)}
                      style={{ height: 40, width: "100%" }}
                    />
                  ) : (
                    <Select
                      options={TRANSPORT_TYPES}
                      placeholder="เลือกประเภทขนส่ง"
                      showSearch
                      optionFilterProp="label"
                      style={{ height: 40, width: "100%" }}
                      onChange={handleTransportChange}
                      value={transportValue}
                      listHeight={160} 
                      virtual 
                    />
                  )}
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="Customer account"
                  name="Customer_account"
                  rules={[{ required: true }]}
                  style={{ color: "#657589" }}
                >
                  <Select
                    showSearch
                    placeholder="Select Customer Account"
                    onChange={handleAccountChange}
                    loading={loading}
                    listHeight={160} 
                    virtual 
                    style={{ height: 40, width: "100%" }}
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
            </Row>

            <Row gutter={16} style={{ marginTop: "10px" }}>
              <Col span={8}>
                <Form.Item
                  label="Customer Name"
                  name="Customer_name"
                  rules={[{ required: true }]}
                >
                  <Input
                    style={{ height: 40, width: "100%" }}
                    placeholder="ข้อมูลชื่อลูกค้า"
                    value={selectedAccount?.customerName || "-"}
                    disabled
                  />
                </Form.Item>
              </Col>
              <Col span={8}>
                <Form.Item
                  label="Invoice Name"
                  name="Invoice_name"
                  style={{ color: "#657589" }}
                  rules={[{ required: true }]}
                >
                  <Select
                    showSearch
                    value={selectedAccount?.customerName || "-"}
                    placeholder="Select Invoice Name"
                    onSearch={handleInvoiceSearch} 
                    onChange={handleInvoiceChange}
                    loading={loading}
                    listHeight={160}
                    virtual
                    style={{ height: 40, width: "100%" }}
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
                <Form.Item 
                  label="Tax ID" 
                  name="Tax"
                  style={{ color: "#657589" }}
                  rules={[{ required: true }]}
                >
                  <Input 
                    style={{ height: 40, width: "100%" }}
                    placeholder="ข้อมูลเลขภาษี"
                    value={selectedAccount?.taxID || "-"} 
                    disabled 
                  />
                </Form.Item>
              </Col>
            </Row>

            <Row gutter={16} style={{ marginTop: "10px" }}>
              <Col span={18}>
                <Form.Item
                  label="Invoice Address"
                  name="Address"
                  style={{ color: "#657589" }}
                  rules={[{ required: true }]}
                  // rules={[{ required: true }]}
                >
                  <Input
                    style={{ width: "100%", height: "40px"}}
                    placeholder="ข้อมูลที่อยู่ลูกค้า"
                    value={selectedAccount?.address || "-"} 
                    disabled 
                  />
                </Form.Item>
              </Col>
              <Col span={6}>
                <Button
                  id="NewInvoiceAddress"
                  type="primary"
                  onClick={handleOpen}
                  style={{ width: "100%", height: "40px", marginTop: "30px" }}
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
            <Row gutter={16} style={{ marginTop: "10px", width: "100%", justifyContent: "center"}}>
              <Col span={7}>
                <Form.Item
                  id="SKU"
                  label={<span style={{ color: "#657589" }}>รหัสสินค้า</span>}
                  name="SKU"
                  // rules={[{ required: true, message: "กรุณากรอก SKU" }]}
                >
                  <Select
                    showSearch
                    style={{ width: "100%", height: "40px" }}
                    dropdownStyle={{ minWidth: 200 }}
                    listHeight={160}
                    placeholder="Search by SKU"
                    value={selectedSKU} // ใช้ค่าที่เลือก
                    onSearch={handleSearchSKU} // ใช้สำหรับค้นหา SKU
                    onChange={handleSKUChange} // เมื่อเลือก SKU
                    loading={loading}
                    virtual
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
                    <span style={{ color: "#657589" }}>ชื่อสินค้า</span>
                  }
                  name="SKU_Name"
                  // rules={[{ required: true, message: "กรุณาเลือก SKU Name" }]}
                >
                  <Select
                    showSearch
                    style={{ width: "100%", height: "40px" }}
                    dropdownStyle={{ minWidth: 300 }}
                    listHeight={160}
                    placeholder="Search by Product Name"
                    value={selectedName} // ใช้ค่าที่เลือก
                    onSearch={handleSearchNameAlias} // ใช้สำหรับค้นหา Name Alias
                    onChange={handleNameChange} // เมื่อเลือก Name Alias
                    loading={loading}
                    virtual 
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
            </Row>
            <Row gutter={16} style={{ marginTop: "10px", width: "100%", justifyContent: "center" }}>
              <Col span={4}>
                <Form.Item
                  id="qty"
                  label={<span style={{ color: "#657589" }}>จำนวนเริ่มต้น</span>}
                  name="QTY"
                  // rules={[{ required: true, message: "กรุณากรอก QTY" }]}
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
                  id="returnQTY"
                  label={<span style={{ color: "#657589" }}>จำนวนที่คืน</span>}
                  name="ReturnQTY"
                  // rules={[{ required: true, message: "กรุณากรอก QTY" }]}
                >
                  <InputNumber
                    min={1}
                    max={100}
                    value={returnQty}
                    onChange={(value) => setReturnQty(value)}
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
                  id="pricePerUnit"
                  label={<span style={{ color: "#657589" }}>ราคาต่อหน่วย</span>}
                  name="PricePerUnit"
                  // rules={[{ required: true, message: "กรุณากรอก Price" }]}
                >
                  <InputNumber
                    min={1}
                    max={100000}
                    value={pricePerUnit}
                    onChange={(value) => setPricePerUnit(value)}
                    step={0.01}
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
                  label={<span style={{ color: "#657589" }}>ราคารวม</span>}
                  name="Price"
                  // rules={[{ required: true, message: "กรุณากรอก Price" }]}
                >
                  <InputNumber
                    min={1}
                    max={100000}
                    value={price}
                    // onChange={(value) => setPrice(value)}
                    step={0.01}
                    disabled
                    style={{
                      width: "100%",
                      height: "40px",
                      lineHeight: "40px",
                    }}
                  />
                </Form.Item>
              </Col>
              <Col span={4}>
                <Button
                  id="add"
                  type="primary"
                  style={{ width: "100%", height: "40px", marginTop: 30 }}
                  onClick={handleAdd} 
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
                      <span style={{ color: "#657589" }}>Invoice name</span>
                    }
                    name="Invoice_Name"
                    rules={[{ required: true, message: "Please Select Invoice name" }]}
                  >
                    <Input
                      style={{ width: "400px", height: "40px" }}
                      placeholder="ข้อมูลลูกค้าตามใบสั่งซื้อ"
                      value={form.getFieldValue("Invoice_Name")} // ใช้ค่าจาก form
                      disabled
                    />
                  </Form.Item>
                </Col>

                <Col>
                  <Form.Item
                    id="้houseno"
                    label={
                      <span style={{ color: "#657589" }}>บ้านเลขที่</span>
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
                    label={<span style={{ color: "#657589" }}>จังหวัด</span>}
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
                    label={<span style={{ color: "#657589" }}>เขต</span>}
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
                    label={<span style={{ color: "#657589" }}>แขวง</span>}
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
                      <span style={{ color: "#657589" }}>รหัสไปรษณีย์</span>
                    }
                    name="PostalCode"
                    rules={[{ required: true, message: "Please Select Postcode" }]}
                  >
                    <Input
                      style={{ width: "400px", height: "40px" }}
                      placeholder="ข้อมูลเลขไปรษณีย์"
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

          <Row gutter={16} style={{ marginBottom: 20 }}>
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
        <div >
          <Table
            dataSource={dataSource.slice((currentPage - 1) * pageSize, currentPage * pageSize)} // แสดงเฉพาะจำนวนรายการที่เลือก
            columns={columns}
            rowKey="key"
            pagination={false} // ปิด pagination ใน Table
            style={{
              width: "100%",
              tableLayout: "auto",
              border: "1px solid #ddd",
              borderRadius: "8px",
            }}
            scroll={{ x: "max-content" }}
            bordered={false}
          />
        
        {showPagination && (
          <div>
            {/* showTotal แสดงอยู่เหนือ showPagination */}
            <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20 }}>
                <span style={{ fontSize: '14px', fontWeight: 'bold', color: '#555' }}>
                    ทั้งหมด <span style={{ color: '#007bff' }}>{dataSource.length}</span> รายการ
                </span>
            </div>
            <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginTop: 20, gap: 10 }}>
              {/* ปุ่มไปหน้าแรก */}
              <button
                  onClick={() => handlePageChange(1, pageSize)}
                  disabled={currentPage === 1}
                  style={{
                      fontSize: "14px",
                      // fontWeight: "bold",
                      padding: "4px 10px",
                      border: "1px solid #ddd",
                      borderRadius: "6px",
                      background: currentPage === 1 ? "#f5f5f5" : "#fff",
                      cursor: currentPage === 1 ? "not-allowed" : "pointer",
                  }}
              >
                  {"<<"}
              </button>

              {/* ปุ่มไปหน้าก่อน */}
              <button
                  onClick={() => handlePageChange(currentPage - 1, pageSize)}
                  disabled={currentPage === 1}
                  style={{
                      fontSize: "14px",
                      // fontWeight: "bold",
                      padding: "4px 10px",
                      border: "1px solid #ddd",
                      borderRadius: "6px",
                      background: currentPage === 1 ? "#f5f5f5" : "#fff",
                      cursor: currentPage === 1 ? "not-allowed" : "pointer",
                  }}
              >
                  {"<"}
              </button>

              {/* แสดงเลขหน้าแบบ [ 1 / 9 ] */}
              <span style={{ fontSize: "14px", fontWeight: 'bold' }}>
                  [ {currentPage} to {Math.ceil(dataSource.length / pageSize)} ]
              </span>

              {/* ปุ่มไปหน้าถัดไป */}
              <button
                  onClick={() => handlePageChange(currentPage + 1, pageSize)}
                  disabled={currentPage === Math.ceil(dataSource.length / pageSize)}
                  style={{
                      fontSize: "14px",
                      // fontWeight: "bold",
                      padding: "4px 10px",
                      border: "1px solid #ddd",
                      borderRadius: "6px",
                      background: currentPage === Math.ceil(dataSource.length / pageSize) ? "#f5f5f5" : "#fff",
                      cursor: currentPage === Math.ceil(dataSource.length / pageSize) ? "not-allowed" : "pointer",
                  }}
              >
                  {">"}
              </button>

              {/* ปุ่มไปหน้าสุดท้าย */}
              <button
                  onClick={() => handlePageChange(Math.ceil(dataSource.length / pageSize), pageSize)}
                  disabled={currentPage === Math.ceil(dataSource.length / pageSize)}
                  style={{
                      fontSize: "14px",
                      // fontWeight: "bold",
                      padding: "4px 10px",
                      border: "1px solid #ddd",
                      borderRadius: "6px",
                      background: currentPage === Math.ceil(dataSource.length / pageSize) ? "#f5f5f5" : "#fff",
                      cursor: currentPage === Math.ceil(dataSource.length / pageSize) ? "not-allowed" : "pointer",
                  }}
              >
                  {">>"}
              </button>

              {/* เลือกจำนวนรายการต่อหน้า */}
              <select
                  value={pageSize}
                  onChange={(e) => handlePageChange(1, Number(e.target.value))}
                  className="paginate"
                  style={{
                      fontSize: "14px",
                      fontWeight: "bold",
                      padding: "4px 10px",
                      border: "1px solid #ddd",
                      borderRadius: "6px",
                      cursor: "pointer",
                  }}
              >
                  <option value="5">5 รายการ</option>
                  <option value="10">10 รายการ</option>
                  <option value="20">20 รายการ</option>
              </select>
            </div>
            </div>
          )}
        </div>
          <Row justify="center" gutter={16}>
              <Button
                id="Submit"
                onClick={showModal}
                className="submit-trade"
              >
                Submit
              </Button>
              <Modal
                title="คุณแน่ใจหรือไม่ว่าต้องการส่งข้อมูล?"
                open={isModalVisible}
                onOk={handleOk}
                onCancel={handleCancel}
                okText="ใช่"
                cancelText="ไม่"
                centered
                style={{ textAlign: 'center'}}
                footer={
                  <div style={{ textAlign: "center" }}> {/* ทำให้ปุ่มอยู่ตรงกลาง */}
                    <Button key="ok" type="default" onClick={handleOk} style={{ marginRight: 8 }} className="button-yes">
                      Yes
                    </Button>
                    <Button key="cancel" type="dashed" onClick={handleCancel} className="button-no">
                      No
                    </Button>
                  </div>
                }

              >
              </Modal>
          </Row>
        </Layout.Content>
      </Layout>
    </ConfigProvider>
  );
};

export default CreateTradeReturn;
