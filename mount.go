package common

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"fmt"

	"golang.org/x/sys/unix"
)

func isMounted(targetDir string) bool {
	findmntCmd := exec.Command("findmnt", "-n", targetDir)
	findmntStdout, err := findmntCmd.StdoutPipe()
	if err != nil {
		Fail("Could not get findmount stdout pipe: ", err.Error())
	}

	if err = findmntCmd.Start(); err != nil {
		Fail("findmnt failed to start: ", err.Error())
	}

	findmntScanner := bufio.NewScanner(findmntStdout)
	findmntScanner.Split(bufio.ScanWords)
	findmntScanner.Scan()
	if findmntScanner.Err() != nil {
		Fail("Couldn't read findnmnt output: ", findmntScanner.Err().Error())
	}

	findmntText := findmntScanner.Text()
	if err = findmntCmd.Wait(); err != nil {
		_, isExitError := err.(*exec.ExitError)
		if !isExitError {
			Fail("findmnt failed: ", err.Error())
		}
	}

	return findmntText == targetDir
}

func currentFormat(device string) (string, error) {

	lsblkCmd := exec.Command("lsblk", "-n", "-o", "FSTYPE", device)
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

func Mount(targetDir string, device string, opt DefaultOptions) Result {
	fsType := opt.FsType
	if fsType == "" {
		return Fail("No filesystem type specified")
	}

	var res unix.Stat_t
	if err := unix.Stat(device, &res); err != nil {
		return Fail("Could not stat ", device, ": ", err.Error())
	}

	if res.Mode&unix.S_IFMT != unix.S_IFBLK {
		return Fail("Not a block device: ", device)
	}

	if isMounted(targetDir) {
		return Succeed()
	}

	format, err := currentFormat(device)
	if err != nil {
		return Fail("Can not get current format ", err.Error())
	}
	if format != fsType {
		mkfsCmd := exec.Command("mkfs", "-t", fsType, device)
		if mkfsOut, err := mkfsCmd.CombinedOutput(); err != nil {
			return Fail("Could not mkfs: ", err.Error(), " Output: ", string(mkfsOut), "Target: ", targetDir, "Device ", device)
		}
	}

	if err := os.MkdirAll(targetDir, 0777); err != nil {
		return Fail("Could not create directory: ", err.Error())
	}

	mountCmd := exec.Command("mount", device, targetDir)
	if mountOut, err := mountCmd.CombinedOutput(); err != nil {
		return Fail("Could not mount: ", err.Error(), " Output: ", string(mountOut))
	}

	return Succeed("Target: ", targetDir, "Device ", device)
}

func Unmount(targetDir string) Result {
	if !isMounted(targetDir) {
		return Succeed()
	}

	umountCmd := exec.Command("umount", targetDir)
	if umountOut, err := umountCmd.CombinedOutput(); err != nil {
		return Fail("Could not umount: ", err.Error(), " Output: ", string(umountOut))
	}

	return Succeed("Target: ", targetDir)
}
