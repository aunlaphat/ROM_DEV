CREATE TABLE [BeforeReturnOrder](
	[RecID] [int] IDENTITY(1,1) NOT NULL,
	[OrderNo] [varchar](50) NOT NULL,
	[SaleOrder] [varchar](20) NOT NULL,
	[SaleReturn] [varchar](50) NOT NULL,
	[ChannelID] [int] NOT NULL,
	[ReturnType] [varchar](50) NOT NULL,
	[CustomerID] [varchar](50) NOT NULL,
	[TrackingNo] [varchar](50) NOT NULL,
	[Logistic] [varchar](100) NOT NULL,
	[WarehouseID] [int] NOT NULL,
	[SoStatusID] [int] NULL,
	[MkpStatusID] [int] NULL,
	[ReturnDate] [datetime] NULL,
	[StatusReturnID] [int] NOT NULL,
	[StatusConfID] [int] NOT NULL,
	[ConfirmBy] [varchar](100) NULL,
	[CreateBy] [varchar](100) NOT NULL,
	[CreateDate] [datetime] NOT NULL,
	[UpdateBy] [varchar](100) NULL,
	[UpdateDate] [datetime] NULL,
	[CancelID] [int] NULL,
 CONSTRAINT [PK_BeforeReturnOrder] PRIMARY KEY CLUSTERED 
(
	[RecID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [BeforeReturnOrder] ADD  CONSTRAINT [DF_BeforeReturnOrder_CreateDate]  DEFAULT (getdate()) FOR [CreateDate]
GO

ALTER TABLE [BeforeReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_BeforeReturnOrder_StatusConfirm] FOREIGN KEY([StatusConfID])
REFERENCES [StatusConfirm] ([StatusConfID])
GO

ALTER TABLE [BeforeReturnOrder] CHECK CONSTRAINT [FK_BeforeReturnOrder_StatusConfirm]
GO

ALTER TABLE [BeforeReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_BeforeReturnOrder_StatusReturn] FOREIGN KEY([StatusReturnID])
REFERENCES [StatusReturn] ([StatusReturnID])
GO

ALTER TABLE [BeforeReturnOrder] CHECK CONSTRAINT [FK_BeforeReturnOrder_StatusReturn]
GO

ALTER TABLE [BeforeReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_BeforeReturnOrder_Warehouse] FOREIGN KEY([WarehouseID])
REFERENCES [Warehouse] ([WarehouseID])
GO

ALTER TABLE [BeforeReturnOrder] CHECK CONSTRAINT [FK_BeforeReturnOrder_Warehouse]
GO

---------------------

CREATE TABLE [BeforeReturnOrderLine](
	[RecID] [int] IDENTITY(1,1) NOT NULL,
	[OrderNo] [varchar](50) NULL,
	[SKU] [varchar](50) NOT NULL,
	[QTY] [int] NULL,
	[ReturnQTY] [int] NULL,
	[Price] [numeric](18, 2) NULL,
	[CreateBy] [varchar](100) NOT NULL,
	[CreateDate] [datetime] NOT NULL,
	[AlterSKU] [varchar](50) NULL,
	[UpdateBy] [varchar](100) NULL,
	[UpdateDate] [datetime] NULL,
	[TrackingNo] [varchar](50) NOT NULL
) ON [PRIMARY]
GO

ALTER TABLE [BeforeReturnOrderLine] ADD  CONSTRAINT [DF_BeforeReturnOrderLine_CreateDate]  DEFAULT (getdate()) FOR [CreateDate]
GO

---------------------

CREATE TABLE [ReturnOrder](
	[ReturnID] [varchar](50) NOT NULL,
	[OrderNo] [varchar](50) NOT NULL,
	[SaleOrder] [varchar](50) NOT NULL,
	[SaleReturn] [varchar](50) NOT NULL,
	[TrackingNo] [varchar](50) NOT NULL,
	[PlatfID] [int] NULL,
	[ChannelID] [int] NULL,
	[OptStatusID] [int] NULL,
	[AxStatusID] [int] NULL,
	[PlatfStatusID] [int] NULL,
	[Remark] [varchar](255) NULL,
	[CreateBy] [varchar](100) NOT NULL,
	[CreateDate] [datetime] NOT NULL,
	[UpdateBy] [varchar](100) NULL,
	[UpdateDate] [datetime] NULL,
	[CancelID] [int] NULL,
	[StatusCheckID] [int] NULL,
	[CheckBy] [varchar](100) NULL,
	[Description] [varchar](255) NULL,
 CONSTRAINT [PK_ReturnOrder_1] PRIMARY KEY CLUSTERED 
(
	[ReturnID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [ReturnOrder] ADD  CONSTRAINT [DF_ReturnOrder_CreateDate]  DEFAULT (getdate()) FOR [CreateDate]
GO

ALTER TABLE [ReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_ReturnOrder_Channel] FOREIGN KEY([ChannelID])
REFERENCES [Channel] ([ChannelID])
GO

ALTER TABLE [ReturnOrder] CHECK CONSTRAINT [FK_ReturnOrder_Channel]
GO

ALTER TABLE [ReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_ReturnOrder_Platforms] FOREIGN KEY([PlatfID])
REFERENCES [Platforms] ([PlatfID])
GO

ALTER TABLE [ReturnOrder] CHECK CONSTRAINT [FK_ReturnOrder_Platforms]
GO

ALTER TABLE [ReturnOrder]  WITH CHECK ADD  CONSTRAINT [FK_ReturnOrder_StatusCheck] FOREIGN KEY([StatusCheckID])
REFERENCES [StatusCheck] ([StatusCheckID])
GO

ALTER TABLE [ReturnOrder] CHECK CONSTRAINT [FK_ReturnOrder_StatusCheck]
GO

---------------------

CREATE TABLE [ReturnOrderLine](
	[RecID] [int] IDENTITY(1,1) NOT NULL,
	[ReturnID] [varchar](50) NOT NULL,
	[OrderNo] [varchar](50) NOT NULL,
	[TrackingNo] [varchar](50) NOT NULL,
	[SKU] [varchar](50) NOT NULL,
	[ReturnQTY] [int] NOT NULL,
	[CheckQTY] [int] NULL,
	[Price] [numeric](18, 2) NOT NULL,
	[CreateBy] [varchar](100) NOT NULL,
	[CreateDate] [datetime] NOT NULL,
	[AlterSKU] [varchar](50) NULL,
	[UpdateBy] [varchar](100) NULL,
	[UpdateDate] [datetime] NULL,
 CONSTRAINT [PK_ReturnOrderLine_1] PRIMARY KEY CLUSTERED 
(
	[RecID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO

----------------------------

CREATE VIEW [dbo].[V_DetailOrder]
AS
SELECT    Data_WEBChecker.dbo.[Order].OrderNo, Data_WEBChecker.dbo.[Order].SoNo, Data_WEBChecker.dbo.OperationStatus.StatusName AS StatusMKP, 
                      CASE WHEN DFIXAX63LIVE_2019.dbo.SALESTABLE.SALESSTATUS = 1 THEN 'open order' WHEN DFIXAX63LIVE_2019.dbo.SALESTABLE.SALESSTATUS = 3 THEN 'invoiced' WHEN DFIXAX63LIVE_2019.dbo.SALESTABLE.SALESSTATUS
                       = 2 THEN 'Delivered' ELSE 'unknown' END AS SalesStatus, Data_WEBChecker.dbo.OrderLine.SKU, Data_WEBChecker.dbo.OrderLine.ItemName, Data_WEBChecker.dbo.OrderLine.QTY, 
                      Data_WEBChecker.dbo.OrderLine.Price, Data_WEBChecker.dbo.[Order].CreateDate
FROM         Data_WEBChecker.dbo.[Order] INNER JOIN
                      Data_WEBChecker.dbo.OrderLine ON Data_WEBChecker.dbo.[Order].OrderNo = Data_WEBChecker.dbo.OrderLine.OrderNo INNER JOIN
                      Data_WEBChecker.dbo.OperationStatus ON Data_WEBChecker.dbo.[Order].OptStatusID = Data_WEBChecker.dbo.OperationStatus.OptStatusID LEFT OUTER JOIN
                      DFIXAX63LIVE_2019.dbo.SALESTABLE ON Data_WEBChecker.dbo.[Order].SoNo = DFIXAX63LIVE_2019.dbo.SALESTABLE.SALESID
WHERE     (Data_WEBChecker.dbo.[Order].CreateDate > CONVERT(DATETIME, '2024-11-01 00:00:00', 102)) AND (Data_WEBChecker.dbo.[Order].OptStatusID <> 8)
GO

--------------------------

CREATE VIEW [db_owner].[ROM_V_UserPermission]
AS
SELECT        TOP (100) PERCENT db_owner.UserRole.UserID, dbo.ROM_V_User.Username, dbo.ROM_V_User.Password, dbo.ROM_V_User.NickName, dbo.ROM_V_User.FullNameTH, dbo.ROM_V_User.DepartmentNo, 
                         db_owner.UserRole.RoleID, db_owner.Role.RoleName, db_owner.Role.Description, db_owner.Role.Permission
FROM            db_owner.UserRole INNER JOIN
                         db_owner.Role ON db_owner.UserRole.RoleID = db_owner.Role.RoleID INNER JOIN
                         dbo.ROM_V_User ON db_owner.UserRole.UserID = dbo.ROM_V_User.UserID
ORDER BY db_owner.UserRole.UserID
GO