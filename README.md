# sidecar

[![GoReportCard Widget]][GoReportCard] [![Travis Widget]][Travis]

[GoReportCard]: https://goreportcard.com/report/github.com/amalgam8/sidecar
[GoReportCard Widget]: https://goreportcard.com/badge/github.com/amalgam8/sidecar
[Travis]: https://travis-ci.org/amalgam8/sidecar
[Travis Widget]: https://travis-ci.org/amalgam8/sidecar.svg?branch=master

A language agnostic sidecar for building microservice applications with
automatic service registration, and load-balancing

### Architecture

![Sidecar architecture](https://github.com/amalgam8/sidecar/blob/master/sidecar.jpg)

### Environment variables needed to run sidecar
    
* ENDPOINT_HOST, ENDPOINT_PORT -- IP and port of service instance to register
* SERVICE -- Name of service to register
* SD_URL, SD_TOKEN -- URL and auth token for use with service discovery
* RE_URL -- Service proxy control plane URL
* SP_TENANT_ID, SP_TENANT_TOKEN - ID and auth token for use with service  proxy
  
#### IBM MessageHub integration - environment variables
* VCAP_SERVICES_MESSAGEHUB_0_CREDENTIALS_KAFKA_REST_URL
* VCAP_SERVICES_MESSAGEHUB_0_CREDENTIALS_API_KEY
* VCAP_SERVICES_MESSAGEHUB_0_CREDENTIALS_KAFKA_BROKERS_SASL_[0,1,2,3..]
* VCAP_SERVICES_MESSAGEHUB_0_CREDENTIALS_USER
* VCAP_SERVICES_MESSAGEHUB_0_CREDENTIALS_PASSWORD

### Running sidecar

Command line arguments
* -register - enable automatic service registration
* -proxy - enable nginx service proxy
* -log - use filebeat to propagate nginx logs to logstash
* -supervise - invoke and monitor application process

Usage:
```bash
sidecar -register -proxy -log -supervise myapp arg1 arg2 -arg3=3 -arg4=4
```
## License
Copyright 2016 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
