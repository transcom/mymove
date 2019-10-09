import React from 'react';
import { startCase } from 'lodash';

/* eslint-disable security/detect-object-injection */
const TitleizedField = ({ source, record = {} }) => {
  return <span>{startCase(record[source])}</span>;
};
/* eslint-enable security/detect-object-injection */

export default TitleizedField;
