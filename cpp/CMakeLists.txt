cmake_minimum_required(VERSION 3.10)
project(niagara CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include)

# Find required packages
find_package(Threads REQUIRED)

# Add websocketpp as header-only library
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/external/websocketpp)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/external/asio/asio/include)

# Add source files
set(SOURCES
    src/main.cpp
    src/RandomDataGenerator.cpp
    src/DataHub.cpp
    src/CommandLineParser.cpp
)

# Create executable
add_executable(niagara ${SOURCES})

# Link libraries
target_link_libraries(niagara Threads::Threads)

# Add install rules
install(TARGETS niagara DESTINATION bin)
