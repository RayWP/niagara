#pragma once

#include <string>

struct CommandLineOptions {
    std::string host;
    int port;
    int rate;
    int workers;
};

class CommandLineParser {
public:
    static CommandLineOptions parse(int argc, char* argv[]);
    static void printHelp();
};
