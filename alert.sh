cat $1 | (curl http://mwc16.azurewebsites.net/distance/update -s -S --header 'Content-Type: application/json' -H 'Accept: application/json' -d @- ) 
