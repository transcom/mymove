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
          label="I’ll move things myself"
          value={SHIPMENT_OPTIONS.PPM}
          name="moveType"
          defaultChecked
        />
        <ul>
          <li>This is a PPM - “personally procured move”</li>
          <li>You arrange to move some or all of your belongings</li>
          <li>The government pays you an incentive based on weight</li>
          <li>DIY or hire your own movers</li>
        </ul>
        <Radio
          id={SHIPMENT_OPTIONS.HHG}
          label="The government packs for me and moves me"
          value={SHIPMENT_OPTIONS.HHG}
          name="moveType"
          disabled
        />
        <ul>
          <li>This is an HHG shipment — “household goods”</li>
          <li>The most popular kind of shipment</li>
          <li>Professional movers take care of the whole shipment</li>
          <li>They pack and move it for you</li>
        </ul>
      </div>
    </div>
  </div>
);

export default SelectMoveType;
