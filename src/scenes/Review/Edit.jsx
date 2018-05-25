import React from 'react';
import Summary from './Summary';

export default () => (
  <div className="usa-grid">
    <div className="usa-width-one-whole">
      <h1>Edit Move</h1>
      <p>
        Changes to your orders or shipments could impact your move, including
        the estimated PPM incentive.
      </p>
      <Summary />
    </div>
  </div>
);
