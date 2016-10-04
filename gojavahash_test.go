package gojavahash

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestGoJavaHash(t *testing.T) {
	Convey("With a new JavaServerSelector", t, func() {
		js := &JavaServerSelector{}
		a, err := js.PickServer("test")
		So(a, ShouldBeNil)
		So(err, ShouldNotBeNil)
		Convey("Add Single Server should work", func() {
			err := js.AddServer("processdata-mc-lab2--aaaa.int.signalfuse.com:11211")
			So(err, ShouldBeNil)
			a, err := js.PickServer("test")
			So(err, ShouldBeNil)
			So(a.String(), ShouldEqual, "10.1.2.50:11211")
			Convey("Add More Servers should work", func() {
				err = js.AddServer("processdata-mc-lab3--aaab.int.signalfuse.com:11211")
				So(err, ShouldBeNil)
				err = js.AddServer("processdata-mc-lab4--aaac.int.signalfuse.com:11211")
				So(err, ShouldBeNil)
				Convey("Pick server should pick the correct server", func() {
					a, err = js.PickServer("PT:BqDQY5OAAAA:i-9caabe23")
					So(err, ShouldBeNil)
					So(a.String(), ShouldEqual, "10.1.71.12:11211")
					a, err = js.PickServer("PT:BqDQY5OAAAA:i-567480fd")
					So(err, ShouldBeNil)
					So(a.String(), ShouldEqual, "10.1.40.112:11211")
					a, err = js.PickServer("PT:BuStxseAAAA:i-316cc2b1")
					So(err, ShouldBeNil)
					So(a.String(), ShouldEqual, "10.1.2.50:11211")
				})
				Convey("For each should address all servers and return errors returned", func() {
					addrs := []net.Addr{}
					err = js.Each(func(a net.Addr) error {
						addrs = append(addrs, a)
						return nil
					})
					So(len(addrs), ShouldEqual, 3)
					So(err, ShouldBeNil)
					err = js.Each(func(a net.Addr) error {
						return errors.New("nope")
					})
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
	Convey("An invalid server name should fail", t, func() {
		_, err := New("blarg:11211")
		So(err, ShouldNotBeNil)
	})
	Convey("An invalid port should fail", t, func() {
		_, err := New("blarg11211")
		So(err, ShouldNotBeNil)
	})
	Convey("A unix name should work", t, func() {
		_, err := New("localhost/127.0.0.1:11211")
		So(err, ShouldBeNil)
	})
	Convey("With a new JavaServerSelector with small reps", t, func() {
		js := &JavaServerSelector{numReps: 1}
		err := js.AddServer("processdata-mc-lab2--aaaa.int.signalfuse.com:11211")
		So(err, ShouldBeNil)
		err = js.AddServer("processdata-mc-lab3--aaab.int.signalfuse.com:11211")
		So(err, ShouldBeNil)
		So(len(js.ring), ShouldEqual, 2)
		a, err := js.PickServer("UndWAeWGmN")
		So(err, ShouldBeNil)
		So(a.String(), ShouldEqual, "10.1.40.112:11211")
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
