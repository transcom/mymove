import React from 'react';
import PropTypes from 'prop-types';

import styles from './SectionWrapper.module.scss';

const SectionWrapper = ({ children }) => <div className={styles.sectionWrapper}>{children}</div>;

SectionWrapper.propTypes = {
  children: PropTypes.oneOfType([PropTypes.arrayOf(PropTypes.element), PropTypes.element.isRequired]),
};

SectionWrapper.defaultProps = {
  children: ' ',
};

export default SectionWrapper;
