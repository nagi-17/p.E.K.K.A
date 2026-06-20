# P.E.K.K.A's PlayHouse

You should have Docker and Docker Compose installed.

SETUP INSTRUCTIONS :
1. Clone the Repository
   ```bash
   git clone https://github.com/nagi-17/p.E.K.K.A.git
   cd p.E.K.K.A
   ```

2. Copy the sample env file

   ```bash
   cp .env.sample .env
   ```

3. Run the following command in the project root folder to start the containers:
   ```bash
   docker compose up --build
   ```

To stop the game services, press `Ctrl+C` in your terminal or run `docker compose down`.