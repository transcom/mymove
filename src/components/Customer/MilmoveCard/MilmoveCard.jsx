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

export default MilmoveCard;
