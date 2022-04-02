# Creamy Miner

## Usage

Configured through environment variables.

| Env Var                        | Description                                                                      | Default                                  | Example Value                            |
|--------------------------------|----------------------------------------------------------------------------------|------------------------------------------|------------------------------------------|
| CREAMY_MINER_ENV_FILE          | A file to load env vars from  Should contain one KEY=VALUE pair per line         | (none)                                   | /etc/creamy-miner.env                    |
| CREAMY_MINER_POOLS             | A list of pools to work for  Will try to get shares in as many pools as possible | (none)                                   | localhost:23380,127.0.0.1:23380          |
| CREAMY_MINER_POOL              | A single pool to work for  Ignored if CREAMY_MINER_POOLS is set                  | localhost:23380                          | localhost:23380                          |
| CREAMY_MINER_ADDRESS           | Address to mine to                                                               | c04rt84spfjc9xy88snx5r256qv0tmy664zcdrnc | c04rt84spfjc9xy88snx5r256qv0tmy664zcdrnc |
| CREAMY_MINER_FIELDS_DIR        | Location of snowblossom.1-10 directories                                         | fields                                   | /etc/snowblossom/fields/mainnet          |
| CREAMY_MINER_THREAD_MULTIPLIER | Multiplier of threads to run based on os.NumCPU                                  | 10                                       | 1                                        |     

## Configuring Systemd Service

- Download creamy-miner to `~/creamy-miner`
- Ensure it is executable by running `chmod +x ~/creamy-miner`
- Create a file at `~/creamy-miner.env` containing your config vars
- Make sure the Systemd user folder exists: `mkdir -p ~/.config/systemd/user/`
- Write the below service file to `~/.config/systemd/user/creamy-miner.service`:

```
[Unit]
Description=Creamy Miner
After=network.target

[Service]
Type=simple
Environment="CREAMY_MINER_ENV_FILE=creamy-miner.env"
WorkingDirectory=%h
ExecStart=%h/creamy-miner
Restart=on-failure

[Install]
WantedBy=default.target
```

- Reload, enable, and start the service:

```
systemctl --user daemon-reload
systemctl --user enable creamy-miner
systemctl --user start creamy-miner
```

- Done!
- To view logs, run `journalctl --user-unit creamy-miner`
- To update creamy-miner, stop the service with `systemctl --user stop creamy-miner`, overwrite `~/creamy-miner` with the latest version, and restart the service with `systemctl --user start creamy-miner`
