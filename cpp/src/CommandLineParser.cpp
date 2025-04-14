#include "CommandLineParser.h"
#include <iostream>
#include <cstring>

CommandLineOptions CommandLineParser::parse(int argc, char* argv[]) {
    CommandLineOptions options = {
        "localhost",  // default host
        8080,         // default port
        60,           // default rate (messages per minute)
        1             // default worker count
    };
    
    for (int i = 1; i < argc; ++i) {
        std::string arg = argv[i];
        
        if (arg.find("--host=") == 0) {
            options.host = arg.substr(7);
        }
        else if (arg.find("--port=") == 0) {
            options.port = std::stoi(arg.substr(7));
        }
        else if (arg.find("--rate=") == 0) {
            options.rate = std::stoi(arg.substr(7));
        }
        else if (arg.find("--worker=") == 0) {
            options.workers = std::stoi(arg.substr(9));
        }
        else if (arg == "--help" || arg == "-h") {
            printHelp();
            exit(0);
        }
    }
    
    // Print the configuration
    std::cout << "Server configuration:" << std::endl;
    std::cout << "Host: " << options.host << std::endl;
    std::cout << "Port: " << options.port << std::endl;
    std::cout << "Rate: " << options.rate << " messages per minute" << std::endl;
    std::cout << "Workers: " << options.workers << std::endl;
    
    return options;
}

void CommandLineParser::printHelp() {
    std::cout << "Niagara - WebSocket server for stock trading data simulation" << std::endl;
    std::cout << "Usage: niagara [options]" << std::endl;
    std::cout << std::endl;
    std::cout << "Options:" << std::endl;
    std::cout << "  --host=<string>   Host address (default: localhost)" << std::endl;
    std::cout << "  --port=<number>   Port number (default: 8080)" << std::endl;
    std::cout << "  --rate=<number>   Data messages per minute (default: 60)" << std::endl;
    std::cout << "  --worker=<number> Number of worker threads (default: 1)" << std::endl;
    std::cout << "  --help, -h        Show this help message" << std::endl;
    std::cout << std::endl;
    std::cout << "Example:" << std::endl;
    std::cout << "  niagara --host=localhost --port=8080 --rate=10 --worker=2" << std::endl;
}
