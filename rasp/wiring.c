#include <wiringPi.h>

// sudo apt-get install wiringpi
// gcc -o wiring wiring.c -lwiringPi

int main(void)
{

    wiringPiSetup();

    pinMode(25, OUTPUT);

    for (;;)
    {

        digitalWrite(25, HIGH);
        delay(1000);

        digitalWrite(25, LOW);
        delay(1000);
    }
}
