import React from 'react';
import PropTypes from 'prop-types';

import { HistoryLogValuesShape, eventNamePlainTextToDisplay } from 'constants/historyLogUIDisplayName';

const PlainTextDetails = ({ eventName, changedValues, oldValues }) => {
  let textToDisplay = '';
  if (eventNamePlainTextToDisplay[eventName]) {
    textToDisplay = eventNamePlainTextToDisplay[eventName](changedValues, oldValues);
  }
  return <div>{textToDisplay}</div>;
};

PlainTextDetails.propTypes = {
  eventName: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  oldValues: HistoryLogValuesShape,
};

PlainTextDetails.defaultProps = {
  eventName: '',
  changedValues: {},
  oldValues: {},
};

export default PlainTextDetails;
