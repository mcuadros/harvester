#!/bin/bash

# collectord - Startup script for collectord

# chkconfig: 35 85 15
# description: low footprint collector and parser for events and logs
# processname: collectord
# config: /etc/collectord.conf
# pidfile: /var/run/collectord.pid

. /etc/rc.d/init.d/functions

# things from mongod.conf get there by mongod reading it

CONFIGFILE="/etc/collectord.conf"
OPTIONS=" -f $CONFIGFILE"
PIDFILE=/var/run/collectord.pid
LOCKFILE=/var/lock/subsys/collectord
USER=collectord
exec=/usr/bin/collectord


start()
{
  echo -n $"Starting collectord: "
  daemon --pidfile="$PIDFILE" --user "$MONGO_USER" $exec $OPTIONS
  RETVAL=$?
  echo
  [ $RETVAL -eq 0 ] && touch $LOCKFILE
}

stop()
{
  echo -n $"Stopping collectord: "
  killproc -p "$PIDFILE" $exec
  RETVAL=$?
  echo
  [ $RETVAL -eq 0 ] && rm -f $LOCKFILE
}

restart () {
    stop
    start
}

ulimit -n 12000
RETVAL=0

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  restart|reload|force-reload)
    restart
    ;;
  condrestart)
    [ -f /var/lock/subsys/mongod ] && restart || :
    ;;
  status)
    status -p "$PIDFILE" -l "$LOCKFILE" $exec
    RETVAL=$?
    ;;
  *)
    echo "Usage: $0 {start|stop|status|restart|reload|force-reload|condrestart}"
    RETVAL=1
esac

exit $RETVAL
