import { useState } from "react";
import { Button, Checkbox, ConfigProvider, Descriptions, DescriptionsProps, Divider, Form, Layout, Select, Table } from "antd";
import Popup from "reactjs-popup";

const CreateExpense = () => {
    const [selectedVandorNo, setSelectedVandorNo] = useState<string | null>(null);
    const [open, setOpen] = useState(false);
    const [accountData, setAccountData] = useState([
        { Account_Name: 'อินทรา เขียวทน', Account_Number: '6566639875', Account_Bank: 'ธนาคารกรุงเทพ', Name_Contact: 'อินทรา เขียวทน', Email: 'aintra@dcom.co.th', Tel: '0896875423', selected: false },
    ]);
    const [selectedAccountName, setSelectedAccountName] = useState<string>("");

    const data = [
        { VanderNo: "V0001", VenderName: "ปิยวลี ศรีสุวรรณ", PostingDate: "2 ต.ค. 2024", Month1: "10-2024", Month2: "10-2024" },
        { VanderNo: "V0002", VenderName: "ธีรยุทธ นุชผดุง", PostingDate: "3 ต.ค. 2024", Month1: "10-2024", Month2: "10-2024" },
        { VanderNo: "V0003", VenderName: "กฤษกร ฤทธิเดช", PostingDate: "4 ต.ค. 2024", Month1: "10-2024", Month2: "10-2024" },
        { VanderNo: "V0004", VenderName: "เขนณิสสา ตีะสุ", PostingDate: "5 ต.ค. 2024", Month1: "10-2024", Month2: "10-2024" },
    ];

    const selectedVendorData = data.find((vendor) => vendor.VanderNo === selectedVandorNo);
    
    const handleOpen = () => setOpen(true);
    const handleClose = () => setOpen(false);
    
    const handleConfirm = () => {
        const selectedAccounts = accountData.filter(account => account.selected).map(account => account.Account_Name);
        setSelectedAccountName(selectedAccounts.join(", "));
        handleClose();
    };

    const items: DescriptionsProps['items'] = [
        {
            label: 'Vandor No',
            children: (
                <Select
                    style={{ width: 150 }}
                    showSearch
                    placeholder="Select a vendor"
                    optionFilterProp="children"
                    onChange={setSelectedVandorNo}
                    value={selectedVandorNo}
                    filterOption={(input, option) =>
                        (option?.children as unknown as string).toLowerCase().includes(input.toLowerCase())
                    }
                >
                    {data.map((vendor) => (
                        <Select.Option key={vendor.VanderNo} value={vendor.VanderNo}>
                            {vendor.VanderNo}
                        </Select.Option>
                    ))}
                </Select>
            ),
        },
        {
            label: 'Vandor Name',
            children: (
                <Select
                    value={selectedVendorData?.VenderName || ''}
                    disabled
                    style={{ width: 150 }}
                >
                    {selectedVendorData && (
                        <Select.Option value={selectedVendorData.VenderName}>
                            {selectedVendorData.VenderName}
                        </Select.Option>
                    )}
                </Select>
            ),
        },
        {
            label: 'Posting Date',
            children: (
                <Select
                    value={selectedVendorData?.PostingDate || ''}
                    disabled
                    style={{ width: 150 }}
                >
                    {selectedVendorData && (
                        <Select.Option value={selectedVendorData.PostingDate}>
                            {selectedVendorData.PostingDate}
                        </Select.Option>
                    )}
                </Select>
            ),
        },
        {
            label: 'Month',
            children: (
                <div style={{ display: 'flex', gap: '10px' }}>
                    <Select
                        value={selectedVendorData?.Month1 || ''}
                        disabled
                        style={{ width: 150 }}
                    >
                        {selectedVendorData && (
                            <Select.Option value={selectedVendorData.Month1}>
                                {selectedVendorData.Month1}
                            </Select.Option>
                        )}
                    </Select>
                    <div>ถึง</div>
                    <Select
                        value={selectedVendorData?.Month2 || ''}
                        disabled
                        style={{ width: 150 }}
                    >
                        {selectedVendorData && (
                            <Select.Option value={selectedVendorData.Month2}>
                                {selectedVendorData.Month2}
                            </Select.Option>
                        )}
                    </Select>
                </div>
            ),
        },
        {
            label: 'ชื่อบัญชี',
            children: (
                <>
                    <Form style={{ display: 'flex', alignItems: 'center' }}>
                        <Form.Item style={{ marginBottom: 0 }}>
                            {/* ใช้ <span> แทน <Select> เพื่อแสดงชื่อบัญชี */}
                            <span style={{ width: 150 ,border: '1px solid #d9d9d9', padding: '4px 11px', borderRadius: '5px', display: 'inline-block' }} onClick={handleOpen}>
                                {selectedAccountName || 'เลือกชื่อบัญชี'} {/* แสดงชื่อบัญชีที่เลือก */}
                            </span>
                        </Form.Item>
                    </Form>
                    <Popup 
                        open={open} 
                        onClose={handleClose} 
                        modal 
                        overlayStyle={{ background: 'rgba(0, 0, 0, 0.5)' }} 
                        contentStyle={{ 
                            background: 'white', 
                            borderRadius: '8px', 
                            boxShadow: '0 4px 8px rgba(0, 0, 0, 0.2)',
                            width: '600px', 
                            height: '600px', 
                            overflow: 'auto', 
                        }}
                    >
                        <div>
                        <div style={{ marginLeft: "28px", fontSize: "25px",  marginTop: "20px", fontWeight: "bold", color: "DodgerBlue" }}>
                        Select Account
            </div>
                            <Table 
                            style={{marginTop:'20px'}}
                                dataSource={accountData}
                                columns={[
                                    {
                                        title: 'เลือก',
                                        render: (_, record, index) => (
                                            <Checkbox
                                                checked={accountData[index]?.selected || false}
                                                onChange={(e) => {
                                                    const newData = [...accountData];
                                                    newData[index].selected = e.target.checked;
                                                    setAccountData(newData);
                                                }}
                                            />
                                        ),
                                    },
                                    {
                                        title: 'ชื่อบัญชี',
                                        dataIndex: 'Account_Name',
                                    },
                                    {
                                        title: 'หมายเลขบัญชี',
                                        dataIndex: 'Account_Number',
                                    },
                                    {
                                        title: 'ธนาคาร',
                                        dataIndex: 'Account_Bank',
                                    },
                                    {
                                        title: 'ชื่อผู้ติดต่อ',
                                        dataIndex: 'Name_Contact',
                                    },
                                    {
                                        title: 'อีเมล',
                                        dataIndex: 'Email',
                                    },
                                    {
                                        title: 'โทรศัพท์',
                                        dataIndex: 'Tel',
                                    },
                                ]}
                                pagination={false} 
                                rowKey="Account_Number" 
                            />
                            <Button 
                                type="primary" 
                                onClick={handleConfirm}
                                style={{ marginTop: '10px' }} 
                            >
                                ยืนยัน
                            </Button>
                        </div>
                    </Popup>
                </>
            ),
        }
    ];

    return (
        <ConfigProvider>
            <div style={{ marginLeft: "28px", fontSize: "25px", fontWeight: "bold", color: "DodgerBlue" }}>
                Create Expense
            </div>
            <Layout>
                <Layout.Content style={{ margin: "24px", padding: 36, minHeight: 360, background: "#fff", borderRadius: "8px" }}>
                    <Descriptions
                        bordered
                        column={{ xs: 1, sm: 2, md: 2, lg: 2, xl: 2, xxl: 2 }}
                        items={items}
                    />
                </Layout.Content>
            </Layout>
        </ConfigProvider>
    );
};

export default CreateExpense;
