import React from 'react';
import { Card } from '@trussworks/react-uswds';
import { bool, node } from 'prop-types';

import styles from './MilmoveCard.module.scss';

const MilmoveCard = ({ children, headerFirst }) => (
  <Card className={styles.MilmoveCard} containerProps={{ className: styles.CardContainer }} headerFirst={headerFirst}>
    {children}
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
