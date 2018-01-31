package cloud

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

var ExecCommand = exec.Command
var UnixStat = unix.Stat

func isMounted(targetDir string) (bool, error) {
	findmntCmd := ExecCommand("findmnt", "-n", targetDir)
	findmntStdout, err := findmntCmd.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("could not get findmount stdout pipe: %v", err.Error())
	}

	if err = findmntCmd.Start(); err != nil {
		return false, fmt.Errorf("findmnt failed to start: %v", err.Error())
	}

	findmntScanner := bufio.NewScanner(findmntStdout)
	findmntScanner.Split(bufio.ScanWords)
	findmntScanner.Scan()
	if findmntScanner.Err() != nil {
		return false, fmt.Errorf("couldn't read findnmnt output: %v", findmntScanner.Err().Error())
	}

	findmntText := findmntScanner.Text()
	if err = findmntCmd.Wait(); err != nil {
		_, isExitError := err.(*exec.ExitError)
		if !isExitError {
			return false, fmt.Errorf("findmnt failed: %v", err.Error())
		}
	}

	return findmntText == targetDir, nil
}

func currentFormat(device string) (string, error) {

	lsblkCmd := ExecCommand("lsblk", "-n", "-o", "FSTYPE", device)
	lsblkOut, err := lsblkCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("lsblk -n -o FSTYPE %s: output[%s] error[%s]", device, string(lsblkOut), err.Error())
	}

	output := strings.TrimSuffix(string(lsblkOut), "\n")
	lines := strings.Split(output, "\n")
	if lines[0] != "" {
		// The device is formatted
		return lines[0], nil
	}

	if len(lines) == 1 {
		// The device is unformatted and has no dependent devices
		return "", nil
	}

	// The device has dependent devices, most probably partitions (LVM, LUKS
	// and MD RAID are reported as FSTYPE and caught above).
	return "unknown data, probably partitions", nil
}

func Mount(targetDir string, device string, opt DefaultOptions) error {
	fsType := opt.FsType
	if fsType == "" {
		return errors.New("No filesystem type specified")
	}

	var res unix.Stat_t
	if err := UnixStat(device, &res); err != nil {
		return err
	}

	if res.Mode&unix.S_IFMT != unix.S_IFBLK {
		return fmt.Errorf("not a block device: %v", device)
	}

	if ok, err := isMounted(targetDir); err != nil || ok {
		return err
	}

	format, err := currentFormat(device)
	if err != nil {
		return err
	}
	if format != fsType {
		mkfsCmd := ExecCommand("mkfs", "-t", fsType, device)
		if mkfsOut, err := mkfsCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("could not mkfs: %v Output: %v", err.Error(), string(mkfsOut))
		}
	}

	if err := os.MkdirAll(targetDir, 0777); err != nil {
		return err
	}

	mountCmd := ExecCommand("mount", device, targetDir)
	if mountOut, err := mountCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("could not mount: %v Output: %v", err.Error(), string(mountOut))
	}

	return nil
}

func Unmount(targetDir string) error {
	ok, err := isMounted(targetDir)
	if err != nil || !ok {
		return err
	}

	umountCmd := ExecCommand("umount", targetDir)
	if umountOut, err := umountCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("could not umount: %v Output: %v", err.Error(), string(umountOut))
	}

	return nil
}
