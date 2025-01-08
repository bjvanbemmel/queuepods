FROM python:3.13.1-alpine3.21
WORKDIR /app
RUN pip install pyfirmata2 pika pytz
COPY ./* ./
CMD [ "python", "main.py" ]
