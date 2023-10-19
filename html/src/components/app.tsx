import { Component, h } from 'preact';
import { Form } from './form';

interface Props {}

interface State {
    space: string;
    user: string;
}

export class App extends Component<Props, State> {
    constructor(props) {
        super(props);
        this.state = { space: '', user: '' };
    }

    componentDidMount() {
        this.fetchData();
    }

    fetchData = async () => {
        try {
            const response = await fetch('/switch-space?auth=fe');
            if (response.status < 300) {
                const jsonData = await response.json();
                const info = jsonData.cookie.split('_');
                let user,
                    space = '';
                for (const f of info) {
                    const map = f.split(':');
                    switch (map[0]) {
                        case 'space':
                            space = map[1];
                            break;
                        case 'u':
                            user = map[1];
                            break;
                    }
                }
                this.setState({ space, user });
            } else {
                this.setState({ space: '', user: '' });
            }
        } catch (error) {
            this.setState({ space: '', user: '' });
            console.error('Error fetching data:', error);
        }
    };
    render() {
        const header = this.state.space?.length ? (
            <div class="badge">
                Hi <div class="current-space">{this.state.user}</div> you are in{' '}
                <div class="current-space space">{this.state.space}</div> box.
            </div>
        ) : (
            ''
        );

        return (
            <div class="login">
                {header}
                <h1>Switch to space</h1>
                <Form name={this.state.space} />
            </div>
        );
    }
}
