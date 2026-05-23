#!/usr/bin/python3
from lib.core.ANSIColor import Color

xbanner = f"""
{Color.Bold}{Color.Red}
________ ____________________.___________________  .____      ____   _____________   ________________
\\_____  \\\\______   \\______   \\   \\__    ___/  _  \\ |    |     \\   \\ /   /   _____/  /  _  \\__    ___/
 /   |   \\|       _/|    |  _/   | |    | /  /_\\  \\|    |      \\   Y   /\\_____  \\  /  /_\\  \\|    |   
/    |    \\    |   \\|    |   \\   | |    |/    |    \\    |___    \\     / /        \\/    |    \\    |   
\\_______  /____|_  /|______  /___| |____|\\____|__  /_______ \\    \\___/ /_______  /\\____|__  /____|   
        \\/       \\/        \\/                    \\/        \\/                  \\/         \\/         
{Color.Reset}
"""

xhelper = f"""
    {Color.White}[{Color.Red}M{Color.White}]: {Color.Cyan} Available Methods:
        {Color.White}[{Color.Red}+{Color.White}] {Color.White} LAYER 7: {Color.Orange} APPLICATION
            {Color.White}     ->{Color.Orange} GET                   \t: {Color.Orange} HTTP GET Flood With Cache Bypass
            {Color.White}     ->{Color.Orange} POST                  \t: {Color.Orange} HTTP POST Flood [ 64KB Payloads ]
            {Color.White}     ->{Color.Orange} PUT                   \t: {Color.Orange} HTTP PUT Flood
            {Color.White}     ->{Color.Orange} HEAD                  \t: {Color.Orange} HTTP HEAD Flood
            {Color.White}     ->{Color.Orange} DELETE                   \t: {Color.Orange} HTTP DELETE Flood
            {Color.White}     ->{Color.Orange} PATCH                 \t: {Color.Orange} HTTP PATCH Flood
            {Color.White}     ->{Color.Orange} OPTIONS                   \t: {Color.Orange} HTTP OPTIONS Flood

        {Color.White}[{Color.Red}+{Color.White}]{Color.White} LAYER 7: {Color.Orange} ADVANCED METHODS
            {Color.White}     ->{Color.Orange} XMLRPC                    \t: {Color.Orange} XML-RPC Pingback attack
            {Color.White}     ->{Color.Orange} RANDOM                    \t: {Color.Orange} Random HTTP Methods Flood
            {Color.White}     ->{Color.Orange} SLOWLORIS                 \t: {Color.Orange} Slowloris Attack [ Keep Alive ]
            {Color.White}     ->{Color.Orange} SLOW-POST                 \t: {Color.Orange} Slow-POST Body Attack
            {Color.White}     ->{Color.Orange} CACHE                 \t: {Color.Orange} Cache Bypass Flood
            {Color.White}     ->{Color.Orange} BYPASS                    \t: {Color.Orange} WAF Bypass Techniques
            {Color.White}     ->{Color.Orange} CONNECT                   \t: {Color.Orange} HTTP CONNECT Flood
            {Color.White}     ->{Color.Orange} TRACE                 \t: {Color.Orange} HTTP TRACE Flood
            {Color.White}     ->{Color.Orange} SLOW-READ                 \t: {Color.Orange} Slow-READ Body Attack
            {Color.White}     ->{Color.Orange} RUDY                  \t: {Color.Orange} ARE-YOU-DEAD-YET Attack

        {Color.White}[{Color.Red}+{Color.White}]{Color.White} LAYER 7: {Color.Orange} HTTP/2 | HTTP/3
            {Color.White}     ->{Color.Orange} H2-GET                    \t: {Color.Orange} HTTP/2 GET With Priority
            {Color.White}     ->{Color.Orange} H2-POST                   \t: {Color.Orange} HTTP/2 POST With Multiplexing
            {Color.White}     ->{Color.Orange} H2-RAPID                  \t: {Color.Orange} HTTP/2 RAPID Reset
            {Color.White}     ->{Color.Orange} H2-PING                   \t: {Color.Orange} HTTP/2 PING Flood
            {Color.White}     ->{Color.Orange} H3-GET                    \t: {Color.Orange} HTTP/3 QUIC GET
            {Color.White}     ->{Color.Orange} H3-POST                   \t: {Color.Orange} HTTP/3 QUIC POST

        {Color.White}[{Color.Red}+{Color.White}]{Color.White} LAYER 4: {Color.Orange} TRANSPORT
            {Color.White}     ->{Color.Orange} TCP                   \t: {Color.Orange} TCP Connection Flood
            {Color.White}     ->{Color.Orange} UDP                   \t: {Color.Orange} UDP Packet Flood [ 64KB Payloads ]
            {Color.White}     ->{Color.Orange} SYN                   \t: {Color.Orange} TCP SYN Flood {Color.White} [{Color.Orange} REQUIRES ROOT {Color.White}]
            {Color.White}     ->{Color.Orange} ACK                   \t: {Color.Orange} TCP ACK Flood {Color.White} [{Color.Orange} REQUIRES ROOT {Color.White}]
            {Color.White}     ->{Color.Orange} RST                   \t: {Color.Orange} TCP RST Flood {Color.White} [{Color.Orange} REQUIRES ROOT {Color.White}]
            {Color.White}     ->{Color.Orange} FIN                   \t: {Color.Orange} TCP FIN Flood {Color.White} [{Color.Orange} REQUIRES ROOT {Color.White}]
            {Color.White}     ->{Color.Orange} SYNACK                    \t: {Color.Orange} TCP SYN-ACK Reflection
            {Color.White}     ->{Color.Orange} PSH                   \t: {Color.Orange} TCP PSH + ACK Flood
            {Color.White}     ->{Color.Orange} URG                   \t: {Color.Orange} TCP URG Flood
            {Color.White}     ->{Color.Orange} XMAS                  \t: {Color.Orange} TCP XMAS SCAN Flood
            {Color.White}     ->{Color.Orange} NULL                  \t: {Color.Orange} TCP NULL SCAN Flood

        {Color.White}[{Color.Red}+{Color.White}]{Color.White} LAYER 4: {Color.Orange} AMPLIFICATIONS
            {Color.White}     ->{Color.Orange} UDP-FRAG                  \t: {Color.Orange} UDP FRAGMENTATION Flood
            {Color.White}     ->{Color.Orange} DNS-AMP                   \t: {Color.Orange} DNS AMPLIFICATION
            {Color.White}     ->{Color.Orange} NTP-AMP                   \t: {Color.Orange} NTP AMPLIFICATION
            {Color.White}     ->{Color.Orange} SSDP-AMP                  \t: {Color.Orange} SSDP AMPLIFICATION
            {Color.White}     ->{Color.Orange} MEMCACHED                 \t: {Color.Orange} MEMCACHED AMPLIFICATION
            {Color.White}     ->{Color.Orange} CHARGEN                   \t: {Color.Orange} CHARGEN AMPLIFICATION

        {Color.White}[{Color.Red}+{Color.White}]{Color.White} LAYER 3: {Color.Orange} NETWORKS
            {Color.White}     ->{Color.Orange} ICMP                  \t: {Color.Orange} ICMP Ping Flood {Color.White} [{Color.Orange} REQUIRES ROOT {Color.White}]
            {Color.White}     ->{Color.Orange} PING                  \t: {Color.Orange} PING Flood
            {Color.White}     ->{Color.Orange} SMURF                 \t: {Color.Orange} SMURF Attack
            {Color.White}     ->{Color.Orange} FRAGGLE                   \t: {Color.Orange} FRAGGLE Attack [ UDP + ECHO ]

    {Color.White}[{Color.Red}C{Color.White}]: {Color.Cyan} coresuration Supports:
            {Color.White}[{Color.Red}+{Color.White}]: {Color.White} User-Agent Headers Randomization
            {Color.White}[{Color.Red}+{Color.White}]: {Color.White} Proxy Address Randomization/Proxychaining
            {Color.White}[{Color.Red}+{Color.White}]: {Color.White} Referers Randomization/Requests Sources Randomization
            {Color.White}[{Color.Red}+{Color.White}]: {Color.White} HTTP1 | HTTP2 | HTTP3 coresurations Support
            {Color.White}[{Color.Red}+{Color.White}]: {Color.White} Autofingerprinting JA3, TLS, Browser Like Requests [ chrome | firefox | safari ]

    {Color.White}[{Color.Red}EXIT{Color.White}]: {Color.Red} CTRL + C To Stop.

{Color.Reset}
"""


def Banner():
    print(xbanner)


def Helper():
    print(xhelper)
