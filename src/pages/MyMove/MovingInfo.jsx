/* eslint-disable camelcase */
import React from 'react';

import SectionWrapper from 'components/Customer/SectionWrapper';
import ConnectedWizardPage from 'shared/WizardPage/index';
import { WizardPageShape } from 'types/customerShapes';
import { no_op } from 'shared/utils';

export const MovingInfo = ({ wizardProps }) => {
  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <ConnectedWizardPage handleSubmit={no_op} {...wizardProps} hideBackBtn showFinishLaterBtn>
      <h1 data-testid="shipmentsHeader">Figure out your shipments</h1>
      <p>Handy tips as you decide how to move</p>
      <SectionWrapper>
        <h2 data-testid="shipmentsSubHeader">Move in one shipment or more</h2>
        <p>
          It’s common to move in a few shipments. Everything can go in one batch, or you can divide your belongings into
          several shipments.
        </p>
        <hr />
        <h2 className="margin-top-2" data-testid="shipmentsSubHeader">
          Keep important things with you
        </h2>
        <p>
          It’s a good idea to move things you’ll need right away and prized possessions yourself. Select a PPM
          (personally procured move) shipment to do that.
        </p>
        <hr />
        <h2 className="margin-top-2" data-testid="shipmentsSubHeader">
          Spread out your pickup dates
        </h2>
        <p>It’s easier to coordinate multiple shipments if you don’t schedule all the pickups on the same day.</p>
      </SectionWrapper>
    </ConnectedWizardPage>
  );
};

MovingInfo.propTypes = {
  wizardProps: WizardPageShape.isRequired,
};

export default MovingInfo;
