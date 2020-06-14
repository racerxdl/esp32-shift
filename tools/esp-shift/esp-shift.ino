#include <SPI.h>

#include "shift.pb.h"


#include "pb_common.h"
#include "pb.h"
#include "pb_decode.h"

#define ONBOARD_LED 2
#define OUTPUT_ENABLE 14
#define RESET 27
#define LATCH 13
#define MAX_BYTES 128
#define NUM_REGISTERS 8
#define NUM_BITS (NUM_REGISTERS * 8)

uint8_t states[NUM_BITS];
uint8_t curPos = 0;
uint8_t bstate[NUM_REGISTERS];

uint8_t connstatus[NUM_REGISTERS/2];
int     connpins[NUM_REGISTERS/2] = { 34, 35, 32, 33 };

SPIClass * vspi = NULL;
static const int spiClk = 10000000; // 10 MHz


// ProtoBuf Part
CmdMsg message = CmdMsg_init_default;
uint8_t buffer[MAX_BYTES];
uint16_t numBytes = 0;


// mapIO maps the correct pins from Shift Register to Output
// 
// The Relay Board expects this pinout:
//  --------------------------------------------------
// | 5V |  2 |  4 |  6 |  8 | 10 | 12 | 14 | 16 | GND |
// | 5V |  1 |  3 |  5 |  7 |  9 | 11 | 13 | 15 | GND |
//  --------------------------------------------------
//
// But instead we mapped (by ease of soldering) like this:
//  --------------------------------------------------
// | 5V |  7 |  6 |  5 |  4 |  3 |  2 |  1 |  0 | GND |
// | 5V |  8 |  9 | 10 | 11 | 12 | 13 | 14 | 15 | GND |
//  --------------------------------------------------
// 
// So this function receives a IO Pin number and converts to Shift Register Bit number
uint8_t mapIO(uint8_t iopin) {
  uint8_t ioPort = iopin / 16;                    // Find which IO Port this pin belongs
  uint8_t localIoPinNum = iopin - (ioPort * 16);  // Find the local pin number inside this IO Port
  uint8_t shiftRegisterBitNum;                    // Store the result bit

  // if the IOPort > 1, the shiftregisters are inverted
  if (ioPort >= 2) {
    localIoPinNum = 15 - localIoPinNum;
  }
  // Say we have for each I/O the Shift Registers A and B
  // The output will be mapped like this
  //  --------------------------------------------------
  // | 5V | A7 | A6 | A5 | A4 | A3 | A2 | A1 | A0 | GND |
  // | 5V | B0 | B1 | B2 | B3 | B4 | B5 | B6 | B7 | GND |
  //  --------------------------------------------------
  //                        to
  //  --------------------------------------------------
  // | 5V |  1 |  3 |  5 |  7 |  9 | 11 | 13 | 15 | GND |
  // | 5V |  0 |  2 |  4 |  6 |  8 | 10 | 12 | 14 | GND |
  //  --------------------------------------------------
  //  So the expected IO to Shift mapping is:
  //    Shift A has odd numbers increasing
  //    Shift B has even numbers decreasing
  uint8_t shiftRegisterNum = localIoPinNum % 2; // If odd, shift register 0, if even shift register 1
  uint8_t inShiftPin = localIoPinNum / 2;       // Now we can map directly

  if ( shiftRegisterNum == 0 ) {                // Register A
    shiftRegisterBitNum = 7 - inShiftPin;       // Register A just reverse
  } else {
    inShiftPin += 1;                            // One bit shifted
    shiftRegisterBitNum = 7 + inShiftPin;       // Add offset of the second shiftregister
  }
  
  shiftRegisterBitNum += ioPort * 16;           // Add the IO Port offset

  return shiftRegisterBitNum;
}

// Code
void setup() {
  Serial.begin(115200);
  Serial.println("(STS) Initializing SPI");
  vspi = new SPIClass(VSPI);
  vspi->begin();

  Serial.println("(STS) Setting pin modes");
  pinMode(LATCH, OUTPUT);
  pinMode(OUTPUT_ENABLE, OUTPUT);
  pinMode(RESET, OUTPUT);
  pinMode(ONBOARD_LED, OUTPUT);

  for (int i = 0; i < NUM_REGISTERS/2; i++) {
    pinMode(connpins[i], INPUT);
    connstatus[i] = digitalRead(connpins[i]);
  }
  
  Serial.println("(STS) Resetting flags");
  clear();
  digitalWrite(OUTPUT_ENABLE, HIGH);
  digitalWrite(RESET, HIGH);
  digitalWrite(ONBOARD_LED, LOW);
  Serial.println("(STS) READY");
}

