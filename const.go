// Package line2go provides low level interfaces to LINE Thrift protocol and servers
//
// "Icecream" is a placeholder name until a better suit name have met
package line2go

// LINE service URLs
const (
	LineThriftServer    = "gd2.line.naver.jp"
	LineOSServer        = "os.line.naver.jp"
	LineOSURL           = "https://os.line.naver.jp/"
	LineStickerURL      = "http://dl.stickershop.line.naver.jp/products/"
	LineVerificationURL = "https://gd2.line.naver.jp/Q"
)

// LINE service paths
const (
	LineLoginPath   = "/api/v4/TalkService.do"
	LineCommandPath = "/S4"
	LinePollPath    = "/P4"
	LineShopPath    = "/SHOP4"
)

// Header Constants
const (
	LineUserAgent   = "purple-line (LINE for libpurple/Pidgin)"
	LineApplication = "DESKTOPWIN\t3.2.1.83\tWINDOWS\t5.1.2600-XP-x64"
)

// Toggle parameters
const (
	HTTPPrefix        = "https://"
	AppUserAgent      = "LINE2Go"
	DefaultFetchCount = 50
)
