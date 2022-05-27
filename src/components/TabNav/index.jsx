import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './index.module.scss';

const TabNav = ({ items, role }) => (
  <nav className={classNames(styles.tabNav)} role={role}>
    <div>
      <ul className={classNames(styles.tabList)}>
        {items.map((item /* , index */) => (
          <li key={item.id} className={classNames(styles.tabItem)}>
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

TabNav.defaultProps = {
  role: null,
};

export default TabNav;
