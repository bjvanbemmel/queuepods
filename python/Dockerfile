FROM python:3.13.1-alpine3.21
WORKDIR /app
RUN pip install pyfirmata2 pika
COPY ./* ./
CMD [ "python", "main.py" ]
