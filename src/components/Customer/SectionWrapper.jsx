/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './SectionWrapper.module.scss';

const SectionWrapper = ({ children, className, ...props }) => (
  <div className={classnames(styles.sectionWrapper, className)} {...props}>
    {children}
  </div>
);

SectionWrapper.propTypes = {
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
};

SectionWrapper.defaultProps = {
  className: '',
};

export default SectionWrapper;
