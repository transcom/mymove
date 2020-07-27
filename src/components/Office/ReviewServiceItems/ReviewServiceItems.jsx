import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import sortServiceItemsByGroup from '../../../utils/serviceItems';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { SERVICE_ITEM_STATUS } from 'shared/constants';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import ServiceItemCard from 'components/Office/ReviewServiceItems/ServiceItemCard';
import { toDollarString } from 'shared/formatters';

const ReviewServiceItems = ({ header, serviceItemCards, handleClose }) => {
  const [curCardIndex, setCardIndex] = useState(0);
  const [sortedCards] = useState(sortServiceItemsByGroup(serviceItemCards));
  const totalCards = serviceItemCards.length;

  const { APPROVED, REJECTED } = SERVICE_ITEM_STATUS;

  const handleClick = (index) => {
    setCardIndex(index);
  };

  const calculateTotals = (values) => {
    let approvedSum = 0;
    let rejectedSum = 0;

    serviceItemCards.forEach((serviceItem) => {
      const itemValues = values[`${serviceItem.id}`];
      if (itemValues?.status === APPROVED) approvedSum += serviceItem.amount;
      else if (itemValues?.status === REJECTED) rejectedSum += serviceItem.amount;
    });

    return {
      approved: approvedSum,
      rejected: rejectedSum,
    };
  };

  //  let requestedSum = 0; // TODO - use in Complete review screen
  const formValues = {};
  // TODO - preset these based on existing values
  serviceItemCards.forEach((serviceItem) => {
    formValues[serviceItem.id] = {
      status: serviceItem.status,
      rejectionReason: serviceItem.rejectionReason,
    };

    // requestedSum += serviceItem.amount; // TODO - use in Complete review screen
  });

  const currentCard = sortedCards[parseInt(curCardIndex, 10)];

  return (
    <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
      <Formik initialValues={formValues}>
        {({ values, handleChange, setValues }) => {
          const clearServiceItemValues = (id) => {
            setValues({
              ...values,
              [`${id}`]: {
                status: undefined,
                rejectionReason: undefined,
              },
            });
          };

          return (
            <Form className={styles.form}>
              <div className={styles.top}>
                <Button data-testid="closeSidebar" type="button" onClick={handleClose} unstyled>
                  <XLightIcon />
                </Button>
                <div data-testid="itemCount" className={styles.eyebrowTitle}>
                  {curCardIndex + 1} OF {totalCards} ITEMS
                </div>
                <h2 className={styles.header}>{header}</h2>
              </div>
              <div className={styles.body}>
                <ServiceItemCard
                  key={`serviceItemCard_${currentCard.id}`}
                  // eslint-disable-next-line react/jsx-props-no-spreading
                  {...currentCard}
                  value={values[currentCard.id]}
                  onChange={handleChange}
                  clearValues={clearServiceItemValues}
                />
              </div>
              <div className={styles.bottom}>
                <Button
                  data-testid="prevServiceItem"
                  type="button"
                  onClick={() => handleClick(curCardIndex - 1)}
                  secondary
                  disabled={curCardIndex === 0}
                >
                  Previous
                </Button>
                <Button
                  data-testid="nextServiceItem"
                  type="button"
                  onClick={() => handleClick(curCardIndex + 1)}
                  disabled={curCardIndex + 1 === totalCards}
                >
                  Next
                </Button>
                <div className={styles.totalApproved}>
                  <div className={styles.totalLabel}>Total approved</div>
                  <div data-testid="approvedAmount" className={styles.totalAmount}>
                    {toDollarString(calculateTotals(values).approved)}
                  </div>
                </div>
              </div>
            </Form>
          );
        }}
      </Formik>
    </div>
  );
};

ReviewServiceItems.propTypes = {
  header: PropTypes.string,
  serviceItemCards: ServiceItemCardsShape.isRequired,
  handleClose: PropTypes.func.isRequired,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
};

export default ReviewServiceItems;
