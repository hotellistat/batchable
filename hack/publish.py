import pika
import uuid
import json
import time


connection = pika.BlockingConnection(
    pika.ConnectionParameters(
        'localhost',
        5672,
        # credentials=pika.PlainCredentials("main", "OGQXDO2I39")
    ))
channel = connection.channel()


for i in range(1):
    message = {
        "specversion": "1.0",
        "type": "com.hotellistat.scraping.synxis.apisession",
        "id": str(uuid.uuid4()),
        "source": "testing",
        "nopublish": True,
        "data": {
            "job": str(uuid.uuid4()),
            "group": str(uuid.uuid4()),
            "id_ota": 1,
            "id_hotel": 2074,
            "hotel_ota_id": "58053",
            "crawl_date": "2021-05-19",
            "days_to_crawl": 2,
            "length_of_stay": 1,
            "max_persons": 2,
            "country_code": "de",
            "currency": "EUR",
            "closures": []

            # "type": "auto",
            # "ota_id": 1,
            # "hotel_id": 17,
            # "hotel_ota_id": "de/rocco-forte-the-charles.de.html",
            # "offset": 0,
            # "crawl_date": "2021-05-01",
            # "days_to_crawl": 1,
            # "length_of_stay": 1,
            # "max_persons": 2,
            # "country_code": "de",
            # "currency": "EUR",
            # "closures": []
        }

    }
    channel.basic_publish(
        exchange='',
        properties=pika.BasicProperties(
            delivery_mode=2,
        ),
        routing_key='com.hotellistat.scraping.synxis.apisession',
        body=json.dumps(message)
    )


connection.close()
