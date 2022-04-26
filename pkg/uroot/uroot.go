// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uroot creates root file systems from Go programs.
//
// uroot will appropriately compile the Go programs, create symlinks for their
// names, and assemble an initramfs with additional files as specified.
package uroot

import (
	"debug/elf"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	gbbgolang "github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/golang"
	"github.com/u-root/u-root/pkg/ldd"
	"github.com/u-root/u-root/pkg/uflag"
	"github.com/u-root/u-root/pkg/ulog"
	"github.com/u-root/u-root/pkg/uroot/builder"
	"github.com/u-root/u-root/pkg/uroot/initramfs"
)

// These constants are used in DefaultRamfs.
const (
	// This is the literal timezone file for GMT-0. Given that we have no
	// idea where we will be running, GMT seems a reasonable guess. If it
	// matters, setup code should download and change this to something
	// else.
	gmt0 = "TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00TZif2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x00\x04\xf8\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00GMT\x00\x00\x00\nGMT0\n"

	nameserver = "nameserver 8.8.8.8\n"
)

// DefaultRamRamfs returns a cpio.Archive for the target OS.
// If an OS is not known it will return a reasonable u-root specific
// default.
func DefaultRamfs() *cpio.Archive {
	switch golang.Default().GOOS {
	case "linux":
		return cpio.ArchiveFromRecords([]cpio.Record{
			cpio.Directory("bin", 0o755),
			cpio.Directory("dev", 0o755),
			cpio.Directory("env", 0o755),
			cpio.Directory("etc", 0o755),
			cpio.Directory("lib64", 0o755),
			cpio.Directory("proc", 0o755),
			cpio.Directory("sys", 0o755),
			cpio.Directory("tcz", 0o755),
			cpio.Directory("tmp", 0o777),
			cpio.Directory("ubin", 0o755),
			cpio.Directory("usr", 0o755),
			cpio.Directory("usr/lib", 0o755),
			cpio.Directory("var/log", 0o777),
			cpio.CharDev("dev/console", 0o600, 5, 1),
			cpio.CharDev("dev/tty", 0o666, 5, 0),
			cpio.CharDev("dev/null", 0o666, 1, 3),
			cpio.CharDev("dev/port", 0o640, 1, 4),
			cpio.CharDev("dev/urandom", 0o666, 1, 9),
			cpio.StaticFile("etc/resolv.conf", nameserver, 0o644),
			cpio.StaticFile("etc/localtime", gmt0, 0o644),
		})
	default:
		return cpio.ArchiveFromRecords([]cpio.Record{
			cpio.Directory("ubin", 0o755),
			cpio.Directory("bbin", 0o755),
		})
	}
}

// Commands specifies a list of Golang packages to build with a builder, e.g.
// in busybox mode, source mode, or binary mode.
//
// See Builder for an explanation of build modes.
type Commands struct {
	// Builder is the Go compiler mode.
	Builder builder.Builder

	// Packages are the Go commands to include (compiled or otherwise) and
	// add to the archive.
	//
	// Currently allowed formats:
	//
	//   - package imports; e.g. github.com/u-root/u-root/cmds/ls
	//   - globs of package imports; e.g. github.com/u-root/u-root/cmds/*
	//   - paths to package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
	//   - globs of paths to package directories; e.g. ./cmds/*
	//
	// Directories may be relative or absolute, with or without globs.
	// Globs are resolved using filepath.Glob.
	Packages []string

	// BinaryDir is the directory in which the resulting binaries are
	// placed inside the initramfs.
	//
	// BinaryDir may be empty, in which case Builder.DefaultBinaryDir()
	// will be used.
	BinaryDir string
}

// TargetDir returns the initramfs binary directory for these Commands.
func (c Commands) TargetDir() string {
	if len(c.BinaryDir) != 0 {
		return c.BinaryDir
	}
	return c.Builder.DefaultBinaryDir()
}

