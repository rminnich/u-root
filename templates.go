// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var templates = map[string][]string{
	"all": {
		"github.com/u-root/u-root/cmds/*",
	},
	// Core should be things you don't want to live without.
	"core": {
		"github.com/u-root/u-root/cmds/ansi",
		"github.com/u-root/u-root/cmds/boot",
		"github.com/u-root/u-root/cmds/cat",
		"github.com/u-root/u-root/cmds/cbmem",
		"github.com/u-root/u-root/cmds/chmod",
		"github.com/u-root/u-root/cmds/chroot",
		"github.com/u-root/u-root/cmds/cmp",
		"github.com/u-root/u-root/cmds/console",
		"github.com/u-root/u-root/cmds/cp",
		"github.com/u-root/u-root/cmds/cpio",
		"github.com/u-root/u-root/cmds/date",
		"github.com/u-root/u-root/cmds/dd",
		"github.com/u-root/u-root/cmds/df",
		"github.com/u-root/u-root/cmds/dhclient",
		"github.com/u-root/u-root/cmds/dirname",
		"github.com/u-root/u-root/cmds/dmesg",
		"github.com/u-root/u-root/cmds/echo",
		"github.com/u-root/u-root/cmds/elvish",
		"github.com/u-root/u-root/cmds/false",
		"github.com/u-root/u-root/cmds/field",
		"github.com/u-root/u-root/cmds/find",
		"github.com/u-root/u-root/cmds/free",
		"github.com/u-root/u-root/cmds/freq",
		"github.com/u-root/u-root/cmds/gpgv",
		"github.com/u-root/u-root/cmds/gpt",
		"github.com/u-root/u-root/cmds/grep",
		"github.com/u-root/u-root/cmds/gzip",
		"github.com/u-root/u-root/cmds/hexdump",
		"github.com/u-root/u-root/cmds/hostname",
		"github.com/u-root/u-root/cmds/id",
		"github.com/u-root/u-root/cmds/init",
		"github.com/u-root/u-root/cmds/insmod",
		"github.com/u-root/u-root/cmds/installcommand",
		"github.com/u-root/u-root/cmds/io",
		"github.com/u-root/u-root/cmds/ip",
		"github.com/u-root/u-root/cmds/kexec",
		"github.com/u-root/u-root/cmds/kill",
		"github.com/u-root/u-root/cmds/lddfiles",
		"github.com/u-root/u-root/cmds/ln",
		"github.com/u-root/u-root/cmds/losetup",
		"github.com/u-root/u-root/cmds/ls",
		"github.com/u-root/u-root/cmds/lsmod",
		"github.com/u-root/u-root/cmds/mkdir",
		"github.com/u-root/u-root/cmds/mkfifo",
		"github.com/u-root/u-root/cmds/mknod",
		"github.com/u-root/u-root/cmds/modprobe",
		"github.com/u-root/u-root/cmds/mount",
		"github.com/u-root/u-root/cmds/msr",
		"github.com/u-root/u-root/cmds/mv",
		"github.com/u-root/u-root/cmds/netcat",
		"github.com/u-root/u-root/cmds/ntpdate",
		"github.com/u-root/u-root/cmds/pci",
		"github.com/u-root/u-root/cmds/ping",
		"github.com/u-root/u-root/cmds/printenv",
		"github.com/u-root/u-root/cmds/ps",
		"github.com/u-root/u-root/cmds/pwd",
		"github.com/u-root/u-root/cmds/pxeboot",
		"github.com/u-root/u-root/cmds/readlink",
		"github.com/u-root/u-root/cmds/rm",
		"github.com/u-root/u-root/cmds/rmmod",
		"github.com/u-root/u-root/cmds/rsdp",
		"github.com/u-root/u-root/cmds/seq",
		"github.com/u-root/u-root/cmds/shutdown",
		"github.com/u-root/u-root/cmds/sleep",
		"github.com/u-root/u-root/cmds/sort",
		"github.com/u-root/u-root/cmds/stty",
		"github.com/u-root/u-root/cmds/switch_root",
		"github.com/u-root/u-root/cmds/sync",
		"github.com/u-root/u-root/cmds/tail",
		"github.com/u-root/u-root/cmds/tee",
		"github.com/u-root/u-root/cmds/true",
		"github.com/u-root/u-root/cmds/truncate",
		"github.com/u-root/u-root/cmds/umount",
		"github.com/u-root/u-root/cmds/uname",
		"github.com/u-root/u-root/cmds/uniq",
		"github.com/u-root/u-root/cmds/unshare",
		"github.com/u-root/u-root/cmds/validate",
		"github.com/u-root/u-root/cmds/vboot",
		"github.com/u-root/u-root/cmds/wc",
		"github.com/u-root/u-root/cmds/wget",
		"github.com/u-root/u-root/cmds/which",
	},
	// Minimal should be things you can't live without.
	"minimal": {
		"github.com/u-root/u-root/cmds/cat",
		"github.com/u-root/u-root/cmds/chmod",
		"github.com/u-root/u-root/cmds/cmp",
		"github.com/u-root/u-root/cmds/console",
		"github.com/u-root/u-root/cmds/cp",
		"github.com/u-root/u-root/cmds/date",
		"github.com/u-root/u-root/cmds/dd",
		"github.com/u-root/u-root/cmds/df",
		"github.com/u-root/u-root/cmds/dhclient",
		"github.com/u-root/u-root/cmds/dmesg",
		"github.com/u-root/u-root/cmds/echo",
		"github.com/u-root/u-root/cmds/elvish",
		"github.com/u-root/u-root/cmds/find",
		"github.com/u-root/u-root/cmds/free",
		"github.com/u-root/u-root/cmds/gpgv",
		"github.com/u-root/u-root/cmds/grep",
		"github.com/u-root/u-root/cmds/gzip",
		"github.com/u-root/u-root/cmds/hostname",
		"github.com/u-root/u-root/cmds/id",
		"github.com/u-root/u-root/cmds/init",
		"github.com/u-root/u-root/cmds/insmod",
		"github.com/u-root/u-root/cmds/io",
		"github.com/u-root/u-root/cmds/ip",
		"github.com/u-root/u-root/cmds/kexec",
		"github.com/u-root/u-root/cmds/kill",
		"github.com/u-root/u-root/cmds/ln",
		"github.com/u-root/u-root/cmds/losetup",
		"github.com/u-root/u-root/cmds/ls",
		"github.com/u-root/u-root/cmds/lsmod",
		"github.com/u-root/u-root/cmds/mkdir",
		"github.com/u-root/u-root/cmds/mknod",
		"github.com/u-root/u-root/cmds/modprobe",
		"github.com/u-root/u-root/cmds/mount",
		"github.com/u-root/u-root/cmds/msr",
		"github.com/u-root/u-root/cmds/mv",
		"github.com/u-root/u-root/cmds/pci",
		"github.com/u-root/u-root/cmds/ping",
		"github.com/u-root/u-root/cmds/printenv",
		"github.com/u-root/u-root/cmds/ps",
		"github.com/u-root/u-root/cmds/pwd",
		"github.com/u-root/u-root/cmds/readlink",
		"github.com/u-root/u-root/cmds/rm",
		"github.com/u-root/u-root/cmds/rmmod",
		"github.com/u-root/u-root/cmds/seq",
		"github.com/u-root/u-root/cmds/shutdown",
		"github.com/u-root/u-root/cmds/sleep",
		"github.com/u-root/u-root/cmds/sync",
		"github.com/u-root/u-root/cmds/tail",
		"github.com/u-root/u-root/cmds/tee",
		"github.com/u-root/u-root/cmds/truncate",
		"github.com/u-root/u-root/cmds/umount",
		"github.com/u-root/u-root/cmds/uname",
		"github.com/u-root/u-root/cmds/unshare",
		"github.com/u-root/u-root/cmds/wc",
		"github.com/u-root/u-root/cmds/wget",
		"github.com/u-root/u-root/cmds/which",
	},
	// coreboot-app minimal environment
	"coreboot-app": {
		"github.com/u-root/u-root/cmds/cat",
		"github.com/u-root/u-root/cmds/cbmem",
		"github.com/u-root/u-root/cmds/chroot",
		"github.com/u-root/u-root/cmds/console",
		"github.com/u-root/u-root/cmds/cp",
		"github.com/u-root/u-root/cmds/dd",
		"github.com/u-root/u-root/cmds/dhclient",
		"github.com/u-root/u-root/cmds/dmesg",
		"github.com/u-root/u-root/cmds/elvish",
		"github.com/u-root/u-root/cmds/find",
		"github.com/u-root/u-root/cmds/grep",
		"github.com/u-root/u-root/cmds/id",
		"github.com/u-root/u-root/cmds/init",
		"github.com/u-root/u-root/cmds/insmod",
		"github.com/u-root/u-root/cmds/ip",
		"github.com/u-root/u-root/cmds/kill",
		"github.com/u-root/u-root/cmds/ls",
		"github.com/u-root/u-root/cmds/modprobe",
		"github.com/u-root/u-root/cmds/mount",
		"github.com/u-root/u-root/cmds/pci",
		"github.com/u-root/u-root/cmds/ping",
		"github.com/u-root/u-root/cmds/ps",
		"github.com/u-root/u-root/cmds/pwd",
		"github.com/u-root/u-root/cmds/rm",
		"github.com/u-root/u-root/cmds/rmmod",
		"github.com/u-root/u-root/cmds/shutdown",
		"github.com/u-root/u-root/cmds/sshd",
		"github.com/u-root/u-root/cmds/switch_root",
		"github.com/u-root/u-root/cmds/tail",
		"github.com/u-root/u-root/cmds/tee",
		"github.com/u-root/u-root/cmds/uname",
		"github.com/u-root/u-root/cmds/wget",
	},
	// bloat are things known to build with 1.10 or newer
	// we use it to try to understand bloat
	// git checkout abc607
	// go build .
	// ./u-root -build=bb bloat
	//
	// go version go1.10 linux/amd64
	// -rwxr-xr-x 1 rminnich rminnich 11402356 Feb  7 11:07 /tmp/initramfs.linux_amd64.cpio
	// go version go1.11.5 linux/amd64
	// -rwxr-xr-x 1 rminnich rminnich 13021172 Feb  7 11:12 /tmp/initramfs.linux_amd64.cpio
	// go version devel +4b3f04c63b Thu Jan 10 18:15:48 2019 +0000 linux/amd64
	// -rwxr-xr-x 1 rminnich rminnich 13791444 Feb  7 11:15 /tmp/initramfs.linux_amd64.cpio
	// The trend is not good.
	"bloat": {
		"github.com/u-root/u-root/cmds/ansi",
		"github.com/u-root/u-root/cmds/boot",
		"github.com/u-root/u-root/cmds/cat",
		"github.com/u-root/u-root/cmds/cbmem",
		"github.com/u-root/u-root/cmds/chmod",
		"github.com/u-root/u-root/cmds/chroot",
		"github.com/u-root/u-root/cmds/cmp",
		"github.com/u-root/u-root/cmds/console",
		"github.com/u-root/u-root/cmds/cp",
		"github.com/u-root/u-root/cmds/cpio",
		"github.com/u-root/u-root/cmds/date",
		"github.com/u-root/u-root/cmds/dd",
		"github.com/u-root/u-root/cmds/df",
		"github.com/u-root/u-root/cmds/dhclient",
		"github.com/u-root/u-root/cmds/dirname",
		"github.com/u-root/u-root/cmds/dmesg",
		"github.com/u-root/u-root/cmds/echo",
		"github.com/u-root/u-root/cmds/elvish",
		"github.com/u-root/u-root/cmds/false",
		"github.com/u-root/u-root/cmds/field",
		"github.com/u-root/u-root/cmds/find",
		"github.com/u-root/u-root/cmds/free",
		"github.com/u-root/u-root/cmds/freq",
		"github.com/u-root/u-root/cmds/gpgv",
		"github.com/u-root/u-root/cmds/gpt",
		"github.com/u-root/u-root/cmds/grep",
		"github.com/u-root/u-root/cmds/gzip",
		"github.com/u-root/u-root/cmds/hexdump",
		"github.com/u-root/u-root/cmds/hostname",
		"github.com/u-root/u-root/cmds/id",
		"github.com/u-root/u-root/cmds/init",
		"github.com/u-root/u-root/cmds/insmod",
		"github.com/u-root/u-root/cmds/installcommand",
		"github.com/u-root/u-root/cmds/io",
		"github.com/u-root/u-root/cmds/ip",
		"github.com/u-root/u-root/cmds/kill",
		"github.com/u-root/u-root/cmds/lddfiles",
		"github.com/u-root/u-root/cmds/ln",
		"github.com/u-root/u-root/cmds/losetup",
		"github.com/u-root/u-root/cmds/ls",
		"github.com/u-root/u-root/cmds/lsmod",
		"github.com/u-root/u-root/cmds/mkdir",
		"github.com/u-root/u-root/cmds/mkfifo",
		"github.com/u-root/u-root/cmds/mknod",
		"github.com/u-root/u-root/cmds/modprobe",
		"github.com/u-root/u-root/cmds/mount",
		"github.com/u-root/u-root/cmds/msr",
		"github.com/u-root/u-root/cmds/mv",
		"github.com/u-root/u-root/cmds/netcat",
		"github.com/u-root/u-root/cmds/ntpdate",
		"github.com/u-root/u-root/cmds/pci",
		"github.com/u-root/u-root/cmds/ping",
		"github.com/u-root/u-root/cmds/printenv",
		"github.com/u-root/u-root/cmds/ps",
		"github.com/u-root/u-root/cmds/pwd",
		"github.com/u-root/u-root/cmds/pxeboot",
		"github.com/u-root/u-root/cmds/readlink",
		"github.com/u-root/u-root/cmds/rm",
		"github.com/u-root/u-root/cmds/rmmod",
		"github.com/u-root/u-root/cmds/rsdp",
		"github.com/u-root/u-root/cmds/seq",
		"github.com/u-root/u-root/cmds/shutdown",
		"github.com/u-root/u-root/cmds/sleep",
		"github.com/u-root/u-root/cmds/sort",
		"github.com/u-root/u-root/cmds/stty",
		"github.com/u-root/u-root/cmds/switch_root",
		"github.com/u-root/u-root/cmds/sync",
		"github.com/u-root/u-root/cmds/tail",
		"github.com/u-root/u-root/cmds/tee",
		"github.com/u-root/u-root/cmds/true",
		"github.com/u-root/u-root/cmds/truncate",
		"github.com/u-root/u-root/cmds/umount",
		"github.com/u-root/u-root/cmds/uname",
		"github.com/u-root/u-root/cmds/uniq",
		"github.com/u-root/u-root/cmds/unshare",
		"github.com/u-root/u-root/cmds/validate",
		"github.com/u-root/u-root/cmds/vboot",
		"github.com/u-root/u-root/cmds/wc",
		"github.com/u-root/u-root/cmds/wget",
		"github.com/u-root/u-root/cmds/which",
	},
}
