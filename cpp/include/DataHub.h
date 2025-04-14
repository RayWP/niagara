    #pragma once

#include <vector>
#include <memory>
#include <mutex>
#include "RandomDataGenerator.h"

#include <websocketpp/config/asio_no_tls.hpp>
#include <websocketpp/server.hpp>

class DataHub {
public:
    using WebsocketServer = websocketpp::server<websocketpp::config::asio>;
    using ConnectionHandle = websocketpp::connection_hdl;
    
private:
    std::vector<ConnectionHandle> connections;
    int rate;
    RandomDataGenerator randomDataGenerator;
    std::mutex connectionsMutex;
    bool isRunning;
    std::thread broadcastThread;
    WebsocketServer* serverPtr;
    
public:
    DataHub(int rate, WebsocketServer* server);
    ~DataHub();
    
    void addConnection(ConnectionHandle hdl);
    void removeConnection(ConnectionHandle hdl);
    void broadcastData();
    void startBroadcasting();
    void stopBroadcasting();
};
