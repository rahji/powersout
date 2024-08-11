# powersout

I have an Arduino with a datalogger shield from the olden days.

I plug it into an outlet when I go on vacation.

It records a timestamp in a file called `DATALOG.TXT` on its SD card, about every 5 minutes.

It records the string `REBOOT` every time it is powered on.

This program analyzes the `DATALOG.TXT` file and outputs the approximate durations of any power outages.

![Made with VHS](https://vhs.charm.sh/vhs-7ySglDHtnlT17aZmzHRZ9p.gif)

## Arduino Code

I don't have the simple Arduino source code that I wrote anymore, but it probably looks something like this:

```cpp
#include <SPI.h>
#include <SD.h>

const int chipSelect = SDCARD_SS_PIN;
unsigned long secondsDelay = 300000L

void writeString(String s) {
  File dataFile = SD.open("DATALOG.TXT", FILE_WRITE);

  if (dataFile) {
    dataFile.println(s);
    dataFile.close();
  } else {
    Serial.println("error opening datalog.txt");
  }
}

void setup() {
  Serial.begin(9600);
  while (!Serial) {
  }

  Serial.print("Initializing SD card...");

  if (!SD.begin(chipSelect)) {
    Serial.println("Card failed, or not present");
    while (1);
  }
  Serial.println("card initialized.");

  writeString("REBOOT");
}

void loop() {
  writeString( String(millis()/1000) )
  delay(secondsDelay);
}
```