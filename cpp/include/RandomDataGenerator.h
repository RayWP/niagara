#pragma once

#include <vector>
#include <string>
#include <random>

class RandomDataGenerator {
private:
    std::vector<std::string> stockCode;
    std::vector<int> stockPrice;
    std::vector<std::string> buyOrSell;
    std::mt19937 rng;
    
public:
    RandomDataGenerator();
    std::string getData();
};
