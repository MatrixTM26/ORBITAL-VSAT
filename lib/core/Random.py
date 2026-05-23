#!/usr/bin/env python

import random


class RandomGenerator:
    @staticmethod
    def RandomString(Length=10):
        return "".join(
            random.choice("abcdefghijklmnopqrstuvwxyz0123456789") for _ in range(Length)
        )

    @staticmethod
    def RandomIP():
        return f"{random.randint(1, 254)}.{random.randint(1, 254)}.{random.randint(1, 254)}.{random.randint(1, 254)}"

    @staticmethod
    def CalculateChecksum(Data):
        Sum = 0
        for I in range(0, len(Data), 2):
            if I + 1 < len(Data):
                Sum += (Data[I] << 8) + Data[I + 1]
            else:
                Sum += Data[I] << 8
        Sum = (Sum >> 16) + (Sum & 0xFFFF)
        Sum += Sum >> 16
        return ~Sum & 0xFFFF
