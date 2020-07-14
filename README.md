# librebread

Librebread is server for mock SMS sender services and smtp service for testing purposes. SmsRu and Devino telecom has beed imaplemended. Librebread is just random name

## Docker

https://hub.docker.com/r/vasyahuyasa/librebread

## API

### HTTP 80 port

| URL                    | DESCRIPTION |
|------------------------|-------------|
| `/`                    | index, sent messages |

__LibeSMS__

| URL                | DESCRIPTION |
|--------------------|-------------|
| `/libre/send`      | send sms    |
| `/libew/check`     | not implemented |

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

__LibeSMS__

| URL                | DESCRIPTION |
|--------------------|-------------|
| `/libre/send`      | send sms    |
| `/libew/check`     | not implemented |

### SMTP 25 port

Plain SMTP server

### POP3 110 port

Mock of pop3 server