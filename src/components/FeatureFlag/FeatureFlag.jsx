import React from 'react';
import PropTypes from 'prop-types';

import { milmoveLogger } from 'utils/milmoveLog';
import { getFeatureFlagForUser } from 'services/internalApi';

export const ENABLED_VALUE = 'enabled';
export const DISABLED_VALUE = 'disabled';

export function featureIsEnabled(val) {
  return val && val === ENABLED_VALUE;
}

// Example of how we might have a FeatureFlag component
export const FeatureFlag = ({ flagKey, flagContext, render }) => {
  const [flagValue, setFlagValue] = React.useState('');

  React.useEffect(() => {
    getFeatureFlagForUser(flagKey, flagContext)
      .then((result) => {
        if (result.match) {
          setFlagValue(result.value);
        } else {
          setFlagValue(DISABLED_VALUE);
        }
      })
      .catch((error) => {
        milmoveLogger.error(error);
        setFlagValue(DISABLED_VALUE);
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
  flagContext: {},
  render: null,
};

export default FeatureFlag;
