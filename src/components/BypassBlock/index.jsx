import React from 'react';
import styles from './index.module.scss';

const BypassBlock = ({ anchorLink }) => (
  <div className={styles.bypassBlock}>
    <a href={anchorLink} >Skip to Content</a>
  </div>
);

BypassBlock.defaultProps = {
  anchorLink: '#main',
};

BypassBlock.propTypes = {
  anchorLink: string,
};

export default BypassBlock;