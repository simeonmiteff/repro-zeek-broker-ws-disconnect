module Test;

redef Broker::disable_ssl = T;

global cnt: count = 0;

event zeek_init()
    {
    Broker::listen_websocket("0.0.0.0", 16666/tcp);
    Broker::subscribe("/simeonmiteff/test");
    }

event Test::evt()
    {
    ++cnt;
    }

event Broker::peer_added(endpoint: Broker::EndpointInfo, msg: string)
    {
    print fmt("peer added, endpoint=%s, msg=%s", endpoint, msg);
    }

event Broker::peer_lost(endpoint: Broker::EndpointInfo, msg: string)
    {
    print fmt("peer lost at cnt=%d, endpoint=%s, msg=%s", cnt, endpoint, msg);
    terminate();
    }
