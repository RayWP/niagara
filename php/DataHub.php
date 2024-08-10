<?php
use Workerman\Connection\TcpConnection;
require_once __DIR__ . '/RandomDataGenerator.php';
class DataHub
{
    private array $connections = [];
    private int $rate;
    private RandomDataGenerator $randomDataGenerator;
    public function __construct(int $rate)
    {
        $this->rate = $rate;
        $this->randomDataGenerator = new RandomDataGenerator();
    }

    public function addConnection(TcpConnection $connection)
    {
        echo "New Connection added, id: $connection->id\n";
        $this->connections[] = $connection;
        if(count($this->connections) == 1) {
            // start data when at least one connection is connected
            $this->generateAndSendData();
        }
    }

    public function removeConnection(TcpConnection $connection)
    {
        $index = array_search($connection, $this->connections);
        if ($index !== false) {
            echo "Connection removed, id: $connection->id\n";
            unset($this->connections[$index]);
        }
        if(count($this->connections) == 0) {
            echo "No more connections, stop generating data\n";
            \Workerman\Lib\Timer::delAll();
        }
    }

    public function generateAndSendData() {
        $timeout_between_data = 60 / $this->rate;
        echo "Start data generation";
        \Workerman\Lib\Timer::add($timeout_between_data, function()
        {
            $data = $this->randomDataGenerator->getData();
            echo "Data: $data\n";
            foreach ($this->connections as $connection) {
                $connection->send($data);
            }
        });
    }
}