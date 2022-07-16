package money

import (
	"sync"

	"github.com/rafaelsanzio/go-flashscore/pkg/applog"
	"github.com/rafaelsanzio/go-flashscore/pkg/errs"
)

var (
	// all currencies listed at https://en.wikipedia.org/wiki/ISO_4217
	// except BTN (Bhutanese ngultrum) and MKD (Macedonian denar), which
	// happen to overlap with country codes -- rather than rename all of these something like
	// USD_CURRENCY (which would work) I'm omitting until it become relevant...

	AED          = currency("United Arab Emirates dirham", "AED", "784", 2)
	AFN          = currency("Afghan afghani", "AFN", "971", 2)
	ALL          = currency("Albanian lek", "ALL", "008", 2)
	AMD          = currency("Armenian dram", "AMD", "051", 2)
	ANG          = currency("Netherlands Antillean guilder", "ANG", "532", 2)
	AOA          = currency("Angolan kwanza", "AOA", "973", 2)
	ARS          = currency("Argentine peso", "ARS", "032", 2)
	AUD          = currency("Australian dollar", "AUD", "036", 2)
	AWG          = currency("Aruban florin", "AWG", "533", 2)
	AZN          = currency("Azerbaijani manat", "AZN", "944", 2)
	BAM          = currency("Bosnia and Herzegovina", "BAM", "977", 2)
	BBD          = currency("Barbados dollar", "BBD", "052", 2)
	BDT          = currency("Bangladeshi taka", "BDT", "050", 2)
	BGN          = currency("Bulgarian lev", "BGN", "975", 2)
	BHD          = currency("Bahraini dinar", "BHD", "048", 3)
	BIF          = currency("Burundian franc", "BIF", "108", 0)
	BMD          = currency("Bermudian dollar", "BMD", "060", 2)
	BND          = currency("Brunei dollar", "BND", "096", 2)
	BOB          = currency("Boliviano", "BOB", "068", 2)
	BRL          = currency("Brazilian real", "BRL", "986", 2)
	BSD          = currency("Bahamian dollar", "BSD", "044", 2)
	BTN_CURRENCY = currency("Ngultrum", "BTN", "064", 2)
	BWP          = currency("Botswana pula", "BWP", "072", 2)
	BYN          = currency("Belarusian ruble 2016", "BYN", "933", 2)
	BYR          = currency("Belarusian ruble", "BYR", "974", 0)
	BZD          = currency("Belize dollar", "BZD", "084", 2)
	CAD          = currency("Canadian dollar", "CAD", "124", 2)
	CDF          = currency("Congolese franc", "CDF", "976", 2)
	CHF          = currency("Swiss franc", "CHF", "756", 2)
	CLP          = currency("Chilean peso", "CLP", "152", 0)
	CNY          = currency("Chinese yuan", "CNY", "156", 2)
	CNH          = currency("Chinese yuan", "CNH", "157", 2)
	CNX          = currency("Chinese Renminbi", "CNX", "158", 2)
	COP          = currency("Colombian peso", "COP", "170", 2)
	CRC          = currency("Costa Rican colon", "CRC", "188", 2)
	CUC          = currency("Cuban convertible peso", "CUC", "931", 2)
	CUP          = currency("Cuban peso", "CUP", "192", 2)
	CVE          = currency("Cape Verde escudo", "CVE", "132", 2)
	CZK          = currency("Czech koruna", "CZK", "203", 2)
	DJF          = currency("Djiboutian franc", "DJF", "262", 0)
	DKK          = currency("Danish krone", "DKK", "208", 2)
	DOP          = currency("Dominican peso", "DOP", "214", 2)
	DZD          = currency("Algerian dinar", "DZD", "012", 2)
	EGP          = currency("Egyptian pound", "EGP", "818", 2)
	ERN          = currency("Eritrean nakfa", "ERN", "232", 2)
	ESA          = currency("Bolivar Fuerte", "ESA", "996", 2)
	ETB          = currency("Ethiopian birr", "ETB", "230", 2)
	EUR          = currency("Euro", "EUR", "978", 2)
	FJD          = currency("Fiji dollar", "FJD", "242", 2)
	FKP          = currency("Falkland Islands pound", "FKP", "238", 2)
	GBP          = currency("Pound sterling", "GBP", "826", 2)
	GEL          = currency("Georgian lari", "GEL", "981", 2)
	GHS          = currency("Ghanaian cedi", "GHS", "936", 2)
	GIP          = currency("Gibraltar pound", "GIP", "292", 2)
	GMD          = currency("Gambian dalasi", "GMD", "270", 2)
	GNF          = currency("Guinean franc", "GNF", "324", 0)
	GTQ          = currency("Guatemalan quetzal", "GTQ", "320", 2)
	GYD          = currency("Guyanese dollar", "GYD", "328", 2)
	HKD          = currency("Hong Kong dollar", "HKD", "344", 2)
	HNL          = currency("Honduran lempira", "HNL", "340", 2)
	HRK          = currency("Croatian kuna", "HRK", "191", 2)
	HTG          = currency("Haitian gourde", "HTG", "332", 2)
	HUF          = currency("Hungarian forint", "HUF", "348", 2)
	IDR          = currency("Indonesian rupiah", "IDR", "360", 2)
	ILS          = currency("Israeli new shekel", "ILS", "376", 2)
	INR          = currency("Indian rupee", "INR", "356", 2)
	IQD          = currency("Iraqi dinar", "IQD", "368", 3)
	IRR          = currency("Iranian rial", "IRR", "364", 2)
	ISK          = currency("Icelandic krona", "ISK", "352", 2)
	JMD          = currency("Jamaican dollar", "JMD", "388", 2)
	JOD          = currency("Jordanian dinar", "JOD", "400", 3)
	JPY          = currency("Japanese yen", "JPY", "392", 0)
	KES          = currency("Kenyan shilling", "KES", "404", 2)
	KGS          = currency("Kyrgyzstani som", "KGS", "417", 2)
	KHR          = currency("Cambodian riel", "KHR", "116", 2)
	KMF          = currency("Comoro franc", "KMF", "174", 0)
	KPW          = currency("North Korean won", "KPW", "408", 2)
	KRW          = currency("South Korean won", "KRW", "410", 0)
	KWD          = currency("Kuwaiti dinar", "KWD", "414", 3)
	KYD          = currency("Cayman Islands dollar", "KYD", "136", 2)
	KZT          = currency("Kazakhstani tenge", "KZT", "398", 2)
	LAK          = currency("Lao kip", "LAK", "418", 2)
	LBP          = currency("Lebanese pound", "LBP", "422", 2)
	LKR          = currency("Sri Lankan rupee", "LKR", "144", 2)
	LRD          = currency("Liberian dollar", "LRD", "430", 2)
	LSL          = currency("Lesotho loti", "LSL", "426", 2)
	LYD          = currency("Libyan dinar", "LYD", "434", 3)
	MAD          = currency("Moroccan dirham", "MAD", "504", 2)
	MDL          = currency("Moldovan leu", "MDL", "498", 2)
	MGA          = currency("Malagasy Ariary", "MGA", "969", 2)
	MKD_CURRENCY = currency("Denar", "MKD", "807", 2)
	MMK          = currency("Myanmar kyat", "MMK", "104", 2)
	MNT          = currency("Mongolian togrog", "MNT", "496", 2)
	MOP          = currency("Macanese pataca", "MOP", "446", 2)
	MRO          = currency("Ouguiya", "MRO", "478", 2)
	MRU          = currency("Ouguiya", "MRU", "929", 2)
	MUR          = currency("Mauritian rupee", "MUR", "480", 2)
	MVR          = currency("Maldivian rufiyaa", "MVR", "462", 2)
	MWK          = currency("Malawian kwacha", "MWK", "454", 2)
	MXN          = currency("Mexican peso", "MXN", "484", 2)
	MYR          = currency("Malaysian ringgit", "MYR", "458", 2)
	MZN          = currency("Mozambican metical", "MZN", "943", 2)
	NAD          = currency("Namibian dollar", "NAD", "516", 2)
	NGN          = currency("Nigerian naira", "NGN", "566", 2)
	NIO          = currency("Nicaraguan cordoba", "NIO", "558", 2)
	NOK          = currency("Norwegian krone", "NOK", "578", 2)
	NPR          = currency("Nepalese rupee", "NPR", "524", 2)
	NZD          = currency("New Zealand dollar", "NZD", "554", 2)
	OMR          = currency("Omani rial", "OMR", "512", 3)
	PAB          = currency("Panamanian balboa", "PAB", "590", 2)
	PEN          = currency("Peruvian Sol", "PEN", "604", 2)
	PGK          = currency("Papua New Guinean kina", "PGK", "598", 2)
	PHP          = currency("Philippine peso", "PHP", "608", 2)
	PKR          = currency("Pakistani rupee", "PKR", "586", 2)
	PLN          = currency("Polish zloty", "PLN", "985", 2)
	PYG          = currency("Paraguayan guarani", "PYG", "600", 0)
	QAR          = currency("Qatari riyal", "QAR", "634", 2)
	RON          = currency("Romanian leu", "RON", "946", 2)
	RSD          = currency("Serbian dinar", "RSD", "941", 2)
	RUB          = currency("Russian ruble", "RUB", "643", 2)
	RWF          = currency("Rwandan franc", "RWF", "646", 0)
	SAR          = currency("Saudi riyal", "SAR", "682", 2)
	SBD          = currency("Solomon Islands dollar", "SBD", "090", 2)
	SCR          = currency("Seychelles rupee", "SCR", "690", 2)
	SDG          = currency("Sudanese pound", "SDG", "938", 2)
	SEK          = currency("Swedish krona/kronor", "SEK", "752", 2)
	SGD          = currency("Singapore dollar", "SGD", "702", 2)
	SHP          = currency("Saint Helena pound", "SHP", "654", 2)
	SLL          = currency("Sierra Leonean leone", "SLL", "694", 2)
	SOS          = currency("Somali shilling", "SOS", "706", 2)
	SRD          = currency("Surinamese dollar", "SRD", "968", 2)
	SSP          = currency("South Sudanese pound", "SSP", "728", 2)
	STD          = currency("Sao Tome and Principe dobra", "STD", "678", 2)
	STN          = currency("Dobra", "STN", "930", 2)
	SVC          = currency("El Salvador Colon", "SVC", "222", 2)
	SYP          = currency("Syrian pound", "SYP", "760", 2)
	SZL          = currency("Swazi lilangeni", "SZL", "748", 2)
	THB          = currency("Thai baht", "THB", "764", 2)
	TJS          = currency("Tajikistani somoni", "TJS", "972", 2)
	TMT          = currency("Turkmenistani manat", "TMT", "934", 2)
	TND          = currency("Tunisian dinar", "TND", "788", 3)
	TOP          = currency("Tongan paanga", "TOP", "776", 2)
	TRY          = currency("Turkish lira", "TRY", "949", 2)
	TTD          = currency("Trinidad and Tobago dollar", "TTD", "780", 2)
	TWD          = currency("New Taiwan dollar", "TWD", "901", 2)
	TZS          = currency("Tanzanian shilling", "TZS", "834", 2)
	UAH          = currency("Ukrainian hryvnia", "UAH", "980", 2)
	UGX          = currency("Ugandan shilling", "UGX", "800", 0)
	USD          = currency("United States dollar", "USD", "840", 2)
	UYU          = currency("Uruguayan peso", "UYU", "858", 2)
	UZS          = currency("Uzbekistan som", "UZS", "860", 2)
	VEF          = currency("Venezuelan bolivar", "VEF", "937", 2)
	VES          = currency("Venezuelan bolivar soberano", "VES", "928", 2)
	VND          = currency("Vietnamese dong", "VND", "704", 0)
	VUV          = currency("Vanuatu vatu", "VUV", "548", 0)
	WST          = currency("Samoan tala", "WST", "882", 2)
	XAF          = currency("CFA franc BEAC", "XAF", "950", 0)
	XCD          = currency("East Caribbean dollar", "XCD", "951", 2)
	XOF          = currency("CFA franc BCEAO", "XOF", "952", 0)
	XPF          = currency("CFP franc", "XPF", "953", 0)
	YER          = currency("Yemeni rial", "YER", "886", 2)
	ZAR          = currency("South African rand", "ZAR", "710", 2)
	ZMW          = currency("Zambian kwacha", "ZMW", "967", 2)
	ZWL          = currency("Zimbabwean Dollar", "ZWL", "932", 2)
)

