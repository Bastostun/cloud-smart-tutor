package pve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type NoVNCConfig struct {
	VMID      int    `json:"vmid"`
	Node      string `json:"node"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	DPI       int    `json:"dpi"`
}

type NoVNCToken struct {
	Ticket    string `json:"ticket"`
	CSRFToken string `json:"CSRFPreventionToken"`
	User      string `json:"user"`
	VNCTicket string `json:"vncticket"`
	VMID      int    `json:"vmid"`
	Node      string `json:"node"`
	Port      int    `json:"port"`
}

type PVENetworkInfo struct {
	Iface       string `json:"iface"`
	Type        string `json:"type"`
	Active      bool   `json:"active"`
	Address     string `json:"address"`
	Address6    string `json:"address6"`
	Gateway     string `json:"gateway"`
	Netmask     string `json:"netmask"`
	Bridge      string `json:"bridge,omitempty"`
	BridgePorts string `json:"bridge_ports,omitempty"`
	VLANID      int    `json:"vlan_id,omitempty"`
}

type PVEStorageInfo struct {
	Storage   string  `json:"storage"`
	Type      string  `json:"type"`
	Active    bool    `json:"active"`
	Enabled   bool    `json:"enabled"`
	Shared    bool    `json:"shared"`
	Content   string  `json:"content"`
	Used      uint64  `json:"used"`
	Avail     uint64  `json:"avail"`
	Total     uint64  `json:"total"`
}

type ZFSPoolInfo struct {
	Pool        string  `json:"pool"`
	Health      string  `json:"health"`
	Allocated   uint64  `json:"allocated"`
	Free        uint64  `json:"free"`
	Size        uint64  `json:"size"`
	DedupRatio  float64 `json:"dedup_ratio"`
	Compression string  `json:"compression"`
}

type ZFSDatasetInfo struct {
	Name        string `json:"name"`
	Used        uint64 `json:"used"`
	Available   uint64 `json:"available"`
	Referenced  uint64 `json:"referenced"`
	Mountpoint  string `json:"mountpoint"`
	Type        string `json:"type"`
}

type KSMInfo struct {
	PagesShared     uint64 `json:"pages_shared"`
	PagesSharing    uint64 `json:"pages_sharing"`
	PagesUnshared   uint64 `json:"pages_unshared"`
	PagesVolatile   uint64 `json:"pages_volatile"`
	FullScans       uint64 `json:"full_scans"`
	MemorySaved     uint64 `json:"memory_saved"`
	MergeAcrossNodes int   `json:"merge_across_nodes"`
	Run             uint64 `json:"run"`
}

type VNCConsoleURL struct {
	URL       string `json:"url"`
	Ticket    string `json:"ticket"`
	VMID      int    `json:"vmid"`
	Node      string `json:"node"`
	Port      int    `json:"port"`
	Password  string `json:"password"`
}

func (c *PVEClient) GenerateVNCTicket(node string, vmID int) (*VNCConsoleURL, error) {
	payload := map[string]interface{}{
		"vmid": vmID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/vncproxy", node, vmID), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("vncproxy request failed: %w", err)
	}

	var response struct {
		Data struct {
			Port     int    `json:"port"`
			Ticket   string `json:"ticket"`
			User     string `json:"user"`
			Password string `json:"password"`
			UPID     string `json:"upid"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("decode vncproxy response: %w", err)
	}

	return &VNCConsoleURL{
		URL:      fmt.Sprintf("wss://%s:%d", strings.TrimPrefix(c.Endpoint, "https://"), response.Data.Port),
		Ticket:   response.Data.Ticket,
		VMID:     vmID,
		Node:     node,
		Port:     response.Data.Port,
		Password: response.Data.Password,
	}, nil
}

func (c *PVEClient) GetVMNetworks(node string, vmID int) ([]string, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var networks []string
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("net%d", i)
		if val, ok := response.Data[key]; ok {
			networks = append(networks, fmt.Sprintf("%v", val))
		}
	}

	return networks, nil
}

