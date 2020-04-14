import React from 'react';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const PageHeader = () => (
  <div className={cx('page-header')}>
    <div>
      <h2>
        Smith, Kerry
        <span data-testid="tag" className="usa-tag usa-tag--cyan usa-tag--large">
          #ABC123K
        </span>
      </h2>
    </div>
    <div>right</div>
  </div>
);

export default PageHeader;
