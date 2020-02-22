import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import classNames from 'classnames/bind';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const TabNav = ({ options }) => (
  <div className={cx('tab-nav')}>
    {options.map(({ title, active, notice }, index) => (
      <span tabIndex={index.toString()} key={index.toString()} className={cx('tab-item')}>
        <span className={cx('tab-title', { 'tab-title-active': active })}>{title}</span>
        {notice && <Tag>{notice}</Tag>}
      </span>
    ))}
  </div>
);

TabNav.propTypes = {
  options: PropTypes.arrayOf(
    PropTypes.shape({
      title: PropTypes.string,
      active: PropTypes.bool,
      notice: PropTypes.string,
    }),
  ).isRequired,
};

export default TabNav;
