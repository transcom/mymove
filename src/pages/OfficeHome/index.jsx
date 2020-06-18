import { connect } from 'react-redux';

import OfficeHome from './OfficeHome';

import { selectCurrentUser } from 'shared/Data/users';

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);

  return {
    activeRole: state.auth.activeRole,
    userRoles: user.roles,
  };
};

export default connect(mapStateToProps)(OfficeHome);
