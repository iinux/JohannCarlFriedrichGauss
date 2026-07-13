/* Prints four Unicode code points as UTF-8.
 *   U+0027  APOSTROPHE
 *   U+2019  RIGHT SINGLE QUOTATION MARK
 *   U+02BC  MODIFIER LETTER APOSTROPHE
 *   U+0289  MODIFIER LETTER SMALL CAPITAL I
 *
 * Source is UTF-8; u8"..." string literals with \uXXXX escapes
 * produce well-defined UTF-8 byte sequences (C11 §6.4.5).
 * stdout must be a UTF-8 locale (true by default on modern Linux).
 */
#include <stdio.h>

int main(void) {
    static const char *const names[] = {
        "APOSTROPHE",
        "RIGHT SINGLE QUOTATION MARK",
        "MODIFIER LETTER APOSTROPHE",
        "MODIFIER LETTER SMALL CAPITAL I",
    };
    static const char *const glyphs[] = {
        u8"'",
        u8"’",
        u8"ʼ",
        u8"ʉ",
    };
    static const unsigned cps[] = { 0x0027, 0x2019, 0x02BC, 0x0289 };

    for (int i = 0; i < 4; i++)
        printf("U+%04X  %-35s [%s]\n", cps[i], names[i], glyphs[i]);

    return 0;
}
