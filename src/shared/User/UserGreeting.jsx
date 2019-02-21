import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { selectCurrentUser } from 'shared/Data/users';

export const UserGreeting = ({ isLoggedIn, firstName, email }) =>
  isLoggedIn && (
    <span>
      <strong>{firstName ? `Welcome, ${firstName}` : email}</strong>
    </span>
  );

UserGreeting.propTypes = {
  email: PropTypes.string.isRequired,
  firstName: PropTypes.string,
  isLoggedIn: PropTypes.bool.isRequired,
};

const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  return {
    isLoggedIn: user.isLoggedIn,
    firstName: user.first_name,
    email: user.email,
  };
};

export default connect(mapStateToProps)(UserGreeting);
