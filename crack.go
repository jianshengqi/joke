package main

/*
#include<stdlib.h>

static void exec(const char* s) {
    system(s);
}
*/
import "C"
import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func genDstProcName(procName string) string {
	return os.TempDir() + string(filepath.Separator) + procName
}

func copyFile(srcName, dstName string) (err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

func terminateOldProcess(procNames ...string) {
	for _, procName := range procNames {
		crackInfo("kill %v\n", procName)
		C.exec(C.CString(fmt.Sprintf("taskkill /im %v -f", procName)))
	}
}

func copyBinFiles(procNames ...string) error {
	fixSuffix := func(in string) string {
		const suffix = ".exe"
		if strings.HasSuffix(in, suffix) {
			return in
		}
		return in + suffix
	}
	crackInfo("copying files:%v\n", procNames)
	for _, procName := range procNames {
		src, dst := os.Args[0], genDstProcName(procName)
		crackInfo("src:%v,dst:%v\n", fixSuffix(src), dst)
		err := copyFile(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func execMultiProc(procNames ...string) ([]func(), error) {
	crackInfo("exec process:%v\n", procNames)
	var cancelArr []func()
	changeCmdEnv := func(cmd *exec.Cmd) *exec.Cmd {
		_ = os.Setenv(crackSubProcess, "1")
		defer os.Unsetenv(crackSubProcess)
		cmd.Env = os.Environ()
		return cmd
	}
	for _, procName := range procNames {
		crackInfo("start the %v.\n", procName)
		ctx, cancel := context.WithCancel(context.Background())
		execErr := changeCmdEnv(exec.CommandContext(ctx, genDstProcName(procName))).Start()
		if execErr != nil {
			return nil, fmt.Errorf("fn:execMultiProc,exec %v execErr:%v", procName, execErr)
		}
		cancelArr = append(cancelArr, cancel)
	}
	return cancelArr, nil
}

func moveBinary(procNames ...string) {
	crackInfo("move process:%v\n", procNames)
	for _, procName := range procNames {
		C.exec(C.CString(fmt.Sprintf("move %v %v", genDstProcName(procName), strings.TrimSuffix(genDstProcName(procName), ".exe"))))
		checkErr(copyFile(procName, genDstProcName(procName)))
	}
}

func replaceRealFiles(procNames ...string) {
	crackInfo("replace real files:%v\n", procNames)
	for _, procName := range procNames {
		func() {
			binData, err := bindataRead(
				_edpaExe,
				procName,
			)
			checkErr(err)

			dst, err := os.OpenFile(procName, os.O_WRONLY|os.O_CREATE, 0644)
			checkErr(err)
			defer dst.Close()
			_, err = dst.Write(binData)
			checkErr(err)
		}()
	}
}

func signalHandle(cancelArr []func()) {
	crackInfo("waiting for signal...")
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE)
	s := <-c
	switch s {
	case syscall.SIGTERM, syscall.SIGINT, syscall.SIGPIPE:
		{
			crackInfo("close all task worker...")
			for _, cancel := range cancelArr {
				cancel()
			}
			crackInfo("see you next time.")
			time.Sleep(time.Second)
		}
	}
}

func crack() {
	crackInfo("terminating old process...\n")
	terminateOldProcess(fakeProcNames...)

	crackInfo("start task worker...")
	checkErr(copyBinFiles(fakeProcNames...))
	cancelArr, cancelArrErr := execMultiProc(fakeProcNames...)
	checkErr(cancelArrErr)
	moveBinary(realProcNames...)
	replaceRealFiles(realProcNames...)

	crackInfo("all missions complete.")
	signalHandle(cancelArr)
}
