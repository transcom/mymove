import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

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

const mapStateToProps = ({ user }) => ({
  isLoggedIn: user.isLoggedIn,
  firstName: user.firstName,
  email: user.email,
});

export default connect(mapStateToProps)(UserGreeting);
