cmake_minimum_required(VERSION 3.0)

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY "./bin")
project("test")
# add_executable(test ./test/test.cpp)
# target_link_libraries(test gtest gmock pthread)
add_executable(solution ./solution/g.cpp)