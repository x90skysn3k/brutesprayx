# BruteSprayX

![Version](https://img.shields.io/badge/Version-2.0-red)[![goreleaser](https://github.com/x90skysn3k/brutesprayx/actions/workflows/release.yml/badge.svg)](https://github.com/x90skysn3k/brutesprayx/actions/workflows/release.yml)

Created by: Shane Young/@t1d3nio && Jacob Robles/@shellfail 

Inspired by: Leon Johnson/@sho-luv

# Description
BruteSprayX takes Nmap GNMAP/XML output, newline separated JSON, Nexpose `XML Export` output or Nessus `.nessus` exports and automatically brute-forces services with default credentials. BruteSpray finds non-standard ports, make sure to use `-sV` with Nmap.

BrutesprayX is Brutespray but written in Go!

<img src="https://i.imgur.com/ZTS5be9.png" width="500">

# Installation

TODO

# Usage

If using Nmap, scan with ```-oG nmap.gnmap``` or ```-oX nmap.xml```.

If using Nexpose, export the template `XML Export`. 

If using Nessus, export your `.nessus` file.

Command: ```brutesprayx -h```

Command: ```brutesprayx -f nmap.gnmap -u userlist -p passlist```

Command: ```brutesprayx -f nmap.xml -u userlist -p passlist```

Command: ```brutesprayx -H ssh://127.0.0.1:22 -u userlist -p passlist```


## Examples

#### Using Custom Wordlists:

```brutesprayx -f nmap.gnmap -u /usr/share/wordlist/user.txt -p /usr/share/wordlist/pass.txt -t 5 ```

#### Brute-Forcing Specific Services:

```brutesprayx -f nmap.gnmap -u admin -p password -s ftp,ssh,telnet -t 5 ```

#### Specific Credentials:
   
```brutesprayx -f nmap.gnmap -u admin -p password -t 5 ```

#### Use Nmap XML Output

```brutesprayx -f nmap.xml -u admin -p password -t 5 ```

#### Use JSON Output

```brutesprayx -f out.json -u admin -p password -t 5 ```

# Supported Services

* ssh
* ftp
* telnet
* mssql
* postgresql
* imap
* pop3
* smbnt
* smtp
* snmp

# Services in Progress

* rdp
* mysql
* 

# Data Specs
```json
{"host":"127.0.0.1","port":"3306","service":"mysql"}
{"host":"127.0.0.10","port":"3306","service":"mysql"}
...
```
If using Nexpose, export the template `XML Export`. 

If using Nessus, export your `.nessus` file.

List example
```
ssh:127.0.0.1:22
ftp:127.0.0.1:21
```

# Changelog
Changelog notes are available at [CHANGELOG.md](https://github.com/x90skysn3k/brutesprayx/blob/master/CHANGELOG.md)
