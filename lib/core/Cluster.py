#!/usr/bin/env python

from time import sleep
from concurrent.futures import ThreadPoolExecutor


class Cluster:
    def __init__(self, MethodsInstance, Threads, Running):
        self.Methods = MethodsInstance
        self.Threads = Threads
        self.Running = Running

    def GetMethodExecutor(self, MethodName):
        MethodOptions = {
            "GET": self.Methods.HTTPExecutor,
            "POST": self.Methods.HTTPExecutor,
            "PUT": self.Methods.HTTPExecutor,
            "HEAD": self.Methods.HTTPExecutor,
            "DELETE": self.Methods.HTTPExecutor,
            "PATCH": self.Methods.HTTPExecutor,
            "OPTIONS": self.Methods.HTTPExecutor,
            "CONNECT": self.Methods.HTTPExecutor,
            "TRACE": self.Methods.HTTPExecutor,
            "RANDOM": self.Methods.HTTPExecutor,
            "SLOWLORIS": self.Methods.SlowlorisExecutor,
            "SLOWPOST": self.Methods.SlowPostExecutor,
            "H2GET": self.Methods.H2Executor,
            "H2POST": self.Methods.H2Executor,
            "H2PING": self.Methods.H2PingExecutor,
            "TCP": self.Methods.TCPExecutor,
            "SYN": self.Methods.SYNExecutor,
            "ACK": self.Methods.ACKExecutor,
            "RST": self.Methods.RSTExecutor,
            "FIN": self.Methods.FINExecutor,
            "XMAS": self.Methods.XMASExecutor,
            "UDP": self.Methods.UDPExecutor,
            "UDPFRAG": self.Methods.UDPFragExecutor,
            "DNSAMP": self.Methods.DNSAmpExecutor,
            "NTPAMP": self.Methods.NTPAmpExecutor,
            "ICMP": self.Methods.ICMPExecutor,
        }
        return MethodOptions.get(MethodName, self.Methods.HTTPExecutor)

    def ProcessingTask(self, ProcessID, MethodName):
        ThreadsExecutor = self.GetMethodExecutor(MethodName)
        with ThreadPoolExecutor(max_workers=self.Threads) as Executor:
            Futures = [Executor.submit(ThreadsExecutor, I) for I in range(self.Threads)]
            for Future in Futures:
                Future.result()
            while self.Running.value:
                sleep(1)
