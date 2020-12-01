import React from 'react';
import classNames from 'classnames/bind';

import { ReactComponent as MmLogo } from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const MilMoveHeader = () => (
  <div className={cx('mm-header')}>
    <MmLogo />
    <div className={cx('links')}>
      <span>
        <a href="#">Navigation Link</a>
      </span>
      <span>
        <a href="#">Navigation Link</a>
      </span>
      <span>
        <a href="#">Navigation Link</a>
      </span>
      <span className={cx('line-add')}>&nbsp;</span>
      <span>Baker, Riley</span>
      <span>
        <a href="#">Sign out</a>
      </span>
    </div>
  </div>
);

export default MilMoveHeader;
