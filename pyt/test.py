#!/usr/bin/env python

import struct, sys, time, serial

from shift_pb2 import *


ser = serial.Serial('/dev/ttyUSB0', 115200, timeout=0.1)

addr = 0x22


def SendMessage(msg):
  b = msg.SerializeToString()
  ser.write(struct.pack("H", len(b)))
  ser.write(b)

def CheckInput():
  if ser.in_waiting > 0:
    lines = ser.readlines()
    for line in lines:
      print line.split("\n")[0]

time.sleep(1)
CheckInput()

print "Resetting"
x = CmdMsg()
x.cmd = CmdMsg.Reset
SendMessage(x)

time.sleep(1)
CheckInput()

# x.cmd = CmdMsg.SetPin
# x.devAddr = addr
# x.data = "\x00\x00"
# SendMessage(x)

# lastState = False
# print "Looping"
# while True:
#   x.cmd = CmdMsg.HealthCheck
#   x.devAddr = addr
#   x.data = "A"
#   SendMessage(x)

#   # x.cmd = CmdMsg.SetGPIOAB
#   # x.data = "\xFF\xFF"
#   # x.data = "\xFF" if lastState else "\x00"
#   # x.data = "\xFF"
#   # x.data += "\xFF" if lastState else "\x00"
#   # SendMessage(x)

#   x.cmd = CmdMsg.SetPin
#   x.data = chr(12)
#   x.data += "\x01" if lastState else "\x00"
#   SendMessage(x)

#   # for i in range(9):
#   #   x.cmd = CmdMsg.SetPin
#   #   x.devAddr = addr
#   #   x.data = chr(i)
#   #   if lastState:
#   #     x.data += "\x01"
#   #   else:
#   #     x.data += "\x00"
#   #   SendMessage(x)

#   lastState = not lastState

#   time.sleep(0.5)
