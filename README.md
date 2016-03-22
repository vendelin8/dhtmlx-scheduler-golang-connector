# golang database connector for dhtmlx scheduler
dhtmlx-scheduler-golang-connector

* Tested for MySQL and SQLite.
* Not Tested for PostgreSQL, but probably works.

### To download the package use:
go get github.com/vendelin8/dhtmlx-scheduler-golang-connector

### For MySql databases, make sure to initialize database with utf-8 collation:
CREATE DATABASE test
DEFAULT CHARACTER SET utf8
DEFAULT COLLATE utf8\_general\_ci;

### To run the examples, download dhtmlx from here:
* http://dhtmlx.com/docs/products/dhtmlxScheduler/download.shtml
* and extract the codebase directory to examples/static/
* Then call: go install github.com/vendelin8/dhtmlx-scheduler-golang-connector/examples && $GOPATH/bin/examples
* Open http://0.0.0.0:1212/static/ in your browser


Contributions are welcome.
