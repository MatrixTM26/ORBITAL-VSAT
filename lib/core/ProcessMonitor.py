#!/usr/bin/env python

from time import sleep
from lib.core.ANSIColor import Color


class ProcessMonitor:
    def __init__(
        self, Running, RequestsCount, BytesSent, StatsLock, Host, Method, IP, Port
    ):
        self.Running = Running
        self.RequestsCount = RequestsCount
        self.BytesSent = BytesSent
        self.StatsLock = StatsLock
        self.Host = Host
        self.Method = Method
        self.IP = IP
        self.Port = Port

    def Monitor(self):
        LastCount = 0
        LastBytes = 0
        while self.Running.value:
            sleep(1)
            with self.StatsLock:
                Count = self.RequestsCount.value
                TotalBytes = self.BytesSent.value
            Diff = Count - LastCount
            ByteDiff = TotalBytes - LastBytes
            RPS = Diff
            Mbps = (ByteDiff * 8) / (1024 * 1024)
            LastCount = Count
            LastBytes = TotalBytes
            print(
                f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} REQUESTS: {Color.Green}[{Color.Red} {Count:,} {Color.Green}] {Color.White} TARGET: {Color.DarkGreen} {self.Host} {Color.White} METHODS: {Color.DarkGreen} {self.Method} {Color.White} IP: {Color.DarkGreen} {self.IP}:{self.Port} {Color.White} RPS: {Color.DarkGreen} {RPS:,} {Color.White} BW: {Color.DarkGreen} {Mbps:.1f} Mbps {Color.Reset}"
            )
