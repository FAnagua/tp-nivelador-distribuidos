import sys


def generar_docker_compose(nombre_archivo: str, n_clientes: int):

    text_compose = f"""services:
  server:
    build:
      context: ./services/server
      dockerfile: Dockerfile
    container_name: server
    environment:
      - PYTHONUNBUFFERED=1
      - SERVER_HOST=server
      - SERVER_PORT=5678
      - AGENCY_QUORUM_MIN={n_clientes}
    ports:
      - "5678:5678"
"""

    for i in range(0, n_clientes):
        text_compose += f"""
  client_{i}:
    build:
      context: ./services/client
      dockerfile: Dockerfile
    container_name: client_{i}
    depends_on:
      - server
    environment:
      - AGENCY_ID={i}
      - SERVER_HOST=server
      - SERVER_PORT=5678
      - INPUT_FILE=/input/input-{i}.csv
      - OUTPUT_FILE=/output/output-{i}.csv
      - BATCH_SIZE=30
    volumes:
      - ./input:/input
      - ./output:/output
"""
    

    with open(nombre_archivo, 'w') as file:
        file.write(text_compose)

def main():
    if len(sys.argv) != 2:
        print("Se epera: python generador_docker_compose.py <cantidad_clientes>")
        sys.exit(1)

    nombre_archivo = "docker-compose.yaml"

    try:
        n_clientes = int(sys.argv[1])
    except ValueError:
        print("Cantidad de clientes debe ser un número entero.")
        sys.exit(1)

    generar_docker_compose(nombre_archivo, n_clientes)
    print(f"Se creo el archivo {nombre_archivo}")

if __name__ == "__main__":
    main()