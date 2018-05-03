var http = require("http")
var Log = require("log")
  , log = new Log("info");

log.info("Setting up server")

http.createServer(function(request, response) {
	response.writeHead(200, {"Content-Type": "text/json"});
	response.end();
}).listen(8888);
