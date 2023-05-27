import json
import sys
import requests
import cf_sensitive

from colorama import Fore, Style

zone_id = cf_sensitive.zone_id
token = cf_sensitive.token


# https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records

def get_all_record():
    resp = requests.get(
        'https://api.cloudflare.com/client/v4/zones/{}/dns_records'.format(zone_id),
        headers={
            'Authorization': 'Bearer ' + token,
            'Content-Type': 'application/json'
        })
    if not json.loads(resp.text)['success']:
        return None
    domains = json.loads(resp.text)['result']
    print_format = '{:<36s} {:<32s} {:<4s} {:<50s} {:<9s} {:<7s}'
    print(print_format.format('id', 'name', 'type', 'content', 'proxiable', 'proxied'))
    print("-" * 150)
    domains = sorted(domains, key=lambda x: x['name'])
    for domain in domains:
        if domain['proxied']:
            proxied = Fore.GREEN + str(True) + Style.RESET_ALL
        else:
            proxied = str(False)
        print(print_format.format(domain['id'], domain['name'], domain['type'], domain['content'],
                                  str(domain['proxiable']), proxied))
    return None


def update_record(dns_id, dns_name, dns_type, dns_content, proxied=False):
    resp = requests.put(
        'https://api.cloudflare.com/client/v4/zones/{}/dns_records/{}'.format(
            zone_id, dns_id),
        json={
            'type': dns_type,
            'name': dns_name,
            'content': dns_content,
            'proxied': proxied
        },
        headers={
            'Authorization': 'Bearer ' + token,
            'Content-Type': 'application/json'
        })
    if not json.loads(resp.text)['success']:
        print(resp.text)
        return False
    return True


def add_record(dns_name, dns_type, dns_content, proxied=False):
    resp = requests.post(
        'https://api.cloudflare.com/client/v4/zones/{}/dns_records'.format(zone_id),
        json={
            'type': dns_type,
            'name': dns_name,
            'content': dns_content,
            'proxied': proxied
        },
        headers={
            'Authorization': 'Bearer ' + token,
            'Content-Type': 'application/json'
        })
    if not json.loads(resp.text)['success']:
        print(resp.text)
        return False
    return True


if __name__ == '__main__':
    # update_record('cdf095d296853b3b884d375e46a9904f', 'love', 'A', '1.1.1.1')
    # add_record('new_love', 'A', '8.8.8.8')
    # get_all_record()
    argv_len = len(sys.argv)
    action = 'list'
    if argv_len == 1:
        pass
    elif argv_len == 2:
        pass
    elif argv_len < 5 or argv_len > 6:
        print('argv error')
        sys.exit(-1)
    else:
        action = sys.argv[1]

    if action == 'list':
        get_all_record()
    elif action == 'add':
        add_record(sys.argv[2], sys.argv[3], sys.argv[4])
        get_all_record()
    elif action == 'edit':
        update_record(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5])
        get_all_record()
    elif action == 'edit_proxy':
        update_record(sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5], True)
        get_all_record()
    else:
        print('error action')
