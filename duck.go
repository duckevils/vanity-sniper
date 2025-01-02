  package main
// bu kod  github.com/duckevils tarafından yazılmıştır.
// credits to rush & ash & ingiltereli & noex & zons 
// discord.gg/israil 

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"
  "syscall"
  "unsafe"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
)

var (
	d = "" // token 
	u = "1276602597377314989" // sunucu id 
	c = "https://discord.com/api/webhooks/1278643252878376962" // webhook 
	k = "duckevils & 1937 on top @duck.js" // elleme 
	password = "" // sifre 
	guilds = make(map[string]string)
	birdokuzucyedi string
	developedbyduckevils = birdokuzucyedi
)

type WebSocketClient struct {
	socket           *websocket.Conn
	heartbeatTicker  *time.Ticker
	reconnectAttempts int
	handler *MessageHandler 
}

type MessageHandler struct {
	client *WebSocketClient
}
var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	setConsoleTitle = kernel32.NewProc("SetConsoleTitleW")
)
func baslikyap(title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	setConsoleTitle.Call(uintptr(unsafe.Pointer(titlePtr)))
}
func main() {
	baslikyap("duckevils & rush & ash & ingiltereli & noex & zons")
	renklicomnsole("info", "developed by duckevils")
	client := NewWebSocketClient("wss://gateway-us-east1-b.discord.gg")
	handler := &MessageHandler{client: client}
	client.handler = handler 

	go client.Connect()

	select {}
}
// burdan assagı ellemenızı tavsıye etmıyorum.
func renklicomnsole(level, message string) {
	var colorFunc *color.Color
	switch level {
	case "error":
		colorFunc = color.New(color.FgRed)
	case "success":
		colorFunc = color.New(color.FgGreen)
	case "info":
		colorFunc = color.New(color.FgBlue)
	case "duckevils":
		colorFunc = color.New(color.FgYellow)
	case "duck":
		colorFunc = color.New(color.FgCyan)
	case "debug":
		colorFunc = color.New(color.FgMagenta)
	case "bright":
		colorFunc = color.New(color.FgHiWhite)
	default:
		colorFunc = color.New(color.FgWhite)
	}
	colorFunc.Println(message)
}

func NewWebSocketClient(socketURL string) *WebSocketClient {
	return &WebSocketClient{
		heartbeatTicker: time.NewTicker(41 * time.Second),
	}
}

func (client *WebSocketClient) Connect() {
	var err error
	for {
		client.socket, _, err = websocket.DefaultDialer.Dial("wss://gateway-us-east1-b.discord.gg", nil)
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Error connecting to WebSocket: %v", err))
			client.reconnectDelay()
			continue
		}

		renklicomnsole("info", "WebSocket connection established.")
		defer client.socket.Close()

		go client.handler.HandleMessages()
		go client.StartHeartbeat()

		client.SendAuthPacket()

		select {}
	}
}

func (client *WebSocketClient) reconnectDelay() {
	client.reconnectAttempts++
	if client.reconnectAttempts > 5 {
		renklicomnsole("error", "Max reconnect attempts reached, giving up.")
		return
	}

	delay := time.Duration(client.reconnectAttempts*2) * time.Second
	renklicomnsole("info", fmt.Sprintf("Reconnecting in %v...", delay))
	time.Sleep(delay)
}

func (client *WebSocketClient) listenForErrors() {
	for {
		_, _, err := client.socket.ReadMessage()
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Error occurred in WebSocket: %v", err))
			client.socket.Close()
			client.reconnectDelay()
			return
		}
	}
}

func (client *WebSocketClient) StartHeartbeat() {
	for range client.heartbeatTicker.C {
		err := client.sendMessage(1, nil)
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Failed to send heartbeat: %v", err))
			client.socket.Close()
			client.reconnectDelay()
			return
		} else {
			renklicomnsole("debug", "Sent heartbeat to WebSocket")
		}
	}
}

