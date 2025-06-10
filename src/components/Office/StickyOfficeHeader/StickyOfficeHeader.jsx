import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './StickyOfficeHeader.module.scss';

import CUIHeader from 'components/CUIHeader/CUIHeader';
// import { OfficeUserInfoShape } from 'types/index';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';

const StickyOfficeHeader = ({ displayChangeRole }) => {
  // may need to add these back to the useTXOquery params isLoading, isError, errors
  // const { move, order, customerData, isLoading, isError, errors } = useTXOMoveInfoQueries('DP3QXQ');
  return (
    <div className={styles.stickyHeader}>
      <CUIHeader />
      <OfficeLoggedInHeader />
      <div className={styles.changeRole}>
        {displayChangeRole && <Link to="/select-application">Change user role</Link>}
      </div>
    </div>
  );
};

StickyOfficeHeader.propTypes = {
  displayChangeRole: PropTypes.string,
};

StickyOfficeHeader.defaultProps = {
  displayChangeRole: false,
};

export default StickyOfficeHeader;
