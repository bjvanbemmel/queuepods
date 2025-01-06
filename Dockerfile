FROM python:3.13.1-alpine3.21
WORKDIR /app
RUN pip install pyfirmata2
COPY ./main.py ./main.py
CMD [ "python", "main.py" ]
