import { Component, h } from 'preact';

interface Props {
    name: string;
}

interface State {
    u: string;
    p: string;
    space: string;
    error: string;
}

export class Form extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { u: '', space: '', p: '', error: '' };

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        console.log(event.target.name);
        const state = { ...this.state };
        state[event.target.name] = event.target.value;
        this.setState({ ...state, error: '' });
    }

    async handleSubmit(event) {
        event.preventDefault();
        const opt: RequestInit = {
            headers: {
                'Content-Type': 'application/json;charset=UTF-8',
            },
            mode: 'same-origin',
            method: 'POST',
            body: JSON.stringify({
                url: '',
                p: { skip: true, v: btoa(this.state.p) },
                u: { skip: false, v: btoa(this.state.u) },
                space: { skip: false, v: btoa(this.state.space) },
            }),
        };
        const resp = await fetch('/switch-space', opt);
        if (resp.status < 300) {
            const j = await resp.json();
            window.location.href = '/' + j.url;
            this.setState({ ...this.state, error: '' });
            console.log('json resp', j);
        } else {
            this.setState({ ...this.state, error: 'Wrong username or password!' });
        }
    }

    render() {
        const errorMsg = this.state.error?.length ? <div class="badge error">{this.state.error}</div> : '';
        return (
            <form onSubmit={this.handleSubmit}>
                <input
                    type="text"
                    name="u"
                    value={this.state.u}
                    placeholder="username"
                    required={true}
                    onChange={this.handleChange}
                />
                <input
                    type="password"
                    value={this.state.p}
                    name="p"
                    placeholder="password"
                    required={true}
                    onChange={this.handleChange}
                />
                <input
                    type="text"
                    value={this.state.space}
                    name="space"
                    placeholder="login space"
                    required={true}
                    onChange={this.handleChange}
                />
                {errorMsg}
                <button type="submit" class="btn btn-primary btn-block btn-large">
                    Let me in.
                </button>
            </form>
        );
    }
}
