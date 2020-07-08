import React from 'react';
import { Radio } from '@trussworks/react-uswds';

import { SHIPMENT_TYPE } from 'shared/constants';

const SelectMoveType = () => (
  <div className="grid-container usa-prose">
    <div className="usa-grid">
      <div className="grid-row grid-gap">
        <h1 className="sm-heading">How do you want to move your belongings?</h1>
        <div className="grid-col-9 desktop:grid-col-12">
          <Radio
            id={SHIPMENT_TYPE.PPM}
            label="Arrange it all yourself"
            value={SHIPMENT_TYPE.PPM}
            name="moveType"
            defaultChecked
          />
        </div>
      </div>
      <div className="grid-row grid-gap">
        <div className="grid-col-9 desktop:grid-col-12">
          <Radio
            id={SHIPMENT_TYPE.HHG}
            label="Have professionals pack and move it all"
            value={SHIPMENT_TYPE.HHG}
            name="moveType"
            disabled
          />
        </div>
      </div>
    </div>
  </div>
);

export default SelectMoveType;
