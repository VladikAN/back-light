#include <Adafruit_NeoPixel.h>

#define SERIAL_SPEED 9600
#define BUFFER_SIZE 256
#define PIXELS_PIN 5
#define PIXELS_BRIGHTLESS 10
#define PIXELS_NUM 12
#define ANIM_DELAY 100

Adafruit_NeoPixel pixels = Adafruit_NeoPixel(
  PIXELS_NUM,
  PIXELS_PIN,
  NEO_GRB + NEO_KHZ800);

enum Animations {
  None,
  Started,
};

Animations _currentAnimation = Animations::None;
String inputString;
int inputLed;

void setup() {
  Serial.begin(SERIAL_SPEED);
  
  inputString.reserve(BUFFER_SIZE);
  inputLed = 0;

  pixels.begin();
  pixels.show();
  pixels.setBrightness(PIXELS_BRIGHTLESS);

  setAnimationStarted();
}

void loop() {
  runAnimations();
}

void serialEvent() {
  if (_currentAnimation != Animations::None) {
    return;
  }

  while (Serial.available()) { 
    char ch = (char)Serial.read();
    if (ch == '\n' || ch == '\0') {
      // The rest of the leds
      while (inputLed != 0 && inputLed < PIXELS_NUM) {
        pixels.setPixelColor(inputLed, pixels.Color(0, 0, 0));
        pixels.show();
        inputLed++;
      }

      // reset vars
      inputString = "";
      inputLed = 0;
    } else if (ch == ';') {
      long number = inputString.length() > 0
        ? strtol(&inputString[0], NULL, 16)
        : 0;
      inputString = "";

      long r = number >> 16;
      long g = number >> 8 & 0xFF;
      long b = number & 0xFF;
      
      pixels.setPixelColor(inputLed, pixels.Color(r, g, b));
      pixels.show();

      inputLed = inputLed < PIXELS_NUM - 1
        ? inputLed + 1
        : 0;
    } else {
      inputString.concat(ch);
    }
  }
}

void setAnimationStarted() {
  _currentAnimation = Animations::Started;
}

void runAnimations() {
  if (_currentAnimation == Animations::None) {
    return;
  }

  if (_currentAnimation == Animations::Started) {
    runStartedAnimation();
    return;
  }
}

void runStartedAnimation() {
  int frame = 0;
  while (frame < PIXELS_NUM) {
      pixels.setPixelColor(frame, pixels.Color(255, 0, 0));
      pixels.show();

      frame += 1;
      delay(ANIM_DELAY);
    }

    pixels.clear();
    pixels.show();
    _currentAnimation = Animations::None;
}