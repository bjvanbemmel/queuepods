from pyfirmata2 import Arduino, Pin
from enum import Enum

queue_population: int = 0

board: Arduino = Arduino('/dev/ttyACM0')
board.samplingOn(100)

class Lights():
    RED: Pin = board.digital[12]
    YELLOW: Pin = board.digital[11]
    GREEN: Pin = board.digital[10]

    @staticmethod
    def reset() -> None:
        for light in [Lights.RED, Lights.YELLOW, Lights.GREEN]:
            light.write(0)

    @staticmethod
    def update() -> None:
        Lights.reset()

        global queue_population
        if queue_population < 0:
            queue_population = 0

        match queue_population:
            case 0:
                Lights.GREEN.write(1)
            case 1:
                Lights.YELLOW.write(1)
            case _:
                Lights.RED.write(1)

class Buttons(Enum):
    ENTRY = board.get_pin('d:3:u')
    EXIT = board.get_pin('d:2:u')

def entry_callback(value: bool) -> None:
    if not value:
        return

    global queue_population
    queue_population += 1

def exit_callback(value: bool) -> None:
    if not value:
        return

    global queue_population
    queue_population -= 1

Buttons.ENTRY.value.register_callback(entry_callback)
Buttons.EXIT.value.register_callback(exit_callback)

previous_population: int = 0
while True:
    if previous_population != queue_population:
        print(f"New queue population: {queue_population}")

    previous_population = queue_population

    Lights.update()
