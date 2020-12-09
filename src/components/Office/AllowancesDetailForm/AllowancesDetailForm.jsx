import React from 'react';

import { TextMaskedInput } from 'components/form/fields/TextInput';
import styles from 'components/Office/AllowancesDetailForm/AllowancesDetailForm.module.scss';

const AllowancesDetailForm = () => {
  return (
    <div className={styles.AllowancesDetailForm}>
      <TextMaskedInput
        name="authorizedWeight"
        label="Authorized weight"
        id="authorizedWeightInput"
        mask="NUM lbs" // Nested masking imaskjs
        lazy={false} // immediate masking evaluation
        blocks={{
          // our custom masking key
          NUM: {
            mask: Number,
            thousandsSeparator: ',',
            scale: 0, // whole numbers
            signed: false, // positive numbers
          },
        }}
      />
    </div>
  );
};

AllowancesDetailForm.propTypes = {};

export default AllowancesDetailForm;
