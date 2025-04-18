<?php
require_once __DIR__ . '/vendor/autoload.php';
require_once __DIR__ . '/DataHub.php';
require_once __DIR__ . '/RandomDataGenerator.php';

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
$rate = $options['rate'] ?? 60;
$host = $options['host'] ?? 'localhost';
$port = $options['port'] ?? 8000;
$worker_count = $options['worker'] ?? 1;

$dataRunner = new DataHub($rate);

$worker = new Worker('websocket://' . $host . ':' . $port);
$worker->count = $worker_count;
$worker->onWebSocketConnect = function (TcpConnection $connection) use ($dataRunner) {
    // bind the connection
    $dataRunner->addConnection($connection);
};

$worker->onClose = function (TcpConnection $connection) use ($dataRunner) {
    echo "connection closed on id $connection->id\n";
    $dataRunner->removeConnection($connection);
};

Worker::runAll();

