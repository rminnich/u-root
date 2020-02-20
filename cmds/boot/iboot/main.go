// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/boot/diskboot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/cmdline"
	"github.com/u-root/u-root/pkg/mount"
)

var (
	v       = flag.Bool("v", false, "Print debug messages")
	verbose = func(string, ...interface{}) {}
	dryrun  = flag.Bool("dryrun", false, "Only print out kexec commands")

	devGlob           = flag.String("dev", "/sys/class/block/*", "Device glob")
	sDeviceIndex      = flag.String("d", "", "Device index")
	sConfigIndex      = flag.String("c", "", "Config index")
	sEntryIndex       = flag.String("n", "", "Entry index")
	removeCmdlineItem = flag.String("remove", "console", "comma separated list of kernel params value to remove from parsed kernel configuration (default to console)")
	reuseCmdlineItem  = flag.String("reuse", "console", "comma separated list of kernel params value to reuse from current kernel (default to console)")
	appendCmdline     = flag.String("append", "", "Additional kernel params")

	devices []*diskboot.Device
)

func getDevice() ([]*diskboot.Device, error) {
	devices = diskboot.FindDevices(*devGlob)
	if len(devices) == 0 {
		return nil, errors.New("No devices found")
	}

	verbose("Got devices: %#v", devices)
	return devices, nil
}

func getConfig(device *diskboot.Device) ([]*diskboot.Config, error) {
	configs := device.Configs
	if len(configs) == 0 {
		return nil, errors.New("No config found")
	}

	verbose("Got configs: %#v", configs)
	return configs, nil
}

func bootEntry(config *diskboot.Config, entry *diskboot.Entry) error {
	verbose("Booting entry: %v", entry)
	filter := cmdline.NewUpdateFilter(*appendCmdline, strings.Split(*removeCmdlineItem, ","), strings.Split(*reuseCmdlineItem, ","))
	err := entry.KexecLoad(config.MountPath, filter, *dryrun)
	if err != nil {
		return fmt.Errorf("wrror doing kexec load: %v", err)
	}

	if *dryrun {
		return nil
	}

	err = kexec.Reboot()
	if err != nil {
		return fmt.Errorf("error doing kexec reboot: %v", err)
	}
	return nil
}

func cleanDevices() {
	for _, device := range devices {
		if err := device.Unmount(mount.MNT_FORCE); err != nil {
			log.Printf("Error unmounting device %s: %v", device, err)
		}
	}
}

func main() {
	flag.Parse()
	if *v {
		verbose = log.Printf
	}
	defer cleanDevices()

	devs, err := getDevice()
	if err != nil {
		log.Panic(err)
	}
	var tab []string
	for _, d := range devs {
		configs, err := getConfig(d)
		if err != nil {
			log.Print(err)
			continue
		}
		for _, c := range configs {
			for i, e := range c.Entries {
				def := ""
				if i == c.DefaultEntry {
					def = "Default"
				}
				tab = append(tab, filepath.Join(d.Device, def, e.Name))
			}
		}
	}
	log.Printf("%q", tab)
	f := NewStringCompleter(tab)

	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := t.Raw()
	if err != nil {
		log.Printf("non-fatal cannot get tty: %v", err)
	}
	defer t.Set(r)

	if false {
		//		if err := bootEntry(config, entry); err != nil {
		//			log.Panic(err)
		//		}
	}
}
