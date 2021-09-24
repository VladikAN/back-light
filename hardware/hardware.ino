#include <Adafruit_NeoPixel.h>

#define BUFFER_SIZE 256
#define PIN 5
#define NUMPIXELS 16
Adafruit_NeoPixel pixels = Adafruit_NeoPixel(NUMPIXELS, PIN, NEO_GRB + NEO_KHZ800);

String inputString;
int inputLed;

void setup() {
  Serial.begin(9600);
  
  inputString.reserve(BUFFER_SIZE);
  inputLed = 0;

  pixels.begin();
  pixels.show();
  pixels.setBrightness(200);

  Serial.println("Ready\n");
}

void loop() {}

void serialEvent() {
  while (Serial.available()) { 
    char ch = (char)Serial.read();
    if (ch == '\n' || ch == '\0') {
      inputString = "";
      inputLed = 0;
    } else if (ch == ';') {
      long number = strtol( &inputString[0], NULL, 16);
      long r = number >> 16;
      long g = number >> 8 & 0xFF;
      long b = number & 0xFF;

      pixels.setPixelColor(inputLed, pixels.Color(r, g, b));
      pixels.show();

      inputString = "";
      inputLed += 1;
      if (inputLed >= NUMPIXELS) {
        inputLed = 0;
      }
    } else {
      inputString.concat(ch);
    }
  }
}
