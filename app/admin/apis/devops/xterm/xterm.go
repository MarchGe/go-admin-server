package xterm

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"time"
)

var _xterm = &Xterm{}

type Xterm struct {
}

const readBufferSize = 1024 * 16    // websocket read buffer size
const writeBufferSize = 1024 * 1024 // websocket write buffer size
const readLimit = 1024 * 16         // websocket connection read limit, eg. conn.ReadMessage(), exceed readLimit limitation will return err
var wsUpgrader = &websocket.Upgrader{
	HandshakeTimeout: 30 * time.Second,
	ReadBufferSize:   readBufferSize,
	WriteBufferSize:  writeBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetXterm() *Xterm {
	return _xterm
}

// Connect godoc
//
//	@Summary	websocket连接
//	@Tags		Web Shell
//	@Param		rows	query		int		true	"pty的行数"
//	@Param		cols	query		int		true	"pty的列数"
//	@Param		token	query		string	true	"Web Shell连接使用的token"
//	@Success	200		{object}	nil
//	@Router		/terminal/ws [get]
func (a *Xterm) Connect(c *gin.Context) {
	winSize := &WinSize{}
	if e := c.ShouldBindQuery(winSize); e != nil {
		E.PanicErr(e)
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("websocket upgrade error", slog.Any("err", err))
		return
	}
	conn.SetReadLimit(readLimit)
	defer func() { _ = conn.Close() }()

	if err = a.handleConn(conn, winSize); err != nil {
		var closeErr *websocket.CloseError
		needNotHandleCodes := []int{websocket.CloseNormalClosure, websocket.CloseGoingAway}
		if !errors.As(err, &closeErr) {
			_ = conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
			slog.Error("handleConnection error", slog.Any("err", err))
		} else if !slices.Contains(needNotHandleCodes, closeErr.Code) {
			slog.Error("handleConnection error", slog.Any("err", closeErr))
		}
	}
}

// ConnectWithRemoteSSH godoc
//
//	@Summary	websocket连接，并通过ssh连接其他主机
//	@Tags		Web Shell With SSH
//	@Param		rows	query		int		true	"pty的行数"
//	@Param		cols	query		int		true	"pty的列数"
//	@Param		token	query		string	true	"Web Shell连接使用的token"
//	@Param		id		path		int64	true	"主机ID"
//	@Success	200		{object}	nil
//	@Router		/terminal/ws/ssh/:id [get]
func (a *Xterm) ConnectWithRemoteSSH(c *gin.Context) {
	winSize := &WinSize{}
	if e := c.ShouldBindQuery(winSize); e != nil {
		E.PanicErr(e)
	}
	hostId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	host, _ := dvservice.GetHostService().FindOneById(hostId)
	if host == nil {
		R.Fail(c, "连接的主机信息不存在", http.StatusBadRequest)
		return
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("websocket upgrade error", slog.Any("err", err))
		return
	}
	conn.SetReadLimit(readLimit)
	defer func() { _ = conn.Close() }()

	if err = a.handleConnWithSSH(conn, winSize, host); err != nil {
		var closeErr *websocket.CloseError
		needNotHandleCodes := []int{websocket.CloseNormalClosure, websocket.CloseGoingAway}
		if !errors.As(err, &closeErr) {
			_ = conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
			if !errors.Is(err, io.EOF) {
				slog.Error("handleConnection error", slog.Any("err", err))
			}
		} else if !slices.Contains(needNotHandleCodes, closeErr.Code) {
			slog.Error("handleConnection error", slog.Any("err", closeErr))
		}
	}
}

func resizeTTY(ptmx *os.File, wz *WinSize) {
	err := pty.Setsize(ptmx, &pty.Winsize{
		Rows: uint16(wz.Rows),
		Cols: uint16(wz.Cols),
	})
	if err != nil {
		slog.Error("applies pty size to tty error", slog.Any("err", err))
	}
}

func (a *Xterm) handleConn(conn *websocket.Conn, winSize *WinSize) error {
	bash, err := getBash()
	if err != nil {
		return fmt.Errorf("get bash error, %w", err)
	}
	shell := exec.Command(bash)
	ptmx, err := pty.Start(shell)
	if err != nil {
		return fmt.Errorf("pty start error, %w", err)
	}

	go resizeTTY(ptmx, winSize)
	defer func() {
		_ = ptmx.Close()
		_ = shell.Process.Kill()
	}()

	go func() {
		n, buf := 0, make([]byte, 32*1024)
		for {
			if n, err = ptmx.Read(buf); err != nil {
				return
			}
			if err = conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				if err != websocket.ErrCloseSent {
					slog.Error("write data to websocket error", slog.Any("err", err))
				}
				return
			}
		}
	}()

	for {
		_, cmdBytes, e := conn.ReadMessage()
		if e != nil {
			return fmt.Errorf("read data from websocket error, %w", e)
		}
		if _, e = ptmx.Write(cmdBytes); e != nil {
			return fmt.Errorf("write data to ptmx error, %w", e)
		}
	}
}

func (a *Xterm) handleConnWithSSH(conn *websocket.Conn, winSize *WinSize, host *dvmodel.Host) error {
	addr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	password, err := utils.DecryptString(config.GetConfig().EncryptKey, host.Password, "")
	if err != nil {
		return err
	}
	clientConfig := ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         constant.SshEstablishTimeoutInSeconds * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, &clientConfig)
	if err != nil {
		return fmt.Errorf("ssh connect failed, %w", err)
	}
	defer func() { _ = client.Close() }()
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("ssh client.NewSession() error, %w", err)
	}
	defer func() { _ = session.Close() }()

	if err = session.RequestPty("xterm", winSize.Rows, winSize.Cols, ssh.TerminalModes{}); err != nil {
		return fmt.Errorf("ssh request pty error, %w", err)
	}
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get stdout pipe of ssh error, %w", err)
	}
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("get stdin pipe of ssh error, %w", err)
	}
	// login
	if err = session.Shell(); err != nil {
		return fmt.Errorf("ssh login error, %w", err)
	}
	go func() {
		n, buf := 0, make([]byte, 32*1024)
		for {
			if n, err = stdoutPipe.Read(buf); err != nil {
				if err != io.EOF {
					slog.Error("read data from ssh stdout pipe error", slog.Any("err", err))
				}
				return
			}
			if err = conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				if err != websocket.ErrCloseSent {
					slog.Error("write data to websocket error", slog.Any("err", err))
				}
				return
			}
		}
	}()

	for {
		_, cmdBytes, e := conn.ReadMessage()
		if e != nil {
			return fmt.Errorf("read data from websocket error, %w", e)
		}
		if _, e = stdinPipe.Write(cmdBytes); e != nil {
			return fmt.Errorf("write data to ssh stdin error, %w", e)
		}
	}
}

func getBash() (string, error) {
	unsupportedPlatforms := []string{"windows"} // pty not support windows platform
	if slices.Contains(unsupportedPlatforms, runtime.GOOS) {
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	if _, err := exec.LookPath("bash"); err != nil {
		return "", fmt.Errorf("bash command not found, %w", err)
	}
	return "bash", nil
}
