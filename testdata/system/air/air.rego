package system.air

import future.keywords.in

default allow := false

allow {
	input.sacNo="4613810890"
}

action := x {
    response := http.send({"url":"https://www.baidu.com","method":"GET"})
    print(response)

	x := {"test":1,"sacNo":input.sacNo}
}

order := y {
	y := 1
}