// Copyright 2018 Google LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strace

import (
	"golang.org/x/sys/unix"
)

const archWidth = 64

// This is the arm64 syscall map. One might think that this one map could be used for all Linux
// flavors on all architectures. Ah, no. It's Linux, not Plan 9. Every arch has a different
// system call set.
var syscalls = SyscallMap{
	unix.SYS_READ:  makeSyscallInfo("read", Hex, ReadBuffer, Hex),
	unix.SYS_WRITE: makeSyscallInfo("write", Hex, WriteBuffer, Hex),

	unix.SYS_CLOSE: makeSyscallInfo("close", Hex),

	unix.SYS_FSTAT: makeSyscallInfo("fstat", Hex, Stat),

	unix.SYS_LSEEK:          makeSyscallInfo("lseek", Hex, Hex, Hex),
	unix.SYS_MMAP:           makeSyscallInfo("mmap", Hex, Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MPROTECT:       makeSyscallInfo("mprotect", Hex, Hex, Hex),
	unix.SYS_MUNMAP:         makeSyscallInfo("munmap", Hex, Hex),
	unix.SYS_BRK:            makeSyscallInfo("brk", Hex),
	unix.SYS_RT_SIGACTION:   makeSyscallInfo("rt_sigaction", Hex, Hex, Hex),
	unix.SYS_RT_SIGPROCMASK: makeSyscallInfo("rt_sigprocmask", Hex, Hex, Hex, Hex),
	unix.SYS_RT_SIGRETURN:   makeSyscallInfo("rt_sigreturn"),
	unix.SYS_IOCTL:          makeSyscallInfo("ioctl", Hex, Hex, Hex),
	unix.SYS_PREAD64:        makeSyscallInfo("pread64", Hex, ReadBuffer, Hex, Hex),
	unix.SYS_PWRITE64:       makeSyscallInfo("pwrite64", Hex, WriteBuffer, Hex, Hex),
	unix.SYS_READV:          makeSyscallInfo("readv", Hex, ReadIOVec, Hex),
	unix.SYS_WRITEV:         makeSyscallInfo("writev", Hex, WriteIOVec, Hex),

	unix.SYS_SCHED_YIELD: makeSyscallInfo("sched_yield"),
	unix.SYS_MREMAP:      makeSyscallInfo("mremap", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MSYNC:       makeSyscallInfo("msync", Hex, Hex, Hex),
	unix.SYS_MINCORE:     makeSyscallInfo("mincore", Hex, Hex, Hex),
	unix.SYS_MADVISE:     makeSyscallInfo("madvise", Hex, Hex, Hex),
	unix.SYS_SHMGET:      makeSyscallInfo("shmget", Hex, Hex, Hex),
	unix.SYS_SHMAT:       makeSyscallInfo("shmat", Hex, Hex, Hex),
	unix.SYS_SHMCTL:      makeSyscallInfo("shmctl", Hex, Hex, Hex),
	unix.SYS_DUP:         makeSyscallInfo("dup", Hex),

	unix.SYS_NANOSLEEP: makeSyscallInfo("nanosleep", Timespec, PostTimespec),
	unix.SYS_GETITIMER: makeSyscallInfo("getitimer", ItimerType, PostItimerVal),

	unix.SYS_SETITIMER:   makeSyscallInfo("setitimer", ItimerType, ItimerVal, PostItimerVal),
	unix.SYS_GETPID:      makeSyscallInfo("getpid"),
	unix.SYS_SENDFILE:    makeSyscallInfo("sendfile", Hex, Hex, Hex, Hex),
	unix.SYS_SOCKET:      makeSyscallInfo("socket", SockFamily, SockType, SockProtocol),
	unix.SYS_CONNECT:     makeSyscallInfo("connect", Hex, SockAddr, Hex),
	unix.SYS_ACCEPT:      makeSyscallInfo("accept", Hex, PostSockAddr, SockLen),
	unix.SYS_SENDTO:      makeSyscallInfo("sendto", Hex, Hex, Hex, Hex, SockAddr, Hex),
	unix.SYS_RECVFROM:    makeSyscallInfo("recvfrom", Hex, Hex, Hex, Hex, PostSockAddr, SockLen),
	unix.SYS_SENDMSG:     makeSyscallInfo("sendmsg", Hex, SendMsgHdr, Hex),
	unix.SYS_RECVMSG:     makeSyscallInfo("recvmsg", Hex, RecvMsgHdr, Hex),
	unix.SYS_SHUTDOWN:    makeSyscallInfo("shutdown", Hex, Hex),
	unix.SYS_BIND:        makeSyscallInfo("bind", Hex, SockAddr, Hex),
	unix.SYS_LISTEN:      makeSyscallInfo("listen", Hex, Hex),
	unix.SYS_GETSOCKNAME: makeSyscallInfo("getsockname", Hex, PostSockAddr, SockLen),
	unix.SYS_GETPEERNAME: makeSyscallInfo("getpeername", Hex, PostSockAddr, SockLen),
	unix.SYS_SOCKETPAIR:  makeSyscallInfo("socketpair", SockFamily, SockType, SockProtocol, Hex),
	unix.SYS_SETSOCKOPT:  makeSyscallInfo("setsockopt", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_GETSOCKOPT:  makeSyscallInfo("getsockopt", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_CLONE:       makeSyscallInfo("clone", CloneFlags, Hex, Hex, Hex, Hex),

	unix.SYS_EXECVE:    makeSyscallInfo("execve", Path, ExecveStringVector, ExecveStringVector),
	unix.SYS_EXIT:      makeSyscallInfo("exit", Hex),
	unix.SYS_WAIT4:     makeSyscallInfo("wait4", Hex, Hex, Hex, Rusage),
	unix.SYS_KILL:      makeSyscallInfo("kill", Hex, Hex),
	unix.SYS_UNAME:     makeSyscallInfo("uname", Uname),
	unix.SYS_SEMGET:    makeSyscallInfo("semget", Hex, Hex, Hex),
	unix.SYS_SEMOP:     makeSyscallInfo("semop", Hex, Hex, Hex),
	unix.SYS_SEMCTL:    makeSyscallInfo("semctl", Hex, Hex, Hex, Hex),
	unix.SYS_SHMDT:     makeSyscallInfo("shmdt", Hex),
	unix.SYS_MSGGET:    makeSyscallInfo("msgget", Hex, Hex),
	unix.SYS_MSGSND:    makeSyscallInfo("msgsnd", Hex, Hex, Hex, Hex),
	unix.SYS_MSGRCV:    makeSyscallInfo("msgrcv", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MSGCTL:    makeSyscallInfo("msgctl", Hex, Hex, Hex),
	unix.SYS_FCNTL:     makeSyscallInfo("fcntl", Hex, Hex, Hex),
	unix.SYS_FLOCK:     makeSyscallInfo("flock", Hex, Hex),
	unix.SYS_FSYNC:     makeSyscallInfo("fsync", Hex),
	unix.SYS_FDATASYNC: makeSyscallInfo("fdatasync", Hex),
	unix.SYS_TRUNCATE:  makeSyscallInfo("truncate", Path, Hex),
	unix.SYS_FTRUNCATE: makeSyscallInfo("ftruncate", Hex, Hex),

	unix.SYS_GETCWD: makeSyscallInfo("getcwd", PostPath, Hex),
	unix.SYS_CHDIR:  makeSyscallInfo("chdir", Path),
	unix.SYS_FCHDIR: makeSyscallInfo("fchdir", Hex),

	unix.SYS_FCHMOD: makeSyscallInfo("fchmod", Hex, Mode),

	unix.SYS_FCHOWN: makeSyscallInfo("fchown", Hex, Hex, Hex),

	unix.SYS_UMASK:        makeSyscallInfo("umask", Hex),
	unix.SYS_GETTIMEOFDAY: makeSyscallInfo("gettimeofday", Timeval, Hex),
	unix.SYS_GETRLIMIT:    makeSyscallInfo("getrlimit", Hex, Hex),
	unix.SYS_GETRUSAGE:    makeSyscallInfo("getrusage", Hex, Rusage),
	unix.SYS_SYSINFO:      makeSyscallInfo("sysinfo", Hex),
	unix.SYS_TIMES:        makeSyscallInfo("times", Hex),
	unix.SYS_PTRACE:       makeSyscallInfo("ptrace", PtraceRequest, Hex, Hex, Hex),
	unix.SYS_GETUID:       makeSyscallInfo("getuid"),
	unix.SYS_SYSLOG:       makeSyscallInfo("syslog", Hex, Hex, Hex),
	unix.SYS_GETGID:       makeSyscallInfo("getgid"),
	unix.SYS_SETUID:       makeSyscallInfo("setuid", Hex),
	unix.SYS_SETGID:       makeSyscallInfo("setgid", Hex),
	unix.SYS_GETEUID:      makeSyscallInfo("geteuid"),
	unix.SYS_GETEGID:      makeSyscallInfo("getegid"),
	unix.SYS_SETPGID:      makeSyscallInfo("setpgid", Hex, Hex),
	unix.SYS_GETPPID:      makeSyscallInfo("getppid"),

	unix.SYS_SETSID:          makeSyscallInfo("setsid"),
	unix.SYS_SETREUID:        makeSyscallInfo("setreuid", Hex, Hex),
	unix.SYS_SETREGID:        makeSyscallInfo("setregid", Hex, Hex),
	unix.SYS_GETGROUPS:       makeSyscallInfo("getgroups", Hex, Hex),
	unix.SYS_SETGROUPS:       makeSyscallInfo("setgroups", Hex, Hex),
	unix.SYS_SETRESUID:       makeSyscallInfo("setresuid", Hex, Hex, Hex),
	unix.SYS_GETRESUID:       makeSyscallInfo("getresuid", Hex, Hex, Hex),
	unix.SYS_SETRESGID:       makeSyscallInfo("setresgid", Hex, Hex, Hex),
	unix.SYS_GETRESGID:       makeSyscallInfo("getresgid", Hex, Hex, Hex),
	unix.SYS_GETPGID:         makeSyscallInfo("getpgid", Hex),
	unix.SYS_SETFSUID:        makeSyscallInfo("setfsuid", Hex),
	unix.SYS_SETFSGID:        makeSyscallInfo("setfsgid", Hex),
	unix.SYS_GETSID:          makeSyscallInfo("getsid", Hex),
	unix.SYS_CAPGET:          makeSyscallInfo("capget", Hex, Hex),
	unix.SYS_CAPSET:          makeSyscallInfo("capset", Hex, Hex),
	unix.SYS_RT_SIGPENDING:   makeSyscallInfo("rt_sigpending", Hex),
	unix.SYS_RT_SIGTIMEDWAIT: makeSyscallInfo("rt_sigtimedwait", Hex, Hex, Timespec, Hex),
	unix.SYS_RT_SIGQUEUEINFO: makeSyscallInfo("rt_sigqueueinfo", Hex, Hex, Hex),
	unix.SYS_RT_SIGSUSPEND:   makeSyscallInfo("rt_sigsuspend", Hex),
	unix.SYS_SIGALTSTACK:     makeSyscallInfo("sigaltstack", Hex, Hex),

	unix.SYS_PERSONALITY: makeSyscallInfo("personality", Hex),

	unix.SYS_STATFS:  makeSyscallInfo("statfs", Path, Hex),
	unix.SYS_FSTATFS: makeSyscallInfo("fstatfs", Hex, Hex),

	unix.SYS_GETPRIORITY:            makeSyscallInfo("getpriority", Hex, Hex),
	unix.SYS_SETPRIORITY:            makeSyscallInfo("setpriority", Hex, Hex, Hex),
	unix.SYS_SCHED_SETPARAM:         makeSyscallInfo("sched_setparam", Hex, Hex),
	unix.SYS_SCHED_GETPARAM:         makeSyscallInfo("sched_getparam", Hex, Hex),
	unix.SYS_SCHED_SETSCHEDULER:     makeSyscallInfo("sched_setscheduler", Hex, Hex, Hex),
	unix.SYS_SCHED_GETSCHEDULER:     makeSyscallInfo("sched_getscheduler", Hex),
	unix.SYS_SCHED_GET_PRIORITY_MAX: makeSyscallInfo("sched_get_priority_max", Hex),
	unix.SYS_SCHED_GET_PRIORITY_MIN: makeSyscallInfo("sched_get_priority_min", Hex),
	unix.SYS_SCHED_RR_GET_INTERVAL:  makeSyscallInfo("sched_rr_get_interval", Hex, Hex),
	unix.SYS_MLOCK:                  makeSyscallInfo("mlock", Hex, Hex),
	unix.SYS_MUNLOCK:                makeSyscallInfo("munlock", Hex, Hex),
	unix.SYS_MLOCKALL:               makeSyscallInfo("mlockall", Hex),
	unix.SYS_MUNLOCKALL:             makeSyscallInfo("munlockall"),
	unix.SYS_VHANGUP:                makeSyscallInfo("vhangup"),

	unix.SYS_PIVOT_ROOT: makeSyscallInfo("pivot_root", Hex, Hex),

	unix.SYS_PRCTL: makeSyscallInfo("prctl", Hex, Hex, Hex, Hex, Hex),

	unix.SYS_ADJTIMEX:      makeSyscallInfo("adjtimex", Hex),
	unix.SYS_SETRLIMIT:     makeSyscallInfo("setrlimit", Hex, Hex),
	unix.SYS_CHROOT:        makeSyscallInfo("chroot", Path),
	unix.SYS_SYNC:          makeSyscallInfo("sync"),
	unix.SYS_ACCT:          makeSyscallInfo("acct", Hex),
	unix.SYS_SETTIMEOFDAY:  makeSyscallInfo("settimeofday", Timeval, Hex),
	unix.SYS_MOUNT:         makeSyscallInfo("mount", Path, Path, Path, Hex, Path),
	unix.SYS_UMOUNT2:       makeSyscallInfo("umount2", Path, Hex),
	unix.SYS_SWAPON:        makeSyscallInfo("swapon", Hex, Hex),
	unix.SYS_SWAPOFF:       makeSyscallInfo("swapoff", Hex),
	unix.SYS_REBOOT:        makeSyscallInfo("reboot", Hex, Hex, Hex, Hex),
	unix.SYS_SETHOSTNAME:   makeSyscallInfo("sethostname", Hex, Hex),
	unix.SYS_SETDOMAINNAME: makeSyscallInfo("setdomainname", Hex, Hex),

	unix.SYS_INIT_MODULE:   makeSyscallInfo("init_module", Hex, Hex, Hex),
	unix.SYS_DELETE_MODULE: makeSyscallInfo("delete_module", Hex, Hex),

	//	unix.SYS_QUERY_MODULE:query_module (only present in Linux < 2.6)
	unix.SYS_QUOTACTL:   makeSyscallInfo("quotactl", Hex, Hex, Hex, Hex),
	unix.SYS_NFSSERVCTL: makeSyscallInfo("nfsservctl", Hex, Hex, Hex),
	// 	unix.SYS_GETPMSG:getpmsg (not implemented in the Linux kernel)
	// 	unix.SYS_PUTPMSG:putpmsg (not implemented in the Linux kernel)
	// 	unix.SYSCALL:afs_syscall (not implemented in the Linux kernel)
	// 	unix.SYS_TUXCALL:tuxcall (not implemented in the Linux kernel)
	// 	unix.SYS_SECURITY:security (not implemented in the Linux kernel)
	unix.SYS_GETTID:       makeSyscallInfo("gettid"),
	unix.SYS_READAHEAD:    makeSyscallInfo("readahead", Hex, Hex, Hex),
	unix.SYS_SETXATTR:     makeSyscallInfo("setxattr", Path, Path, Hex, Hex, Hex),
	unix.SYS_LSETXATTR:    makeSyscallInfo("lsetxattr", Path, Path, Hex, Hex, Hex),
	unix.SYS_FSETXATTR:    makeSyscallInfo("fsetxattr", Hex, Path, Hex, Hex, Hex),
	unix.SYS_GETXATTR:     makeSyscallInfo("getxattr", Path, Path, Hex, Hex),
	unix.SYS_LGETXATTR:    makeSyscallInfo("lgetxattr", Path, Path, Hex, Hex),
	unix.SYS_FGETXATTR:    makeSyscallInfo("fgetxattr", Hex, Path, Hex, Hex),
	unix.SYS_LISTXATTR:    makeSyscallInfo("listxattr", Path, Path, Hex),
	unix.SYS_LLISTXATTR:   makeSyscallInfo("llistxattr", Path, Path, Hex),
	unix.SYS_FLISTXATTR:   makeSyscallInfo("flistxattr", Hex, Path, Hex),
	unix.SYS_REMOVEXATTR:  makeSyscallInfo("removexattr", Path, Path),
	unix.SYS_LREMOVEXATTR: makeSyscallInfo("lremovexattr", Path, Path),
	unix.SYS_FREMOVEXATTR: makeSyscallInfo("fremovexattr", Hex, Path),
	unix.SYS_TKILL:        makeSyscallInfo("tkill", Hex, Hex),

	unix.SYS_FUTEX:             makeSyscallInfo("futex", Hex, FutexOp, Hex, Timespec, Hex, Hex),
	unix.SYS_SCHED_SETAFFINITY: makeSyscallInfo("sched_setaffinity", Hex, Hex, Hex),
	unix.SYS_SCHED_GETAFFINITY: makeSyscallInfo("sched_getaffinity", Hex, Hex, Hex),

	unix.SYS_IO_SETUP:     makeSyscallInfo("io_setup", Hex, Hex),
	unix.SYS_IO_DESTROY:   makeSyscallInfo("io_destroy", Hex),
	unix.SYS_IO_GETEVENTS: makeSyscallInfo("io_getevents", Hex, Hex, Hex, Hex, Timespec),
	unix.SYS_IO_SUBMIT:    makeSyscallInfo("io_submit", Hex, Hex, Hex),
	unix.SYS_IO_CANCEL:    makeSyscallInfo("io_cancel", Hex, Hex, Hex),

	unix.SYS_LOOKUP_DCOOKIE: makeSyscallInfo("lookup_dcookie", Hex, Hex, Hex),

	// 	unix.SYS_EPOLL_CTL_OLD:epoll_ctl_old (not implemented in the Linux kernel)
	// 	unix.SYS_EPOLL_WAIT_OLD:epoll_wait_old (not implemented in the Linux kernel)
	unix.SYS_REMAP_FILE_PAGES: makeSyscallInfo("remap_file_pages", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_GETDENTS64:       makeSyscallInfo("getdents64", Hex, Hex, Hex),
	unix.SYS_SET_TID_ADDRESS:  makeSyscallInfo("set_tid_address", Hex),
	unix.SYS_RESTART_SYSCALL:  makeSyscallInfo("restart_syscall"),
	unix.SYS_SEMTIMEDOP:       makeSyscallInfo("semtimedop", Hex, Hex, Hex, Hex),
	unix.SYS_FADVISE64:        makeSyscallInfo("fadvise64", Hex, Hex, Hex, Hex),
	unix.SYS_TIMER_CREATE:     makeSyscallInfo("timer_create", Hex, Hex, Hex),
	unix.SYS_TIMER_SETTIME:    makeSyscallInfo("timer_settime", Hex, Hex, ItimerSpec, PostItimerSpec),
	unix.SYS_TIMER_GETTIME:    makeSyscallInfo("timer_gettime", Hex, PostItimerSpec),
	unix.SYS_TIMER_GETOVERRUN: makeSyscallInfo("timer_getoverrun", Hex),
	unix.SYS_TIMER_DELETE:     makeSyscallInfo("timer_delete", Hex),
	unix.SYS_CLOCK_SETTIME:    makeSyscallInfo("clock_settime", Hex, Timespec),
	unix.SYS_CLOCK_GETTIME:    makeSyscallInfo("clock_gettime", Hex, PostTimespec),
	unix.SYS_CLOCK_GETRES:     makeSyscallInfo("clock_getres", Hex, PostTimespec),
	unix.SYS_CLOCK_NANOSLEEP:  makeSyscallInfo("clock_nanosleep", Hex, Hex, Timespec, PostTimespec),
	unix.SYS_EXIT_GROUP:       makeSyscallInfo("exit_group", Hex),

	unix.SYS_EPOLL_CTL: makeSyscallInfo("epoll_ctl", Hex, Hex, Hex, Hex),
	unix.SYS_TGKILL:    makeSyscallInfo("tgkill", Hex, Hex, Hex),

	// 	unix.SYS_VSERVER:vserver (not implemented in the Linux kernel)
	unix.SYS_MBIND:           makeSyscallInfo("mbind", Hex, Hex, Hex, Hex, Hex, Hex),
	unix.SYS_SET_MEMPOLICY:   makeSyscallInfo("set_mempolicy", Hex, Hex, Hex),
	unix.SYS_GET_MEMPOLICY:   makeSyscallInfo("get_mempolicy", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MQ_OPEN:         makeSyscallInfo("mq_open", Hex, Hex, Hex, Hex),
	unix.SYS_MQ_UNLINK:       makeSyscallInfo("mq_unlink", Hex),
	unix.SYS_MQ_TIMEDSEND:    makeSyscallInfo("mq_timedsend", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MQ_TIMEDRECEIVE: makeSyscallInfo("mq_timedreceive", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_MQ_NOTIFY:       makeSyscallInfo("mq_notify", Hex, Hex),
	unix.SYS_MQ_GETSETATTR:   makeSyscallInfo("mq_getsetattr", Hex, Hex, Hex),
	unix.SYS_KEXEC_LOAD:      makeSyscallInfo("kexec_load", Hex, Hex, Hex, Hex),
	unix.SYS_WAITID:          makeSyscallInfo("waitid", Hex, Hex, Hex, Hex, Rusage),
	unix.SYS_ADD_KEY:         makeSyscallInfo("add_key", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_REQUEST_KEY:     makeSyscallInfo("request_key", Hex, Hex, Hex, Hex),
	unix.SYS_KEYCTL:          makeSyscallInfo("keyctl", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_IOPRIO_SET:      makeSyscallInfo("ioprio_set", Hex, Hex, Hex),
	unix.SYS_IOPRIO_GET:      makeSyscallInfo("ioprio_get", Hex, Hex),

	unix.SYS_INOTIFY_ADD_WATCH: makeSyscallInfo("inotify_add_watch", Hex, Hex, Hex),
	unix.SYS_INOTIFY_RM_WATCH:  makeSyscallInfo("inotify_rm_watch", Hex, Hex),
	unix.SYS_MIGRATE_PAGES:     makeSyscallInfo("migrate_pages", Hex, Hex, Hex, Hex),
	unix.SYS_OPENAT:            makeSyscallInfo("openat", Hex, Path, OpenFlags, Mode),
	unix.SYS_MKDIRAT:           makeSyscallInfo("mkdirat", Hex, Path, Hex),
	unix.SYS_MKNODAT:           makeSyscallInfo("mknodat", Hex, Path, Mode, Hex),
	unix.SYS_FCHOWNAT:          makeSyscallInfo("fchownat", Hex, Path, Hex, Hex, Hex),

	unix.SYS_UNLINKAT:        makeSyscallInfo("unlinkat", Hex, Path, Hex),
	unix.SYS_RENAMEAT:        makeSyscallInfo("renameat", Hex, Path, Hex, Path),
	unix.SYS_LINKAT:          makeSyscallInfo("linkat", Hex, Path, Hex, Path, Hex),
	unix.SYS_SYMLINKAT:       makeSyscallInfo("symlinkat", Path, Hex, Path),
	unix.SYS_READLINKAT:      makeSyscallInfo("readlinkat", Hex, Path, ReadBuffer, Hex),
	unix.SYS_FCHMODAT:        makeSyscallInfo("fchmodat", Hex, Path, Mode),
	unix.SYS_FACCESSAT:       makeSyscallInfo("faccessat", Hex, Path, Oct, Hex),
	unix.SYS_PSELECT6:        makeSyscallInfo("pselect6", Hex, Hex, Hex, Hex, Hex, Hex),
	unix.SYS_PPOLL:           makeSyscallInfo("ppoll", Hex, Hex, Timespec, Hex, Hex),
	unix.SYS_UNSHARE:         makeSyscallInfo("unshare", Hex),
	unix.SYS_SET_ROBUST_LIST: makeSyscallInfo("set_robust_list", Hex, Hex),
	unix.SYS_GET_ROBUST_LIST: makeSyscallInfo("get_robust_list", Hex, Hex, Hex),
	unix.SYS_SPLICE:          makeSyscallInfo("splice", Hex, Hex, Hex, Hex, Hex, Hex),
	unix.SYS_TEE:             makeSyscallInfo("tee", Hex, Hex, Hex, Hex),
	unix.SYS_SYNC_FILE_RANGE: makeSyscallInfo("sync_file_range", Hex, Hex, Hex, Hex),
	unix.SYS_VMSPLICE:        makeSyscallInfo("vmsplice", Hex, Hex, Hex, Hex),
	unix.SYS_MOVE_PAGES:      makeSyscallInfo("move_pages", Hex, Hex, Hex, Hex, Hex, Hex),
	unix.SYS_UTIMENSAT:       makeSyscallInfo("utimensat", Hex, Path, UTimeTimespec, Hex),
	unix.SYS_EPOLL_PWAIT:     makeSyscallInfo("epoll_pwait", Hex, Hex, Hex, Hex, Hex, Hex),

	unix.SYS_TIMERFD_CREATE: makeSyscallInfo("timerfd_create", Hex, Hex),

	unix.SYS_FALLOCATE:         makeSyscallInfo("fallocate", Hex, Hex, Hex, Hex),
	unix.SYS_TIMERFD_SETTIME:   makeSyscallInfo("timerfd_settime", Hex, Hex, ItimerSpec, PostItimerSpec),
	unix.SYS_TIMERFD_GETTIME:   makeSyscallInfo("timerfd_gettime", Hex, PostItimerSpec),
	unix.SYS_ACCEPT4:           makeSyscallInfo("accept4", Hex, PostSockAddr, SockLen, SockFlags),
	unix.SYS_SIGNALFD4:         makeSyscallInfo("signalfd4", Hex, Hex, Hex, Hex),
	unix.SYS_EVENTFD2:          makeSyscallInfo("eventfd2", Hex, Hex),
	unix.SYS_EPOLL_CREATE1:     makeSyscallInfo("epoll_create1", Hex),
	unix.SYS_DUP3:              makeSyscallInfo("dup3", Hex, Hex, Hex),
	unix.SYS_PIPE2:             makeSyscallInfo("pipe2", PipeFDs, Hex),
	unix.SYS_INOTIFY_INIT1:     makeSyscallInfo("inotify_init1", Hex),
	unix.SYS_PREADV:            makeSyscallInfo("preadv", Hex, ReadIOVec, Hex, Hex),
	unix.SYS_PWRITEV:           makeSyscallInfo("pwritev", Hex, WriteIOVec, Hex, Hex),
	unix.SYS_RT_TGSIGQUEUEINFO: makeSyscallInfo("rt_tgsigqueueinfo", Hex, Hex, Hex, Hex),
	unix.SYS_PERF_EVENT_OPEN:   makeSyscallInfo("perf_event_open", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_RECVMMSG:          makeSyscallInfo("recvmmsg", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_FANOTIFY_INIT:     makeSyscallInfo("fanotify_init", Hex, Hex),
	unix.SYS_FANOTIFY_MARK:     makeSyscallInfo("fanotify_mark", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_PRLIMIT64:         makeSyscallInfo("prlimit64", Hex, Hex, Hex, Hex),
	unix.SYS_NAME_TO_HANDLE_AT: makeSyscallInfo("name_to_handle_at", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_OPEN_BY_HANDLE_AT: makeSyscallInfo("open_by_handle_at", Hex, Hex, Hex),
	unix.SYS_CLOCK_ADJTIME:     makeSyscallInfo("clock_adjtime", Hex, Hex),
	unix.SYS_SYNCFS:            makeSyscallInfo("syncfs", Hex),
	unix.SYS_SENDMMSG:          makeSyscallInfo("sendmmsg", Hex, Hex, Hex, Hex),
	unix.SYS_SETNS:             makeSyscallInfo("setns", Hex, Hex),
	unix.SYS_GETCPU:            makeSyscallInfo("getcpu", Hex, Hex, Hex),
	unix.SYS_PROCESS_VM_READV:  makeSyscallInfo("process_vm_readv", Hex, ReadIOVec, Hex, IOVec, Hex, Hex),
	unix.SYS_PROCESS_VM_WRITEV: makeSyscallInfo("process_vm_writev", Hex, IOVec, Hex, WriteIOVec, Hex, Hex),
	unix.SYS_KCMP:              makeSyscallInfo("kcmp", Hex, Hex, Hex, Hex, Hex),
	unix.SYS_FINIT_MODULE:      makeSyscallInfo("finit_module", Hex, Hex, Hex),
	unix.SYS_SCHED_SETATTR:     makeSyscallInfo("sched_setattr", Hex, Hex, Hex),
	unix.SYS_SCHED_GETATTR:     makeSyscallInfo("sched_getattr", Hex, Hex, Hex),
	unix.SYS_RENAMEAT2:         makeSyscallInfo("renameat2", Hex, Path, Hex, Path, Hex),
	unix.SYS_SECCOMP:           makeSyscallInfo("seccomp", Hex, Hex, Hex),
}

// FillArgs pulls the correct registers to populate system call arguments
// and the system call number into a TraceRecord. Note that the system
// call number is not technically an argument. This is good, in a sense,
// since it makes the function arguments end up in "the right place"
// from the point of view of the caller. The performance improvement is
// negligible, as you can see by a look at the GNU runtime.
func (s *SyscallEvent) FillArgs() {
	s.Args = SyscallArguments{
		{uintptr(s.Regs.Regs[1])},
		{uintptr(s.Regs.Regs[2])},
		{uintptr(s.Regs.Regs[3])},
		{uintptr(s.Regs.Regs[4])},
		{uintptr(s.Regs.Regs[5])},
		{uintptr(s.Regs.Regs[6])}}
	s.Sysno = int(uint32(s.Regs.Regs[0]))
}

// FillRet fills the TraceRecord with the result values from the registers.
func (s *SyscallEvent) FillRet() {
	s.Ret = [2]SyscallArgument{{uintptr(s.Regs.Regs[0])}}
	if errno := int(s.Regs.Regs[0]); errno < 0 {
		s.Errno = unix.Errno(-errno)
	}
}
