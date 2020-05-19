import http from 'k6/http';
import { check } from "k6";
import { Rate, Counter } from "k6/metrics";

export let errorTrx = new Counter("Error Transaction"); // Counter for Error non 200 HTTP code
export let errorRate = new Rate("Error Rate");

function makeid(length) {
    var result           = '';
    var characters       = 'abcdefghijklmnopqrstuvwxyz';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

export default function() {

    let id = makeid(20)
    let url = `http://localhost:2222/users/${id}`;

    let hit2 =  http.get(url);

    //parse and check response
    let checkRes = check(hit2, {
        "Status is 200": (r) => r.status === 200
    });

    if (hit2.status !== 200) {
        console.log(hit2.body)
    }

    errorTrx.add(!checkRes);
    errorRate.add(!checkRes);
}