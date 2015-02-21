package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	SyslogSocket = "/dev/log"
)

func sysReboot() error {
	log.Printf("rebooting")
	syscall.Sync()
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func sysHalt() error {
	syscall.Sync()
	log.Printf("halting")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
}

func sysPoweroff() error {
	syscall.Sync()
	log.Printf("power off")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func sysReinit() error {
	log.Printf("restarting init(pid=1), sending signals for all children")
	if err := syscall.Kill(-1, syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send SIGTERM for all processes: %s", err)
	}
	time.Sleep(1 * time.Second)
	syscall.Kill(-1, syscall.SIGKILL)
	if err := collectChildren(false); err != nil {
		return err
	}
	return syscall.Exec(os.Args[0], os.Args, os.Environ())
}

var signalHandlers = map[os.Signal]func() error{
	syscall.SIGHUP:  sysReboot,
	syscall.SIGINT:  sysReboot,
	syscall.SIGPWR:  sysHalt,
	syscall.SIGQUIT: sysReinit,
	syscall.SIGTERM: sysReboot,
	syscall.SIGUSR1: sysHalt,
	syscall.SIGUSR2: sysPoweroff,
}

func setupSignal() chan os.Signal {
	ch := make(chan os.Signal)
	for sig := range signalHandlers {
		signal.Notify(ch, sig)
	}
	return ch
}

func handleSignal(sig os.Signal) error {
	if handler := signalHandlers[sig]; handler != nil {
		return handler()
	}
	return nil
}

func handleSyslogConn(conn net.Conn) {
	defer conn.Close()
	for {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			if err == syscall.EAGAIN {
				continue
			}
			log.Printf("error while reading data from syslog socket: %s", err)
		}
		break
	}
}

func handleSyslog(ln net.Listener) {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept syslog connection request: %s", err)
			continue
		}
		go handleSyslogConn(conn)
	}
}

func serviceSyslog() error {
	ln, err := net.Listen("unix", SyslogSocket)
	if err != nil {
		return fmt.Errorf("failed to create syslog socket %s: %s", SyslogSocket, err)
	}
	go handleSyslog(ln)
	return nil
}

func serviceInitialService(init []string) error {
	cmd := exec.Command(init[0], init[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start initial service %s: %s", init[0], err)
	}
	return nil
}

func collectChildren(block bool) error {
	for {
		var waitStatus syscall.WaitStatus
		wpid, err := syscall.Wait4(-1, &waitStatus, 0, nil)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			if err == syscall.ECHILD {
				if block {
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					break
				}
			}
			return fmt.Errorf("error while waiting child exit: %s", err)
		}
		log.Printf("child %d exit with status %d", wpid, waitStatus.ExitStatus())
	}
	return nil
}

func startInit(init []string) error {
	chSignal := setupSignal()

	if opts.Syslog {
		if err := serviceSyslog(); err != nil {
			return err
		}
	}

	if err := serviceInitialService(init); err != nil {
		return err
	}

	chError := make(chan error, 1)
	go func() {
		chError <- collectChildren(true)
	}()

	for {
		select {
		case err := <-chError:
			return err
		case sig := <-chSignal:
			log.Printf("received signal %d", sig)
			if err := handleSignal(sig); err != nil {
				return err
			}
		}
	}
}

var opts struct {
	Syslog bool `long:"syslog" description:"Enable syslog(/dev/log)"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = os.Args[0]
	parser.Usage = "[OPTIONS] INIT"

	args, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if len(args) < 1 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if err := startInit(args); err != nil {
		log.Fatalln(err)
	}
}