func (c *PVEClient) GetNodeNetworks(node string) ([]PVENetworkInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/network", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []PVENetworkInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) GetNodeStorage(node string) ([]PVEStorageInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/storage", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []PVEStorageInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) GetZFSStatus(node string) ([]ZFSPoolInfo, []ZFSDatasetInfo, error) {
	pools, err := c.getZFSPools(node)
	if err != nil {
		return nil, nil, err
	}

	datasets, err := c.getZFSDatasets(node)
	if err != nil {
		return nil, nil, err
	}

	return pools, datasets, nil
}

func (c *PVEClient) getZFSPools(node string) ([]ZFSPoolInfo, error) {
	respBody, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/api2/json", node), nil)
	if err != nil {
		return nil, err
	}

	var result []ZFSPoolInfo
	_ = json.Unmarshal(respBody, &result)

	if len(result) == 0 {
		return c.getStorageAsFallback(node)
	}

	return result, nil
}

func (c *PVEClient) getStorageAsFallback(node string) ([]ZFSPoolInfo, error) {
	storages, err := c.GetNodeStorage(node)
	if err != nil {
		return nil, err
	}

	var pools []ZFSPoolInfo
	for _, s := range storages {
		if s.Type == "zfspool" || s.Type == "lvm" {
			pools = append(pools, ZFSPoolInfo{
				Pool:      s.Storage,
				Health:    "online",
				Allocated: s.Used,
				Free:      s.Avail,
				Size:      s.Total,
			})
		}
	}

	return pools, nil
}

func (c *PVEClient) getZFSDatasets(node string) ([]ZFSDatasetInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/disks/zfs", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []ZFSDatasetInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) GetKSMInfo(node string) (*KSMInfo, error) {
	script := `cat /sys/kernel/mm/ksm/pages_shared && echo && cat /sys/kernel/mm/ksm/pages_sharing && echo && cat /sys/kernel/mm/ksm/pages_unshared && echo && cat /sys/kernel/mm/ksm/pages_volatile && echo && cat /sys/kernel/mm/ksm/full_scans && echo && cat /sys/kernel/mm/ksm/pages_to_scan && echo && cat /sys/kernel/mm/ksm/merge_across_nodes && echo && cat /sys/kernel/mm/ksm/run`

	payload := map[string]interface{}{
		"script": script,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu", node), bytes.NewReader(body))
	if err != nil {
		return c.getKSMFromNodeStatus(node)
	}

	var result KSMInfo
	_ = json.Unmarshal(respBody, &result)
	return &result, nil
}

func (c *PVEClient) getKSMFromNodeStatus(node string) (*KSMInfo, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/status", node), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			KSM struct {
				Shared  uint64 `json:"shared"`
				Sharing uint64 `json:"sharing"`
				Saved   uint64 `json:"saved"`
			} `json:"ksm"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &KSMInfo{
		PagesShared:  response.Data.KSM.Shared,
		PagesSharing: response.Data.KSM.Sharing,
		MemorySaved:  response.Data.KSM.Saved * 4096,
	}, nil
}

func (c *PVEClient) ResizeVMDisk(node string, vmID int, size string) error {
	payload := map[string]interface{}{
		"disk": "scsi0",
		"size": size,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", fmt.Sprintf("/nodes/%s/qemu/%d/resize", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) SetVMConfig(node string, vmID int, config map[string]interface{}) error {
	body, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) GetVMConfig(node string, vmID int) (map[string]interface{}, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) MigrateVM(node string, vmID int, targetNode string, online bool) error {
	payload := map[string]interface{}{
		"target": targetNode,
		"online": online,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/migrate", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) DeleteVM(node string, vmID int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/nodes/%s/qemu/%d", node, vmID), nil)
	return err
}

func (c *PVEClient) RebootVM(node string, vmID int) error {
	payload := map[string]interface{}{
		"forceStop": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/reboot", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) ResetVM(node string, vmID int) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/reset", node, vmID), nil)
	return err
}

func (c *PVEClient) ShutdownVM(node string, vmID int) error {
	payload := map[string]interface{}{
		"forceStop": false,
		"timeout":   30,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/shutdown", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) ForceStopVM(node string, vmID int) error {
	_, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/status/stop", node, vmID), nil)
	return err
}

func (c *PVEClient) GetVMIPAddresses(node string, vmID int) ([]string, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu/%d/agent/network-get-interfaces", node, vmID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Result []struct {
				Name        string `json:"name"`
				IPAddresses []struct {
					Address string `json:"ip-address"`
					Prefix  int    `json:"prefix"`
					Type    string `json:"ip-address-type"`
				} `json:"ip-addresses"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var ips []string
	for _, iface := range response.Data.Result {
		for _, ip := range iface.IPAddresses {
			if ip.Type == "ipv4" {
				ips = append(ips, fmt.Sprintf("%s/%d", ip.Address, ip.Prefix))
			}
		}
	}

	return ips, nil
}

