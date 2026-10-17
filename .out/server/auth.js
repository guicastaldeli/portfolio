export class Auth {
    form = null;
    email = null;
    password = null;
    constructor() {
        console.log('Auth constructed');
        this.form = document.querySelector('#login-form');
        console.log('form found:', !!this.form);
        this.form = document.querySelector('#login-form');
        this.email = document.querySelector('#email');
        this.password = document.querySelector('#password');
        this.init();
    }
    /**
     * Set
     */
    set() {
        if (!this.form)
            return;
        this.form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const err = document.getElementById('err');
            err.textContent = '';
            try {
                const res = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    credentials: 'include',
                    body: JSON.stringify({
                        email: this.email?.value,
                        password: this.password?.value
                    })
                });
                if (res.ok) {
                    window.location.href = '/';
                }
                else {
                    err.textContent = 'Invalid credentials';
                }
            }
            catch (e) {
                err.textContent = 'Network error';
            }
        });
    }
    /**
     *
     * Init
     *
     */
    init() {
        this.set();
    }
}
