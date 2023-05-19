package rpc

import (
	"net"
	"net/http"
	"time"

	"github.com/chuanmoon/utils/mtz"
	"github.com/phuslu/iploc"
)

type BodyArgs struct {
	Header http.Header
	Body   []byte
}

type EmptyArgs struct {
}

type CommonArgs struct {
	Platform   string
	AppVersion string
	OsVersion  string
	Lang       string
	TimeZone   string
	DeviceId   string
	Token      string
	ClientIp   string

	countryCode string
}

func (args *CommonArgs) getCountryCode() string {
	if args.countryCode == "" {
		args.countryCode = string(iploc.Country(net.ParseIP(args.ClientIp)))
	}
	return args.countryCode
}

func (args *CommonArgs) GetTimeZone() string {
	return mtz.TimeZone(args.TimeZone, args.getCountryCode())
}

func (args *CommonArgs) LoadLocation() *time.Location {
	return mtz.LoadLocation(args.GetTimeZone())
}
