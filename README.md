# COVID Simulator Documentation

This documentation provides instructions on how to set up and use the COVID simulator, which simulates patient flow during a COVID outbreak. The simulator uses RabbitMQ for messaging and involves two scripts written in Go: `producer.go` and `consumer.go`. Additionally, a Python script `chart.py` is provided to visualize the output data.

## Prerequisites

Before running the simulator, ensure that the following software is installed on your system:

1. **RabbitMQ**: A message broker that facilitates communication between the producer and consumer scripts.
   - Installation instructions: [RabbitMQ Installation Guide](https://www.rabbitmq.com/download.html)

2. **Go**: A programming language used to write the producer and consumer scripts.
   - Installation instructions: [Go Installation Guide](https://golang.org/doc/install)

3. **Python 3**: Required for running the `chart.py` script and data visualization with the `matplotlib` library.
   - Installation instructions: [Python 3 Installation Guide](https://www.python.org/downloads/)

4. **Matplotlib**: A Python library for creating visualizations.
   - Installation:

     ```bash
     pip install matplotlib
     ```

## Getting Started

### 1. Setting Up RabbitMQ

Ensure RabbitMQ is installed and running on your local machine or server. You can start RabbitMQ by running the following command:

```bash
rabbitmq-server
```

### 2. Running the Consumer Script

The consumer.go script should be run first. This script listens for messages from the producer and processes them.

To run the consumer.go script, execute the following command:

```bash
go run consumer.go
```

### 3. Running the Producer Script

After the consumer script is running, you can start the producer.go script. This script generates events and sends them to RabbitMQ, where they will be consumed and processed by the consumer.go script.

To run the producer.go script, execute the following command:

```bash
go run producer.go
```

### 4. Processing Output

After all events have been processed, the consumer.go script will output the results in JSON format. The output will look something like this:

```json
[
    {
        "Day": 1,
        "NotAttendedPacients": 236,
        "AttendedPacients": 202,
        "AttendedPacientsMorning": 91,
        "AttendedPacientsAfternoon": 111,
        "AttentionMedian": 50,
        "AttentionMedianMorning": 45,
        "AttentionMedianAfternoon": 55,
        "WaitingTimeMedian": 153
    },
    {
        "Day": 2,
        "NotAttendedPacients": 253,
        "AttendedPacients": 208,
        "AttendedPacientsMorning": 98,
        "AttendedPacientsAfternoon": 110,
        "AttentionMedian": 52,
        "AttentionMedianMorning": 49,
        "AttentionMedianAfternoon": 55,
        "WaitingTimeMedian": 164
    }
]
```

### 5. Visualizing the Data

To visualize the output data, you can use the provided chart.py script. Follow these steps:

Open the chart.py script in a text editor.
Replace the content of the data variable with the JSON data generated from the consumer.go script.

```python
data = [
    {
        "Day": 1,
        "NotAttendedPacients": 236,
        "AttendedPacients": 202,
        "AttendedPacientsMorning": 91,
        "AttendedPacientsAfternoon": 111,
        "AttentionMedian": 50,
        "AttentionMedianMorning": 45,
        "AttentionMedianAfternoon": 55,
        "WaitingTimeMedian": 153
    },
    {
        "Day": 2,
        "NotAttendedPacients": 253,
        "AttendedPacients": 208,
        "AttendedPacientsMorning": 98,
        "AttendedPacientsAfternoon": 110,
        "AttentionMedian": 52,
        "AttentionMedianMorning": 49,
        "AttentionMedianAfternoon": 55,
        "WaitingTimeMedian": 164
    }
]
```

Save the changes and run the chart.py script to generate the visualization:

```bash
python chart.py
```

### 6. Output Visualization

Running chart.py will generate a visualization based on the JSON data. This can include graphs or charts that help you analyze the performance metrics like the number of attended/not attended patients, median attention time, and waiting time for each day.
