#include <iostream>
#include <string>
#include <memory>
#include <websocketpp/config/asio_no_tls.hpp>
#include <websocketpp/server.hpp>

#include "CommandLineParser.h"
#include "DataHub.h"

using WebsocketServer = websocketpp::server<websocketpp::config::asio>;

int main(int argc, char* argv[]) {
    try {
        // Parse command line arguments
        CommandLineOptions options = CommandLineParser::parse(argc, argv);
        
        // Create and configure WebSocket server
        WebsocketServer server;
        
        // Set logging settings
        server.set_access_channels(websocketpp::log::alevel::all);
        server.clear_access_channels(websocketpp::log::alevel::frame_payload);
        
        // Initialize ASIO
        server.init_asio();
        
        // Create data hub
        DataHub dataHub(options.rate, &server);
        
        // Set WebSocket callbacks
        server.set_open_handler([&dataHub](websocketpp::connection_hdl hdl) {
            dataHub.addConnection(hdl);
        });
        
        server.set_close_handler([&dataHub](websocketpp::connection_hdl hdl) {
            dataHub.removeConnection(hdl);
        });
        
        // Set up endpoint
        std::string uri = options.host + ":" + std::to_string(options.port);
        server.listen(options.port);
        
        // Start the server accept loop
        server.start_accept();
        
        std::cout << "WebSocket server started at ws://" << uri << std::endl;
        std::cout << "Press Ctrl+C to quit" << std::endl;
        
        // Start the ASIO io_service run loop
        server.run();
    } 
    catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
    
    return 0;
}