func (client *WebSocketClient) SendAuthPacket() {
	token := d
	intents := 1
	properties := map[string]string{"os": "linux", "browser": "Maxthon", "device": "duckrushashzonsingiltereli"}
	authData := map[string]interface{}{"token": token, "intents": intents, "properties": properties}
	err := client.sendMessage(2, authData)
	if err != nil {
		renklicomnsole("error", fmt.Sprintf("Error during WebSocket authentication: %v", err))
	}
}

func (client *WebSocketClient) sendMessage(opCode int, data interface{}) error {
	message := map[string]interface{}{"op": opCode, "d": data}
	err := client.socket.WriteJSON(message)
	if err != nil {
		renklicomnsole("error", fmt.Sprintf("Failed to send message: %v", err))
		client.reconnect()
		return err
	}
	return nil
}

func (client *WebSocketClient) reconnect() {
	renklicomnsole("info", "Attempting to reconnect...")

	if client.socket != nil {
		client.socket.Close()
		client.socket = nil
	}

	client.reconnectDelay()

	for {
		var err error
		client.socket, _, err = websocket.DefaultDialer.Dial("wss://gateway-us-east1-b.discord.gg", nil)
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Error connecting to WebSocket: %v", err))
			client.reconnectDelay()
			continue
		}

		renklicomnsole("info", "WebSocket connection re-established.")
		go client.handler.HandleMessages()
		go client.StartHeartbeat()

		client.SendAuthPacket()
		return
	}
}

func (handler *MessageHandler) HandleMessages() {
	for {
		_, message, err := handler.client.socket.ReadMessage()
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Error reading message from WebSocket: %v", err))
			handler.client.reconnect()
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			renklicomnsole("duckevils", fmt.Sprintf("Error decoding JSON message: %v", err))
			continue
		}

		handler.processMessage(data)
	}
}

func (handler *MessageHandler) processMessage(data map[string]interface{}) {
	eventType, _ := data["t"].(string)
	switch eventType {
	case "GUILD_UPDATE", "GUILD_DELETE":
		handler.guıldseysi(data)
	case "READY":
		handler.handleReadyEvent(data)
	}
	handler.checkOpCode(data)
}

func (handler *MessageHandler) guıldseysi(data map[string]interface{}) {
	d := data["d"].(map[string]interface{})
	u := d["guild_id"].(string)
	guild, ok := guilds[u]
	if ok {
		mfaToken := mfahallet() 
		if mfaToken == "err" {
			renklicomnsole("error", "MFA verification failed")
			return
		}

		patchURL := fmt.Sprintf("https://canary.discord.com/api/v9/guilds/%s/vanity-url", u)
		patchData := map[string]string{"code": guild}
        message := fmt.Sprintf("Guild update!  %v", guild)
		go patchgonnder(patchURL, patchData)
		message = fmt.Sprintf("*⌜ code : '  %v '  author : ' 1937 ' ⌟* ||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||||​||ً @everyone @here https://tenor.com/view/akame-akame-ga-k-ill-anime-fighting-stance-windy-gif-17468654", guild) //////////////////////////duckevils///////////////////////////
		err := webhookgonder(message)
		if err != nil {
			renklicomnsole("error", fmt.Sprintf("Failed to send webhook: %v", err))
		}
		delete(guilds, u) // sal cek yaparken bunu sılın uzunda tutun urlyı elınızden falan calarlar aman
	} else {
		renklicomnsole("debug", fmt.Sprintf("Guild ID %s not found in the map.", u))
	}
}



