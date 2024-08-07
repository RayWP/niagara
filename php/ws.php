<?php
require_once __DIR__ . '/vendor/autoload.php';

use Workerman\Connection\TcpConnection;
use Workerman\Protocols\Http\Request;
use Workerman\Worker;

$longopts = [
    'rate::', // Required value -rate value
    'host::', // Optional value --host=value
    'port::', // Optional value --port=value
    'worker::', // Optional value --worker=value
];

$options = getopt('', $longopts);
echo json_encode($options) . PHP_EOL;
$rate = $options['rate'] ?? 1;
$host = $options['host'] ?? 'localhost';
$port = $options['port'] ?? 8000;
$worker = $options['worker'] ?? 1;

echo "rate: $rate, host: $host, port: $port\n";

$worker = new Worker('websocket://' . $host . ':' . $port);
$worker->count = 1;
$worker->onWebSocketConnect = function (TcpConnection $connection) use ($rate) {
//        // You can verify the legitimacy of the connection here and close it if it is not valid
//        // $_SERVER['HTTP_ORIGIN'] indicates the site from which the page initiated the websocket connection
//        if($_SERVER['HTTP_ORIGIN'] != 'https://www.workerman.net')
//        {
//            $connection->close();
//        }

    // debug the connection object
    var_dump($connection);

        $timeout_between_data = 60 / $rate;
        $connection->timeout_id = \Workerman\Lib\Timer::add($timeout_between_data, function() use ($connection)
        {
            $connection->send('hello');
        });
};
Worker::runAll();