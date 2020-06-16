import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const TabNav = ({ items }) => (
  <nav className={cx('usa-nav', 'tab-nav')}>
    <div className="usa-nav__inner">
      <ul className={cx('usa-nav__primary', 'tab-list')}>
        {items.map((item, index) => (
          <li key={index.toString()} className={cx('usa-nav__primary-item', 'tab-item')} role="tab">
            {item}
          </li>
        ))}
      </ul>
    </div>
  </nav>
);

TabNav.propTypes = {
  items: PropTypes.arrayOf(PropTypes.node),
};

export default TabNav;