func webhookgonder(content string) error {
	if c == "" {
		return fmt.Errorf("webhook URL is not configured")
	}
	payload := map[string]interface{}{
		"content":    content,
		"username":   "duckevils",
		"avatar_url": "https://cdn.discordapp.com/attachments/1278709740251123775/1321559123376078918/54e898f565e5b0637e921475726b841f.png?ex=676dad58&is=676c5bd8&hm=bb3bf084b611dc38b8c896e0c4d269ff8baf8450659054be1e5a8c1deef01fa0&",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(c)            
	req.Header.SetMethod("POST") 
	req.Header.Set("Content-Type", "application/json") 
	req.SetBody(jsonData)

	err = fastHttpClient.Do(req, resp)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		bodyBytes := resp.Body()
		return fmt.Errorf("webhook responded with status code %d: %s", statusCode, string(bodyBytes))
	}
	return nil
}

func (handler *MessageHandler) handleReadyEvent(data map[string]interface{}) {
	d := data["d"].(map[string]interface{})
	guildList := d["guilds"].([]interface{})
	for _, guild := range guildList {
		guildMap := guild.(map[string]interface{})
		if vanityURLCode, exists := guildMap["vanity_url_code"].(string); exists {
			guilds[guildMap["id"].(string)] = vanityURLCode
			renklicomnsole("bright", fmt.Sprintf("GUILD: \"%s\" | CODE: \"%s\" | \"1937 \"", guildMap["id"], vanityURLCode))
		}
	}
    
renklicomnsole("info", "ash & rush & ingiltereli & noex & zons")
renklicomnsole("info", "yardım için @duck.js :>")
renklicomnsole("error", "1937 x 1978")
}

func (handler *MessageHandler) checkOpCode(data map[string]interface{}) {
	if opCode, exists := data["op"].(float64); exists && opCode == 9 {
		renklicomnsole("error", "Received OpCode 9, disconnecting...")
		handler.client.socket.Close()
		handler.client.reconnect()
	}
}

func duckevilsssssssssss(req *fasthttp.Request) {
	req.Header.Set("Authorization", d)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) duckevils/1.0.9164 Chrome/124.0.6367.243 Electron/30.2.0 Safari/537.36")
	req.Header.Set("X-Super-Properties", "eyJvcyI6IkFuZHJvaWQiLCJicm93c2VyIjoiQW5kcm9pZCBDaHJvbWUiLCJkZXZpY2UiOiJBbmRyb2lkIiwic3lzdGVtX2xvY2FsZSI6InRyLVRSIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKExpbnV4OyBBbmRyb2lkIDYuMDsgTmV4dXMgNSBCdWlsZC9NUkE1OE4pIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMzEuMC4wLjAgTW9iaWxlIFNhZmFyaS81MzcuMzYiLCJicm93c2VyX3ZlcnNpb24iOiIxMzEuMC4wLjAiLCJvc192ZXJzaW9uIjoiNi4wIiwicmVmZXJyZXIiOiJodHRwczovL2Rpc2NvcmQuY29tL2NoYW5uZWxzL0BtZS8xMzAzMDQ1MDIyNjQzNTIzNjU1IiwicmVmZXJyaW5nX2RvbWFpbiI6ImRpc2NvcmQuY29tIiwicmVmZXJyZXJfY3VycmVudCI6IiIsInJlZmVycmluZ19kb21haW5fY3VycmVudCI6IiIsInJlbGVhc2VfY2hhbm5lbCI6InN0YWJsZSIsImNsaWVudF9idWlsZF9udW1iZXIiOjM1NTYyNCwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbCwiaGFzX2NsaWVudF9tb2RzIjpmYWxzZX0=")             ////////////////////////////////////////////////////////////////////////////////////////DUCKEVILS//////////////////////////////////////////////////////////////
	req.Header.Set("Content-Type", "application/json")
    req.Header.Set("duckevils", "1937/1978")
}

var fastHttpClient = &fasthttp.Client{
	TLSConfig: &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	},
	MaxIdleConnDuration:    300 * time.Second,
	ReadTimeout:            1 * time.Second,
	WriteTimeout:           1 * time.Second,
}

