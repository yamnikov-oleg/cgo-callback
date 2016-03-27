// Offset of each register in the register map,
// as it's saved by entry_386.s

#pragma once

#define EAX 0x00
#define EBX 0x04
#define ECX 0x08
#define EDX 0x0C
#define ESI 0x10
#define EDI 0x14
// Must be set if ST0 was modified and should be loaded
#define ST0_SET_FLAG 0x18
#define ST0 0x19
