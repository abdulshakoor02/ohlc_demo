# OHLC (open high low close candle data)

The the application fetches the live trade data of btcusdt,ethusdt & pepeusdt pairs from binance and aggreagates the data to form OHLC and exposes each tick data to grpc server.Application also stores each closed candle for 1m to postgres database.

# Prerequisite

Docker should be installed.

# Installation

=> clone the repository
=> setup the .env file based on the env.example
=> run docker compose up.
