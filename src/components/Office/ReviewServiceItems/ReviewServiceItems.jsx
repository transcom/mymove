import React, { useState } from 'react';
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
  // TODO - preset these based on existing values
  serviceItemCards.forEach((serviceItem) => {
    formValues[serviceItem.id] = {
      status: serviceItem.status,
      rejectionReason: undefined,
    };
  });

  const currentCard = sortedCards[parseInt(curCardIndex, 10)];

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
