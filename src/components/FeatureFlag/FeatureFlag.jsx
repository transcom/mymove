import React from 'react';
import PropTypes from 'prop-types';

import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import { getFeatureFlagForUser } from 'services/internalApi';

// Example of how we might have a FeatureFlag component
// This is probably not production ready
export const FeatureFlag = ({ flagKey, flagContext, render }) => {
  const [flagValue, setFlagValue] = React.useState('');

  React.useEffect(() => {
    getFeatureFlagForUser(flagKey, flagContext)
      .then((result) => {
        if (!result.enabled) {
          setFlagValue('disabled');
        } else {
          setFlagValue(result.value);
        }
      })
      .catch((error) => {
        milmoveLog(MILMOVE_LOG_LEVEL.ERROR, error);
        setFlagValue('disabled');
      });
  });

  return render(flagValue);
};

FeatureFlag.propTypes = {
  flagKey: PropTypes.string.isRequired,
  flagContext: PropTypes.object,
  render: PropTypes.func,
};

FeatureFlag.defaultProps = {
  flagKey: '',
  flagContext: {},
  render: null,
};

export default FeatureFlag;