func (c *PVEClient) ExecuteInVM(node string, vmID int, command []string) (string, error) {
	payload := map[string]interface{}{
		"command": command,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/agent/exec", node, vmID), bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var response struct {
		Data struct {
			Exited int    `json:"exited"`
			OutData string `json:"out-data"`
			ErrData string `json:"err-data"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	if response.Data.Exited != 0 {
		return "", fmt.Errorf("command exited with code %d: %s", response.Data.Exited, response.Data.ErrData)
	}

	return response.Data.OutData, nil
}

func (c *PVEClient) CreateFirewallRule(node string, vmID int, direction string, action string, protocol string, dport string) error {
	payload := map[string]interface{}{
		"direction": direction,
		"action":    action,
		"protocol":  protocol,
		"dport":     dport,
		"enable":    1,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/firewall/rules", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) DeleteFirewallRule(node string, vmID int, ruleIndex int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/nodes/%s/qemu/%d/firewall/rules/%d", node, vmID, ruleIndex), nil)
	return err
}

func (c *PVEClient) GetFirewallRules(node string, vmID int) ([]map[string]interface{}, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu/%d/firewall/rules", node, vmID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *PVEClient) ReadFileInVM(node string, vmID int, filePath string) (string, error) {
	payload := map[string]interface{}{
		"command": []string{"cat", filePath},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	respBody, err := c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/agent/file-read", node, vmID), bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var response struct {
		Data struct {
			Result string `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	return response.Data.Result, nil
}

func (c *PVEClient) WriteFileInVM(node string, vmID int, filePath string, content string) error {
	payload := map[string]interface{}{
		"file": filePath,
		"content": content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", fmt.Sprintf("/nodes/%s/qemu/%d/agent/file-write", node, vmID), bytes.NewReader(body))
	return err
}

func (c *PVEClient) GetNoVNCURL(node string, vmID int, width int, height int) (string, error) {
	vnc, err := c.GenerateVNCTicket(node, vmID)
	if err != nil {
		return "", err
	}

	noVNCURL := fmt.Sprintf("%s/?console=kvm&novnc=1&vmid=%d&node=%s&resize=scale&width=%d&height=%d&websocket=&path=/api2/json/nodes/%s/qemu/%d/vncwebsocket&port=%d&password=%s",
		c.Endpoint,
		vmID,
		node,
		width,
		height,
		node,
		vmID,
		vnc.Port,
		vnc.Password,
	)

	return noVNCURL, nil
}

func (c *PVEClient) GetNoVNCTicket(vmID int, node string) (*NoVNCToken, error) {
	vnc, err := c.GenerateVNCTicket(node, vmID)
	if err != nil {
		return nil, err
	}

	return &NoVNCToken{
		Ticket:    vnc.Ticket,
		CSRFToken: c.CSRFToken,
		VNCTicket: vnc.Password,
		VMID:      vmID,
		Node:      node,
		Port:      vnc.Port,
	}, nil
}

func (c *PVEClient) GetLogs(node string, vmID int, start int, limit int) ([]string, error) {
	body, err := c.doRequest("GET", fmt.Sprintf("/nodes/%s/qemu/%d/status/current", node, vmID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Log string `json:"log"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Data.Log == "" {
		return []string{}, nil
	}

	return strings.Split(response.Data.Log, "\n"), nil
}
