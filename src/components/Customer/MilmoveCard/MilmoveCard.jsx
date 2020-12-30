import React from 'react';
import { CardGroup, Card, CardHeader, CardBody } from '@trussworks/react-uswds';
import { node, string } from 'prop-types';

import * as styles from './MilmoveCard.module.scss';

const MilmoveCard = ({ headerText, children }) => (
  <CardGroup>
    <Card containerProps={{ className: styles.CardContainer }}>
      <CardHeader>
        <h3>{headerText}</h3>
      </CardHeader>
      <CardBody className={styles.CardBody}>{children}</CardBody>
    </Card>
  </CardGroup>
);

MilmoveCard.propTypes = {
  headerText: string.isRequired,
  children: node.isRequired,
};

/*
const PPMExplanation = () => (
  <div className={styles.Section}>
    <h3>Hold on to things you’ll need quickly</h3>
    <img src={PPMShipmentImg} alt="PPM Shipment" />
    <p>Hand-carry important documents — ID, medical info, orders, school records, etc.</p>
    <p>
      Pack a set of things that you’ll need when you arrive — clothes, electronics, chargers, cleaning supplies, etc.
      Valuables that can’t be replaced are also a good idea.
    </p>
    <p>To be paid for moving these things, select a PPM shipment.</p>
  </div>
);

const HHGExplanation = () => (
  <div className={styles.Section}>
    <h3>One move, several parts</h3>
    <img src={HHGShipmentImg} alt="HHG Shipment" />
    <p>It’s common to move a few things yourself and have professional movers pack and move the rest.</p>
    <p>
      You can have things picked up or delivered to more than one place — your home and an office, for example. But
      multiple shipments make it easier to go over weight and end up paying for part of your move yourself.
    </p>
  </div>
);

const MoveCounselorExplanation = () => (
  <div className={styles.Section}>
    <h3>Talk to your move counselor</h3>
    <img src={MoveCounselorImg} alt="Move counselor meeting" />
    <p>
      A session with a move counselor is free. Counselors have a lot of experience with military moves and can steer you
      through complicated situations.
    </p>
    <p className={styles.CounselorListHeader}>Your counselor can identify:</p>
    <ul>
      <li>belongings that won&apos;t count against your weight allowance</li>
      <li>excess weight, excess distance, and other things that can cost you money</li>
      <li>things to make your move easier</li>
    </ul>
  </div>
);

const MoversExplanation = () => (
  <div className={styles.Section}>
    <h3>Talk to your movers</h3>
    <img src={MovingTruckImg} alt="Moving truck" />
    <p>
      If you have any shipments using professional movers, you&apos;ll be referred to a point of contact for your move.
    </p>
    <p>When things get complicated or you have questions during your move, they are there to help.</p>
    <p>
      It’s OK if things change after you submit your move info. Your movers or your counselor will make things work.
    </p>
  </div>
);
*/
export default MilmoveCard;
