package p2p

import (
	"bufio"
	"core-blockchain/common/utils"
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
	GeneralChannel   *Channel
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

func NewCLIUI(generalChannel, miningChannel, fullNodesChannel *Channel) *CLIUI {
	app := tview.NewApplication()

	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf(" HOST (%s) ", strings.ToUpper(ShortID(generalChannel.self))))

	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(strings.ToUpper(ShortID(generalChannel.self) + " > ")).
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

		if line == "/quit" {
			app.Stop()
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
		GeneralChannel:   generalChannel,
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

func (ui *CLIUI) Run(net *Network, callback func(*Network)) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("UI Crashed: %v", r)
		}
	}()

	go ui.handleEvents(net, callback)
	defer ui.end()

	return ui.app.Run()
}

func (ui *CLIUI) end() {
	ui.doneCh <- struct{}{}
}

func (ui *CLIUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", strings.ToUpper(ShortID(ui.GeneralChannel.self))))
	fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, msg)
}

func (ui *CLIUI) refreshPeers() {
	peers := ui.GeneralChannel.ListPeers()
	minerPeers := ui.MiningChannel.ListPeers()
	fullnodes := ui.FullNodesChannel.ListPeers()
	idStrs := make([]string, 0, len(peers))

	for _, pId := range peers {
		peerId := strings.ToUpper(ShortID(pId))
		if len(minerPeers) != 0 && slices.Contains(minerPeers, pId) {
			idStrs = append(idStrs, "MINER: "+peerId)
		} else if len(fullnodes) != 0 && slices.Contains(fullnodes, pId) {
			idStrs = append(idStrs, "FNODE: "+peerId)
		} else {
			idStrs = append(idStrs, "NODE: "+peerId)
		}
	}

	ui.peerList.SetText(strings.Join(idStrs, "\n"))
	ui.app.Draw()
}

func (ui *CLIUI) handleEvents(net *Network, callback func(*Network)) {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	go ui.readFromLogs(net, callback)
	log.Info("HOST ADDR: ", net.Host.Addrs())

	for {
		select {
		case input := <-ui.inputCh:
			err := ui.GeneralChannel.Publish(input, nil, "")
			if err != nil {
				log.Errorf("Publish error: %s", err)
			}
			ui.displaySelfMessage(input)

		case <-peerRefreshTicker.C:
			ui.refreshPeers()

		case content := <-ui.GeneralChannel.Content:
			ui.HandleStream(net, content)

		case content := <-ui.MiningChannel.Content:
			ui.HandleStream(net, content)

		case content := <-ui.FullNodesChannel.Content:
			ui.HandleStream(net, content)

		case <-ui.doneCh:
			log.Info("Stopping CLI UI")

			return
		}
	}
}

func (net *Network) handleLimitedOperation(limiter chan struct{}, content *ChannelContent, handler func(*ChannelContent)) {
	select {
	case limiter <- struct{}{}:
		go func() {
			defer func() {
				<-limiter
			}()
			handler(content)
		}()
	default:
		log.Warnf("WARNING: Processing is busy, skipping operation for content from %s", content.SendFrom)
	}
}

func (ui *CLIUI) HandleStream(net *Network, content *ChannelContent) {
	if content.Payload != nil {
		command := BytesToCmd(content.Payload[:commandLength])
		log.Infof("Received %s command\n", command)

		switch command {
		case PREFIX_TX_FROM_POOL:
			net.HandleGetTxFromPool(content)
		case PREFIX_HEADER:
			net.HandleGetHeader(content)
		case PREFIX_BLOCKS_HEADER:
			net.handleLimitedOperation(net.blockProcessingLimiter, content, net.HandleGetBlocksHeader)
		case PREFIX_INVENTORY:
			net.HandleGetInventory(content)
		case PREFIX_BLOCK:
			net.handleLimitedOperation(net.blockProcessingLimiter, content, net.HandleGetBlocksData)
		case PREFIX_TX:
			if net.syncCompleted {
				net.handleLimitedOperation(net.txProcessingLimiter, content, net.HandleTx)
			}
		case PREFIX_GET_DATA:
			net.handleLimitedOperation(net.blockProcessingLimiter, content, net.HandleGetData)
		default:
			log.Warning("Unknown command received: ", command)
		}
	}
}

func (ui *CLIUI) readFromLogs(net *Network, callback func(*Network)) {
	instanceId := net.Blockchain.InstanceId

	filename := "/logs/console.log"
	if instanceId != "" {
		filename = fmt.Sprintf("/logs/console_%s.log", instanceId)
	}

	logFile := path.Join(Root, filename)
	err := os.WriteFile(logFile, []byte(""), 0644)
	utils.ErrorHandle(err)

	log.SetOutput(io.Discard)

	f, err := os.Open(logFile)
	utils.ErrorHandle(err)
	defer f.Close()

	r := bufio.NewReader(f)
	info, err := f.Stat()
	utils.ErrorHandle(err)

	callback(net)

	logLevels := map[string]string{
		"info":    "green",
		"warn":    "yellow",
		"warning": "yellow",
		"error":   "red",
		"fatal":   "red",
		"panic":   "red",
	}

	oldSize := info.Size()

	for {
		for line, _, err := r.ReadLine(); err != io.EOF; line, _, err = r.ReadLine() {
			var data Log

			if err := json.Unmarshal(line, &data); err != nil {
				panic(err)
			}

			prompt := fmt.Sprintf("[%s]:", withColor(logLevels[data.Level], strings.ToUpper(data.Level)))
			fmt.Fprintf(ui.hostWindow, "%s %s\n", prompt, data.Msg)
			ui.hostWindow.ScrollToEnd()
		}

		pos, err := f.Seek(0, io.SeekCurrent)
		utils.ErrorHandle(err)

		for {
			time.Sleep(time.Second)

			newInfo, err := f.Stat()
			utils.ErrorHandle(err)
			newSize := newInfo.Size()

			if newSize != oldSize {
				if newSize < oldSize {
					f.Seek(0, io.SeekStart)
				} else {
					f.Seek(pos, io.SeekStart)
				}
				r = bufio.NewReader(f)
				oldSize = newSize
				break
			}
		}

	}
}

func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-:-:-]", color, msg)
}
