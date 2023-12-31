# BruteSprayX

![Version](https://img.shields.io/badge/Version-2.0.1-red)[![goreleaser](https://github.com/x90skysn3k/brutesprayx/actions/workflows/release.yml/badge.svg)](https://github.com/x90skysn3k/brutesprayx/actions/workflows/release.yml)

Created by: Shane Young/@t1d3nio && Jacob Robles/@shellfail 

Inspired by: Leon Johnson/@sho-luv

# Description
BruteSprayx is a golang version of the original BruteSpray. Without needing to rely on other tools this version will be extensible to bruteforce many different services and is way faster than it's Python counterpart. Currently BruteSprayX takes Nmap GNMAP/XML output, newline separated JSON, Nexpose `XML Export` output, Nessus `.nessus` exports, and lists. It will bruteforce supported servics found in those files. This tool is for research purposes and not intended for illegal use. 

<img src="https://imgur.com/HL5jP5W.png" width="500">

# Installation

[Release Binaries](https://github.com/x90skysn3k/brutesprayx/releases)

To Build:

```go build -o brutesprayx main.go```

# Usage

If using Nmap, scan with ```-oG nmap.gnmap``` or ```-oX nmap.xml```.

If using Nexpose, export the template `XML Export`. 

If using Nessus, export your `.nessus` file.

Command: ```brutesprayx -h```

Command: ```brutesprayx -f nmap.gnmap -u userlist -p passlist```

Command: ```brutesprayx -f nmap.xml -u userlist -p passlist```

Command: ```brutesprayx -H ssh://127.0.0.1:22 -u userlist -p passlist```


## Examples

<img src="brutesprayx.gif" width="512">

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
* mysql
* vmauthd

# Services in Progress

* rdp
* and more

# Data Specs
```json
{"host":"127.0.0.1","port":"3306","service":"mysql"}
{"host":"127.0.0.10","port":"3306","service":"mysql"}
```
If using Nexpose, export the template `XML Export`. 

If using Nessus, export your `.nessus` file.

List example
```
ssh:127.0.0.1:22
ftp:127.0.0.1:21
```
