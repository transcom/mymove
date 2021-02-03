import React from 'react';
import { string } from 'prop-types';

import styles from './index.module.scss';

const BypassBlock = ({ anchorLink }) => (
  <div>
    <a className={styles.bypassBlock} href={anchorLink}>
      Skip to content
    </a>
  </div>
);

BypassBlock.propTypes = {
  anchorLink: string,
};

BypassBlock.defaultProps = {
  anchorLink: '#main',
};

export default BypassBlock;
