import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export let options = {
  stages: [
    { duration: '1m', target: 500 },
    { duration: '1m', target: 1000 },
    { duration: '2m', target: 2000 },
    { duration: '2m', target: 3500 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<800'],
    http_req_failed: ['rate<0.01'],
  },
};

const actions = ['login', 'click', 'purchase', 'logout'];

export default function () {
  const payload = JSON.stringify({
    user_id: randomIntBetween(1, 10000),
    action: actions[Math.floor(Math.random() * actions.length)],
  });

  const headers = { 'Content-Type': 'application/json' };

  const res = http.post('http://api:8080/event', payload, { headers });

  check(res, {
    'status is 202': (r) => r.status === 202,
    'not 500': (r) => r.status !== 500,
    'response time OK': (r) => r.timings.duration < 800,
  });

  sleep(1);
}
