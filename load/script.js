import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 100,
    duration: '30s',
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    const email = `user_${__VU}_${__ITER}_${Date.now()}@example.com`;
    const password = 'password123';

    // 1. Register
    const registerRes = http.post(`${BASE_URL}/api/v1/auth/register`, JSON.stringify({
        email: email,
        password: password,
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(registerRes, {
        'register status 201': (r) => r.status === 201,
    });

    // 2. Login
    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
        email: email,
        password: password,
    }), {
        headers: { 'Content-Type': 'application/json' },
    });
    check(loginRes, {
        'login status 200': (r) => r.status === 200,
    });

    const accessToken = loginRes.json('access_token');

    // 3. 10 раз: Create URL → Redirect
    for (let i = 0; i < 10; i++) {
        const createRes = http.post(`${BASE_URL}/api/v1/urls`, JSON.stringify({
            original_url: 'https://example.com',
        }), {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`,
            },
        });
        check(createRes, {
            'create url status 201': (r) => r.status === 201,
        });

        const shortCode = createRes.json('short_code');

        const redirectRes = http.get(`${BASE_URL}/${shortCode}`, {
            redirects: 0,
        });
        check(redirectRes, {
            'redirect status 302': (r) => r.status === 302,
        });
    }
}