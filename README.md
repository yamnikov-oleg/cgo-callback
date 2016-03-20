# cgo-callback
[![Build Status](https://travis-ci.org/yamnikov-oleg/cgo-callback.svg?branch=master)](https://travis-ci.org/yamnikov-oleg/cgo-callback)

That's gonna be a golang package for dynamic Cgo callbacks.

WIP

**Current limitations:**
+ Supported calling conventions:
  - [x] **System V x64** (Unix 64-bit)
  - [ ] **Microsoft x64** (Windows 64-bit)
  - [x] **cdecl** (Practically any OS 32-bit)
  - [x] **stdcall** (Windows callbacks on 32-bit)
+ Only integers, pointers and floats as arguments and return values.
