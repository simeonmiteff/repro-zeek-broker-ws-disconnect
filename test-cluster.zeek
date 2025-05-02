module Test;

global cnt: count = 0;

event zeek_init()
	{
	Cluster::listen_websocket([ $listen_host="0.0.0.0", $listen_port=16666/tcp ]);
	Cluster::subscribe("/simeonmiteff/test");
	}

event zeek_done()
	{
	print "Last count was", cnt;
	}

event Test::evt()
	{
	++cnt;
	}

event Cluster::websocket_client_lost(info: Cluster::EndpointInfo)
	{
	print fmt("peer lost at cnt=%d, endpoint=%s", cnt, info);
	terminate();
	}

event Cluster::websocket_client_added(info: Cluster::EndpointInfo,
    subscriptions: string_vec)
	{
	print fmt("peer added, endpoint=%s, subscriptions=%s", info, subscriptions);
	}
