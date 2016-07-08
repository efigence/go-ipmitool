package ipmitool

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInstance(t *testing.T) {
	i := Instance{
		IP:       "172.16.121.12",
		AuthType: "MD5",
		User:     "USERID",
		Password: "PASSW0RD",
	}
	str, err := i.GetChassisStatus()
	Convey("", t, func() {
		So(err, ShouldEqual, nil)
		So(str, ShouldEqual, nil)
	})
}
