import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './index.module.scss';

const TabNav = ({ items, role, className }) => (
  <nav className={classNames(styles.tabNav, className)} role={role}>
    <div>
      <ul className={classNames(styles.tabList)}>
        {items.map((item, index) => (
          <li key={index.toString()} className={classNames(styles.tabItem)}>
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
  className: PropTypes.string,
};

TabNav.defaultProps = {
  role: null,
  className: null,
};

export default TabNav;
