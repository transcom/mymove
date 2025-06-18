import React, { useState } from 'react';

import styles from './FeatureFlags.module.scss'; // Importing styles from SCSS module

import { BOOLEAN_FLAG_TYPE, VARIANT_FLAG_TYPE } from 'components/FeatureFlag/FeatureFlag';
import { getBooleanFeatureFlagForUser, getVariantFeatureFlagForUser } from 'services/internalApi';
import { milmoveLogger } from 'utils/milmoveLog';

const displayVariant = (result) => {
  if ('variant' in result) {
    return result.variant === '' ? '<none>' : result.variant;
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
      .then((result) => setFlagResult(result))
      .catch((error) => {
        milmoveLogger.error(error);
        setFlagResult(null);
      });
  };

  const handleTypeChange = (e) => {
    setFlagType(e.target.value);
  };

  const handleKeyChange = (e) => {
    setFlagKey(e.target.value);
  };

  return (
    <div className={styles.container}>
      <p>Show feature flag for current user. This tests the feature flag API.</p>
      <hr />
      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.formGroup}>
          <label htmlFor="adminFeatureFlagType" className={styles.label}>
            Feature Flag Type
          </label>
          <select
            id="adminFeatureFlagType"
            name="featureFlagType"
            value={flagType}
            onChange={handleTypeChange}
            className={styles.select}
          >
            <option value={BOOLEAN_FLAG_TYPE}>Boolean</option>
            <option value={VARIANT_FLAG_TYPE}>Variant</option>
          </select>
        </div>
        <div className={styles.formGroup}>
          <label htmlFor="adminFeatureFlagKey" className={styles.label}>
            Feature Flag Key
          </label>
          <input
            id="adminFeatureFlagKey"
            name="featureFlagKey"
            type="text"
            value={flagKey || ''}
            onChange={handleKeyChange}
            className={styles.input}
          />
        </div>
        <button type="submit" className={styles.button} aria-label="Evaluate Feature Flag">
          Evaluate
        </button>
      </form>
      <hr />
      {flagResult && (
        <div className={styles.result}>
          <div className={styles.resultRow}>
            <strong>Entity:</strong>
            <span>{flagResult.entity}</span>
          </div>
          <div className={styles.resultRow}>
            <strong>Key:</strong>
            <span>{flagResult.key}</span>
          </div>
          <div className={styles.resultRow}>
            <strong>Match:</strong>
            <span>{flagResult.match.toString()}</span>
          </div>
          <div className={styles.resultRow}>
            <strong>Value:</strong>
            <span>{displayVariant(flagResult)}</span>
          </div>
          <div className={styles.resultRow}>
            <strong>Namespace:</strong>
            <span>{flagResult.namespace}</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default FeatureFlagEvaluate;
