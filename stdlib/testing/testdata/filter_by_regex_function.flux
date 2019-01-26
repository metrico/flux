package main
 
import "testing"

option now = () => (2030-01-01T00:00:00Z)

inData = "
#datatype,string,long,dateTime:RFC3339,long,string,string,string,string
#group,false,false,false,false,true,true,true,true
#default,_result,,,,,,,
,result,table,_time,_value,_field,_measurement,host,name
,,0,2018-05-22T19:53:26Z,15204688,io_time,diskio,host.local,disk0
,,0,2018-05-22T19:53:36Z,15204894,io_time,diskio,host.local,disk0
,,0,2018-05-22T19:53:46Z,15205102,io_time,diskio,host.local,disk0
,,0,2018-05-22T19:53:56Z,15205226,io_time,diskio,host.local,disk0
,,0,2018-05-22T19:54:06Z,15205499,io_time,diskio,host.local,disk0
,,0,2018-05-22T19:54:16Z,15205755,io_time,diskio,host.local,disk0
,,1,2018-05-22T19:53:26Z,648,ixo_time,diskio,host.local,disk2
,,1,2018-05-22T19:53:36Z,648,ixo_time,diskio,host.local,disk2
,,1,2018-05-22T19:53:46Z,648,ixo_time,diskio,host.local,disk2
,,1,2018-05-22T19:53:56Z,648,ixo_time,diskio,host.local,disk2
,,1,2018-05-22T19:54:06Z,648,ixo_time,diskio,host.local,disk2
,,1,2018-05-22T19:54:16Z,648,ixo_time,diskio,host.local,disk2
"

outData = "
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,long,string,string,string,string
#group,false,false,true,true,false,false,true,true,true,true
#default,_result,,,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement,host,name
,,0,2018-05-20T19:53:26Z,2030-01-01T00:00:00Z,2018-05-22T19:54:16Z,15205755,io_time,diskio,host.local,disk0
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,long,string,string,string,string
#group,false,false,true,true,false,false,true,true,true,true
#default,_result,1,2018-05-20T19:53:26Z,2030-01-01T00:00:00Z,,,ixo_time,diskio,host.local,disk2
,result,table,_start,_stop,_time,_value,_field,_measurement,host,name
"
regexFunc = (table=<-, regLiteral) =>
	(table
		|> range(start: 2018-05-20T19:53:26Z)
		|> filter(fn: (r) =>
			(r._field =~ regLiteral))
		|> max())
t_filter_by_regex_function = (table=<-) =>
	(table
		|> regexFunc(regLiteral: /io.*/))

test _filter_by_regex_function = () =>
	({input: testing.loadStorage(csv: inData), want: testing.loadMem(csv: outData), fn: t_filter_by_regex_function})

testing.run(case: _filter_by_regex_function)