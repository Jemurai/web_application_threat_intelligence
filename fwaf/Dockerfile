FROM python:latest

WORKDIR /code

RUN git clone https://github.com/faizann24/Fwaf-Machine-Learning-driven-Web-Application-Firewall fwaf

WORKDIR /code/fwaf

COPY . .

RUN pip install --trusted-host pypi.python.org -r requirements.txt

CMD ["python", "script.py"]