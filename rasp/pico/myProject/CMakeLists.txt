cmake_minimum_required(VERSION 3.12)
include(pico_sdk_import.cmake)
project(test_project)
pico_sdk_init()
add_executable(myProject    
blink.c    
)
pico_add_extra_outputs(myProject)
target_link_libraries(myProject pico_stdlib)
