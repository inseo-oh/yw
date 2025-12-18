/*
 * This file is part of YW project. Copyright 2025 Oh Inseo (YJK)
 * SPDX-License-Identifier: BSD-3-Clause
 * See LICENSE for details, and LICENSE_WHATWG_SPECS for WHATWG license
 * information.
 */
#include "yw_encoding.h"
#include <stdint.h>

typedef struct YW_Context YW_Context;
struct YW_Context
{
    uint32_t codepoint;
    int bytesSeen;
    int bytesNeeded;
    uint8_t lowerBoundary;
    uint8_t upperBoundary;
};
