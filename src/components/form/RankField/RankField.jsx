import React, { useEffect } from 'react';
import { useFormikContext } from 'formik';
import { Label } from '@trussworks/react-uswds';

import RequiredAsterisk from '../RequiredAsterisk';
import { DropdownInput } from '../fields';

import styles from './RankField.module.scss';

const RankField = ({ rankOptions, handleChange }) => {
  const { setFieldValue } = useFormikContext();

  useEffect(() => {
    if (rankOptions.length === 1) {
      setFieldValue('rank', rankOptions[0].key);
    }
  }, [rankOptions, setFieldValue]);
  if (rankOptions.length > 1) {
    return (
      <DropdownInput
        label="Rank"
        name="rank"
        id="rank"
        required
        options={rankOptions}
        showRequiredAsterisk
        onChange={(e) => handleChange(e)}
        data-testid="RankDropdown"
      />
    );
  }
  return (
    <div className={styles.rankContainer}>
      <Label>
        <span>
          Rank <RequiredAsterisk />
        </span>
      </Label>
      <span className={styles.rankText}>{rankOptions[0]?.value}</span>
    </div>
  );
};

export default RankField;
