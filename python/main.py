import datetime
import json
from typing import List
from pika.adapters import BlockingConnection
from pika.adapters.blocking_connection import BlockingChannel
from pika.credentials import PlainCredentials
from pyfirmata2 import Arduino, Pin, time
from enum import Enum
import pika
from pyfirmata2.util import os

queue_population: int = 0
queue_max_capacity: int = 5
attraction_name: str = os.environ['ATTRACTION_NAME']

board: Arduino = Arduino('/dev/ttyACM0')
board.samplingOn(100)

credentials: PlainCredentials = pika.PlainCredentials(username=os.environ['MQ_USERNAME'], password=os.environ['MQ_PASSWORD'])
connection: BlockingConnection = pika.BlockingConnection(pika.ConnectionParameters('rabbitmq', credentials=credentials))
channel: BlockingChannel = connection.channel()
channel.queue_declare(queue='queuepods')

message_queue: List[str] = [] # We now have a queue that we use to send messages to the other queue from an IoT device that itself is used to monitor queues...

class Levels(Enum):
    INFO  = "info"
    WARN  = "warning"
    ERROR = "error"

class Events(Enum):
    PERSON_ENTERED_QUEUE  = "person_entered_queue"
    PERSON_LEFT_QUEUE     = "person_left_queue"
    QUEUE_EMPTY           = "queue_empty"
    QUEUE_ALMOST_EMPTY    = "queue_almost_empty"
    QUEUE_ALMOST_FULL     = "queue_almost_full"
    QUEUE_FULL            = "queue_full"
    POPULATION_MONITORING = "population_monitoring"

previous_event: Events = Events.QUEUE_EMPTY

class Logger():
    def serialize(self, event: Events, level: Levels, body: str) -> str:
        return json.dumps({
            "attraction": attraction_name,
            "event": event.value,
            "level": level.value,
            "timestamp": datetime.datetime.now().isoformat(),
            "value": body,
        })

    def info(self, event: Events, body: str, publish_to_queue: bool = False) -> None:
        message = self.serialize(event=event, level=Levels.INFO, body=body)
        if publish_to_queue and previous_event != event:
            message_queue.append(message)

    def warn(self, event: Events, body: str, publish_to_queue: bool = False) -> None:
        message = self.serialize(event=event, level=Levels.WARN, body=body)
        if publish_to_queue and previous_event != event:
            message_queue.append(message)

    def error(self, event: Events, body: str, publish_to_queue: bool = False) -> None:
        message = self.serialize(event=event, level=Levels.ERROR, body=body)
        if publish_to_queue and previous_event != event:
            message_queue.append(message)

logger = Logger()

class Lights():
    RED: Pin = board.digital[12]
    YELLOW: Pin = board.digital[11]
    GREEN: Pin = board.digital[10]

    @staticmethod
    def reset() -> None:
        for light in [Lights.RED, Lights.YELLOW, Lights.GREEN]:
            light.write(0)

class Buttons(Enum):
    ENTRY = board.get_pin('d:3:u')
    EXIT = board.get_pin('d:2:u')

def update_state() -> None:
    global previous_event
    global queue_population
    if queue_population < 0:
        queue_population = 0

    Lights.reset()
    if queue_population == 0:
        logger.warn(event=Events.QUEUE_EMPTY, body=f"There are no people present within the queue for attraction with name {attraction_name}.", publish_to_queue=True)
        previous_event = Events.QUEUE_EMPTY
        Lights.GREEN.write(1)
    elif queue_population >= queue_max_capacity:
        logger.warn(event=Events.QUEUE_FULL, body=f"The queue for attraction with name {attraction_name} has reached full occupancy.", publish_to_queue=True)
        previous_event = Events.QUEUE_FULL
        Lights.RED.write(1)
    elif queue_population / queue_max_capacity * 100 < 70:
        logger.info(event=Events.QUEUE_ALMOST_EMPTY, body=f"The queue for attraction with name {attraction_name} has an occupancy of less than 70%.", publish_to_queue=True)
        previous_event = Events.QUEUE_ALMOST_EMPTY
        Lights.GREEN.write(1)
    elif queue_population / queue_max_capacity* 100 > 70:
        logger.warn(event=Events.QUEUE_ALMOST_FULL, body=f"The queue for attraction with name {attraction_name} has an occupancy of at least 70%.", publish_to_queue=True)
        previous_event = Events.QUEUE_ALMOST_FULL
        Lights.YELLOW.write(1)

def entry_callback(value: bool) -> None:
    if not value:
        return

    global queue_population
    queue_population += 1

    logger.info(event=Events.PERSON_ENTERED_QUEUE, body=f"Person has entered the queue for attraction with name {attraction_name}")

def exit_callback(value: bool) -> None:
    if not value:
        return

    global queue_population
    queue_population -= 1
    
    logger.info(event=Events.PERSON_LEFT_QUEUE, body=f"Person has left the queue for attraction with name {attraction_name}")

Buttons.ENTRY.value.register_callback(entry_callback)
Buttons.EXIT.value.register_callback(exit_callback)

previous_population: int = 0
update_state()

monitoring_previous_timestamp: float = time.time()
message_queue_previous_timestamp: float = time.time()
while True:
    current_timestamp: float = time.time()
    if current_timestamp - monitoring_previous_timestamp > 3:
        monitoring_previous_timestamp = current_timestamp
        logger.info(event=Events.POPULATION_MONITORING, body=str(queue_population), publish_to_queue=True)

    if current_timestamp - message_queue_previous_timestamp > 0.01:
        message_queue_previous_timestamp = current_timestamp
        if len(message_queue) != 0:
            channel.basic_publish(exchange='', routing_key='queuepods', body=message_queue.pop())

    if previous_population != queue_population:
        update_state()

    previous_population = queue_population
