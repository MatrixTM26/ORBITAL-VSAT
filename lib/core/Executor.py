#!/usr/bin/env python

import multiprocessing as MP
import threading
from time import sleep
from concurrent.futures import ThreadPoolExecutor
from lib.core.ANSIColor import Color
from lib.core.StdIO import Logging
from lib.core.Methods import Methods
from lib.core.Cluster import Cluster
from lib.core.ProcessMonitor import ProcessMonitor


class Executor:
    def __init__(self, Config):
        self.Config = Config
        self.Running = Config.get("Running")
        self.RequestsCount = Config.get("RequestsCount")
        self.BytesSent = Config.get("BytesSent")
        self.StatsLock = Config.get("StatsLock")
        self.Method = Config.get("Method")
        self.Threads = Config.get("Threads")
        self.Duration = Config.get("Duration")
        self.ClusterMode = Config.get("ClusterMode")
        self.Processes = Config.get("Processes")
        self.Host = Config.get("Host")
        self.IP = Config.get("IP")
        self.Port = Config.get("Port")

    def StartMonitor(self):
        Monitor = ProcessMonitor(
            self.Running,
            self.RequestsCount,
            self.BytesSent,
            self.StatsLock,
            self.Host,
            self.Method,
            self.IP,
            self.Port,
        )
        StatsThread = threading.Thread(target=Monitor.Monitor, daemon=True)
        StatsThread.start()

    def StartClusterMode(self):
        MethodsInstance = Methods(self.Config)
        ClusterInstance = Cluster(MethodsInstance, self.Threads, self.Running)

        Processes = []
        for I in range(self.Processes):
            Process = MP.Process(
                target=ClusterInstance.ProcessingTask, args=(I, self.Method)
            )
            Process.start()
            Processes.append(Process)
            sleep(0.02)

        Logging.Typewriter(
            f"{Color.Cyan}[{Color.Red} ORBITAL VSAT {Color.Cyan}] {Color.Cyan} CLUSTER: {Color.Orange} {self.Processes * self.Threads} {Color.Orange} THREADS ACTIVE!\n {Color.Reset}"
        )

        try:
            sleep(self.Duration)
        except KeyboardInterrupt:
            pass

        self.Running.value = 0

        for Process in Processes:
            Process.join(timeout=2)
            if Process.is_alive():
                Process.terminate()

    def StartSingleMode(self):
        MethodsInstance = Methods(self.Config)
        ClusterInstance = Cluster(MethodsInstance, self.Threads, self.Running)

        ThreadsExecutor = ClusterInstance.GetMethodExecutor(self.Method)

        with ThreadPoolExecutor(max_workers=self.Threads) as Executor:
            Futures = [Executor.submit(ThreadsExecutor, I) for I in range(self.Threads)]
            for Future in Futures:
                Future.result()

            Logging.Typewriter(
                f"{Color.Cyan}[{Color.Red} ORBITAL VSAT {Color.Cyan}] {Color.Cyan} RUNNING: {Color.Orange} {self.Threads} {Color.Orange} THREADS!\n {Color.Reset}"
            )

            try:
                sleep(self.Duration)
            except KeyboardInterrupt:
                pass

            self.Running.value = 0

    def Execute(self):
        self.Running.value = 1

        Logging.Typewriter(f"\n{Color.DarkGreen} {'=' * 100}")
        Logging.Typewriter(
            f"{Color.Cyan}[{Color.Red} ORBITAL VSAT {Color.Cyan}] {Color.Cyan} STARTING ATTACK {Color.Orange} {self.Method} {Color.Reset}"
        )
        Logging.Typewriter(f"{Color.DarkGreen} {'=' * 100}\n")

        self.StartMonitor()

        if self.ClusterMode:
            self.StartClusterMode()
        else:
            self.StartSingleMode()

        sleep(2)

        with self.StatsLock:
            FinalCount = self.RequestsCount.value
            FinalBytes = self.BytesSent.value

        Logging.Typewriter(
            f"{Color.Cyan}[{Color.Red} ORBITAL VSAT {Color.Cyan}] {Color.Cyan} FINAL RESULTS {Color.Reset}"
        )
        Logging.Typewriter(f"{Color.DarkGreen} {'=' * 100}")
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} TOTAL REQUESTS: {Color.White} {FinalCount:,} {Color.Reset}"
        )
        Logging.Typewriter(
            f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} TOTAL SENT: {Color.White} {FinalBytes / 1048576:.2f} {Color.Cyan} MB {Color.Reset}"
        )

        if self.Duration > 0:
            Logging.Typewriter(
                f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} AVG RPS: {Color.White} {FinalCount / self.Duration:.0f} {Color.Reset}"
            )
            Logging.Typewriter(
                f"{Color.White}[{Color.Cyan} INFO {Color.White}] {Color.Cyan} AVG BANDWIDTH: {Color.White} {(FinalBytes * 8) / (self.Duration * 1048576):.2f} {Color.Cyan} Mbps {Color.Reset}"
            )