// Opts are the arguments to CreateInitramfs.
//
// Opts contains everything that influences initramfs creation such as the Go
// build environment.
type Opts struct {
	// Env is the Golang build environment (GOOS, GOARCH, etc).
	Env golang.Environ

	// Commands specify packages to build using a specific builder.
	//
	// E.g. the following will build 'ls' and 'ip' in busybox mode, but
	// 'cd' and 'cat' as separate binaries. 'cd', 'cat', 'bb', and symlinks
	// from 'ls' and 'ip' will be added to the final initramfs.
	//
	//   []Commands{
	//     Commands{
	//       Builder: builder.BusyBox,
	//       Packages: []string{
	//         "github.com/u-root/u-root/cmds/ls",
	//         "github.com/u-root/u-root/cmds/ip",
	//       },
	//     },
	//     Commands{
	//       Builder: builder.Binary,
	//       Packages: []string{
	//         "github.com/u-root/u-root/cmds/cd",
	//         "github.com/u-root/u-root/cmds/cat",
	//       },
	//     },
	//   }
	Commands []Commands

	// UrootSource is the filesystem path to the locally checked out
	// u-root source tree. This is needed to resolve templates or
	// import paths of u-root commands.
	UrootSource string

	// TempDir is a temporary directory for builders to store files in.
	TempDir string

	// ExtraFiles are files to add to the archive in addition to the Go
	// packages.
	//
	// Shared library dependencies will automatically also be added to the
	// archive using ldd, unless SkipLDD (below) is true.
	//
	// The following formats are allowed in the list:
	//
	//   - "/home/chrisko/foo:root/bar" adds the file from absolute path
	//     /home/chrisko/foo on the host at the relative root/bar in the
	//     archive.
	//   - "/home/foo" is equivalent to "/home/foo:home/foo".
	ExtraFiles []string

	// If true, do not use ldd to pick up dependencies from local machine for
	// ExtraFiles. Useful if you have all deps revision controlled and wish to
	// ensure builds are repeatable, and/or if the local machine's binaries use
	// instructions unavailable on the emulated cpu.
	//
	// If you turn this on but do not manually list all deps, affected binaries
	// will misbehave.
	SkipLDD bool

	// OutputFile is the archive output file.
	OutputFile initramfs.Writer

	// BaseArchive is an existing initramfs to include in the resulting
	// initramfs.
	BaseArchive initramfs.Reader

	// UseExistingInit determines whether the existing init from
	// BaseArchive should be used.
	//
	// If this is false, the "init" from BaseArchive will be renamed to
	// "inito" (init-original).
	UseExistingInit bool

	// InitCmd is the name of a command to link /init to.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// If this is empty, no init symlink will be created, but a user may
	// still specify a command called init or include an /init file.
	InitCmd string

	// UinitCmd is the name of a command to link /bin/uinit to.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// The u-root init will always attempt to fork/exec a uinit program,
	// and append arguments from both the kernel command-line
	// (uroot.uinitargs) as well as specified in UinitArgs.
	//
	// If this is empty, no uinit symlink will be created, but a user may
	// still specify a command called uinit or include a /bin/uinit file.
	UinitCmd string

	// UinitArgs are the arguments passed to /bin/uinit.
	UinitArgs []string

	// DefaultShell is the default shell to start after init.
	//
	// This can be an absolute path or the name of a command included in
	// Commands.
	//
	// This must be specified to have a default shell.
	DefaultShell string

	// Build options for building go binaries. Ultimate this holds all the
	// args that end up being passed to `go build`.
	BuildOpts *gbbgolang.BuildOpts
}

