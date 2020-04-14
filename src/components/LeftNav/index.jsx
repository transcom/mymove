import React from 'react';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const LeftNav = () => (
  <nav className={cx('left-nav')}>
    <a>Requested Shipments</a>
    <a>Orders</a>
    <a>Allowances</a>
    <a>Customer Info</a>
  </nav>
);

export default LeftNav;
