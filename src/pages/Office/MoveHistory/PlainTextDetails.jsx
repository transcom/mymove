import React from 'react';
import PropTypes from 'prop-types';

import { HistoryLogValuesShape, eventNamePlainTextToDisplay } from 'constants/historyLogUIDisplayName';

const PlainTextDetails = ({ eventName, changedValues }) => {
  let textToDisplay = '';
  if (eventNamePlainTextToDisplay[eventName]) {
    textToDisplay = eventNamePlainTextToDisplay[eventName](changedValues);
  }
  return <div>{textToDisplay}</div>;
};

PlainTextDetails.propTypes = {
  eventName: PropTypes.string,
  changedValues: HistoryLogValuesShape,
};

PlainTextDetails.defaultProps = {
  eventName: '',
  changedValues: [],
};

export default PlainTextDetails;
