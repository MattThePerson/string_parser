# String Parser





## Type verbs


| verb | description |
|---|---|
| %v | default format for the value |
| %#v | Go-syntax representation of the value |
| %T | type of the value |
| %% | literal percent sign |
|    %t | true or false |
|| **Integer:** |
|    %b | base 2 (binary) |
|    %c | the character represented by the corresponding Unicode code point |
|    %d | base 10 (decimal) |
|    %o | base 8 (octal) |
|    %O | base 8 with 0o prefix |
|    %q | single-quoted character literal (Go syntax) |
|    %x | base 16 (hex, lowercase) |
|    %X | base 16 (hex, uppercase) |
|    %U | Unicode format: U+1234 (upper case hex) |
|| **Floating-point and complex:** |
|    %b | decimal with exponent of two, power of two notation (e.g., 123456p-78) |
|    %e | scientific notation (e.g., 1.234456e+78) |
|    %E | scientific notation (e.g., 1.234456E+78) |
|    %f | decimal point but no exponent (e.g., 123.456) |
|    %F | synonym for %f |
|    %g | %e for large exponents, %f otherwise |
|    %G | %E for large exponents, %F otherwise |
|    %x | hexadecimal notation with fractional part, powers of two exponent |
|| **String and slice of bytes:** |
|    %s | uninterpreted bytes of the string or slice |
|    %q | double-quoted string (Go syntax) |
|    %x | base 16, lower-case, two characters per byte |
|    %X | base 16, upper-case, two characters per byte |
    

