# librebread

Librebread is server for mock SMS sender services for testing purposes. SmsRu and Devino telecom has beed imaplemended. Librebread is just random name

## Docker

https://hub.docker.com/r/vasyahuyasa/librebread

## API

### HTTP 80 port

| URL                    | DESCRIPTION |
|------------------------|-------------|
| `/`                    | index, sent messages |
### HTTPS 443 port

__DevinoTelecom__

| URL                    | DESCRIPTION |
|------------------------|-------------|
| `/rest/user/sessionid` |  always return session id MOCK-SESSION-ID |
| `/rest/sms/send`       | send sms |
| `/rest/sms/state`      | check account |

__SmsRU__

| URL                | DESCRIPTION |
|--------------------|-------------|
| `/sms/user/send`   | send sms    |
| `/sms/user/status` | messages status |
