import React from 'react';
import classNames from 'classnames/bind';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const FOUOHeader = () => (
  <div className={cx('fouo-header')}>
    <div className={cx('fouo-header--text')}>Unclassified // For official use only</div>
  </div>
);

export default FOUOHeader;
