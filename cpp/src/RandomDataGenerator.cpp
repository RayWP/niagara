#include "RandomDataGenerator.h"
#include <random>
#include <sstream>

RandomDataGenerator::RandomDataGenerator() 
    : stockCode{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA"},
      stockPrice{100, 200, 300, 400, 500},
      buyOrSell{"B", "S"} {
    
    // Initialize random number generator
    std::random_device rd;
    rng = std::mt19937(rd());
}

std::string RandomDataGenerator::getData() {
    // Generate random indices
    std::uniform_int_distribution<int> stockCodeDist(0, stockCode.size() - 1);
    std::uniform_int_distribution<int> stockPriceDist(0, stockPrice.size() - 1);
    std::uniform_int_distribution<int> buyOrSellDist(0, buyOrSell.size() - 1);
    std::uniform_int_distribution<int> quantityDist(1, 100);
    
    // Get random values
    std::string randomStockCode = stockCode[stockCodeDist(rng)];
    int randomStockPrice = stockPrice[stockPriceDist(rng)];
    std::string randomBuyOrSell = buyOrSell[buyOrSellDist(rng)];
    int randomQuantity = quantityDist(rng);
    
    // Format and return the string
    std::ostringstream oss;
    oss << randomStockCode << "|" << randomStockPrice << "|" << randomBuyOrSell << "|" << randomQuantity;
    
    return oss.str();
}