void clear() {
  for (int i = 0; i < NUM_BITS; i++) {
    states[i] = HIGH;
  }
}

void updateconn() {
  for (int i = 0; i < NUM_REGISTERS/2; i++) {
    connstatus[i] = digitalRead(connpins[i]);
  }
}

void updatebstates() {
  for (int i = 0; i < NUM_REGISTERS; i++) {
    bstate[i] = 0;
  }
  for (int i = 0; i < NUM_BITS; i++) {
    int s = i / 8;
    int o = i % 8;
    bstate[s] |= states[i] * (1 << o);
  }
}

void updateByte(uint8_t num, uint8_t val) {
  uint8_t offset = num * 8;

  for (int i = 0; i < 8; i++) {
    uint8_t shiftpin = mapIO(offset+i);
    states[shiftpin] = val & (1 << i) ? LOW : HIGH;
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
  digitalWrite(OUTPUT_ENABLE, LOW);
}

void CmdReset() {
  clear();
  Serial.println("( OK) All bits reset");
}

void CmdSetPin(uint8_t *data) {
  uint8_t pin = data[0];
  uint8_t val = data[1] ? LOW : HIGH;

  if (pin >= NUM_BITS) {
    Serial.println("( ERR) Pin number cannot be higher than NUM_BITS");
    return;
  }

  uint8_t shiftpin = mapIO(pin); // Remap

  states[shiftpin] = val;
  Serial.print("( OK) Pin ");
  Serial.print(pin);
  Serial.print(" state set to ");
  Serial.println(val);
}

void CmdSetByte(uint8_t *data) {
  uint8_t numByte = data[0];
  uint8_t val = data[1];

  if (numByte >= NUM_REGISTERS) {
    Serial.println("( ERR) Byte number cannot be higher than NUM_REGISTERS");
    return;
  }

  updateByte(numByte, val);
  Serial.print("( OK) Byte ");
  Serial.print(numByte);
  Serial.print(" state set to ");
  Serial.println(val, BIN);
}

long lastHC = millis();

void CmdHealthCheck() {
  Serial.println("( OK) Health Check OK");
  digitalWrite(ONBOARD_LED, HIGH);
  lastHC = millis();
}

void CmdStatus() {
  Serial.print("( OK) BS[");
  for (int i = 0; i < NUM_REGISTERS/2; i++) {
    Serial.print(connstatus[i]);
    if (i < (NUM_REGISTERS/2)-1) {
      Serial.print(", ");
    }
  }
  Serial.println("]");
}

void processPayload() {
  pb_istream_t stream = pb_istream_from_buffer(buffer, numBytes);
  if (!pb_decode(&stream, CmdMsg_fields, &message)) {
    Serial.println("( ERR) Error parsing message");
    return;
  }

  switch (message.cmd) {
    case CmdMsg_Command_HealthCheck:  return CmdHealthCheck();
    case CmdMsg_Command_SetPin:
      if (message.data.size == 2) {
        return CmdSetPin(message.data.bytes);
      }
      Serial.println("( ERR) Expected 2 bytes for SetPin");
      break;
    case CmdMsg_Command_SetByte:
      if (message.data.size == 2) {
        return CmdSetByte(message.data.bytes);
      }
      Serial.println("( ERR) Expected 2 bytes for SetByte");
      break;
    case CmdMsg_Command_Reset: return CmdReset();
    case CmdMsg_Command_Status: return CmdStatus();
  }
}

void receivePayload() {
  int n = Serial.available();
  if (n >= 2) {
    buffer[0] = Serial.read();
    buffer[1] = Serial.read();

    numBytes = *((uint16_t*) buffer);

    if (numBytes >= MAX_BYTES) {
      Serial.print("( ERR) Wanted to receive ");
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
  updateconn();

  if (millis() - lastHC >= 10) {
    digitalWrite(ONBOARD_LED, LOW);
    lastHC = millis();
  }
}
