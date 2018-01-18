package common

import (
	"bufio"
	"os"
	"os/exec"

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

	mkfsCmd := exec.Command("mkfs", "-t", fsType, device)
	if mkfsOut, err := mkfsCmd.CombinedOutput(); err != nil {
		return Fail("Could not mkfs: ", err.Error(), " Output: ", string(mkfsOut))
	}

	if err := os.MkdirAll(targetDir, 0777); err != nil {
		return Fail("Could not create directory: ", err.Error())
	}

	mountCmd := exec.Command("mount", device, targetDir)
	if mountOut, err := mountCmd.CombinedOutput(); err != nil {
		return Fail("Could not mount: ", err.Error(), " Output: ", string(mountOut))
	}

	return Succeed()
}

func Unmount(targetDir string) Result {
	if !isMounted(targetDir) {
		return Succeed()
	}

	umountCmd := exec.Command("umount", targetDir)
	if umountOut, err := umountCmd.CombinedOutput(); err != nil {
		return Fail("Could not umount: ", err.Error(), " Output: ", string(umountOut))
	}

	return Succeed()
}
