import React from 'react';
import PropTypes from 'prop-types';

import {
  HistoryLogValuesShape,
  eventNamePlainTextToDisplay,
  HistoryLogContextShape,
} from 'constants/historyLogUIDisplayName';

const PlainTextDetails = ({ eventName, oldValues, changedValues, context }) => {
  let textToDisplay = '';
  if (eventNamePlainTextToDisplay[eventName]) {
    textToDisplay = eventNamePlainTextToDisplay[eventName](changedValues, oldValues, context);
  }
  return <div>{textToDisplay}</div>;
};

PlainTextDetails.propTypes = {
  eventName: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  oldValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
};

PlainTextDetails.defaultProps = {
  eventName: '',
  changedValues: {},
  oldValues: {},
  context: {},
};

export default PlainTextDetails;
