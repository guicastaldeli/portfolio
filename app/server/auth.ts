export class Auth {
    private form: HTMLFormElement | null = null;

    private email: HTMLInputElement | null = null;
    private password: HTMLInputElement | null = null;

    constructor() {
        this.form = document.querySelector('#login-form');
        this.email = document.querySelector('#email');
        this.password = document.querySelector('#password');
        
        this.init();
    }

    /**
     * Set
     */
    private set(): void {
        if(!this.form) return;

        this.form.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const err = document.getElementById('err');
            err!.textContent = '';
        
            const res = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    email: this.email?.value,
                    password: this.password?.value
                })
            });
            if(res.ok) {
                window.location.href = '/'
            } else {
                err!.textContent = 'Invalid credentials';
            }
        });
    }

    /**
     * 
     * Init
     * 
     */
    private init(): void {
        this.set();
    }
}