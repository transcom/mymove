import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const TabNav = ({ items, role }) => (
  <nav className={cx('usa-nav', 'tabNav')} role={role}>
    <div className="usa-nav__inner">
      <ul className={cx('usa-nav__primary', 'tabList')} role="tablist">
        {items.map((item, index) => (
          <li key={index.toString()} className={cx('usa-nav__primary-item', 'tabItem')}>
            {item}
          </li>
        ))}
      </ul>
    </div>
  </nav>
);

TabNav.propTypes = {
  items: PropTypes.arrayOf(PropTypes.node).isRequired,
  role: PropTypes.string,
};

export default TabNav;
