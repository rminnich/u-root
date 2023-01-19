// Copyright 2015-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// TODO: make templates able to include other templates.
// e.g. "all" below should just say "core" and "boot". Use it to replace
// the old 'systemboot' template.
// Or just call it a day, now that we have the new directory structure, and dump the templates
// completely; that may be our best bet.
var templates = map[string][]string{
	"all": {
		"github.com/u-root/u-root/cmds/core/*",
		"github.com/u-root/u-root/cmds/boot/*boot*",
	},
	"boot": {
		"github.com/u-root/u-root/cmds/boot/*boot*",
	},
	// Absolutely everything, including experimental commands.
	"world": {
		"github.com/u-root/u-root/cmds/*/*",
	},
	// Core should be things you don't want to live without.
	"core": {
		"github.com/u-root/u-root/cmds/core/*",
	},
	// Minimal should be things you can't live without.
	"minimal": {
		"github.com/u-root/u-root/cmds/core/blkid",
		"github.com/u-root/u-root/cmds/core/cat",
		"github.com/u-root/u-root/cmds/core/chmod",
		"github.com/u-root/u-root/cmds/core/cmp",
		"github.com/u-root/u-root/cmds/core/cp",
		"github.com/u-root/u-root/cmds/core/date",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/df",
		"github.com/u-root/u-root/cmds/core/dhclient",
		"github.com/u-root/u-root/cmds/core/dmesg",
		"github.com/u-root/u-root/cmds/core/echo",
		"github.com/u-root/u-root/cmds/core/elvish",
		"github.com/u-root/u-root/cmds/core/find",
		"github.com/u-root/u-root/cmds/core/free",
		"github.com/u-root/u-root/cmds/core/gpgv",
		"github.com/u-root/u-root/cmds/core/grep",
		"github.com/u-root/u-root/cmds/core/gzip",
		"github.com/u-root/u-root/cmds/core/hostname",
		"github.com/u-root/u-root/cmds/core/id",
		"github.com/u-root/u-root/cmds/core/init",
		"github.com/u-root/u-root/cmds/core/insmod",
		"github.com/u-root/u-root/cmds/core/io",
		"github.com/u-root/u-root/cmds/core/ip",
		"github.com/u-root/u-root/cmds/core/kexec",
		"github.com/u-root/u-root/cmds/core/kill",
		"github.com/u-root/u-root/cmds/core/ln",
		"github.com/u-root/u-root/cmds/core/losetup",
		"github.com/u-root/u-root/cmds/core/ls",
		"github.com/u-root/u-root/cmds/core/lsmod",
		"github.com/u-root/u-root/cmds/core/mkdir",
		"github.com/u-root/u-root/cmds/core/mknod",
		"github.com/u-root/u-root/cmds/core/mount",
		"github.com/u-root/u-root/cmds/core/msr",
		"github.com/u-root/u-root/cmds/core/mv",
		"github.com/u-root/u-root/cmds/core/pci",
		"github.com/u-root/u-root/cmds/core/ping",
		"github.com/u-root/u-root/cmds/core/printenv",
		"github.com/u-root/u-root/cmds/core/ps",
		"github.com/u-root/u-root/cmds/core/pwd",
		"github.com/u-root/u-root/cmds/core/readlink",
		"github.com/u-root/u-root/cmds/core/rm",
		"github.com/u-root/u-root/cmds/core/rmmod",
		"github.com/u-root/u-root/cmds/core/seq",
		"github.com/u-root/u-root/cmds/core/shutdown",
		"github.com/u-root/u-root/cmds/core/sleep",
		"github.com/u-root/u-root/cmds/core/sync",
		"github.com/u-root/u-root/cmds/core/tail",
		"github.com/u-root/u-root/cmds/core/tee",
		"github.com/u-root/u-root/cmds/core/truncate",
		"github.com/u-root/u-root/cmds/core/umount",
		"github.com/u-root/u-root/cmds/core/uname",
		"github.com/u-root/u-root/cmds/core/unshare",
		"github.com/u-root/u-root/cmds/core/wc",
		"github.com/u-root/u-root/cmds/core/wget",
		"github.com/u-root/u-root/cmds/core/which",
	},
	// embedded systems, like ARM based gadgets and SBCs
	"embedded": {
		"github.com/u-root/u-root/cmds/core/cat",
		"github.com/u-root/u-root/cmds/core/cp",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/dhclient",
		"github.com/u-root/u-root/cmds/core/dmesg",
		"github.com/u-root/u-root/cmds/core/echo",
		"github.com/u-root/u-root/cmds/core/elvish",
		"github.com/u-root/u-root/cmds/core/free",
		"github.com/u-root/u-root/cmds/core/grep",
		"github.com/u-root/u-root/cmds/core/init",
		"github.com/u-root/u-root/cmds/core/insmod",
		"github.com/u-root/u-root/cmds/core/ip",
		"github.com/u-root/u-root/cmds/core/kexec",
		"github.com/u-root/u-root/cmds/core/ln",
		"github.com/u-root/u-root/cmds/core/ls",
		"github.com/u-root/u-root/cmds/core/mkdir",
		"github.com/u-root/u-root/cmds/core/mount",
		"github.com/u-root/u-root/cmds/core/netcat",
		"github.com/u-root/u-root/cmds/core/ping",
		"github.com/u-root/u-root/cmds/core/rm",
		"github.com/u-root/u-root/cmds/core/rmmod",
		"github.com/u-root/u-root/cmds/core/shutdown",
		"github.com/u-root/u-root/cmds/core/tail",
		"github.com/u-root/u-root/cmds/core/tee",
		"github.com/u-root/u-root/cmds/core/uname",
		"github.com/u-root/u-root/cmds/core/wget",
	},
	// coreboot-app minimal environment
	"coreboot-app": {
		"github.com/u-root/u-root/cmds/core/cat",
		"github.com/u-root/u-root/cmds/exp/cbmem",
		"github.com/u-root/u-root/cmds/core/chroot",
		"github.com/u-root/u-root/cmds/core/cp",
		"github.com/u-root/u-root/cmds/core/dd",
		"github.com/u-root/u-root/cmds/core/dhclient",
		"github.com/u-root/u-root/cmds/core/dmesg",
		"github.com/u-root/u-root/cmds/core/elvish",
		"github.com/u-root/u-root/cmds/core/find",
		"github.com/u-root/u-root/cmds/core/grep",
		"github.com/u-root/u-root/cmds/core/id",
		"github.com/u-root/u-root/cmds/core/init",
		"github.com/u-root/u-root/cmds/core/insmod",
		"github.com/u-root/u-root/cmds/core/ip",
		"github.com/u-root/u-root/cmds/core/kill",
		"github.com/u-root/u-root/cmds/core/ls",
		"github.com/u-root/u-root/cmds/core/mount",
		"github.com/u-root/u-root/cmds/core/pci",
		"github.com/u-root/u-root/cmds/core/ping",
		"github.com/u-root/u-root/cmds/core/ps",
		"github.com/u-root/u-root/cmds/core/pwd",
		"github.com/u-root/u-root/cmds/core/rm",
		"github.com/u-root/u-root/cmds/core/rmmod",
		"github.com/u-root/u-root/cmds/core/shutdown",
		"github.com/u-root/u-root/cmds/core/sshd",
		"github.com/u-root/u-root/cmds/core/switch_root",
		"github.com/u-root/u-root/cmds/core/tail",
		"github.com/u-root/u-root/cmds/core/tee",
		"github.com/u-root/u-root/cmds/core/uname",
		"github.com/u-root/u-root/cmds/core/wget",
	},
	"plan9": {
		"github.com/u-root/u-root/cmds/core/*",
	},
	"tinygo": {
		"github.com/u-root/u-root/cmds/core/backoff",
		"github.com/u-root/u-root/cmds/core/base64",
		"github.com/u-root/u-root/cmds/core/basename",
		////"github.com/u-root/u-root/cmds/core/bind",
		//"github.com/u-root/u-root/cmds/core/blkid",
		"github.com/u-root/u-root/cmds/core/cat",
		//"github.com/u-root/u-root/cmds/core/chmod",
		//"github.com/u-root/u-root/cmds/core/chroot",
		"github.com/u-root/u-root/cmds/core/cmp",
		"github.com/u-root/u-root/cmds/core/comm",
		"github.com/u-root/u-root/cmds/core/cp",
		//////"github.com/u-root/u-root/cmds/core/cpio",
		"github.com/u-root/u-root/cmds/core/date",
		"github.com/u-root/u-root/cmds/core/dd",
		//////"github.com/u-root/u-root/cmds/core/df",
		//////"github.com/u-root/u-root/cmds/core/dhclient",
		"github.com/u-root/u-root/cmds/core/dirname",
		//////"github.com/u-root/u-root/cmds/core/dmesg",
		"github.com/u-root/u-root/cmds/core/echo",
		//////"github.com/u-root/u-root/cmds/core/elvish",
		"github.com/u-root/u-root/cmds/core/false",
		//"github.com/u-root/u-root/cmds/core/find",
		"github.com/u-root/u-root/cmds/core/free",
		//"github.com/u-root/u-root/cmds/core/fusermount",
		//////"github.com/u-root/u-root/cmds/core/gosh",
		//"github.com/u-root/u-root/cmds/core/gpgv",
		//"github.com/u-root/u-root/cmds/core/gpt",
		"github.com/u-root/u-root/cmds/core/grep",
		//"github.com/u-root/u-root/cmds/core/gzip",
		"github.com/u-root/u-root/cmds/core/hexdump",
		//////"github.com/u-root/u-root/cmds/core/hostname",
		//"github.com/u-root/u-root/cmds/core/hwclock",
		//"github.com/u-root/u-root/cmds/core/id",
		//////"github.com/u-root/u-root/cmds/core/init",
		//"github.com/u-root/u-root/cmds/core/insmod",
		//////"github.com/u-root/u-root/cmds/core/io",
		//"github.com/u-root/u-root/cmds/core/ip",
		//"github.com/u-root/u-root/cmds/core/kexec",
		"github.com/u-root/u-root/cmds/core/kill",
		//"github.com/u-root/u-root/cmds/core/lddfiles",
		//"github.com/u-root/u-root/cmds/core/ln",
		//"github.com/u-root/u-root/cmds/core/lockmsrs",
		//"github.com/u-root/u-root/cmds/core/losetup",
		"github.com/u-root/u-root/cmds/core/ls",
		//"github.com/u-root/u-root/cmds/core/lsdrivers",
		//"github.com/u-root/u-root/cmds/core/lsmod",
		//"github.com/u-root/u-root/cmds/core/man",
		//"github.com/u-root/u-root/cmds/core/md5sum",
		"github.com/u-root/u-root/cmds/core/mkdir",
		//////"github.com/u-root/u-root/cmds/core/mkfifo",
		//////"github.com/u-root/u-root/cmds/core/mknod",
		//"github.com/u-root/u-root/cmds/core/mktemp",
		"github.com/u-root/u-root/cmds/core/more",
		//"github.com/u-root/u-root/cmds/core/mount",
		//"github.com/u-root/u-root/cmds/core/msr",
		"github.com/u-root/u-root/cmds/core/mv",
		//"github.com/u-root/u-root/cmds/core/netcat",
		//////"github.com/u-root/u-root/cmds/core/ntpdate",
		//"github.com/u-root/u-root/cmds/core/pci",
		//"github.com/u-root/u-root/cmds/core/ping",
		//"github.com/u-root/u-root/cmds/core/poweroff",
		//"github.com/u-root/u-root/cmds/core/printenv",
		"github.com/u-root/u-root/cmds/core/ps",
		"github.com/u-root/u-root/cmds/core/pwd",
		//"github.com/u-root/u-root/cmds/core/readlink",
		"github.com/u-root/u-root/cmds/core/rm",
		//"github.com/u-root/u-root/cmds/core/rmmod",
		//"github.com/u-root/u-root/cmds/core/rsdp",
		//"github.com/u-root/u-root/cmds/core/scp",
		//"github.com/u-root/u-root/cmds/core/seq",
		//"github.com/u-root/u-root/cmds/core/shasum",
		//"github.com/u-root/u-root/cmds/core/shutdown",
		//"github.com/u-root/u-root/cmds/core/sleep",
		//"github.com/u-root/u-root/cmds/core/sluinit",
		"github.com/u-root/u-root/cmds/core/sort",
		//////"github.com/u-root/u-root/cmds/core/sshd",
		//////"github.com/u-root/u-root/cmds/core/strace",
		//"github.com/u-root/u-root/cmds/core/strings",
		//////"github.com/u-root/u-root/cmds/core/stty",
		//////"github.com/u-root/u-root/cmds/core/switch_root",
		//"github.com/u-root/u-root/cmds/core/sync",
		//"github.com/u-root/u-root/cmds/core/tail",
		//"github.com/u-root/u-root/cmds/core/tar",
		//"github.com/u-root/u-root/cmds/core/tee",
		//"github.com/u-root/u-root/cmds/core/time",
		//"github.com/u-root/u-root/cmds/core/timeout",
		"github.com/u-root/u-root/cmds/core/tr",
		"github.com/u-root/u-root/cmds/core/true",
		//"github.com/u-root/u-root/cmds/core/truncate",
		//"github.com/u-root/u-root/cmds/core/ts",
		//"github.com/u-root/u-root/cmds/core/umount",
		//////"github.com/u-root/u-root/cmds/core/uname",
		"github.com/u-root/u-root/cmds/core/uniq",
		//"github.com/u-root/u-root/cmds/core/unmount",
		//"github.com/u-root/u-root/cmds/core/unshare",
		//"github.com/u-root/u-root/cmds/core/uptime",
		//"github.com/u-root/u-root/cmds/core/watchdog",
		//"github.com/u-root/u-root/cmds/core/watchdogd",
		"github.com/u-root/u-root/cmds/core/wc",
		//"github.com/u-root/u-root/cmds/core/wget",
		//////"github.com/u-root/u-root/cmds/core/which",
		"github.com/u-root/u-root/cmds/core/yes",
		//// One of these commands in exp tickles a busybox bug.
		////"github.com/u-root/u-root/cmds/exp/acpicat",
		////"github.com/u-root/u-root/cmds/exp/acpigrep",
		////"github.com/u-root/u-root/cmds/exp/ansi",
		////"github.com/u-root/u-root/cmds/exp/bootvars",
		////"github.com/u-root/u-root/cmds/exp/bzimage",
		////"github.com/u-root/u-root/cmds/exp/cbmem",
		////"github.com/u-root/u-root/cmds/exp/console",
		////"github.com/u-root/u-root/cmds/exp/crc",
		////"github.com/u-root/u-root/cmds/exp/disk_unlock",
		////"github.com/u-root/u-root/cmds/exp/dmidecode",
		////"github.com/u-root/u-root/cmds/exp/dumpebda",
		////"github.com/u-root/u-root/cmds/exp/ectool",
		////"github.com/u-root/u-root/cmds/exp/ed",
		////"github.com/u-root/u-root/cmds/exp/efivarfs",
		////"github.com/u-root/u-root/cmds/exp/esxiboot",
		////"github.com/u-root/u-root/cmds/exp/fbsplash",
		////"github.com/u-root/u-root/cmds/exp/fdtdump",
		////"github.com/u-root/u-root/cmds/exp/field",
		////"github.com/u-root/u-root/cmds/exp/fixrsdp",
		////"github.com/u-root/u-root/cmds/exp/forth",
		////"github.com/u-root/u-root/cmds/exp/freq",
		////"github.com/u-root/u-root/cmds/exp/getty",
		////"github.com/u-root/u-root/cmds/exp/gosh",
		////"github.com/u-root/u-root/cmds/exp/hdparm",
		////"github.com/u-root/u-root/cmds/exp/ipmidump",
		////"github.com/u-root/u-root/cmds/exp/kconf",
		////"github.com/u-root/u-root/cmds/exp/lsfabric",
		////"github.com/u-root/u-root/cmds/exp/madeye",
		////"github.com/u-root/u-root/cmds/exp/modprobe",
		////"github.com/u-root/u-root/cmds/exp/netbootxyz",
		////"github.com/u-root/u-root/cmds/exp/newsshd",
		////"github.com/u-root/u-root/cmds/exp/nvme_unlock",
		////"github.com/u-root/u-root/cmds/exp/page",
		////"github.com/u-root/u-root/cmds/exp/partprobe",
		////"github.com/u-root/u-root/cmds/exp/pflask",
		////"github.com/u-root/u-root/cmds/exp/pox",
		////"github.com/u-root/u-root/cmds/exp/pxeserver",
		////"github.com/u-root/u-root/cmds/exp/readpe",
		////"github.com/u-root/u-root/cmds/exp/run",
		////"github.com/u-root/u-root/cmds/exp/rush",
		////"github.com/u-root/u-root/cmds/exp/smbios_transfer",
		////"github.com/u-root/u-root/cmds/exp/smn",
		////"github.com/u-root/u-root/cmds/exp/srvfiles",
		////"github.com/u-root/u-root/cmds/exp/ssh",
		//////"github.com/u-root/u-root/cmds/exp/syscallfilter",
		////"github.com/u-root/u-root/cmds/exp/tac",
		////"github.com/u-root/u-root/cmds/exp/tcz",
		////"github.com/u-root/u-root/cmds/exp/uefiboot",
		////"github.com/u-root/u-root/cmds/exp/vboot",
		////"github.com/u-root/u-root/cmds/exp/watch",
		////"github.com/u-root/u-root/cmds/exp/zbi",
		////"github.com/u-root/u-root/cmds/exp/zimage",
	},
}
