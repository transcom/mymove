import React from 'react';
import classNames from 'classnames/bind';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const CUIHeader = () => (
  <div className={cx('cui-header')}>
    <div className={cx('cui-header--text')}>Unclassified // For official use only</div>
  </div>
);

export default CUIHeader;
