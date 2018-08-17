# SADIS Server
-----------
This service acts as a gateway between ONOS and XOS. ONOS apps expect to pull
in deployment information using the SADIS service in a particular format.
The SADIS server listens to REST requests from ONOS, and when it receives a
request it looks up the appropriate object in XOS and delivers it back to ONOS
in the expected format.


### Build
- `docker build -t opencord/sadis-server  .`


### Run
- `docker run --rm --name sadis --env "SADISSERVER_PORT=4245" --env "SADISSERVER_XOS=10.90.0.101:30006" -p 4245:4245 opencord/sadis-server`
- Pass in the location of the XOS REST endpoint using the SADISSERVER_XOS
  environment variable
- The SADIS server will start listening on the configured port for connections
  from the SADIS ONOS app
- Configure the SADIS app with the URL of the SADIS server. See `sadis.json` for
  an example.


### Note:
- Every parameter is set by envconfig (https://github.com/kelseyhightower/envconfig)
- To modify any of the envvars, please set the ENV VAR with the prefix of `SADISSERVER`
- List of parameters set by envconfig:
```bash
Port       int    `default:"8000" desc:"port on which to listen for requests"`
Xos		   string `default:"127.0.0.1:8181" desc:"connection string with which to connect to XOS"`
Username   string `default:"admin@opencord.org" desc:"username with which to connect to XOS"`
Password   string `default:"letmein" desc:"password with which to connect to XOS"`
LogLevel   string `default:"info" envconfig:"LOG_LEVEL" desc:"detail level for logging"`
LogFormat  string `default:"text" envconfig:"LOG_FORMAT" desc:"log output format, text or json"`
```
