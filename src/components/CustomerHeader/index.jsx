import React from 'react';
import classNames from 'classnames/bind';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const CustomerHeader = () => (
  <div className={cx('cust-header')}>
    <div>
      <div className={cx('name-block')}>
        <h2>Smith, Kerry</h2>
        <span className="usa-tag usa-tag--cyan usa-tag--large">#ABC123K</span>
      </div>
      <div>
        <p>
          Navy E-6
          <span className={cx('vertical-bar')}>|</span>
          DoD ID 999999999
        </p>
      </div>
    </div>
    <div className={cx('info-block')}>
      <div>
        <p>Authorized origin</p>
        <h4>JBSA Lackland</h4>
      </div>
      <div>
        <p>Authorized destination</p>
        <h4>JB Lewis-McChord</h4>
      </div>
      <div>
        <p>Report by</p>
        <h4>27 Mar 2020</h4>
      </div>
    </div>
  </div>
);

export default CustomerHeader;
