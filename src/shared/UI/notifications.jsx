import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { getNotifications } from 'shared/UI/ducks';

import './notifications.css';

export function Notifications(props) {
  return (
    <div className="usa-grid notifications">
      {props.notifications.map(notification => (
        <Alert type={notification.severity} heading={notification.title} key={notification.createdAt}>
          {notification.message}
        </Alert>
      ))}
    </div>
  );
}

Notifications.defaultProps = {
  notifications: [],
};
Notifications.propTypes = {
  notifications: PropTypes.arrayOf(PropTypes.object),
};

const mapStateToProps = state => {
  return {
    notifications: getNotifications(state),
  };
};
export default connect(mapStateToProps)(Notifications);
