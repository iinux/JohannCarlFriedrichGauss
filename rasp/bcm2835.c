#include <bcm2835.h>

// #define PIN RPI_GPIO_P1_12
#define PIN RPI_V2_GPIO_P1_37

// http://www.airspayce.com/mikem/bcm2835/
// gcc bcm2835.c -o bcm2835  -l bcm2835

int main(int argc, char **argv)
{

    if (!bcm2835_init())

        return 1;

    bcm2835_gpio_fsel(PIN, BCM2835_GPIO_FSEL_OUTP);

    while (1)
    {

        bcm2835_gpio_write(PIN, HIGH);

        bcm2835_delay(1000);

        bcm2835_gpio_write(PIN, LOW);

        bcm2835_delay(1000);
    }

    bcm2835_close();

    return 0;
}
