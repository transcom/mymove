import React, { useEffect, useState } from 'react';
import { Radio } from '@trussworks/react-uswds';
import { Field, useFormikContext } from 'formik';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';

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
        scale={2} // digits after point, 0 for integers
        signed={false} // disallow negative
        radix="." // fractional delimiter
        mapToRadix={['.']} // symbols to process as radix
        padFractionalZeros // if true, then pads zeros at end to the length of scale
        thousandsSeparator=","
        lazy={false} // immediate masking evaluation
        prefix="$"
      />
      <Hint>
        Note: Any carrier insurance purchased is not a reimbursable expense. Do not add carrier insurance to the total
        above.
      </Hint>
      <TextField label="Tracking number" name="trackingNumber" id="trackingNumber" />
      <legend className="usa-label">Was this pro-gear?</legend>
      <div>
        <Field
          as={Radio}
          id="proGearYes"
          label="Yes"
          name="isProGear"
          value="true"
          checked={values.isProGear === 'true'}
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
          <legend className="usa-label">Who does this pro-gear belong to?</legend>
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
          <TextField label="Brief description of the pro-gear" name="proGearDescription" id="proGearDescription" />
          <MaskedTextField
            name="weightShipped"
            label="Pro-gear weight"
            data-testid="proGearWeight"
            id="weightShipped"
            mask={Number}
            scale={0}
            signed={false}
            thousandsSeparator=","
            lazy={false}
            suffix="lbs"
          />
        </>
      ) : (
        <MaskedTextField
          name="weightShipped"
          label="Weight shipped"
          data-testid="weightShipped"
          id="weightShipped"
          mask={Number}
          scale={0}
          signed={false}
          thousandsSeparator=","
          lazy={false}
          suffix="lbs"
        />
      )}
    </>
  );
};

export default SmallPackageForm;