// CreateInitramfs creates an initramfs built to opts' specifications.
func CreateInitramfs(logger ulog.Logger, opts Opts) error {
	if _, err := os.Stat(opts.TempDir); os.IsNotExist(err) {
		return fmt.Errorf("temp dir %q must exist: %v", opts.TempDir, err)
	}
	if opts.OutputFile == nil {
		return fmt.Errorf("must give output file")
	}

	files := initramfs.NewFiles()

	// Expand commands.
	for index, cmds := range opts.Commands {
		directoryPaths, err := ResolvePackagePaths(logger, opts.Env, opts.UrootSource, cmds.Packages)
		if err != nil {
			return err
		}
		opts.Commands[index].Packages = directoryPaths
	}

	// Make sure this isn't a nil pointer
	if opts.BuildOpts == nil {
		opts.BuildOpts = &gbbgolang.BuildOpts{}
	}

	fmt.Println("DEBUG MARKER BEFORE FAIL")

	// Add each build mode's commands to the archive.
	for i, cmds := range opts.Commands {
		builderTmpDir, err := os.MkdirTemp(opts.TempDir, "builder")
		if err != nil {
			return err
		}

		// Build packages.
		bOpts := builder.Opts{
			Env:       opts.Env,
			BuildOpts: opts.BuildOpts,
			Packages:  cmds.Packages,
			TempDir:   builderTmpDir,
			BinaryDir: cmds.TargetDir(),
		}

		fmt.Printf("DEBUG MARKER %d FOR %#v\n", i, bOpts) // The printed out values look fine to me...

		if err := cmds.Builder.Build(logger, files, bOpts); err != nil {
			fmt.Printf("DEBUG ERROR %d.1 with %#v\n", i, err) // This is never printed, so it times out inside gbb
			return fmt.Errorf("error building: %v", err)
		}
	}

	fmt.Println("DEBUG MARKER AFTER FAIL")

	// Open the target initramfs file.
	archive := &initramfs.Opts{
		Files:           files,
		OutputFile:      opts.OutputFile,
		BaseArchive:     opts.BaseArchive,
		UseExistingInit: opts.UseExistingInit,
	}
	if err := ParseExtraFiles(logger, archive.Files, opts.ExtraFiles, !opts.SkipLDD); err != nil {
		return err
	}
	if err := opts.addSymlinkTo(logger, archive, opts.UinitCmd, "bin/uinit"); err != nil {
		return fmt.Errorf("%v: specify -uinitcmd=\"\" to ignore this error and build without a uinit", err)
	}
	if len(opts.UinitArgs) > 0 {
		if err := archive.AddRecord(cpio.StaticFile("etc/uinit.flags", uflag.ArgvToFile(opts.UinitArgs), 0o444)); err != nil {
			return fmt.Errorf("%v: could not add uinit arguments from UinitArgs (-uinitcmd) to initramfs", err)
		}
	}
	if err := opts.addSymlinkTo(logger, archive, opts.InitCmd, "init"); err != nil {
		return fmt.Errorf("%v: specify -initcmd=\"\" to ignore this error and build without an init (or, did you specify a list, and are you missing github.com/u-root/u-root/cmds/core/init?)", err)
	}
	if err := opts.addSymlinkTo(logger, archive, opts.DefaultShell, "bin/sh"); err != nil {
		return fmt.Errorf("%v: specify -defaultsh=\"\" to ignore this error and build without a shell", err)
	}
	if err := opts.addSymlinkTo(logger, archive, opts.DefaultShell, "bin/defaultsh"); err != nil {
		return fmt.Errorf("%v: specify -defaultsh=\"\" to ignore this error and build without a shell", err)
	}

	// Finally, write the archive.
	if err := initramfs.Write(archive); err != nil {
		return fmt.Errorf("error archiving: %v", err)
	}
	return nil
}

func (o *Opts) addSymlinkTo(logger ulog.Logger, archive *initramfs.Opts, command string, source string) error {
	if len(command) == 0 {
		return nil
	}

	target, err := resolveCommandOrPath(command, o.Commands)
	if err != nil {
		if o.Commands != nil {
			return fmt.Errorf("could not create symlink from %q to %q: %v", source, command, err)
		}
		logger.Printf("Could not create symlink from %q to %q: %v", source, command, err)
		return nil
	}

	// Make a relative symlink from /source -> target
	//
	// E.g. bin/defaultsh -> target, so you need to
	// filepath.Rel(/bin, target) since relative symlinks are
	// evaluated from their PARENT directory.
	relTarget, err := filepath.Rel(filepath.Join("/", filepath.Dir(source)), target)
	if err != nil {
		return err
	}

	if err := archive.AddRecord(cpio.Symlink(source, relTarget)); err != nil {
		return fmt.Errorf("failed to add symlink %s -> %s to initramfs: %v", source, relTarget, err)
	}
	return nil
}

