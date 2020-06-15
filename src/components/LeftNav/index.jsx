import React from 'react';
import classNames from 'classnames/bind';
import { Tag } from '@trussworks/react-uswds';

import { ReactComponent as AlertIcon } from '../../shared/icon/alert.svg';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const LeftNav = () => (
  <div className={cx('sidebar')}>
    <nav className={cx('left-nav')}>
      <a href="#" className={cx('active')}>
        Requested Shipments
        <Tag className="usa-tag--alert usa-tag--alert--small">
          <AlertIcon />
        </Tag>
      </a>
      <a href="#orders-anchor">
        Orders
        <Tag className="usa-tag--teal">INTL</Tag>
      </a>
      <a href="#">Allowances</a>
      <a href="#">
        Customer Info
        <Tag>3</Tag>
      </a>
    </nav>
  </div>
);

export default LeftNav;
