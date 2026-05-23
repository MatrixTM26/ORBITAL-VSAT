#!/usr/bin/env python

import os
from sys import stdout
from time import sleep


def Clear():
    os.system("cls" if os.name == "nt" else "clear")


class Logging:
    @staticmethod
    def Typewriter(text, delay=0.001):
        for char in text:
            stdout.write(char)
            stdout.flush()
            sleep(delay)
        stdout.write("\n")
        stdout.flush()
