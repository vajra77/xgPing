package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"xgPing/probe"
)

type IXFResponse struct {
	Version   string    `json:"version"`
	Generator string    `json:"generator"`
	Timestamp time.Time `json:"timestamp"`
	IxpList   []struct {
		Shortname             string   `json:"shortname"`
		Name                  string   `json:"name"`
		Country               string   `json:"country"`
		URL                   string   `json:"url"`
		PeeringdbID           int      `json:"peeringdb_id"`
		IxfID                 int      `json:"ixf_id"`
		IxpID                 int      `json:"ixp_id"`
		SupportEmail          string   `json:"support_email"`
		SupportPhone          string   `json:"support_phone"`
		SupportContactHours   string   `json:"support_contact_hours"`
		EmergencyEmail        string   `json:"emergency_email"`
		EmergencyPhone        string   `json:"emergency_phone"`
		EmergencyContactHours string   `json:"emergency_contact_hours"`
		BillingContactHours   string   `json:"billing_contact_hours"`
		BillingEmail          string   `json:"billing_email"`
		BillingPhone          string   `json:"billing_phone"`
		PeeringPolicyList     []string `json:"peering_policy_list"`
		Vlan                  []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Ipv4 struct {
				Prefix     string `json:"prefix"`
				MaskLength int    `json:"mask_length"`
			} `json:"ipv4"`
			Ipv6 struct {
				Prefix     string `json:"prefix"`
				MaskLength int    `json:"mask_length"`
			} `json:"ipv6"`
		} `json:"vlan"`
		Switch []struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Colo          string `json:"colo"`
			City          string `json:"city"`
			Country       string `json:"country"`
			PdbFacilityID int    `json:"pdb_facility_id"`
			Manufacturer  string `json:"manufacturer"`
			Model         string `json:"model"`
			Software      string `json:"software"`
		} `json:"switch"`
	} `json:"ixp_list"`
	MemberList []struct {
		Asnum          int       `json:"asnum"`
		MemberSince    time.Time `json:"member_since"`
		URL            string    `json:"url"`
		Name           string    `json:"name"`
		PeeringPolicy  string    `json:"peering_policy"`
		MemberType     string    `json:"member_type"`
		ConnectionList []struct {
			IxpID  int    `json:"ixp_id"`
			State  string `json:"state"`
			IfList []struct {
				SwitchID int `json:"switch_id"`
				IfSpeed  int `json:"if_speed"`
			} `json:"if_list"`
			VlanList []struct {
				VlanID int `json:"vlan_id"`
				Ipv4   struct {
					Address      string   `json:"address"`
					AsMacro      string   `json:"as_macro"`
					Routeserver  bool     `json:"routeserver"`
					MacAddresses []string `json:"mac_addresses"`
					MaxPrefix    int      `json:"max_prefix"`
				} `json:"ipv4"`
				Ipv6 struct {
					Address      string   `json:"address"`
					AsMacro      string   `json:"as_macro"`
					Routeserver  bool     `json:"routeserver"`
					MacAddresses []string `json:"mac_addresses"`
					MaxPrefix    int      `json:"max_prefix"`
				} `json:"ipv6"`
			} `json:"vlan_list"`
		} `json:"connection_list"`
	} `json:"member_list"`
}

func ImportJSONPeers(url string, ixpId int, vlanId int) ([]*probe.Peer, error) {
	result := make([]*probe.Peer, 0)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	var ixfData IXFResponse
	if err := json.Unmarshal(body, &ixfData); err != nil { // Parse []byte to go struct pointer
		return nil, err
	}

	swDb := make(map[int]string)

	for _, ixp := range ixfData.IxpList {
		// store switch names
		if ixp.IxpID == ixpId {
			for _, s := range ixp.Switch {
				swDb[s.ID] = s.Name
			}
		}
	}

	for _, member := range ixfData.MemberList {
		for _, conn := range member.ConnectionList {
			if conn.IxpID == ixpId {
				for _, v := range conn.VlanList {
					if v.VlanID == vlanId {

						result = append(result, probe.NewPeer(
							member.Name,
							swDb[conn.IfList[0].SwitchID],
							v.Ipv4.Address,
							v.Ipv6.Address,
						))

					}
				}
			}
		}
	}

	return result, nil
}

func ImportCSVPeers(filename string) ([]*probe.Peer, error) {
	result := make([]*probe.Peer, 0)

	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("[W] Error closing file: %s\n", err)
		}
	}(file)

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return result, err
	}
	for _, r := range records {
		peer := probe.NewPeer(r[0], r[1], r[2], r[3])
		result = append(result, peer)
	}

	return result, nil
}
