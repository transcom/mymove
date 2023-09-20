import React from 'react';
import PropTypes from 'prop-types';

import { milmoveLogger } from 'utils/milmoveLog';
import { getBooleanFeatureFlagForUser, getVariantFeatureFlagForUser } from 'services/internalApi';

export const BOOLEAN_FLAG_TYPE = 'boolean';
export const VARIANT_FLAG_TYPE = 'variant';

// Example of how we might have a FeatureFlag component
export const FeatureFlag = ({ flagType, flagKey, flagContext, render }) => {
  const [flagValue, setFlagValue] = React.useState('');

  const getFeatureFlagForUser =
    flagType === BOOLEAN_FLAG_TYPE ? getBooleanFeatureFlagForUser : getVariantFeatureFlagForUser;
  const setFlagTypeValue = (result) => {
    switch (flagType) {
      case BOOLEAN_FLAG_TYPE:
        // always set the value to a string, even for boolean flag type
        setFlagValue(result.match.toString());
        break;
      default:
        if (result.match) {
          setFlagValue(result.variant);
        } else {
          setFlagValue('');
        }
    }
  };

  React.useEffect(() => {
    getFeatureFlagForUser(flagKey, flagContext)
      .then((result) => {
        setFlagTypeValue(result);
      })
      .catch((error) => {
        milmoveLogger.error(error);
        setFlagValue('');
      });
  });

  return render(flagValue);
};

FeatureFlag.propTypes = {
  flagType: PropTypes.oneOf([BOOLEAN_FLAG_TYPE, VARIANT_FLAG_TYPE]),
  flagKey: PropTypes.string.isRequired,
  flagContext: PropTypes.object,
  render: PropTypes.func,
};

FeatureFlag.defaultProps = {
  flagContext: {},
  render: null,
};

export default FeatureFlag;
