# detectdupes
A small code to detect duplicates in access log

generateAccessLog.go - generate 10M records randomly and push they in MongoDB

detectDupesM.go - service to detect duplicates without cache (response delay ~7s)

cacheAccessLog.go - generate a new dataBase from 10M records with user ids and unique IPs

detectDupesViaCache - service to detect duplicates with cache (response delay ~1ms)
