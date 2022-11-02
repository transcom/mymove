import React from 'react';
import * as PropTypes from 'prop-types';

import styles from './PreviewRow.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

/** This component renders a row for display in a eval report preview.
 * For the most part it is just doing some style overriding of description list items so that the displayed label/value pairs match the designs for the eval report */
const PreviewRow = ({ isShown, label, data }) => {
  if (!isShown) return null;

  return (
    <div className={descriptionListStyles.row}>
      <dt className={styles.label}>{label}</dt>
      <dd className={styles.data}>{data}</dd>
    </div>
  );
};

export default PreviewRow;

PreviewRow.propTypes = {
  isShown: PropTypes.bool,
  label: PropTypes.node,
  data: PropTypes.node,
};

PreviewRow.defaultProps = {
  isShown: true,
  label: '',
  data: '',
};
