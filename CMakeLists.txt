# Copyright (c) 2024, Oh Inseo (YJK) <dhdlstjtr@gmail.com>
#
# SPDX-License-Identifier: BSD-2-Clause
cmake_minimum_required(VERSION 3.5)
project(yw)

################################################################################
# CMake options
################################################################################
option(WERROR "Enable -Werror"              OFF)
option(UBSAN  "Enable -fsanitize=undefined" OFF)

################################################################################
# Compiler flags (CXXFLAGS)
################################################################################

# Optimization and std flags
set(CMAKE_CXX_FLAGS
    "${CMAKE_CXX_FLAGS} -flto -O3 -std=c++17 -Wall -Wextra")

if (WERROR)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Werror")
endif()

if (UBSAN)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fsanitize=undefined")
endif()

################################################################################
# Executable information
################################################################################

add_executable(
    yw

    src/main.cc

    src/dom/document.cc
    src/dom/documentfragment.cc
    src/dom/element.cc
    src/dom/node.cc
    src/dom/nodeiterator.cc
    src/dom/range.cc
    src/dom/shadowroot.cc
    src/dom/slottable.cc
    src/dom/_debug/node.cc

    src/idl/domexception.cc

    src/_test/testlib.cc
    src/_test/testrun.cc
    src/_test/dom/dom.cc

    src/_utility/logging.cc
)
