#include <stdio.h>
#include <stdlib.h>
#include "crcLib.h"



int main1()
{
    uint8_t data[] = {0xFE,0xE4,0xEF
    //,0xA9,0x94,0xF1
    //,0x00,0x05
    //,0x00
    };

    for (int i = 0; i < sizeof(data); i++) {
        printf("%02x", data[i]);
    }

    printf("\n\n");

    printf("crc4_itu:%02x\n", crc4_itu(data, sizeof(data)));
    printf("crc5_epc:%02x\n", crc5_epc(data, sizeof(data)));
    printf("crc5_itu:%02x\n", crc5_itu(data, sizeof(data)));
    printf("crc5_usb:%02x\n", crc5_usb(data, sizeof(data)));
    printf("crc6_itu:%02x\n", crc6_itu(data, sizeof(data)));
    printf("crc7_mmc:%02x\n", crc7_mmc(data, sizeof(data)));
    printf("crc8:%02x\n", crc8(data, sizeof(data)));
    printf("crc8_itu:%02x\n", crc8_itu(data, sizeof(data)));
    printf("crc8_rohc:%02x\n", crc8_rohc(data, sizeof(data)));
    printf("crc8_maxim:%02x\n", crc8_maxim(data, sizeof(data)));
    printf("crc16_ibm:%02x\n", crc16_ibm(data, sizeof(data)));
    printf("crc16_maxim:%02x\n", crc16_maxim(data, sizeof(data)));
    printf("crc16_usb:%02x\n", crc16_usb(data, sizeof(data)));
    printf("crc16_modbus:%02x\n", crc16_modbus(data, sizeof(data)));
    printf("crc16_ccitt:%02x\n", crc16_ccitt(data, sizeof(data)));
    printf("crc16_ccitt_false:%02x\n", crc16_ccitt_false(data, sizeof(data)));
    printf("crc16_x25:%02x\n", crc16_x25(data, sizeof(data)));
    printf("crc16_xmodem:%02x\n", crc16_xmodem(data, sizeof(data)));
    printf("crc16_dnp:%02x\n", crc16_dnp(data, sizeof(data)));
    printf("crc32:%02x\n", crc32(data, sizeof(data)));
    printf("crc32_mpeg_2:%02x\n", crc32_mpeg_2(data, sizeof(data)));
    
    printf("\n\n");
    
    printf("crc4_itu:%02x\n", ~crc4_itu(data, sizeof(data)));
    printf("crc5_epc:%02x\n", ~crc5_epc(data, sizeof(data)));
    printf("crc5_itu:%02x\n", ~crc5_itu(data, sizeof(data)));
    printf("crc5_usb:%02x\n", ~crc5_usb(data, sizeof(data)));
    printf("crc6_itu:%02x\n", ~crc6_itu(data, sizeof(data)));
    printf("crc7_mmc:%02x\n", ~crc7_mmc(data, sizeof(data)));
    printf("crc8:%02x\n", ~crc8(data, sizeof(data)));
    printf("crc8_itu:%02x\n", ~crc8_itu(data, sizeof(data)));
    printf("crc8_rohc:%02x\n", ~crc8_rohc(data, sizeof(data)));
    printf("crc8_maxim:%02x\n", ~crc8_maxim(data, sizeof(data)));
    printf("crc16_ibm:%02x\n", ~crc16_ibm(data, sizeof(data)));
    printf("crc16_maxim:%02x\n", ~crc16_maxim(data, sizeof(data)));
    printf("crc16_usb:%02x\n", ~crc16_usb(data, sizeof(data)));
    printf("crc16_modbus:%02x\n", ~crc16_modbus(data, sizeof(data)));
    printf("crc16_ccitt:%02x\n", ~crc16_ccitt(data, sizeof(data)));
    printf("crc16_ccitt_false:%02x\n", ~crc16_ccitt_false(data, sizeof(data)));
    printf("crc16_x25:%02x\n", ~crc16_x25(data, sizeof(data)));
    printf("crc16_xmodem:%02x\n", ~crc16_xmodem(data, sizeof(data)));
    printf("crc16_dnp:%02x\n", ~crc16_dnp(data, sizeof(data)));
    printf("crc32:%02x\n", ~crc32(data, sizeof(data)));
    printf("crc32_mpeg_2:%02x\n", ~crc32_mpeg_2(data, sizeof(data)));
    return 0;
}

// 98A9F9A994F100050017D8DB0045EDBC
// FEE4EFA994F100050094780000000000
// FEE47CA994F1000500E82E0000000000

int main()
{
    uint8_t data[] = {0xFE,0xE4,0xEF
    //,0xA9,0x94,0xF1
    //,0x00,0x05
    //,0x00
    };

    int end = 0;
    for (uint16_t init_value = 0x0005; init_value <= 0xffff; init_value++) {
        for (uint16_t poly = 0x0000; poly <= 0xffff; poly++) {
            uint16_t ret = crc16_free(data, sizeof(data), init_value, poly);
            printf("%02x %02x %02x\n", init_value, poly, ret);
            if (ret == 0x9478) {
                printf("success:%02x %02x", init_value, poly);
                end = 1;
                break;
            }

            if (poly == 0xffff) {
                break;
            }
        }

        if (end == 1) {
            break;
        }

        if (init_value == 0xffff) {
            break;
        }
    }
    return 0;
}