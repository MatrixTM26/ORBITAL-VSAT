#!/usr/bin/python


def GetCipherNames(self, CipherCodes):
    ChiperList = {
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
    return [ChiperList.get(c, "") for c in CipherCodes[:8] if c in ChiperList] or [
        "ECDHE+AESGCM"
    ]
