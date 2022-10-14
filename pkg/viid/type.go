package viid

// 视频图像信息类型
type InfoType = int

// 视频图像内容信息类别
const (
	//其他
	InfoType_Other InfoType = iota
	//自动采集
	InfoType_Auto
	// 人工采集
	InfoType_Manual
)

// 是，否，不确定
const (
	// 否
	YesOrNot_NOT int = iota
	// 是
	YesOrNot_Yes
	// 不确定
	YesOrNot_Unknown
)

// B.3.5 密级代码(值类型string) 表B.7
const (
	SecretLevelType_TOP          string = "1" // TOP secret/TS
	SecretLevelType_Confidential string = "2" // Confidential
	SecretLevelType_Secret       string = "3" // Secret
	SecretLevelType_Internal     string = "4" // Internal
	SecretLevelType_Public       string = "5" // Public
	SecretLevelType_Other        string = "9" // Other secret
)

const (
	MaxLength_BasicObjectID                 = 41  // BasicObjectIDType string(41)
	MaxLength_DeviceID                      = 20  // DeviceIDType string(20)
	MaxLength_Title                         = 128 // Title and TitleNote string(128)
	MaxLength_Special                       = 128 // SpecialName string(128)
	MaxLength_KeyWord                       = 200 // KeyWord string(0...200)
	MaxLength_ContentDescription            = 256 // ContentDescription string(256)
	MaxLength_SubjectCharacter              = 256 // SubjectCharacter string(256)
	MaxLength_StoragePath                   = 256 // StoragePath string(256)
	MaxLength_ImageFormat                   = 6   // ImageFormatType string(6)
	MaxLength_PlaceCode                     = 6   //PlaceCodeType string(6)
	MaxLength_PlaceAddress                  = 100 // PlaceFullAddressType string(100)
	MaxLength_DataSourceType                = 2   //DataSourceType string(2)
	MaxLength_SecretLevelType               = 1   // SecretLevelType string(1)
	MaxLength_OrgType                       = 100 // OrgType string(0...100)
	MaxLength_ModelType                     = 100 // ModelType string(0...100)
	MaxLength_QualityGradeType              = 1   //QualityGradeType string(1)
	MaxLength_NameType                      = 50  //NameType string(0...50)
	MaxLength_IDType                        = 50  //IDType string(3)
	MaxLength_IDNumberType                  = 30  // IDNumberType string(0...30)
	MaxLength_DateTimeType                  = 17  // DateTimeFormat (YYYYMMDDHHMMSSMMM)
	MaxLength_HorizontalAndVerticalShotType = 1   // HorizontalShotType VerticalShotType string(1)
	MaxLength_ImageEigenValue               = 256 // eigen value of image,256 dimension
)

type DataSourceType = string

// 视频图像信息对象 video and image information object
const (
	VIIO_Person    = "person"
	VIIO_Face      = "face"
	VIIO_ImageInfo = "iamgeinfo"
)
