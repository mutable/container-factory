#!/usr/bin/env node
require('http').createServer(function(req, res) {
  res.end('goodbye\n')
}).listen(process.env.PORT || 9000)
