import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './SimpleSection.module.scss';

const SimpleSection = ({ children, header, border }) => {
  return (
    <section className={classNames(styles.SimpleSection, { [styles.noBorder]: !border })}>
      <header>{header}</header>
      {children}
    </section>
  );
};

SimpleSection.propTypes = {
  children: PropTypes.node.isRequired,
  header: PropTypes.oneOfType([PropTypes.string, PropTypes.node]).isRequired,
  border: PropTypes.bool,
};

SimpleSection.defaultProps = {
  border: false,
};

export default SimpleSection;
