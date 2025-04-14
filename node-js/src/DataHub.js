class DataHub {
    constructor(rate) {
        this.connections = [];
        this.rate = rate;
        this.randomDataGenerator = new RandomDataGenerator();
    }

    addConnection(connection) {
        console.log(`New Connection added, id: ${connection.id}`);
        this.connections.push(connection);
        if (this.connections.length === 1) {
            // start data when at least one connection is connected
            this.broadcastData();
        }
    }

    removeConnection(connection) {
        const index = this.connections.indexOf(connection);
        if (index !== -1) {
            console.log(`Connection removed, id: ${connection.id}`);
            this.connections.splice(index, 1);
        }
        if (this.connections.length === 0) {
            console.log("No more connections, stop generating data");
            clearInterval(this.timer);
        }
    }

    broadcastData() {
        const timeoutBetweenData = 60 / this.rate * 1000;
        console.log("Start data generation");
        this.timer = setInterval(() => {
            const data = this.randomDataGenerator.getData();
            console.log(`Data: ${data}`);
            this.connections.forEach(connection => {
                connection.send(data);
            });
        }, timeoutBetweenData);
    }
}
