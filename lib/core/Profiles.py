#!/usr/bin/python

"""JA3 Profiles"""

JA3Profiles = {
    "chrome": {
        "ciphers": [
            0x1301,
            0x1302,
            0x1303,
            0xC02B,
            0xC02F,
            0xC02C,
            0xC030,
            0xCCA9,
            0xCCA8,
        ],
        "curves": [29, 23, 24],
    },
    "firefox": {
        "ciphers": [
            0x1301,
            0x1302,
            0x1303,
            0xC02B,
            0xC02F,
            0xCCA9,
            0xCCA8,
            0xC02C,
            0xC030,
        ],
        "curves": [29, 23, 24, 25],
    },
    "safari": {
        "ciphers": [
            0x1301,
            0x1302,
            0x1303,
            0xC02C,
            0xC02B,
            0xC030,
            0xC02F,
            0xCCA9,
            0xCCA8,
        ],
        "curves": [29, 23, 24],
    },
}
