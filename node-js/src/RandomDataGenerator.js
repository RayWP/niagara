class RandomDataGenerator {
    constructor() {
        this.stockCode = ['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'TSLA'];
        this.stockPrice = [100, 200, 300, 400, 500];
        this.buyOrSell = ['B', 'S'];
    }

    getData() {
        const randomStockCode = this.stockCode[Math.floor(Math.random() * this.stockCode.length)];
        const randomStockPrice = this.stockPrice[Math.floor(Math.random() * this.stockPrice.length)];
        const randomBuyOrSell = this.buyOrSell[Math.floor(Math.random() * this.buyOrSell.length)];
        const randomQuantity = Math.floor(Math.random() * 100) + 1;
        return `${randomStockCode}|${randomStockPrice}|${randomBuyOrSell}|${randomQuantity}`;
    }
}