#!/usr/bin/env python

import socket
import ssl
from lib.core.Profiles import JA3Profiles


class Fingerprint:
    @staticmethod
    def GetCipherNames(CipherCodes):
        CipherList = {
            0x1301: "TLS_AES_128_GCM_SHA256",
            0x1302: "TLS_AES_256_GCM_SHA384",
            0x1303: "TLS_CHACHA20_POLY1305_SHA256",
            0xC02B: "ECDHE-ECDSA-AES128-GCM-SHA256",
            0xC02F: "ECDHE-RSA-AES128-GCM-SHA256",
            0xC02C: "ECDHE-ECDSA-AES256-GCM-SHA384",
            0xC030: "ECDHE-RSA-AES256-GCM-SHA384",
            0xCCA9: "ECDHE-ECDSA-CHACHA20-POLY1305",
            0xCCA8: "ECDHE-RSA-CHACHA20-POLY1305",
        }
        return [CipherList.get(C, "") for C in CipherCodes[:8] if C in CipherList] or [
            "ECDHE+AESGCM"
        ]

    @staticmethod
    def CreateJa3Socket(IP, Port, Scheme, Host, Protocol, JAProfile):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            Sock.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)
            Sock.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
            Sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            Sock.settimeout(3)
            Sock.connect((IP, Port))
            if Scheme == "https":
                Context = ssl.SSLContext(ssl.PROTOCOL_TLS_CLIENT)
                Context.check_hostname = False
                Context.verify_mode = ssl.CERT_NONE
                Context.options |= ssl.OP_NO_TLSv1 | ssl.OP_NO_TLSv1_1
                Profile = JA3Profiles[JAProfile]
                CipherNames = Fingerprint.GetCipherNames(Profile["Ciphers"])
                try:
                    Context.set_ciphers(":".join(CipherNames))
                except Exception:
                    Context.set_ciphers("ECDHE+AESGCM:!aNULL")
                if Protocol == "H2":
                    Context.set_alpn_protocols(["h2", "http/1.1"])
                elif Protocol == "H3":
                    Context.set_alpn_protocols(["h3"])
                else:
                    Context.set_alpn_protocols(["http/1.1"])
                Sock = Context.wrap_socket(Sock, server_hostname=Host)
            return Sock
        except Exception:
            return None
