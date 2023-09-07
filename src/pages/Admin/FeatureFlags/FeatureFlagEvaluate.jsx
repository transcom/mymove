import React, { useState } from 'react';

import { BOOLEAN_FLAG_TYPE, VARIANT_FLAG_TYPE } from 'components/FeatureFlag/FeatureFlag';
import { getBooleanFeatureFlagForUser, getVariantFeatureFlagForUser } from 'services/internalApi';
import { milmoveLogger } from 'utils/milmoveLog';

const displayVariant = (result) => {
  if ('variant' in result) {
    if (result.variant === '') {
      return '<none>';
    }
    return result.variant;
  }
  return 'Boolean Flag';
};

export const FeatureFlagEvaluate = () => {
  const [flagType, setFlagType] = useState(BOOLEAN_FLAG_TYPE);
  const [flagKey, setFlagKey] = useState(null);
  const [flagResult, setFlagResult] = useState(null);

  const handleSubmit = (e) => {
    e.preventDefault();
    const getFeatureFlagForUser =
      flagType === BOOLEAN_FLAG_TYPE ? getBooleanFeatureFlagForUser : getVariantFeatureFlagForUser;

    getFeatureFlagForUser(flagKey, {})
      .then((result) => {
        setFlagResult(result);
      })
      .catch((error) => {
        milmoveLogger.error(error);
        setFlagResult(null);
      });
  };

  const handleTypeChange = (e) => {
    if (e.target.value === VARIANT_FLAG_TYPE) {
      setFlagType(VARIANT_FLAG_TYPE);
    } else {
      setFlagType(BOOLEAN_FLAG_TYPE);
    }
  };

  const handleKeyChange = (e) => {
    setFlagKey(e.target.value);
  };

  return (
    <div>
      <span>Show feature flag for current user. This tests the feature flag API.</span>
      <hr />
      <div>
        <form onSubmit={handleSubmit}>
          <div>
            <label htmlFor="adminFeatureFlagType">Feature Flag Type</label>
            <select onChange={handleTypeChange} name="featureFlagType" id="adminFeatureFlagType">
              <option value={BOOLEAN_FLAG_TYPE}>Boolean</option>
              <option value={VARIANT_FLAG_TYPE}>Variant</option>
            </select>
          </div>
          <br />
          <div>
            <label htmlFor="adminFeatureFlagKey">Feature Flag Key</label>
            <input onChange={handleKeyChange} name="featureFlagKey" id="adminFeatureFlagKey" type="text" />
          </div>
          <br />
          <button type="submit">Evaluate</button>
        </form>
      </div>
      <hr />
      <div>
        {flagResult && (
          <dl>
            <dt>Entity</dt>
            <dd>{flagResult.entity}</dd>
            <dt>Key</dt>
            <dd>{flagResult.key}</dd>
            <dt>Match</dt>
            <dd>{flagResult.match.toString()}</dd>
            <dt>Value</dt>
            <dd>{displayVariant(flagResult)}</dd>
            <dt>Namespace</dt>
            <dd>{flagResult.namespace}</dd>
          </dl>
        )}
      </div>
    </div>
  );
};

export default FeatureFlagEvaluate;