// resolvePackagePath finds import paths for a single import path or directory string.
//
// Possible options are:
//
//   ./foobar
//   ./foobar/glob*
//   UROOT_SOURCE/cmds/core/ip
//   UROOT_SOURCE/cmds/core/g*lob
func resolvePackagePath(logger ulog.Logger, env golang.Environ, urootSource string, input string) ([]string, error) {
	// In case the input is a u-root import path we strip the prefix here
	// so that it can get re-added with a proper filepath.
	input = strings.TrimPrefix(input, "github.com/u-root/u-root/")

	// Search the current working directory, as well as the uroot source path if specified
	prefixes := []string{""}
	if len(urootSource) != 0 {
		prefixes = append(prefixes, urootSource)
	}

	// Resolve file system paths to package import paths.
	for _, prefix := range prefixes {
		path := filepath.Join(prefix, input)
		matches, err := filepath.Glob(path)
		if len(matches) == 0 || err != nil {
			continue
		}

		var directories []string
		for _, match := range matches {
			// Only match directories for building.
			// Skip anything that is not a directory
			fileInfo, _ := os.Stat(match)
			if !fileInfo.IsDir() {
				continue
			}
			absPath, _ := filepath.Abs(match)

			if _, err := env.PackageByPath(match); err != nil {
				logger.Printf("Skipping package %q: %v", match, err)
			} else {
				directories = append(directories, absPath)
			}
		}
		return directories, nil
	}

	// No file import paths found. Check if pkg still resolves as a package name.
	pkg, err := env.Package(input)
	if err != nil {
		return nil, fmt.Errorf("%q is neither package or path/glob: %v", input, err)
	}
	absDir, _ := filepath.Abs(pkg.Dir)
	return []string{absDir}, nil
}

func resolveCommandOrPath(cmd string, cmds []Commands) (string, error) {
	if strings.ContainsRune(cmd, filepath.Separator) {
		return cmd, nil
	}

	// Each build mode has its own binary dir (/bbin or /bin or /ubin).
	//
	// Figure out which build mode the shell is in, and symlink to that
	// build mode.
	for _, c := range cmds {
		for _, p := range c.Packages {
			if name := path.Base(p); name == cmd {
				return path.Join("/", c.TargetDir(), cmd), nil
			}
		}
	}

	return "", fmt.Errorf("command or path %q not included in u-root build", cmd)
}

// ResolvePackagePaths takes a list of u-root command paths and directories
// and turns them into exclusively absolute directory paths.
//
// Currently allowed formats:
//
//   - paths to package directories; e.g. $GOPATH/src/github.com/u-root/u-root/cmds/ls
//   - globs of paths to package directories; e.g. ./cmds/*
//   - u-root package paths; e.g. cmds/core/ls or github.com/u-root/u-root/cmds/core/ls
//   - globs of u-root package paths, e.g. cmds/* or github.com/u-root/u-root/cmds/*
//   - if an entry starts with "-" it excludes the matching package(s)
//
// Directories may be relative or absolute, with or without globs.
// Globs are resolved using filepath.Glob.
func ResolvePackagePaths(logger ulog.Logger, env golang.Environ, urootSource string, pkgs []string) ([]string, error) {
	var includes []string
	excludes := map[string]bool{}
	for _, pkg := range pkgs {
		isExclude := strings.HasPrefix(pkg, "-")

		if isExclude {
			pkg = pkg[1:]
		}
		paths, err := resolvePackagePath(logger, env, urootSource, pkg)
		if err != nil {
			return nil, err
		}
		if !isExclude {
			includes = append(includes, paths...)
		} else {
			for _, p := range paths {
				excludes[p] = true
			}
		}
	}
	var result []string
	for _, p := range includes {
		if !excludes[p] {
			result = append(result, p)
		}
	}
	return result, nil
}

