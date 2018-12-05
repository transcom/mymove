import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import './index.css';

export const ProgressTimelineStep = function(props) {
  let classes = classNames({
    step: true,
    completed: props.completed,
    current: props.current,
  });

  return (
    <div className={classes}>
      <div className="dot" />
      <div className="name">{props.name}</div>
    </div>
  );
};

ProgressTimelineStep.propTypes = {
  completed: PropTypes.bool,
  current: PropTypes.bool,
};

// ProgressTimeline renders a subway-map-style timeline. Use ProgressTimelineStep
// components as children to declaritively define the "stops" and their status.
export const ProgressTimeline = function(props) {
  return <div className="progress-timeline">{props.children}</div>;
};

ProgressTimeline.propTypes = {
  children: PropTypes.arrayOf(ProgressTimelineStep).isRequired,
};
