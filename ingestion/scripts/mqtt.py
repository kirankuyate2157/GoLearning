
import signal
import sys
import argparse
import paho.mqtt.client as mqtt
import psycopg2
from datetime import datetime

# Parse command line arguments
parser = argparse.ArgumentParser(description='MQTT to TimescaleDB bridge')
parser.add_argument('--ip', required=True, help='MQTT broker address')
parser.add_argument('--port', type=int, default=1883, help='MQTT broker port')
parser.add_argument('--topic', required=True, help='MQTT topic to subscribe')
parser.add_argument('--username', help='MQTT username (optional)')
parser.add_argument('--password', help='MQTT password (optional)')
parser.add_argument('--qos', type=int, default=0, help='MQTT QoS level')
parser.add_argument('--feed-uid', required=True, help='Feed UID for this MQTT subscription')
parser.add_argument('--db-host', required=True, help='TimescaleDB host')
parser.add_argument('--db-port', type=int, default=5432, help='TimescaleDB port')
parser.add_argument('--db-name', required=True, help='Database name')
parser.add_argument('--db-user', required=True, help='Database user')
parser.add_argument('--db-password', required=True, help='Database password')
parser.add_argument('--db-table', default='mqtt_messages', help='Table name')
args = parser.parse_args()

BROKER = args.ip
PORT = args.port
TOPIC = args.topic
MQTT_USERNAME = args.username
MQTT_PASSWORD = args.password
MQTT_QOS = args.qos
FEED_UID = args.feed_uid
DB_HOST = args.db_host
DB_PORT = args.db_port
DB_NAME = args.db_name
DB_USER = args.db_user
DB_PASSWORD = args.db_password
DB_TABLE = args.db_table

# Connect to TimescaleDB
try:
    conn = psycopg2.connect(
        host=DB_HOST,
        port=DB_PORT,
        database=DB_NAME,
        user=DB_USER,
        password=DB_PASSWORD
    )
    cursor = conn.cursor()
    print(f"Connected to TimescaleDB at {DB_HOST}")
    
    # Create table if not exists
    cursor.execute(f"""
        CREATE TABLE IF NOT EXISTS {DB_TABLE} (
            time TIMESTAMPTZ NOT NULL,
            feed_uid TEXT NOT NULL,
            topic TEXT NOT NULL,
            payload TEXT NOT NULL
        );
    """)
    
    # Create hypertable if not already created
    cursor.execute(f"""
        SELECT create_hypertable('{DB_TABLE}', 'time', if_not_exists => TRUE);
    """)
    conn.commit()
    print(f"Table {DB_TABLE} ready")
except Exception as e:
    print(f"Error connecting to TimescaleDB: {e}")
    sys.exit(1)

def on_connect(client, userdata, flags, rc):
    print(f"Connected to MQTT broker with result code {rc}")
    print(f"Subscribing to topic: {TOPIC} with QoS {MQTT_QOS}")
    client.subscribe(TOPIC, qos=MQTT_QOS)

def on_message(client, userdata, msg):
    payload = msg.payload.decode()
    print(f"Received on {msg.topic}: {payload}")
    
    # Insert into TimescaleDB
    try:
        cursor.execute(
            f"INSERT INTO {DB_TABLE} (time, feed_uid, topic, payload) VALUES (%s, %s, %s, %s)",
            (datetime.now(), FEED_UID, msg.topic, payload)
        )
        conn.commit()
        print(f"Stored in TimescaleDB")
    except Exception as e:
        print(f"Error storing to database: {e}")
        conn.rollback()


client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

# Set authentication if provided
if MQTT_USERNAME and MQTT_PASSWORD:
    client.username_pw_set(MQTT_USERNAME, MQTT_PASSWORD)
    print(f"MQTT authentication configured for user: {MQTT_USERNAME}")

def handle_sigint(signum, frame):
    print("SIGINT received, disconnecting...")
    client.disconnect()
    cursor.close()
    conn.close()
    sys.exit(0)

signal.signal(signal.SIGINT, handle_sigint)

print(f"Connecting to MQTT broker at {BROKER}:{PORT}...")
client.connect(BROKER, PORT, 60)
print(f"Starting ingestion for Feed: (UID: {FEED_UID})")
try:
    client.loop_forever()
except KeyboardInterrupt:
    print("KeyboardInterrupt, disconnecting...")
    client.disconnect()
    cursor.close()
    conn.close()
