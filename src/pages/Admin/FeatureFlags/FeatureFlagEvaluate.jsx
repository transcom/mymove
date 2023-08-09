import React from 'react';

import { getFeatureFlagForUser } from 'services/internalApi';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

export const FeatureFlagEvaluate = () => {
  const [flagKey, setFlagKey] = React.useState(null);
  const [flagResult, setFlagResult] = React.useState(null);
  const [shouldEvaluate, setShouldEvaluate] = React.useState(false);

  const handleSubmit = (e) => {
    e.preventDefault();
    setShouldEvaluate(true);
  };

  const handleChange = (e) => {
    setFlagKey(e.target.value);
  };

  React.useEffect(() => {
    if (shouldEvaluate) {
      getFeatureFlagForUser(flagKey, {})
        .then((result) => {
          setFlagResult(result);
        })
        .catch((error) => {
          milmoveLog(MILMOVE_LOG_LEVEL.ERROR, error);
          setFlagResult(null);
        });
      setShouldEvaluate(false);
    }
  }, [shouldEvaluate, flagKey]);

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
