<?php

class RandomDataGenerator
{
    private $stockCode = ['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'TSLA'];
    private $stockPrice = [100, 200, 300, 400, 500];
    private $buyOrSell = ['B', 'S'];
    public function getData() {
        $randomStockCode = $this->stockCode[rand(0, 4)];
        $randomStockPrice = $this->stockPrice[rand(0, 4)];
        $randomBuyOrSell = $this->buyOrSell[rand(0, 1)];
        $randomQuantity = rand(1, 100);
        return "$randomStockCode|$randomStockPrice|$randomBuyOrSell|$randomQuantity";
    }
}