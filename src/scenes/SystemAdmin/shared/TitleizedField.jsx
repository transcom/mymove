import React from 'react';
import { startCase } from 'lodash';

const TitleizedField = ({ source, record = {} }) => {
  return <span>{startCase(record[source])}</span>;
};

export default TitleizedField;
