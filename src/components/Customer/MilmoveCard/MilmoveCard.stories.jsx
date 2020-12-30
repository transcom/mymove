import React from 'react';

import PPMShipmentImg from '../../../images/ppm-shipment.jpg';

import MilmoveCard from './MilmoveCard';

export default {
  title: 'Components/MilmoveCard',
  component: MilmoveCard,
};

export const Basic = () => (
  <div style={{ maxWidth: 675 }}>
    <MilmoveCard headerText="Hold on to things you’ll need quickly">
      <img src={PPMShipmentImg} alt="PPM Shipment" />
      <p>Hand-carry important documents — ID, medical info, orders, school records, etc.</p>
      <p>
        Pack a set of things that you’ll need when you arrive — clothes, electronics, chargers, cleaning supplies, etc.
        Valuables that can’t be replaced are also a good idea.
      </p>
      <p>To be paid for moving these things, select a PPM shipment.</p>
    </MilmoveCard>
  </div>
);