// ParseExtraFiles adds files from the extraFiles list to the archive.
//
// The following formats are allowed in the extraFiles list:
//
//   - "/home/chrisko/foo:root/bar" adds the file from absolute path
//     /home/chrisko/foo on the host at the relative root/bar in the
//     archive.
//   - "/home/foo" is equivalent to "/home/foo:home/foo".
//
// ParseExtraFiles will also add ldd-listed dependencies if lddDeps is true.
func ParseExtraFiles(logger ulog.Logger, archive *initramfs.Files, extraFiles []string, lddDeps bool) error {
	var err error
	// Add files from command line.
	for _, file := range extraFiles {
		var src, dst string
		parts := strings.SplitN(file, ":", 2)
		if len(parts) == 2 {
			// treat the entry with the new src:dst syntax
			src = filepath.Clean(parts[0])
			dst = filepath.Clean(parts[1])
		} else {
			// plain old syntax
			// filepath.Clean interprets an empty string as CWD for no good reason.
			if len(file) == 0 {
				continue
			}
			src = filepath.Clean(file)
			dst = src
			if filepath.IsAbs(dst) {
				dst, err = filepath.Rel("/", dst)
				if err != nil {
					return fmt.Errorf("cannot make path relative to /: %v: %v", dst, err)
				}
			}
		}
		src, err := filepath.Abs(src)
		if err != nil {
			return fmt.Errorf("couldn't find absolute path for %q: %v", src, err)
		}
		if err := archive.AddFileNoFollow(src, dst); err != nil {
			return fmt.Errorf("couldn't add %q to archive: %v", file, err)
		}

		if lddDeps {
			// Users are frequently naming directories now, not just files.
			// Hence we must use walk here, not just check the one file.
			if err := filepath.Walk(src, func(name string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				// Try to open it as an ELF. If that fails, we can skip the ldd
				// step. The file will still be included from above.
				f, err := elf.Open(name)
				if err != nil {
					return nil
				}
				if err = f.Close(); err != nil {
					logger.Printf("WARNING: Closing ELF file %q: %v", name, err)
				}
				// Pull dependencies in the case of binaries. If `path` is not
				// a binary, `libs` will just be empty.
				libs, err := ldd.List([]string{name})
				if err != nil {
					return fmt.Errorf("WARNING: couldn't add ldd dependencies for %q: %v", name, err)
				}
				for _, lib := range libs {
					// N.B.: we already added information about the src.
					// Don't add it twice. We have to do this check here in
					// case we're renaming the src to a different dest.
					if lib == name {
						continue
					}
					if err := archive.AddFileNoFollow(lib, lib[1:]); err != nil {
						logger.Printf("WARNING: couldn't add ldd dependencies for %q: %v", lib, err)
					}
				}
				return nil
			}); err != nil {
				logger.Printf("Getting dependencies for %q: %v", src, err)
			}
		}
	}
	return nil
}

// AddCommands adds commands to the build.
func (o *Opts) AddCommands(c ...Commands) {
	o.Commands = append(o.Commands, c...)
}

func (o *Opts) AddBusyBoxCommands(pkgs ...string) {
	for i, cmds := range o.Commands {
		if cmds.Builder == builder.BusyBox {
			o.Commands[i].Packages = append(cmds.Packages, pkgs...)
			return
		}
	}

	// Not found? Add first busybox.
	o.AddCommands(BusyBoxCmds(pkgs...)...)
}

// BinaryCmds returns a list of Commands with cmds built as a busybox.
func BinaryCmds(cmds ...string) []Commands {
	if len(cmds) == 0 {
		return nil
	}
	return []Commands{
		{
			Builder:  builder.Binary,
			Packages: cmds,
		},
	}
}

// BusyBoxCmds returns a list of Commands with cmds built as a busybox.
func BusyBoxCmds(cmds ...string) []Commands {
	if len(cmds) == 0 {
		return nil
	}
	return []Commands{
		{
			Builder:  builder.BusyBox,
			Packages: cmds,
		},
	}
}