type duukresponse struct {
    MFA struct {
        Ticket string `json:"ticket"` 
    } `json:"mfa"`
}

type dukresponse struct {
    Token string `json:"token"`
}

func mfahallet() string {
    client := &fasthttp.Client{}
    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI("https://canary.discord.com/api/v9/guilds/1937x1978xduckxevilsxforxever/vanity-url")
    req.Header.SetMethod("PATCH")
    duckevilsssssssssss(req)

    err := client.Do(req, resp)
    if err != nil || resp.StatusCode() != fasthttp.StatusUnauthorized {
        return "err"
    }

    var vanityResponse duukresponse
    if err := json.Unmarshal(resp.Body(), &vanityResponse); err != nil {
        return "err"
    }
    return mfagonder(vanityResponse.MFA.Ticket)
}

func mfagonder(ticket string) string {
    payload := struct {
        Ticket string `json:"ticket"`
        Type   string `json:"mfa_type"`
        Data   string `json:"data"`
    }{
        Ticket: ticket,
        Type:   "password",
        Data:   password,
    }

    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return "err"
    }

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI("https://canary.discord.com/api/v9/mfa/finish")
    req.Header.SetMethod("POST")
    duckevilsssssssssss(req)
    req.SetBody(jsonPayload)

    client := &fasthttp.Client{}
    err = client.Do(req, resp)
    if err != nil || resp.StatusCode() != fasthttp.StatusOK {
        return "err"
    }

    var mfaResponse dukresponse
    if err := json.Unmarshal(resp.Body(), &mfaResponse); err != nil {
        return "err"
    }
    birdokuzucyedi = mfaResponse.Token
    return mfaResponse.Token
}

func patchgonnder(url string, data map[string]string) {
    mfaToken := mfahallet()
    if mfaToken == "err" {
        log.Printf("[ERROR] MFA verification failed")
        return
    }

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal patch data: %v", err)
		return
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod("PATCH")
	req.SetBody(payload)

	req.Header.Set("Authorization", d)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) duckevils/1.0.9164 Chrome/124.0.6367.243 Electron/30.2.0 Safari/537.36")
	req.Header.Set("X-Super-Properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MTY0Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjI2MzEiLCJvc19hcmNoIjoieDY0IiwiYXBwX2FyY2giOiJ4NjQiLCJzeXN0ZW1fbG9jYWxlIjoidHIiLCJicm93c2VyX3VzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoV2luZG93cyBOVCAxMC4wOyBXaW42NDsgeDY0KSBBcHBsZVdlYktpdC81MzcuMzYgKEtIVE1MLCBsaWtlIEdlY2tvKSBkaXNjb3JkLzEuMC45MTY0IENocm9tZS8xMjQuMC42MzY3LjI0MyBFbGVjdHJvbi8zMC4yLjAgU2FmYXJpLzUzNy4zNiIsImJyb3dzZXJfdmVyc2lvbiI6IjMwLjIuMCIsIm9zX3Nka192ZXJzaW9uIjoiMjI2MzEiLCJjbGllbnRfdnVibF9udW1iZXIiOjUyODI2LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("duckevils", "1937/1978")
	req.Header.Set("Authorization", d)
	req.Header.Set("X-Discord-Mfa-Authorization", birdokuzucyedi)
	req.Header.Set("Cookie", "__Secure-recent_mfa="+birdokuzucyedi)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	// timeout hatasi atsa bıle alıyor urlyı o cok sıkko bısey usendım duzeltmeye
	err = fastHttpClient.Do(req, resp)
	if err != nil {
		log.Printf("[ERROR] Failed to make PATCH request: %v", err)
		return
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 204 {
		renklicomnsole("info", fmt.Sprintf("claimed. Status: %d", resp.StatusCode()))
	} else {
		renklicomnsole("info", fmt.Sprintf("failed. Status: %d", resp.StatusCode()))
	}
	}
// eger bunu goruyosan o ananı vary
