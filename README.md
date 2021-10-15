# Back Light

![version - 1 - demo](/pics/demo-min.jpg)

I want to have an extra screen light. Depending on that happens on the screen I will get an extra colorful light behind the monitor or laptop.

This repository consists of two programs:

* hardware - is an arduino program, reading serial input and light the LEDs based on incoming information.

* software - is an Golang command line utility. Keep tracks the screen in order to find dominant colors and transfer this information to arduino.

## Hardware

Is an Arduino Nano with LED ring (CJMCU-2812-16) connected to 5V, GND and PIN 5.