var currencyLookup = sync.Map{}

// Currency contains information about a currency as defined in
// the ISO 4217 specification http://en.wikipedia.org/wiki/ISO_4217
// swagger:model currencyModel
type Currency struct {
	// example: United States dollar
	Description string
	// example: USD
	AlphaCode string // Either alpha or numeric is required
	// example: 840
	NumericCode string
	// example: 2
	Scale int // required
}

func currency(description string, alpha string, numeric string, scale int) Currency {

	c := Currency{description, alpha, numeric, scale}

	currencyLookup.Store(alpha, c)
	currencyLookup.Store(numeric, c)

	return c
}

func SafeCurrencyLookup(code string) (Currency, bool) {
	c, ok := currencyLookup.Load(code)
	if ok {
		return c.(Currency), true
	} else {
		return Currency{}, false
	}
}

// MarshalText marshals a currency to its 3-digit ISO Code
func (c Currency) MarshalText() ([]byte, error) {
	return []byte(c.AlphaCode), nil
}

// UnmarshalText unmarshals a currency struct from its 3-digit ISO Code
func (c *Currency) UnmarshalText(data []byte) error {

	code := string(data)

	curr, ok := SafeCurrencyLookup(code)

	if !ok {
		return errs.ErrInvalidCurrencyCode.Throwf(applog.Log, "Code: %s", code)
	}

	*c = curr

	return nil
}
