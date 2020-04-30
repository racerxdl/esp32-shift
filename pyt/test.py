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

CheckInput()

print "HealthCheck"
x = CmdMsg()
x.cmd = CmdMsg.HealthCheck
SendMessage(x)

time.sleep(0.1)
CheckInput()


lastData = 1

while True:
  x = CmdMsg()
  x.cmd = CmdMsg.HealthCheck
  SendMessage(x)

  x = CmdMsg()
  x.cmd = CmdMsg.Status
  SendMessage(x)

  x.cmd = CmdMsg.SetByte
  x.data = chr(0) + chr(lastData & 0xFF)
  SendMessage(x)

  x.cmd = CmdMsg.SetByte
  x.data = chr(1) + chr(lastData >> 8)
  SendMessage(x)

  lastData <<= 1
  if lastData > 65535:
    lastData = 1
  time.sleep(0.01)
  CheckInput()

# while True:
#   x = CmdMsg()
#   x.cmd = CmdMsg.HealthCheck
#   SendMessage(x)
#   time.sleep(0.01)
#   CheckInput()

#   # Set to 1010 1010
#   x = CmdMsg()
#   x.cmd = CmdMsg.SetByte
#   x.data = chr(0) + chr(0b11110000)
#   SendMessage(x)
#   time.sleep(0.05)
#   CheckInput()
#   time.sleep(1)

#   # Set to 0101 0101
#   x = CmdMsg()
#   x.cmd = CmdMsg.SetByte
#   x.data = chr(0) + chr(0b00001111)
#   SendMessage(x)
#   time.sleep(0.05)
#   CheckInput()
#   time.sleep(1)

#   # Reset
#   # x = CmdMsg()
#   # x.cmd = CmdMsg.Reset
#   # SendMessage(x)
#   # time.sleep(0.1)
#   # CheckInput()
#   # time.sleep(1)
