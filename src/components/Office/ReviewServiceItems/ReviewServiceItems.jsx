import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import sortServiceItemsByGroup from '../../../utils/serviceItems';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import ServiceItemCard from 'components/Office/ReviewServiceItems/ServiceItemCard';

const ReviewServiceItems = ({ header, serviceItemCards, handleClose }) => {
  const [curCardIndex, setCardIndex] = useState(0);
  // eslint-disable-next-line no-unused-vars
  const [totalApproved, setTotalApproved] = useState(0);
  const [sortedCards] = useState(sortServiceItemsByGroup(serviceItemCards));
  const totalCards = serviceItemCards.length;

  const handleClick = (index) => {
    setCardIndex(index);
  };

  const formValues = {};

  let firstBasicIndex = null;
  let lastBasicIndex = null;
  // TODO - preset these based on existing values
  sortedCards.forEach((serviceItem, index) => {
    formValues[serviceItem.id] = {
      status: serviceItem.status,
      rejectionReason: undefined,
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
  });

  const currentCard = sortedCards[parseInt(curCardIndex, 10)];
  const isBasicServiceItem =
    firstBasicIndex !== null && curCardIndex >= firstBasicIndex && curCardIndex <= lastBasicIndex;

  // Similar to componentDidMount and componentDidUpdate
  useEffect(() => {
    const { id } = sortedCards[parseInt(curCardIndex, 10)];
    const element = document.querySelector(`#card-${id}`);
    // scroll into element view
    if (element) {
      element.scrollIntoView();
    }
  });

  return (
    <div data-cy="ReviewServiceItems" className={styles.ReviewServiceItems}>
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
                <Button data-cy="closeSidebar" type="button" onClick={handleClose} unstyled>
                  <XLightIcon />
                </Button>
                <div data-cy="itemCount" className={styles.eyebrowTitle}>
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
                  data-cy="prevServiceItem"
                  type="button"
                  onClick={() => handleClick(curCardIndex - 1)}
                  secondary
                  disabled={curCardIndex === 0}
                >
                  Previous
                </Button>
                <Button
                  data-cy="nextServiceItem"
                  type="button"
                  onClick={() => handleClick(curCardIndex + 1)}
                  disabled={curCardIndex + 1 === totalCards}
                >
                  Next
                </Button>
                <div className={styles.totalApproved}>
                  <div className={styles.totalLabel}>Total approved</div>
                  <div className={styles.totalAmount}>${totalApproved.toFixed(2)}</div>
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
