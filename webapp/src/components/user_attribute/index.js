import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getSteamUserData} from '../../actions';

import {steamProfileInfo} from 'selectors';

import UserAttribute from './user_attribute.jsx';

function mapStateToProps(state, ownProps) {
    const id = ownProps.user ? ownProps.user.id : '';
    const profile = steamProfileInfo(state, id);

    return {
        id,
        profile,
    };
}

function mapDispatchToProps(dispatch) {
    return {
        actions: bindActionCreators({
            getSteamUserData,
        }, dispatch),
    };
}

export default connect(mapStateToProps, mapDispatchToProps)(UserAttribute);
