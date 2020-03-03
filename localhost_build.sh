#!/bin/bash

username="easy"
password="easy"
host="192.167.8.141"
port=22

expect -c "
  spawn rsync -av -e ssh --port=$port --exclude=".git/" --exclude=".git/*" --exclude=".idea/" --exclude=".idea/*" . $username@$host:~/ds-yibasuo-web/
  expect {
    *yes/no* { send -- yes\r;exp_continue; }
    *assword* { send -- $password\r; }
  }
  interact
"