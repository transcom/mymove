import React from 'react';
// import PropTypes from 'prop-types';
import { Link, matchPath, useLocation } from 'react-router-dom';

import styles from './TestingHeader.module.scss';

import CUIHeader from 'components/CUIHeader/CUIHeader';
// import { OfficeUserInfoShape } from 'types/index';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';

const TestingHeader = () => {
  // may need to add these back to the useTXOquery params isLoading, isError, errors
  // const { move, order, customerData, isLoading, isError, errors } = useTXOMoveInfoQueries('DP3QXQ');
  const location = useLocation();
  const displayChangeRole = !matchPath(
    {
      path: '/select-application',
      end: true,
    },
    location.pathname,
  );
  return (
    <div className={styles.test}>
      <CUIHeader />
      <div className={styles.changeRole}>
        {displayChangeRole && <Link to="/select-application">Change user role</Link>}
      </div>
      <OfficeLoggedInHeader />
    </div>
  );
};

export default TestingHeader;
