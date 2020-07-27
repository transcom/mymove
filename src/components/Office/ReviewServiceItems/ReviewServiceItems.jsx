import React, { useState, useEffect } from 'react';
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

const ReviewServiceItems = ({ header, serviceItemCards, handleClose, disableScrollIntoView }) => {
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

  let firstBasicIndex = null;
  let lastBasicIndex = null;
  // TODO - preset these based on existing values
  sortedCards.forEach((serviceItem, index) => {
    formValues[serviceItem.id] = {
      status: serviceItem.status,
      rejectionReason: serviceItem.rejectionReason,
    };

    // here we want to set the first and last index
    // of basic service items to know the bounds
    if (!serviceItem.shipmentType) {
      // no shipemntId, then it is a basic service items
      if (firstBasicIndex === null) {
        // if not set yet, set it the first time we see a basic
        // service item
        firstBasicIndex = index;
      }
      // keep setting the last basic index until the last one
      lastBasicIndex = index;
    }

    // requestedSum += serviceItem.amount; // TODO - use in Complete review screen
  });

  const currentCard = sortedCards[parseInt(curCardIndex, 10)];
  const isBasicServiceItem =
    firstBasicIndex !== null && curCardIndex >= firstBasicIndex && curCardIndex <= lastBasicIndex;

  // Similar to componentDidMount and componentDidUpdate
  useEffect(() => {
    const { id } = sortedCards[parseInt(curCardIndex, 10)];
    const element = document.querySelector(`#card-${id}`);
    // scroll into element view
    if (element && !disableScrollIntoView) {
      element.scrollIntoView();
    }
  });

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
                {
                  // render multiple basic service item cards
                  // otherwise, render only one card for shipment
                  isBasicServiceItem ? (
                    sortedCards.slice(firstBasicIndex, lastBasicIndex + 1).map((curCard) => (
                      <ServiceItemCard
                        key={`serviceItemCard_${curCard.id}`}
                        // eslint-disable-next-line react/jsx-props-no-spreading
                        {...curCard}
                        value={values[curCard.id]}
                        onChange={handleChange}
                        clearValues={clearServiceItemValues}
                      />
                    ))
                  ) : (
                    <ServiceItemCard
                      key={`serviceItemCard_${currentCard.id}`}
                      // eslint-disable-next-line react/jsx-props-no-spreading
                      {...currentCard}
                      value={values[currentCard.id]}
                      onChange={handleChange}
                      clearValues={clearServiceItemValues}
                    />
                  )
                }
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
  disableScrollIntoView: PropTypes.bool,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
  disableScrollIntoView: false,
};

export default ReviewServiceItems;
