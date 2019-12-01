import React from 'react';
import PropTypes from 'prop-types';

export default class UserAttribute extends React.PureComponent {
    static propTypes = {
        id: PropTypes.string.isRequired,
        profile: PropTypes.object.isRequired,
        actions: PropTypes.shape({
            getSteamUserData: PropTypes.func.isRequired,
        }).isRequired,
    };

    componentDidMount() {
        this.props.actions.getSteamUserData(this.props.id);
    }

    render() {
        const profile = this.props.profile;
        if (!profile.personaname) {
            return null;
        }

        return (
            <div style={style.container}>
                <a
                    href={profile.profileurl}
                    target='_blank'
                    rel='noopener noreferrer'
                >
                    <i className='fa fa-steam'/>{' ' + profile.personaname}
                </a>
            </div>
        );
    }
}

const style = {
    container: {
        margin: '5px 0',
    },
};
