import requests
import sys

SWITCHES = {}

if __name__ == '__main__':

    ixf_url = sys.argv[1]
    ixp_id = sys.argv[2]
    vlan_id = sys.argv[3]

    print(f"Requesting: {ixf_url}")
    response = requests.get(ixf_url)
    data = response.json()

    for ixp in data['ixp_list']:
        if ixp['ixp_id'] == ixp_id:
            for s in ixp['switch']:
                SWITCHES.update({ s['id']: s['name'] })


    for member in data['member_list']:
        member_name = member['name']
        for c in member['connection_list']:
            if c['ixp_id'] == ixp_id:
                for v in c['vlan_list']:
                    if v['vlan_id'] == vlan_id:
                        v4address = v.get('ipv4').get('address')
                        v6address = v.get('ipv6').get('address')
                        switch_id = c["if_list"][0]['switch_id']
                        switch_name = SWITCHES[switch_id]
                        print(f"{member_name};{switch_name};{v4address};{v6address}")
