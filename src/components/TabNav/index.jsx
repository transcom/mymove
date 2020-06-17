import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import { Tab, Tabs, TabList } from 'react-tabs';
import classNames from 'classnames/bind';

import styles from './index.module.scss';

const cx = classNames.bind(styles);

const TabNav = ({ options, children }) => (
  <Tabs>
    <TabList className={cx('tab-nav')}>
      {options.map(({ title, notice }, index) => (
        <Tab key={index.toString()} selectedClassName={cx('tab-active')} className={cx('tab-item')}>
          <span className={cx('tab-title')}>{title}</span>
          {notice && <Tag>{notice}</Tag>}
        </Tab>
      ))}
    </TabList>
    {children}
  </Tabs>
);

TabNav.propTypes = {
  options: PropTypes.arrayOf(
    PropTypes.shape({
      title: PropTypes.string,
      notice: PropTypes.string,
    }),
  ).isRequired,
  // eslint-disable-next-line react/require-default-props
  children: (props, propName, componentName) => {
    // eslint-disable-next-line security/detect-object-injection
    const prop = props[propName];
    let error;

    if (React.Children.count(prop) === 0) {
      error = new Error(`\`${componentName}\` requires Children.`);
    }

    React.Children.forEach(prop, (el) => {
      if (error) return;
      if (el.type.name !== 'TabNavPanel') {
        error = new Error(`\`${componentName}\` children must be \`TabNavPanel\`.`);
      }
    });

    return error;
  },
};

export default TabNav;
