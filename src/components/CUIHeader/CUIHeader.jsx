import React from 'react';
import classNames from 'classnames/bind';

import styles from './CUIHeader.module.scss';

const cx = classNames.bind(styles);

const CUIHeader = ({ className }) => {
  return (
    <div className={classNames(styles['cui-header'], className)}>
      <div className={cx('cui-header--text')}>Controlled Unclassified Information</div>
    </div>
  );
};

export default CUIHeader;
