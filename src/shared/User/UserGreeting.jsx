import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const UserGreeting = ({ isLoggedIn, firstName, email }) =>
  isLoggedIn && (
    <span>
      <strong>{firstName ? `Welcome, ${firstName}` : email}</strong>
    </span>
  );

UserGreeting.propTypes = {
  isLoggedIn: PropTypes.bool.isRequired,
  firstName: PropTypes.string.isRequired,
  email: PropTypes.string.isRequired,
};

const mapStateToProps = ({ user }) => ({
  isLoggedIn: user.isLoggedIn,
  firstName: user.firstName,
  email: user.email,
});

export default connect(mapStateToProps)(UserGreeting);
