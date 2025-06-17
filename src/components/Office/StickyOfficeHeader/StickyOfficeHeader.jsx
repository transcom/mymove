import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';

import styles from './StickyOfficeHeader.module.scss';

import CUIHeader from 'components/CUIHeader/CUIHeader';
import OfficeLoggedInHeader from 'containers/Headers/OfficeLoggedInHeader';

const StickyOfficeHeader = ({ displayChangeRole }) => {
  return (
    <div className={styles.stickyHeader}>
      <CUIHeader />
      <div className={styles.changeRole}>
        {displayChangeRole && <Link to="/select-application">Change user role</Link>}
      </div>
      <OfficeLoggedInHeader />
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
