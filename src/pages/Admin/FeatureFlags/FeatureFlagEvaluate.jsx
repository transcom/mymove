import React, { useState } from 'react';

import { getFeatureFlagForUser } from 'services/internalApi';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

export const FeatureFlagEvaluate = () => {
  const [flagKey, setFlagKey] = useState(null);
  const [flagResult, setFlagResult] = useState(null);

  const handleSubmit = (e) => {
    e.preventDefault();
    getFeatureFlagForUser(flagKey, {})
      .then((result) => {
        setFlagResult(result);
      })
      .catch((error) => {
        milmoveLog(MILMOVE_LOG_LEVEL.ERROR, error);
        setFlagResult(null);
      });
  };

  const handleChange = (e) => {
    setFlagKey(e.target.value);
  };

  return (
    <>
      <span>Show feature flag for current user. This tests the feature flag API.</span>
      <form onSubmit={handleSubmit}>
        <input onChange={handleChange} name="featureFlagKey" type="text" />
        <button type="submit">Evaluate</button>
      </form>
      {flagResult && (
        <dl>
          <dt>Entity</dt>
          <dd>{flagResult.entity}</dd>
          <dt>Key</dt>
          <dd>{flagResult.key}</dd>
          <dt>Match</dt>
          <dd>{flagResult.match ? 'true' : 'false'}</dd>
          <dt>Value</dt>
          <dd>{flagResult.value ? flagResult.value : '<none>'}</dd>
          <dt>Namespace</dt>
          <dd>{flagResult.namespace}</dd>
        </dl>
      )}
    </>
  );
};

export default FeatureFlagEvaluate;
