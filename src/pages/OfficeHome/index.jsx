import { connect } from 'react-redux';

import { selectCurrentUser } from 'shared/Data/users';

import OfficeHome from './OfficeHome';

const mapStateToProps = (state) => {
  const user = selectCurrentUser(state);

  return {
    activeRole: state.auth.activeRole,
    userRoles: user.roles,
  };
};

export default connect(mapStateToProps)(OfficeHome);
