import { Button, Checkbox, FormGroup, Label } from '@trussworks/react-uswds';
import { Field, Formik } from 'formik';
import React from 'react';
import PropTypes from 'prop-types';

import { Form } from '../../form/Form';
import formStyles from '../../../styles/form.module.scss';
import { ErrorMessage } from '../../form/ErrorMessage';
import SectionWrapper from '../../Customer/SectionWrapper';
import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import styles from '../../../pages/PrimeUI/CreatePaymentRequest/CreatePaymentRequest.module.scss';
import Hint from '../../Hint';
import { ShipmentShape } from '../../../types/shipment';
import { MTOServiceItemShape } from '../../../types';
import ServiceItem from '../ServiceItem/ServiceItem';
import Shipment from '../Shipment/Shipment';

const CreatePaymentRequestForm = ({
  initialValues,
  onSubmit,
  handleSelectAll,
  createPaymentRequestSchema,
  mtoShipments,
  groupedServiceItems,
}) => (
  <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={createPaymentRequestSchema}>
    {({ isValid, isSubmitting, errors, values, setValues }) => (
      <Form className={formStyles.form}>
        <FormGroup error={errors != null && Object.keys(errors).length > 0}>
          {errors != null && errors.serviceItems && (
            <ErrorMessage display>At least 1 service item must be added when creating a payment request</ErrorMessage>
          )}
          <SectionWrapper className={formStyles.formSection}>
            <dl className={descriptionListStyles.descriptionList}>
              <h2>Move Service Items</h2>
              {groupedServiceItems.basic?.map((mtoServiceItem) => {
                return (
                  <SectionWrapper key={`moveServiceItems${mtoServiceItem.id}`} className={formStyles.formSection}>
                    <div className={styles.serviceItemInputGroup}>
                      <Label htmlFor={mtoServiceItem.id}>Add to payment request</Label>
                      <Field
                        as={Checkbox}
                        type="checkbox"
                        name="serviceItems"
                        value={mtoServiceItem.id}
                        id={mtoServiceItem.id}
                      />
                    </div>
                    <ServiceItem serviceItem={mtoServiceItem} />
                  </SectionWrapper>
                );
              })}
            </dl>
          </SectionWrapper>
          <SectionWrapper className={formStyles.formSection}>
            <dl className={descriptionListStyles.descriptionList}>
              <h2>Shipments</h2>
              {mtoShipments?.map((mtoShipment) => {
                return (
                  <div key={mtoShipment.id}>
                    <Shipment shipment={mtoShipment} />
                    <h2>Shipment Service Items</h2>
                    <Checkbox
                      id={`selectAll-${mtoShipment.id}`}
                      name={`selectAll-${mtoShipment.id}`}
                      label="Add all service items"
                      onClick={(event) => handleSelectAll(mtoShipment.id, values, setValues, event)}
                    />
                    {groupedServiceItems[mtoShipment.id]?.map((mtoServiceItem) => {
                      return (
                        <SectionWrapper
                          key={`shipmentServiceItems${mtoServiceItem.id}`}
                          className={formStyles.formSection}
                        >
                          <div className={styles.serviceItemInputGroup}>
                            <Label htmlFor={mtoServiceItem.id}>Add to payment request</Label>
                            <Field
                              as={Checkbox}
                              type="checkbox"
                              name="serviceItems"
                              value={mtoServiceItem.id}
                              id={mtoServiceItem.id}
                            />
                          </div>
                          <ServiceItem serviceItem={mtoServiceItem} />
                        </SectionWrapper>
                      );
                    })}
                  </div>
                );
              })}
            </dl>
            <Button
              aria-label="Submit Payment Request"
              type="submit"
              disabled={values.serviceItems?.length === 0 || isSubmitting || !isValid}
            >
              Submit Payment Request
            </Button>
            <Hint>
              At least one basic service item or shipment service item is required to create a payment request
            </Hint>
          </SectionWrapper>
        </FormGroup>
      </Form>
    )}
  </Formik>
);

CreatePaymentRequestForm.propTypes = {
  initialValues: PropTypes.shape({
    serviceItems: PropTypes.arrayOf(PropTypes.string),
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  handleSelectAll: PropTypes.func.isRequired,
  createPaymentRequestSchema: PropTypes.shape({
    serviceItems: PropTypes.node,
  }).isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  groupedServiceItems: PropTypes.shape({
    basic: PropTypes.arrayOf(MTOServiceItemShape),
  }),
};

CreatePaymentRequestForm.defaultProps = {
  mtoShipments: undefined,
  groupedServiceItems: undefined,
};

export default CreatePaymentRequestForm;
