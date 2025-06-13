import React, { useEffect, useState } from 'react';
import { Radio } from '@trussworks/react-uswds';
import { Field, useFormikContext } from 'formik';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';
import RequiredAsterisk from 'components/form/RequiredAsterisk';

const SmallPackageForm = () => {
  const { values } = useFormikContext();
  const [showProGear, setShowProGear] = useState(values.isProGear === 'true');

  useEffect(() => {
    if (values.isProGear === 'true') {
      setShowProGear(true);
    } else {
      setShowProGear(false);
    }
  }, [values.isProGear]);

  return (
    <>
      <MaskedTextField
        name="amount"
        label="Package shipment cost"
        id="amount"
        mask={Number}
        scale={2}
        signed={false}
        radix="."
        mapToRadix={['.']}
        padFractionalZeros
        thousandsSeparator=","
        lazy={false}
        prefix="$"
        showRequiredAsterisk
        required
      />
      <Hint>
        Note: Any carrier insurance purchased is not a reimbursable expense. Do not add carrier insurance to the total
        above.
      </Hint>
      <TextField label="Tracking number" name="trackingNumber" id="trackingNumber" showRequiredAsterisk required />
      <legend className="usa-label" aria-label="Required: Was this pro-gear?">
        <span required>
          Was this pro-gear? <RequiredAsterisk />
        </span>
      </legend>
      <div>
        <Field
          as={Radio}
          id="proGearYes"
          label="Yes"
          name="isProGear"
          value="true"
          checked={values.isProGear === 'true'}
          showRequiredAsterisk
          required
        />
        <Field
          as={Radio}
          id="proGearNo"
          label="No"
          name="isProGear"
          value="false"
          checked={values.isProGear === 'false'}
        />
      </div>
      {showProGear ? (
        <>
          <legend className="usa-label" aria-label="Required: Who does this pro-gear belong to?">
            <span required>
              Who does this pro-gear belong to? <RequiredAsterisk />
            </span>
          </legend>
          <div>
            <Field
              as={Radio}
              id="proGearSelf"
              label="Me"
              name="proGearBelongsToSelf"
              value="true"
              checked={values.proGearBelongsToSelf === 'true'}
            />
            <Field
              as={Radio}
              id="proGearSpouse"
              label="My Spouse"
              name="proGearBelongsToSelf"
              value="false"
              checked={values.proGearBelongsToSelf === 'false'}
            />
          </div>
          <TextField
            label="Brief description of the pro-gear"
            name="proGearDescription"
            id="proGearDescription"
            showRequiredAsterisk
            required
          />
          <MaskedTextField
            name="weightShipped"
            label="Package weight"
            data-testid="proGearWeight"
            id="proGearWeight"
            mask={Number}
            scale={0}
            thousandsSeparator=","
            lazy={false}
            suffix="lbs"
            showRequiredAsterisk
            required
          />
        </>
      ) : (
        <MaskedTextField
          name="weightShipped"
          label="Package weight"
          data-testid="weightShipped"
          id="weightShipped"
          mask={Number}
          scale={0}
          thousandsSeparator=","
          lazy={false}
          suffix="lbs"
          showRequiredAsterisk
          required
        />
      )}
    </>
  );
};

export default SmallPackageForm;
