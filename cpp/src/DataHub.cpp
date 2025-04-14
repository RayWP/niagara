#include "DataHub.h"
#include <iostream>
#include <chrono>
#include <thread>

DataHub::DataHub(int rate, WebsocketServer* server) 
    : rate(rate), 
      randomDataGenerator(),
      isRunning(false),
      serverPtr(server) {
}

DataHub::~DataHub() {
    stopBroadcasting();
}

void DataHub::addConnection(ConnectionHandle hdl) {
    std::lock_guard<std::mutex> lock(connectionsMutex);
    std::cout << "New connection added" << std::endl;
    connections.push_back(hdl);
    
    if (connections.size() == 1) {
        // Start broadcasting when first connection is made
        startBroadcasting();
    }
}

void DataHub::removeConnection(ConnectionHandle hdl) {
    std::lock_guard<std::mutex> lock(connectionsMutex);
    
    // Find and remove the connection
    connections.erase(
        std::remove_if(connections.begin(), connections.end(),
            [&](const ConnectionHandle& conn) {
                return !serverPtr->get_con_from_hdl(conn)->get_raw_socket().is_open() ||
                       serverPtr->get_con_from_hdl(hdl).get() == serverPtr->get_con_from_hdl(conn).get();
            }),
        connections.end());
    
    std::cout << "Connection removed, remaining: " << connections.size() << std::endl;
    
    if (connections.empty()) {
        stopBroadcasting();
    }
}

void DataHub::broadcastData() {
    while (isRunning) {
        std::string data = randomDataGenerator.getData();
        std::cout << "Data: " << data << std::endl;
        
        // Calculate sleep time between messages based on rate
        int sleepTimeMs = static_cast<int>(60.0 / rate * 1000);
        
        {
            std::lock_guard<std::mutex> lock(connectionsMutex);
            // Remove closed connections
            connections.erase(
                std::remove_if(connections.begin(), connections.end(),
                    [this](const ConnectionHandle& hdl) {
                        return !serverPtr->get_con_from_hdl(hdl)->get_raw_socket().is_open();
                    }),
                connections.end());
                
            // Send data to all active connections
            for (auto& conn : connections) {
                if (serverPtr->get_con_from_hdl(conn)->get_raw_socket().is_open()) {
                    try {
                        serverPtr->send(conn, data, websocketpp::frame::opcode::text);
                    } catch (const std::exception& e) {
                        std::cerr << "Error sending message: " << e.what() << std::endl;
                    }
                }
            }
        }
        
        // Sleep until next message
        std::this_thread::sleep_for(std::chrono::milliseconds(sleepTimeMs));
    }
}

void DataHub::startBroadcasting() {
    if (!isRunning) {
        isRunning = true;
        std::cout << "Start data generation" << std::endl;
        broadcastThread = std::thread(&DataHub::broadcastData, this);
    }
}

void DataHub::stopBroadcasting() {
    if (isRunning) {
        std::cout << "No more connections, stop generating data" << std::endl;
        isRunning = false;
        if (broadcastThread.joinable()) {
            broadcastThread.join();
        }
    }
}
