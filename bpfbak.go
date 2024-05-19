package bpfbak

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type Pinner interface {
	Pin(string) error
	Unpin() error
	Close() error
}

func detectProg(filepath string) (*ebpf.Program, error) {
	prog, err := ebpf.LoadPinnedProgram(filepath, nil)
	if err != nil {
		return nil, err
	}

	return prog, nil
}

func detectMap(filepath string) (*ebpf.Map, error) {
	m, err := ebpf.LoadPinnedMap(filepath, nil)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func detectLink(filepath string) (link.Link, error) {
	l, err := link.LoadPinnedLink(filepath, nil)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func detectObj(filepath string) (Pinner, error) {
	prog, err := detectProg(filepath)
	if err == nil {
		return prog, nil
	}

	m, err := detectMap(filepath)
	if err == nil {
		return m, nil
	}

	l, err := detectLink(filepath)
	if err == nil {
		return l, nil
	}

	return nil, fmt.Errorf("failed to detect bpf object at %s", filepath)
}

// BackupOpts specifies the options for backing up a bpf object.
type BackupOpts struct {
	// Src is the path to the bpf object to be backed up.
	Src string

	// UnpinSrc specifies whether to unpin the bpf object after backing up.
	UnpinSrc bool

	// Dst is the path to the directory where the bpf object will be pinned to.
	Dst string

	// AutoMountBpffs specifies whether to automatically mount bpffs at:
	// 1. BpffsPath if it is not empty and not already mounted.
	// 2. Dirname(Dst) if BpffsPath is empty.
	AutoMountBpffs bool

	// BpffsPath is the path to the backup directory where bpffs is mounted.
	BpffsPath string
}

func cloneObj(obj Pinner) (Pinner, error) {
	switch obj.(type) {
	case *ebpf.Program:
		prog := obj.(*ebpf.Program)
		return prog.Clone()

	case *ebpf.Map:
		m := obj.(*ebpf.Map)
		return m.Clone()

	case link.Link:
		l := obj.(link.Link)
		info, err := l.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get link info: %v", err)
		}

		return link.NewFromID(info.ID)

	default:
		return nil, fmt.Errorf("unsupported bpf object type")
	}
}

// Backup backs up a bpf object, such as bpf program, bpf map or bpf link.
func Backup(opts BackupOpts) error {
	obj, err := detectObj(opts.Src)
	if err != nil {
		return err
	}
	defer obj.Close()

	dstDir := filepath.Dir(opts.Dst)

	bpffs := opts.BpffsPath
	if bpffs == "" {
		bpffs = dstDir
	}
	bpffs = filepath.Clean(bpffs)

	if err := prepareBpffs(bpffs, opts.AutoMountBpffs); err != nil {
		return err
	}

	if !isexists(dstDir) {
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			return fmt.Errorf("failed to create dir %s: %v", dstDir, err)
		}
	}

	cloned, err := cloneObj(obj)
	if err != nil {
		return fmt.Errorf("failed to clone bpf object: %v", err)
	}
	defer cloned.Close()

	if err := cloned.Pin(opts.Dst); err != nil {
		return fmt.Errorf("failed to pin bpf object to %s: %v", opts.Dst, err)
	}

	if opts.UnpinSrc {
		if err := obj.Unpin(); err != nil {
			_ = cloned.Unpin()
			return fmt.Errorf("failed to unpin bpf object at %s: %v", opts.Src, err)
		}
	}

	return nil
}

// func BackupDir(opts BackupOpts) error {
// 	return filepath.Walk(opts.Src, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return fmt.Errorf("failed to walk %s: %v", path, err)
// 		}

// 		if info.IsDir() {
// 			return nil
// 		}

// 		dst := filepath.Join(opts.Dst, path[len(opts.Src):])
// 		return Backup(BackupOpts{
// 			Src:            path,
// 			UnpinSrc:       opts.UnpinSrc,
// 			Dst:            dst,
// 			AutoMountBpffs: opts.AutoMountBpffs,
// 			BpffsPath:      opts.BpffsPath,
// 		})
// 	})
// }
