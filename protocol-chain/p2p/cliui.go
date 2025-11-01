package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

type CLIUI struct {
	MiningChannel    *Channel
	FullNodesChannel *Channel
	app              *tview.Application
	peerList         *tview.TextView

	hostWindow *tview.TextView
	inputCh    chan string
	doneCh     chan struct{}
}

type Log struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

var (
	_, file, _, _ = runtime.Caller(0)

	Root = filepath.Join(filepath.Dir(file), "../")
)

func NewCLIUI(miningChannel, fullNodesChannel *Channel) *CLIUI {
	app := tview.NewApplication()

	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf(" HOST (%s) ", strings.ToUpper(ShortID(fullNodesChannel.self))))

	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(strings.ToUpper(ShortID(fullNodesChannel.self) + " > ")).
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			return
		}

		line := input.GetText()
		if len(line) == 0 {
			return
		}

		inputCh <- line
		input.SetText("")
	})

	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle(" Peers ")

	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false)

	app.SetRoot(flex, true)

	return &CLIUI{
		MiningChannel:    miningChannel,
		FullNodesChannel: fullNodesChannel,
		app:              app,
		peerList:         peersList,
		hostWindow:       msgBox,
		inputCh:          inputCh,
		doneCh:           make(chan struct{}),
	}
}

func ShortID(p peer.ID) string {
	peerId := p.String()
	return peerId[len(peerId)-8:]
}

func (ui *CLIUI) Run(net *Network) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("UI Crashed: %v", r)
		}
	}()

	go ui.handleEvents(net)
	defer ui.end()

	return ui.app.Run()
}

func (ui *CLIUI) end() {
	ui.doneCh <- struct{}{}
}

func (ui *CLIUI) displaySelfMessage(msg string) {
	ui.app.QueueUpdateDraw(func() {
		prompt := withColor("yellow", fmt.Sprintf("<%s>:", strings.ToUpper(ShortID(ui.FullNodesChannel.self))))
		fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, msg)
		ui.hostWindow.ScrollToEnd()
	})

}

func (ui *CLIUI) refreshPeers() {
	minerPeers := ui.MiningChannel.ListPeers()
	fullnodes := ui.FullNodesChannel.ListPeers()
	idStrs := make([]string, 0, len(fullnodes))

	for _, pId := range fullnodes {
		peerId := strings.ToUpper(ShortID(pId))
		if len(minerPeers) != 0 && slices.Contains(minerPeers, pId) {
			idStrs = append(idStrs, "MINER: "+peerId)
		} else {
			idStrs = append(idStrs, "FNODE: "+peerId)
		}
	}

	ui.app.QueueUpdateDraw(func() {
		ui.peerList.SetText(strings.Join(idStrs, "\n"))
	})
}

func (ui *CLIUI) listenChannels(net *Network) {
	for {
		select {
		case content := <-ui.MiningChannel.Content:
			net.worker.Push(content)

		case content := <-ui.FullNodesChannel.Content:
			net.worker.Push(content)

		case <-ui.doneCh:
			return
		}
	}
}

func (ui *CLIUI) handleEvents(net *Network) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go ui.readFromLogs(net)
	go ui.listenChannels(net)

	for {
		select {
		case input := <-ui.inputCh:
			err := ui.FullNodesChannel.Publish(input, nil, "")
			if err != nil {
				log.Errorf("Publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		case <-ui.doneCh:
			log.Info("Stopping CLI UI")
			return
		case <-ticker.C:
			ui.refreshPeers()
		}
	}
}

func (ui *CLIUI) HandleStream(net *Network, content *ChannelContent) {
	if content.Payload != nil {
		command := BytesToCmd(content.Payload[:commandLength])

		switch command {
		// Sync Block
		case PREFIX_HEADER:
			net.HandleGetHeader(content)
		case PREFIX_BLOCK:
			net.HandleGetBlockData(content)
		case PREFIX_TX_MINING:
			net.HandleTxMining(content)
		case PREFIX_GET_DATA:
			net.HandleGetData(content)
		case PREFIX_HEADER_LOCATOR:
			net.HandleGetHeaderLocator(content)
		case PREFIX_HEADER_SYNC:
			net.HandleGetHeaderSync(content)
		case PREFIX_GET_DATA_SYNC:
			net.HandleGetDataSync(content)
		case PREFIX_BLOCK_SYNC:
			net.HandleGetBlockDataSync(content)

		// Sync Transaction
		case PREFIX_TX:
			net.HandleTx(content)
		case PREFIX_TX_FROM_POOL:
			net.HandleGetTxFromPool(content)
		case PREFIX_INVENTORY:
			net.HandleGetTxPoolInv(content)
		case PREFIX_TXS_Data:
			net.HandleGetTransactions(content)
		case PREFIX_DATA_TX:
			net.HandleGetDataTx(content)
		default:
			log.Warning("Unknown command received: ", command)
		}
	}
}

func (ui *CLIUI) readFromLogs(net *Network) {
	instanceId := net.Blockchain.InstanceId

	filename := "/logs/console.log"
	if instanceId != "" {
		filename = fmt.Sprintf("/logs/console_%s.log", instanceId)
	}

	logFile := path.Join(Root, filename)
	log.SetOutput(io.Discard)

	openLog := func() (*os.File, *bufio.Reader, error) {
		f, err := os.Open(logFile)
		if err != nil {
			return nil, nil, err
		}
		return f, bufio.NewReader(f), nil
	}

	f, r, err := openLog()
	if err != nil {
		log.Errorf("Cannot open log file initially: %v", err)
		return
	}
	defer f.Close()

	logLevels := map[string]string{
		"trace":   "cyan",
		"debug":   "blue",
		"info":    "green",
		"warn":    "yellow",
		"warning": "yellow",
		"error":   "red",
		"fatal":   "red",
		"panic":   "red",
	}

	info, _ := f.Stat()
	oldSize := info.Size()

	for {
		for {
			line, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}

			if err != nil {
				if os.IsNotExist(err) {
					log.Warnf("Log file %s missing (possibly compressed). Waiting for recreation...", logFile)
					break
				}
				log.Errorf("Read line error: %v", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			var data Log
			if err := json.Unmarshal(line, &data); err != nil {
				log.Errorf("Unmarshal log err: %v", err)
				continue
			}

			ui.app.QueueUpdateDraw(func() {
				prompt := fmt.Sprintf("%s:", withColor(logLevels[data.Level], strings.ToUpper(data.Level)))
				fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, data.Msg)
				ui.hostWindow.ScrollToEnd()
			})

		}

		newInfo, err := os.Stat(logFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warnf("Log file %s not found, retrying...", logFile)

				for {
					time.Sleep(1 * time.Second)
					newFile, newReader, reopenErr := openLog()
					if reopenErr == nil {
						log.Infof("Reopened new log file: %s", logFile)
						f.Close()
						f = newFile
						r = newReader
						info, _ = f.Stat()
						oldSize = info.Size()
						break
					}
				}
				continue
			}
			log.Errorf("Stat error: %v", err)
		} else {
			newSize := newInfo.Size()

			if newSize < oldSize {
				log.Warnf("Detected log rotation, reopening %s...", logFile)
				f.Close()
				newFile, newReader, _ := openLog()
				f = newFile
				r = newReader
				oldSize = newSize
				continue
			}

			oldSize = newSize
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func withColor(color, msg string) string {
	return fmt.Sprintf("[%s] %s[-:-:-]", strings.TrimSpace(color), msg)
}
