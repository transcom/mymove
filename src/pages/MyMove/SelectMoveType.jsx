import React from 'react';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_OPTIONS } from 'shared/constants';

export const SelectMoveType = () => (
  <div className="usa-grid">
    <div className="grid-row">
      <div className="grid-col">
        <h1 className="sm-heading">How do you want to move your belongings?</h1>
        <Radio
          id={SHIPMENT_OPTIONS.PPM}
          label="Arrange it all yourself"
          value={SHIPMENT_OPTIONS.PPM}
          name="moveType"
          defaultChecked
        />
        <Radio
          id={SHIPMENT_OPTIONS.HHG}
          label="Have professionals pack and move it all"
          value={SHIPMENT_OPTIONS.HHG}
          name="moveType"
          disabled
        />
      </div>
    </div>
  </div>
);

export default SelectMoveType;
