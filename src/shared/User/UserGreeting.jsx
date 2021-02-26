import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { selectLoggedInUser } from 'store/entities/selectors';
import { selectIsLoggedIn } from '../../store/auth/selectors';

export const UserGreeting = ({ isLoggedIn, firstName, email }) =>
  isLoggedIn && (
    <span className="usa-nav__link">
      <strong>{firstName ? `Welcome, ${firstName}` : email}</strong>
    </span>
  );

UserGreeting.propTypes = {
  email: PropTypes.string.isRequired,
  firstName: PropTypes.string,
  isLoggedIn: PropTypes.bool.isRequired,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);
  return {
    isLoggedIn: selectIsLoggedIn,
    firstName: user.first_name,
    email: user.email,
  };
};

export default connect(mapStateToProps)(UserGreeting);
