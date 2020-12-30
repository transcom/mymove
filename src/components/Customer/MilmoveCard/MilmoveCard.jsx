import React from 'react';
import { Card } from '@trussworks/react-uswds';
import { bool, node } from 'prop-types';

import * as styles from './MilmoveCard.module.scss';

const MilmoveCard = ({ children, headerFirst }) => (
  <Card containerProps={{ className: styles.CardContainer }} headerFirst={headerFirst}>
    <>{children}</>
  </Card>
);

MilmoveCard.propTypes = {
  children: node.isRequired,
  headerFirst: bool,
};

MilmoveCard.defaultProps = {
  headerFirst: true,
};

export default MilmoveCard;
