import React from 'react';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';

const SmallPackageForm = () => {
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
      <MaskedTextField
        name="weightShipped"
        label="Weight shipped"
        data-testid="estimatedWeight"
        id="weightShipped"
        mask={Number}
        scale={0}
        signed={false}
        thousandsSeparator=","
        lazy={false}
        suffix="lbs"
      />
      <TextField label="Tracking number" name="trackingNumber" id="trackingNumber" labelHint="Optional" />
    </>
  );
};

export default SmallPackageForm;
