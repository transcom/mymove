import React, { Component } from 'react';
import { connect } from 'react-redux';

import styles from './Dashboard.module.scss';

import requireCustomerState from 'containers/requireCustomerState/requireCustomerState';
import { profileStates } from 'constants/customerStates';
import { withContext } from 'shared/AppContext';
import withRouter from 'utils/routing';

export class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <div className={styles.dashboardContainer}>
        <div className={`usa-prose grid-container ${styles['grid-container']}`}>
          <h1>MULTIPLE MOVES DASHBOARD</h1>
          <p>This dashboard will allow multiple moves</p>
        </div>
      </div>
    );
  }
}

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});

export default withContext(
  withRouter(connect(mergeProps)(requireCustomerState(Dashboard, profileStates.BACKUP_CONTACTS_COMPLETE))),
);
