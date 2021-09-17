#define BUFFER_SIZE 256

String inputString;
boolean stringComplete;

void setup() {
  Serial.begin(9600);
  inputString.reserve(BUFFER_SIZE);

  Serial.println("Ready\n");
}

void loop() {
  if (stringComplete) {
    draw();
  }
}

void draw() {
  Serial.print("in ");
  Serial.print(inputString.length());
  Serial.print(" bytes: ");
  Serial.println(inputString);

  //
  
  inputString = "";
  stringComplete = false;
}

void serialEvent() {
  if (stringComplete) {
    return;
  }
  
  while (Serial.available()) {
    if (inputString.length() == BUFFER_SIZE) {
      Serial.println("ERROR - BUFFER OVERFLOW");
      break;
    }
    
    char ch = (char)Serial.read();
    if (ch == '\n' || ch == '\0') {
      stringComplete = true;
    } else {
      inputString.concat(ch);
    }
  }
}
