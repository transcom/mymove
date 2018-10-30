import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

const UserGreeting = ({ isLoggedIn, firstName }) =>
  isLoggedIn && (
    <span>
      <strong>Welcome, {firstName}</strong>
    </span>
  );

UserGreeting.propTypes = {
  isLoggedIn: PropTypes.bool.isRequired,
  firstName: PropTypes.string.isRequired,
};

const mapStateToProps = ({ user }) => ({
  isLoggedIn: user.isLoggedIn,
  firstName: user.firstName,
});

export default connect(mapStateToProps)(UserGreeting);
