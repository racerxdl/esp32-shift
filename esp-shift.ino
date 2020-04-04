#include <SPI.h>

#include "shift.pb.h"


#include "pb_common.h"
#include "pb.h"
#include "pb_decode.h"

#define LATCH 13
#define MAX_BYTES 128
#define NUM_REGISTERS 2

uint8_t states[NUM_REGISTERS * 8];
uint8_t curPos = 0;
uint8_t bstate[NUM_REGISTERS];

SPIClass * vspi = NULL;
static const int spiClk = 10000000; // 10 MHz


// ProtoBuf Part
CmdMsg message = CmdMsg_init_default;
uint8_t buffer[MAX_BYTES];
uint16_t numBytes = 0;

// Code

void setup() {
  Serial.begin(115200);
  Serial.println("Initializing SPI");
  vspi = new SPIClass(VSPI);
  vspi->begin();

  // put your setup code here, to run once:
  pinMode(LATCH, OUTPUT);
  clear();
}

void clear() {
  for (int i = 0; i < NUM_REGISTERS * 8; i++) {
    states[i] = HIGH;
  }
}

void updatebstates() {
  for (int i = 0; i < NUM_REGISTERS; i++) {
    bstate[i] = 0;
  }
  for (int i = 0; i < NUM_REGISTERS * 8; i++) {
    int s = i / 8;
    int o = i % 8;
    bstate[s] |= states[i] * (1 << o);
  }
}

void update() {
  updatebstates();
  vspi->beginTransaction(SPISettings(spiClk, MSBFIRST, SPI_MODE0));
  // Since is shifted, use reverse order
  for (int i = NUM_REGISTERS-1; i >= 0; i--) {
    vspi->transfer(bstate[i]);
  }
  vspi->endTransaction();
}

void latch() {
  digitalWrite(LATCH, HIGH);
  delayMicroseconds(1);
  digitalWrite(LATCH, LOW);
  delayMicroseconds(1);
}

void processPayload() {
  pb_istream_t stream = pb_istream_from_buffer(buffer, numBytes);
  if (!pb_decode(&stream, CmdMsg_fields, &message)) {
    Serial.println("( SER) Error parsing message");
    return;
  }

  Serial.print("Received command ");
  Serial.println(message.cmd);
}

void receivePayload() {
  int n = Serial.available();
  if (n >= 2) {
    buffer[0] = Serial.read();
    buffer[1] = Serial.read();

    numBytes = *((uint16_t*) buffer);
    Serial.print("Wants to receive ");
    Serial.println(numBytes);

    if (numBytes >= MAX_BYTES) {
      Serial.print("Wanted to receive ");
      Serial.print(numBytes);
      Serial.print("bytes. But max is MAX_BYTES");
      return;
    }

    for (int i = 0; i < numBytes; i++) {
      buffer[i] = Serial.read();
    }

    processPayload();
  }
}

void loop() {
  receivePayload();
  update();
  latch();
//  clear();
//  curPos++;
//  if (curPos == NUM_REGISTERS * 8) {
//    curPos = 0;
//  }
//  states[curPos] = LOW;
  delay(1);
}
