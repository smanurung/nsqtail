# nsqtail

Tool for tailing NSQ

## Usage
```
smanurung@sonny-macbook nsqtail (master)*$ ./nsqtail --help
usage: nsqtail --topic=TOPIC --lookupd-http-addr=LOOKUPD-HTTP-ADDR [<flags>]

Flags:
  --help               Show context-sensitive help (also try --help-long and --help-man).
  --topic=TOPIC        topic to listen to
  --lookupd-http-addr=LOOKUPD-HTTP-ADDR
                       NSQlookupd address with port, e.g. 127.0.0.1:4161
  --max-in-flight=100  NSQ consumer max-in-flight number
  ```

## Example
```
./nsqtail --topic=topic_name --lookupd-http-addr=lookupd_address:4161
```