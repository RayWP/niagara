<?php
use Workerman\Connection\TcpConnection;
class DataHub
{
    private array $connections = [];
    private int $rate;
    public function __construct(int $rate)
    {
        echo "Create Object\n";
        $this->rate = $rate;
        $this->generateAndSendData();
    }

    public function addConnection(TcpConnection $connection)
    {
        echo "New Connection added, id: $connection->id\n";
        $this->connections[] = $connection;
    }

    public function removeConnection(TcpConnection $connection)
    {
        $index = array_search($connection, $this->connections);
        if ($index !== false) {
            echo "Connection removed, id: $connection->id\n";
            unset($this->connections[$index]);
        }
    }

    public function generateAndSendData() {
        $timeout_between_data = 60 / $this->rate;
        echo "generate function invoked\n";
        echo "Timeout between data: $timeout_between_data\n";
        \Workerman\Lib\Timer::add($timeout_between_data, function()
        {
            echo "Get into generation\n";
            if(count($this->connections) > 0) {
                echo "Generating and sending data\n for " . count($this->connections) . " connections\n";
                $time_in_millis = round(microtime(true) * 1000);
                foreach ($this->connections as $connection) {
                    echo "Sending data to connection $connection->id on worker $connection->worker\n";
                    $connection->send($time_in_millis);
                }
            }
        });
    }
}