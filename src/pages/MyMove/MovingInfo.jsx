import React from 'react';

// eslint-disable-next-line react/prefer-stateless-function
export const MovingInfo = () => {
  return (
    <>
      <h1 data-testid="shipmentsHeader" className="sm-heading">
        Figure out your shipments
      </h1>
      <p>Handy tips as you decide how to move</p>
      <h2 data-testid="shipmentsSubHeader" className="sm-heading">
        Move in one shipment or more
      </h2>
      <p>
        It’s common to move in a few shipments. Everything can go in one batch, or you can divide your belongings into
        several shipments.
      </p>
      <h2 data-testid="shipmentsSubHeader" className="sm-heading">
        Keep important things with you
      </h2>
      <p>
        It’s a good idea to move things you’ll need right away and prized possessions yourself. Select a PPM (personally
        procured move) shipment to do that.
      </p>
      <h2 data-testid="shipmentsSubHeader" className="sm-heading">
        Spread out your pickup dates
      </h2>
      <p>It’s easier to coordinate multiple shipments if you don’t schedule all the pickups on the same day.</p>
    </>
  );
};

export default MovingInfo;
