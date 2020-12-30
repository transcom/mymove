import React from 'react';
import { CardHeader, CardGroup, CardMedia, CardBody } from '@trussworks/react-uswds';

import PPMShipmentImg from '../../../images/ppm-shipment.jpg';

import MilmoveCard from './MilmoveCard';

export default {
  title: 'Components/MilmoveCard',
  component: MilmoveCard,
};

export const IndentedMilmoveCard = () => (
  <div style={{ maxWidth: 675 }}>
    <CardGroup>
      <MilmoveCard>
        <CardHeader>
          <h3>Hold on to things you’ll need quickly</h3>
        </CardHeader>
        <CardMedia inset>
          <img src={PPMShipmentImg} alt="PPM Shipment" />
        </CardMedia>
        <CardBody>
          <p>Hand-carry important documents — ID, medical info, orders, school records, etc.</p>
          <p>
            Pack a set of things that you’ll need when you arrive — clothes, electronics, chargers, cleaning supplies,
            etc. Valuables that can’t be replaced are also a good idea.
          </p>
          <p>To be paid for moving these things, select a PPM shipment.</p>
        </CardBody>
      </MilmoveCard>
    </CardGroup>
  </div>
);

export const ExdentedMilmoveCard = () => (
  <div style={{ maxWidth: 675 }}>
    <CardGroup>
      <MilmoveCard>
        <CardHeader exdent>
          <h3>Hold on to things you’ll need quickly</h3>
        </CardHeader>
        <CardMedia exdent>
          <img src={PPMShipmentImg} alt="PPM Shipment" />
        </CardMedia>
        <CardBody exdent>
          <p>Hand-carry important documents — ID, medical info, orders, school records, etc.</p>
          <p>
            Pack a set of things that you’ll need when you arrive — clothes, electronics, chargers, cleaning supplies,
            etc. Valuables that can’t be replaced are also a good idea.
          </p>
          <p>To be paid for moving these things, select a PPM shipment.</p>
        </CardBody>
      </MilmoveCard>
    </CardGroup>
  </div>
);
