/*!
    \file  main.c
    \brief running led
    
    \version 2019-6-5, V1.0.0, firmware for GD32VF103
*/

/*
    Copyright (c) 2019, GigaDevice Semiconductor Inc.

    Redistribution and use in source and binary forms, with or without modification, 
are permitted provided that the following conditions are met:

    1. Redistributions of source code must retain the above copyright notice, this 
       list of conditions and the following disclaimer.
    2. Redistributions in binary form must reproduce the above copyright notice, 
       this list of conditions and the following disclaimer in the documentation 
       and/or other materials provided with the distribution.
    3. Neither the name of the copyright holder nor the names of its contributors 
       may be used to endorse or promote products derived from this software without 
       specific prior written permission.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" 
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED 
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. 
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, 
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT 
NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR 
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, 
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) 
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY 
OF SUCH DAMAGE.
*/

#include "gd32vf103.h"
#include "systick.h"
#include <stdio.h>

/* BUILTIN LED OF LONGAN BOARDS IS PIN PC13 */
#define LED_R_PIN GPIO_PIN_13
#define LED_R_GPIO_PORT GPIOC

#define LED_G_PIN GPIO_PIN_1
#define LED_G_GPIO_PORT GPIOA

#define LED_B_PIN GPIO_PIN_2
#define LED_B_GPIO_PORT GPIOA

#define LED_C_GPIO_CLK RCU_GPIOC
#define LED_A_GPIO_CLK RCU_GPIOA

void longan_led_init()
{
    /* enable the led clock */
    rcu_periph_clock_enable(LED_C_GPIO_CLK);
    rcu_periph_clock_enable(LED_A_GPIO_CLK);
    /* configure led GPIO port */ 
    gpio_init(LED_R_GPIO_PORT, GPIO_MODE_OUT_PP, GPIO_OSPEED_50MHZ, LED_R_PIN);
    gpio_init(LED_G_GPIO_PORT, GPIO_MODE_OUT_PP, GPIO_OSPEED_50MHZ, LED_G_PIN);
    gpio_init(LED_B_GPIO_PORT, GPIO_MODE_OUT_PP, GPIO_OSPEED_50MHZ, LED_B_PIN);

    // GPIO_BC(LED_GPIO_PORT) = LED_PIN;
}

void longan_led_r_on()
{
    /*
     * LED is hardwired with 3.3V on the anode, we control the cathode
     * (negative side) so we need to use reversed logic: bit clear is on.
     */
    GPIO_BC(LED_R_GPIO_PORT) = LED_R_PIN;
}

void longan_led_g_on()
{
    GPIO_BC(LED_G_GPIO_PORT) = LED_G_PIN;
}

void longan_led_b_on()
{
    GPIO_BC(LED_B_GPIO_PORT) = LED_B_PIN;
}

void longan_led_r_off()
{
    GPIO_BOP(LED_R_GPIO_PORT) = LED_R_PIN;
}

void longan_led_g_off()
{
    GPIO_BOP(LED_G_GPIO_PORT) = LED_G_PIN;
}

void longan_led_b_off()
{
    GPIO_BOP(LED_B_GPIO_PORT) = LED_B_PIN;
}

/*!
    \brief      main function
    \param[in]  none
    \param[out] none
    \retval     none
*/
int main(void)
{
    longan_led_init();

    while(1){
        /* turn on builtin led */
        longan_led_r_on();
        delay_1ms(1000);
        /* turn off uiltin led */
        longan_led_r_off();
        delay_1ms(1000);
        longan_led_g_on();
        delay_1ms(1000);
        longan_led_g_off();
        delay_1ms(1000);
        longan_led_b_on();
        delay_1ms(1000);
        longan_led_b_off();
        delay_1ms(1000);
    }
}
