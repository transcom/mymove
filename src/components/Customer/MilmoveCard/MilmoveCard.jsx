import React from 'react';
import { Card, CardHeader, CardBody } from '@trussworks/react-uswds';

const MilmoveCard = () => (
  <Card gridLayout={{ tablet: { col: 6 } }}>
    <CardHeader>
      <h2 className="usa-card__heading">Hold on to things you’ll need quickly</h2>
    </CardHeader>
    <CardBody>
      <p> This is a standard card with a button in the footer. </p>
    </CardBody>
  </Card>
);

export default MilmoveCard;
